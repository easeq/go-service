package kvstore

import (
	"context"
	"time"
)

// Option for initialization of the store
type Option interface{}

// SetOpts are the additional options passed during the set operation
type SetOpt interface{}

// GetOpts are the additional options passed during the get opertation
type GetOpt interface{}

// SubscribeHandler for the subscribe action
type SubscribeHandler interface {
	// Handle the subscription for the given key
	Handle(key string, args ...interface{}) error
}

// TxnHandler is the interface for handling transactions
type TxnHandler interface {
	// Handle the transaction
	Handle(ctx context.Context, store KVStore) error
}

// KVStore is a key-value data storage interface
type KVStore interface {
	// Init initializes the store with the given options
	Init(opts ...Option) error
	// Put the value for the key
	Put(ctx context.Context, record *Record, opts ...SetOpt) (*Record, error)
	// Get the value for the key
	Get(ctx context.Context, key string, opts ...GetOpt) ([]*Record, error)
	// Delete the key from the store
	Delete(ctx context.Context, key string) error
	// Txn handles transactions
	Txn(ctx context.Context, handler TxnHandler) error
	// Subscribe to the changes made to the given key
	Subscribe(ctx context.Context, key string, handler SubscribeHandler) error
	// Unsubscribe from a subscription
	Unsubscribe(ctx context.Context, key string) error
	// Close the store
	Close() error
	// String returns the name of the store implementation
	String() string
}

// Record set in the store for a specific key
type Record struct {
	// Key to store the record
	Key string `json:"key"`
	// Value for the key set in the store
	Value []byte `json:"value"`
	// Expiry is the time to expire a record
	Expiry time.Duration `json:"expiry,omitempty"`
	// Metadata is the metadata of the record
	Metadata map[string]interface{} `json:"metadata"`
}
