package pool

import "errors"

var (
	ErrConnectionClosed = errors.New("Connections closed")
)

type Pool interface {
	// Init(factory Factory) error
	// Get connection
	Get(address string) (Connection, error)
	// Close the pool
	Close() error
}

type Connection interface {
	// Close connection
	Close() error
}

// Creates and returns a new connection
type Factory func(args ...interface{}) (Connection, error)
