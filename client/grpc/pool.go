package grpc

import (
	"errors"
	"sync"
	"time"

	"github.com/easeq/go-service/pool"
)

const (
	// DefaultTTL is the default value client connection ttl
	DefaultTTL = time.Minute
)

var (
	// ErrGroupNotExist returned when the group hasn't been created
	ErrGroupNotExist = errors.New("Group doesn't exist")
	// ErrConnectionNotExists returned when the connection doesn't exist or has already been closed
	ErrConnectionNotExists = errors.New("Connection doesn't exist")
)

type Pool struct {
	conns   map[string](chan pool.Connection)
	factory pool.Factory
	size    int
	*sync.RWMutex
}

// NewGrpcClientPool - initiates a new grpc client pool
func NewGrpcClientPool(size int, factory pool.Factory) *Pool {
	return &Pool{
		conns:   make(map[string](chan pool.Connection)),
		factory: factory,
		size:    size,
	}
}

// Creates a new connection channel group
func (p *Pool) create(name string) error {
	p.Lock()
	defer p.Unlock()

	if p.conns == nil {
		return pool.ErrConnectionClosed
	}

	p.conns[name] = make(chan pool.Connection, p.size)
	return nil
}

// Get individual connection channel by name and the factory to create connection
func (p *Pool) get(name string) (chan pool.Connection, pool.Factory, error) {
	p.RLock()
	defer p.RLock()

	group, ok := p.conns[name]
	if !ok {
		return nil, p.factory, ErrGroupNotExist
	}

	return group, p.factory, nil
}

// wraps a the connection provided in a standard pool.Connection
func (p *Pool) wrap(address string, conn pool.Connection) pool.Connection {
	gc := &GrpcClient{
		p:       p,
		address: address,
	}
	gc.Connection = conn
	return gc
}

// Get creates or returns an existing gRPC client connection
func (p *Pool) Get(address string) (pool.Connection, error) {
	conns, factory, err := p.get(address)
	if err != nil {
		if err := p.create(address); err != nil {
			return nil, err
		}
	}

	if conns == nil {
		return nil, pool.ErrConnectionClosed
	}

	select {
	case conn := <-conns:
		if conn == nil {
			return nil, pool.ErrConnectionClosed
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

func (p *Pool) add(address string, conn pool.Connection) error {
	if conn == nil {
		return ErrConnectionNotExists
	}

	p.RLock()
	defer p.Unlock()

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
func (p *Pool) Close() error {
	p.Lock()
	defer p.Unlock()

	conns := p.conns
	p.conns = nil
	p.factory = nil

	if conns == nil {
		return pool.ErrConnectionClosed
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
