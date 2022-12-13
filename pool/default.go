package pool

import (
	"errors"
	"sync"

	"github.com/easeq/go-service/logger"
)

var (
	// ErrGroupNotExist returned when the group hasn't been created
	ErrGroupNotExist = errors.New("group doesn't exist")
	// ErrConnectionNotExists returned when the connection doesn't exist or has already been closed
	ErrConnectionNotExists = errors.New("connection doesn't exist")
	// ErrCouldNotAssignConnection when connection assignment fails
	ErrCouldNotAssignConnection = errors.New("assigning connection to gRPC client failed")
	// ErrInvalidConnectionPool returwhen the type assertion to ConnectionPool fails
	ErrInvalidConnectionPool = errors.New("invalid connection pool")
)

// Option to pass as arg while creating new service
type Option func(*ConnectionPool)

// ConnectionPool holds the connections in the pool,
// alongwith the factory to create a new connection,
// size of each type of connection in the pool,
// and the CloseFunc callback used to close a specific connection in the pool
type ConnectionPool struct {
	conns     map[string](chan FactoryConn)
	factory   Factory
	CloseFunc CloseFunc
	size      int
	logger    logger.Logger
	sync.RWMutex
}

// NewPool creates a new pool with size and factory
func NewPool(opts ...Option) *ConnectionPool {
	pool := &ConnectionPool{
		conns: make(map[string](chan FactoryConn)),
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

// WithCloseFunc passes the CloseFunc callback to the pool
func WithCloseFunc(closeFunc CloseFunc) Option {
	return func(p *ConnectionPool) {
		p.CloseFunc = closeFunc
	}
}

// WithLogger sets the logger for the pool
func WithLogger(logger logger.Logger) Option {
	return func(p *ConnectionPool) {
		p.logger = logger
	}
}

// Creates a new connection channel group
func (p *ConnectionPool) create(name string) (chan FactoryConn, Factory, error) {
	p.Lock()

	if p.conns == nil {
		return nil, nil, ErrConnectionClosed
	}

	p.conns[name] = make(chan FactoryConn, p.size)

	p.Unlock()

	return p.get(name)
}

// Get individual connection channel by name and the factory to create connection
func (p *ConnectionPool) get(name string) (chan FactoryConn, Factory, error) {
	p.RLock()
	defer p.RUnlock()

	group, ok := p.conns[name]
	if !ok {
		return nil, p.factory, ErrGroupNotExist
	}

	return group, p.factory, nil
}

// wraps a the connection provided in a standard Connection
func (p *ConnectionPool) wrap(address string, conn FactoryConn) Connection {
	return &ClientConn{
		p:       p,
		address: address,
		conn:    conn,
	}
}

// Get creates or returns an existing gRPC client connection
func (p *ConnectionPool) Get(address string) (Connection, error) {
	conns, factory, err := p.get(address)
	if err != nil {
		if conns, factory, err = p.create(address); err != nil {
			return nil, err
		}
	}

	if conns == nil {
		return nil, ErrConnectionClosed
	}

	select {
	case conn := <-conns:
		if conn == nil {
			return nil, ErrConnectionClosed
		}

		return p.wrap(address, conn), nil
	default:
		conn, err := factory(address)
		if err != nil {
			return nil, err
		}

		return p.wrap(address, conn), nil
	}
}

func (p *ConnectionPool) add(address string, conn FactoryConn) error {
	if conn == nil {
		return ErrConnectionNotExists
	}

	p.RLock()
	defer p.RUnlock()

	if p.conns == nil {
		return p.CloseFunc(conn)
	}

	// add the connection to the pool or close the connection if the pool is full.
	select {
	case p.conns[address] <- conn:
		return nil
	default:
		// Pool group full, close this connection
		return p.CloseFunc(conn)
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
			if err := p.CloseFunc(conn); err != nil {
				return err
			}
		}
	}

	return nil
}

// ClientConn holds the function created by the pool factory method.
// It also keeps the address, and a reference to the parent pool
type ClientConn struct {
	conn    FactoryConn
	address string
	p       *ConnectionPool
}

// Address returns the factory connection address
func (cc *ClientConn) Address() string {
	return cc.address
}

// Conn returns the factory connection
func (cc *ClientConn) Conn() FactoryConn {
	return cc.conn
}

// Close - closes the connection to the gRPC service server
func (cc *ClientConn) Close() error {
	return cc.p.add(cc.address, cc.conn)
}
