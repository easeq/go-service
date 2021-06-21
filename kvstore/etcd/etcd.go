package etcd

import (
	"context"
	"errors"
	"log"

	"github.com/easeq/go-service/kvstore"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	// ErrInvalidLeaseID returned when the leaseID provided is invalid
	ErrInvalidLeaseID = errors.New("invalid etcd leaseID passed")
)

// Etcd holds our etcd instance
type Etcd struct {
	*clientv3.Client
	*Config
}

// NewEtcd returns a new instance of etcd with etcd client and config
func NewEtcd(config *Config) *Etcd {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEndpoints(),
		DialTimeout: config.DialTimeout,
	})

	if err != nil {
		panic("Error creating an kvstore etcd client")
	}

	return &Etcd{client, config}
}

// Init initializes the store with the given options
func (e *Etcd) Init(opts ...kvstore.Option) error {
	log.Printf("Unsupported method %s Init", e.String())
	return nil
}

// GetMetadataLeaseID returns the leaseID from the record metadata
func (e *Etcd) GetMetadataLeaseID(record kvstore.Record) (clientv3.LeaseID, error) {
	lID, ok := record.Metadata["lease_id"]
	if !ok {
		return 0, nil
	}

	leaseID, ok := lID.(clientv3.LeaseID)
	if !ok {
		return 0, ErrInvalidLeaseID
	}

	return leaseID, nil
}

// LeaseID returns the leaseID (if any) to be used by the record
// If exists, it returns the "lease_id" set in the record metadata
// If record expiry is set, then it creates a new leaseID and returns it.
// If it's none of the above, then it returns nil
func (e *Etcd) LeaseID(ctx context.Context, record kvstore.Record) (clientv3.LeaseID, error) {
	leaseID, err := e.GetMetadataLeaseID(record)
	if err != nil {
		return 0, err
	}

	// Create a new lease and return
	if leaseID == 0 && record.Expiry != 0 {
		l, err := e.Client.Lease.Grant(ctx, int64(record.Expiry))
		if err != nil {
			return 0, err
		}

		return l.ID, nil
	}

	// The record doesn't have an expiry
	return 0, nil
}

// RenewLease renews the lease with the given leaseID
// This renews lease if the lease is valid and not 0
func (e *Etcd) RenewLease(ctx context.Context, record kvstore.Record) error {
	leaseID, err := e.GetMetadataLeaseID(record)
	if err != nil {
		return err
	}

	if leaseID == 0 {
		return nil
	}

	if _, err := e.Client.Lease.KeepAliveOnce(ctx, leaseID); err != nil {
		return err
	}

	return nil
}

// Put adds the record into the store
func (e *Etcd) Put(ctx context.Context, record kvstore.Record, opts ...kvstore.SetOpt) (*kvstore.Record, error) {
	leaseID, err := e.LeaseID(ctx, record)
	if err != nil {
		return nil, err
	}

	if err := e.RenewLease(ctx, record); err != nil {
		return nil, err
	}

	// Pass etcd PUT options
	putOpts := []clientv3.OpOption{}
	if leaseID != 0 {
		putOpts = append(putOpts, clientv3.WithLease(leaseID))
	}

	if _, err := e.Client.Put(ctx, record.Key, string(record.Value), putOpts...); err != nil {
		return nil, err
	}

	// Set the leaseID created/renewed
	record.Metadata["lease_id"] = leaseID

	return &record, nil
}

// Get a record by it's key
func (e *Etcd) Get(ctx context.Context, key string, opts ...kvstore.GetOpt) ([]*kvstore.Record, error) {
	response, err := e.Client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	records := make([]*kvstore.Record, response.Count)
	for i, r := range response.Kvs {
		records[i] = &kvstore.Record{
			Key:   string(r.Key),
			Value: r.Value,
			Metadata: map[string]interface{}{
				"lease_id": clientv3.LeaseID(r.Lease),
				"count":    response.Count,
				"header":   response.Header,
				"more":     response.More,
			},
		}
	}

	return records, nil
}

// Delete the key from the store
func (e *Etcd) Delete(ctx context.Context, key string) error {
	_, err := e.Client.Delete(ctx, key)
	return err
}

// Txn handles store transactions
func (e *Etcd) Txn(ctx context.Context, handler kvstore.TxnHandler) error {
	return handler.Handle(ctx, e)
}

// Subscribe to the changes made to the given key
func (e *Etcd) Subscribe(ctx context.Context, key string, handler kvstore.Handler) error {
	cWatch := e.Client.Watch(ctx, key)
	log.Printf("set WATCH on %s", key)

	for {
		watchResp, ok := <-cWatch
		if !ok {
			return nil
		}

		err := handler.Handle(key, watchResp)
		if err != nil {
			return err
		}
	}
}

// Unsubscribe from a subscription
func (e *Etcd) Unsubscribe(ctx context.Context, key string) error {
	log.Printf("Unsupported method %s Unsubscribe", e.String())
	return nil
}

// Close the etcd client
func (e *Etcd) Close() error {
	e.Client.Watcher.Close()
	return e.Client.Close()
}

// String returns the name of the store implementation
func (e *Etcd) String() string {
	return "kvstore-etcd"
}
