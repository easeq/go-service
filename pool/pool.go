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

type FactoryConn interface{}

// Creates and returns a new connection
type Factory func(address string) (FactoryConn, error)

// CloseFunc to close the connection in the pool
type CloseFunc func(conn interface{}) error
