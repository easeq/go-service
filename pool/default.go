package pool

import (
	"errors"
	"fmt"
	"sync"

	"google.golang.org/grpc"
)

var (
	// ErrGroupNotExist returned when the group hasn't been created
	ErrGroupNotExist = errors.New("group doesn't exist")
	// ErrConnectionNotExists returned when the connection doesn't exist or has already been closed
	ErrConnectionNotExists = errors.New("connection doesn't exist")
	// ErrCouldNotAssignConnection when connection assignment fails
	ErrCouldNotAssignConnection = errors.New("assigning connection to gRPC client failed")
)

// Option to pass as arg while creating new service
type Option func(*ConnectionPool)

type ConnectionPool struct {
	conns    map[string](chan Connection)
	factory  Factory
	size     int
	dialOpts []grpc.DialOption
	sync.RWMutex
}

// NewPool creates a new pool with size and factory
func NewPool(opts ...Option) *ConnectionPool {
	pool := &ConnectionPool{
		conns: make(map[string](chan Connection)),
		size:  10,
	}

	for _, opt := range opts {
		opt(pool)
	}

	if pool.factory == nil {
		panic("pool factory cannot be nil")
	}

	return pool
}

// WithSize defines the size of the pool
func WithSize(size int) Option {
	return func(p *ConnectionPool) {
		p.size = size
	}
}

// WithFactory defines the connection creation factory
func WithFactory(factory Factory) Option {
	return func(p *ConnectionPool) {
		p.factory = factory
	}
}

// WithDialOpts defines the grpc dial options
func WithDialOpts(opts ...grpc.DialOption) Option {
	return func(p *ConnectionPool) {
		p.dialOpts = opts
	}
}

// Creates a new connection channel group
func (p *ConnectionPool) create(name string) (chan Connection, Factory, error) {
	p.Lock()

	if p.conns == nil {
		return nil, nil, ErrConnectionClosed
	}

	p.conns[name] = make(chan Connection, p.size)

	p.Unlock()

	return p.get(name)
}

// Get individual connection channel by name and the factory to create connection
func (p *ConnectionPool) get(name string) (chan Connection, Factory, error) {
	p.RLock()
	defer p.RUnlock()

	group, ok := p.conns[name]
	if !ok {
		return nil, p.factory, ErrGroupNotExist
	}

	return group, p.factory, nil
}

// wraps a the connection provided in a standard Connection
func (p *ConnectionPool) wrap(address string, conn Connection) Connection {
	// gc := &GrpcClient{
	// 	p:       p,
	// 	address: address,
	// }
	// gc.Connection = conn
	return nil
}

// Get creates or returns an existing gRPC client connection
func (p *ConnectionPool) Get(address string, opts ...interface{}) (Connection, error) {
	conns, factory, err := p.get(address)
	if err != nil {
		if conns, factory, err = p.create(address); err != nil {
			return nil, err
		}
	}

	if conns == nil {
		return nil, fmt.Errorf("%s", ErrConnectionClosed)
	}

	select {
	case conn := <-conns:
		if conn == nil {
			return nil, fmt.Errorf("%s", ErrConnectionClosed)
		}

		return p.wrap(address, conn), nil
	default:
		conn, err := factory(address, opts...)
		if err != nil {
			return nil, err
		}

		return p.wrap(address, conn), nil
	}
}

func (p *ConnectionPool) add(address string, conn Connection) error {
	if conn == nil {
		return ErrConnectionNotExists
	}

	p.RLock()
	defer p.RUnlock()

	if p.conns == nil {
		return conn.Close()
	}

	// add the connection to the pool or close the connection if the pool is full.
	select {
	case p.conns[address] <- conn:
		return nil
	default:
		// Pool group full, close this connection
		return conn.Close()
	}
}

// Close - closes the connection pool and all it's channels
func (p *ConnectionPool) Close() error {
	p.Lock()
	defer p.Unlock()

	conns := p.conns
	p.conns = nil
	p.factory = nil

	if conns == nil {
		return ErrConnectionClosed
	}

	for key, group := range conns {
		if group == nil {
			continue
		}

		// Close channel
		close(conns[key])

		for conn := range group {
			// Close connection
			if err := conn.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}
