package pool

import "errors"

var (
	ErrConnectionClosed = errors.New("connections closed")
)

type Pool interface {
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
type Factory func(address string, opts ...interface{}) (Connection, error)
