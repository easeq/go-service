package pool

type Pool interface {
	Init(opts ...interface{}) error
	// Get connection
	Get(address string, opts ...interface{}) (Connection, error)
	// Close the pool
	Release(address string) error
}

type Connection interface {
	// Close connection
	Close() error
}
