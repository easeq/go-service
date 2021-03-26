package grpc

import (
	"errors"
	"time"

	"github.com/easeq/go-service/pool"
	"github.com/easeq/go-service/registry"
	"google.golang.org/grpc"
)

const (
	// DefaultTTL is the default value client connection ttl
	DefaultTTL = time.Minute
)

type Pool struct {
	// pool.Pool
	connections map[string]GrpcClient
	ttl         time.Duration
	Registry    registry.ServiceRegistry
	DialOptions []grpc.DialOption
	Scheme      string
}

// Option to pass as arg while creating new service
type Option func(*Pool)

func NewGrpcClientPool(opts ...Option) *Pool {
	p := &Pool{
		connections: make(map[string]GrpcClient),
		ttl:         DefaultTTL,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// WithTTL defines a ttl for the pool
func WithTTL(ttl time.Duration) Option {
	return func(p *Pool) {
		p.ttl = ttl
	}
}

// WithRegistry defines the registry where the services are registered
func WithRegistry(registry registry.ServiceRegistry) Option {
	return func(p *Pool) {
		p.Registry = registry
	}
}

// WithDialOptions defines the gRPC dial options
func WithDialOptions(opts ...grpc.DialOption) Option {
	return func(p *Pool) {
		p.DialOptions = opts
	}
}

// WithScheme defines the gRPC service discovery scheme
func WithScheme(scheme string) Option {
	return func(p *Pool) {
		p.Scheme = scheme
	}
}

// Get creates or returns an existing gRPC client connection
func (p *Pool) Get(client pool.Connection, address string) error {
	conn, ok := p.connections[address]
	if !ok {
		return errors.New("Could not find the requested connection")
	}

	if val, ok := client.(*GrpcClient); ok {
		*val = conn
	}

	return nil
}

// Init gRPC client
func (p *Pool) Init(address string, opts ...interface{}) error {
	// Check whether a connection exists in pool
	if _, ok := p.connections[address]; ok {
		// Refresh TTL timer if the connection is called before TTL deadline
		p.connections[address].timer.Reset(DefaultTTL)
		return nil
	}

	gc, err := NewGrpcClient(p, address, p.ttl, p.DialOptions...)
	if err != nil {
		return err
	}

	p.connections[address] = *gc

	return nil
}

// Release - releases and existing connection
func (p *Pool) Release(name string) error {
	delete(p.connections, name)
	return nil
}
