package pool

import "errors"

var (
	// ErrConnectionClosed is returned when the pool is closed
	ErrConnectionClosed = errors.New("connections closed")
)

// Pool interface to create new pool implementations
type Pool interface {
	// Get connection
	Get(address string) (Connection, error)
	// Close the pool
	Close() error
}

// Connection interface is the connection saved in the pool
type Connection interface {
	// Address returns the connection address
	Address() string
	// Conn returns the factory connection
	Conn() FactoryConn
	// Close connection
	Close() error
}

// FactoryConn interface is the connection returned by the factory passed to the pool
type FactoryConn interface{}

// Factory callback creates and returns a new connection
type Factory func(address string) (FactoryConn, error)

// CloseFunc to close the connection in the pool
type CloseFunc func(conn interface{}) error
