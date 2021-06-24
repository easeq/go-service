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
	// ErrNoResults returned when no results are found
	ErrNoResults = errors.New("no results found")
	// ErrCreatingEtcdClient returned when creating etcd clientv3 fails
	ErrCreatingEtcdClient = errors.New("error creating kvstore etcd client")
)

const (
	// KEY_LEASE_ID points to the lease ID in the etcd record metadata
	KEY_LEASE_ID = "lease_id"
	// KEY_COUNT points to the count in the etcd record metadata
	KEY_COUNT = "count"
	// KEY_HEADER points to the header in the etcd record metadata
	KEY_HEADER = "header"
	// KEY_MORE points to the more boolean var in the etcd record metadata
	KEY_MORE = "more"
	// KEY_REVISION points to the revision in the etcd record metadata
	KEY_REVISION = "revision"
	// KEY_VERSION points to the version in the etcd record metadata
	KEY_VERSION = "version"
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
		panic(ErrCreatingEtcdClient)
	}

	return &Etcd{client, config}
}

// Init initializes the store with the given options
func (e *Etcd) Init(opts ...kvstore.Option) error {
	log.Printf("Unsupported method %s Init", e.String())
	return nil
}

// GetMetadataLeaseID returns the leaseID from the record metadata
func (e *Etcd) GetMetadataLeaseID(record *kvstore.Record) (clientv3.LeaseID, error) {
	lID, ok := record.Metadata[KEY_LEASE_ID]
	if !ok {
		// LeaseID is not a required field, so no error
		return 0, nil
	}

	leaseID, ok := lID.(clientv3.LeaseID)
	if !ok {
		return 0, ErrInvalidLeaseID
	}

	return leaseID, nil
}

// LeaseID returns the leaseID (if any) to be used by the record
// If exists, it renews and returns the "lease_id" set in the record metadata
// If record expiry is set, then it creates a new leaseID and returns it.
// If it's none of the above, then it returns 0
func (e *Etcd) LeaseID(ctx context.Context, record *kvstore.Record) (clientv3.LeaseID, error) {
	leaseID, err := e.GetMetadataLeaseID(record)
	if err != nil {
		return 0, err
	}

	// Renew and use existing lease
	if leaseID != 0 {
		if err := e.RenewLease(ctx, leaseID); err != nil {
			return 0, err
		}

		return leaseID, nil
	}

	// Create a new lease and return
	if record.Expiry != 0 {
		l, err := e.Client.Lease.Grant(ctx, int64(record.Expiry.Seconds()))
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
func (e *Etcd) RenewLease(ctx context.Context, leaseID clientv3.LeaseID) error {
	if leaseID == 0 {
		return nil
	}

	if _, err := e.Client.Lease.KeepAliveOnce(ctx, leaseID); err != nil {
		return err
	}

	return nil
}

// Put adds the record into the store
// Get the lease if lease_id is defined in the record metadata, or create new lease if expiry is defined
// Renew lease using the lease_id in the record metadata, added/used by LeaseID(...)
// Add the record to the store with the lease_id
func (e *Etcd) Put(ctx context.Context, record *kvstore.Record, opts ...kvstore.SetOpt) (*kvstore.Record, error) {
	leaseID, err := e.LeaseID(ctx, record)
	if err != nil {
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

	if record.Metadata == nil {
		record.Metadata = make(map[string]interface{})
	}

	// Set the leaseID created/renewed
	record.Metadata[KEY_LEASE_ID] = leaseID

	return record, nil
}

// Get a record by it's key
func (e *Etcd) Get(ctx context.Context, key string, opts ...kvstore.GetOpt) ([]*kvstore.Record, error) {
	response, err := e.Client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if response.Count == 0 {
		return nil, ErrNoResults
	}

	records := make([]*kvstore.Record, response.Count)
	for i, r := range response.Kvs {
		records[i] = &kvstore.Record{
			Key:   string(r.Key),
			Value: r.Value,
			Metadata: map[string]interface{}{
				KEY_LEASE_ID: clientv3.LeaseID(r.Lease),
				KEY_COUNT:    response.Count,
				KEY_HEADER:   response.Header,
				KEY_MORE:     response.More,
				KEY_REVISION: r.ModRevision,
				KEY_VERSION:  r.Version,
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
func (e *Etcd) Subscribe(ctx context.Context, key string, handler kvstore.SubscribeHandler) error {
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
