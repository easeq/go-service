package jetstream

import (
	"context"
	"fmt"

	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/kvstore"
	"github.com/easeq/go-service/logger"
	"github.com/easeq/go-service/tracer"
	"github.com/nats-io/nats.go"
)

const (
	KEY_REVISION = "revision"
	KEY_OP       = "operation"
	KEY_CREATED  = "created"
	KEY_DELTA    = "delta"
	KEY_BUCKET   = "bucket"
)

// JetStream holds our jetstream instance
type JetStream struct {
	i       component.Initializer
	logger  logger.Logger
	tracer  tracer.Tracer
	wrapper *kvstore.Wrapper
	nc      *nats.Conn
	jsCtx   nats.JetStreamContext
	kv      nats.KeyValue
	Config  *Config
}

// NewJetStream returns a new instance of jetstream with a pre-defined bucket and config
func NewJetStream() *JetStream {
	config := NewConfig()
	nc, err := nats.Connect(config.Address())
	if err != nil {
		panic("error connecting to nats server")
	}

	js, err := nc.JetStream()
	if err != nil {
		panic("error creating JetStreamContext")
	}

	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: config.Bucket,
	})
	if err != nil {
		panic("error creating jetstream bucket")
	}

	j := &JetStream{nc: nc, jsCtx: js, kv: kv, Config: config}
	j.i = NewInitializer(j)
	j.wrapper = kvstore.NewWrapper(j)

	return j
}

// Init initializes the store with the given options
func (j *JetStream) Init(opts ...kvstore.Option) error {
	j.logger.Infof("Unsupported method %s Init", j.String())
	return nil
}

// Put adds the record into the store
func (j *JetStream) Put(ctx context.Context, record *kvstore.Record, opts ...kvstore.SetOpt) (*kvstore.Record, error) {
	cb := func(ctx context.Context, record *kvstore.Record, opts ...kvstore.SetOpt) (*kvstore.Record, error) {
		j.kv.Put(record.Key, record.Value)
		return record, nil
	}

	return j.wrapper.Put(ctx, record, cb, opts...)
}

// Get a record by it's key
func (j *JetStream) Get(ctx context.Context, key string, opts ...kvstore.GetOpt) ([]*kvstore.Record, error) {
	cb := func(ctx context.Context, key string, opts ...kvstore.GetOpt) ([]*kvstore.Record, error) {
		entry, err := j.kv.Get(key)
		if err != nil {
			return nil, err
		}

		return []*kvstore.Record{
			{
				Key:   entry.Key(),
				Value: entry.Value(),
				Metadata: map[string]interface{}{
					KEY_BUCKET:   entry.Bucket(),
					KEY_DELTA:    entry.Delta(),
					KEY_OP:       entry.Operation(),
					KEY_CREATED:  entry.Created(),
					KEY_REVISION: entry.Revision(),
				},
			},
		}, nil
	}

	return j.wrapper.Get(ctx, key, cb, opts...)
}

// Delete the key from the store
func (j *JetStream) Delete(ctx context.Context, key string) error {
	cb := func(ctx context.Context, key string) error {
		return j.kv.Delete(key)
	}

	return j.wrapper.Delete(ctx, key, cb)
}

// Txn handles store transactions (there are no transactions in jetstream
// buckets, manually handle the txn here)
func (j *JetStream) Txn(ctx context.Context, handler kvstore.TxnHandler) error {
	cb := func(ctx context.Context, handler kvstore.TxnHandler) error {
		return handler.Handle(ctx, j)
	}

	return j.wrapper.Txn(ctx, handler, cb)
}

// Subscribe to the changes made to the given key
func (j *JetStream) Subscribe(ctx context.Context, key string, handler kvstore.SubscribeHandler) error {
	cb := func(ctx context.Context, key string, handler kvstore.SubscribeHandler) error {
		watcher, err := j.kv.Watch(key)
		if err != nil {
			return fmt.Errorf("[jetstream] watch (%s) error: %v", key, err)
		}

		for {
			watchResp, ok := <-watcher.Updates()
			if !ok {
				return nil
			}

			done, err := handler.Handle(key, watchResp, watcher)
			if err != nil {
				return err
			}

			if done {
				return watcher.Stop()
			}
		}
	}

	return j.wrapper.Subscribe(ctx, key, handler, cb)
}

// Unsubscribe from a subscription
func (j *JetStream) Unsubscribe(ctx context.Context, key string) error {
	j.logger.Infof("Unsupported method %s Unsubscribe", j.String())
	return nil
}

// String returns the name of the store implementation
func (j *JetStream) String() string {
	return "kvstore-jetstream"
}

func (j *JetStream) HasInitializer() bool {
	return true
}

func (j *JetStream) Initializer() component.Initializer {
	return j.i
}
