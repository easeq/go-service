package pool

type Pool interface {
	Init(address string, opts ...interface{}) error
	// Get connection
	Get(client Connection, address string) error
	// Close the pool
	Release(address string) error
}

type Connection interface {
	// Close connection
	Close() error
}
