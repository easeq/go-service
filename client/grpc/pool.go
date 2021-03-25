package grpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/easeq/go-service/pool"
	"github.com/easeq/go-service/registry"
	"google.golang.org/grpc"
)

const (
	// DefaultTTL is the default value client connection ttl
	DefaultTTL = time.Minute * 5
)

type Pool struct {
	connections map[string]*GrpcClient
	ttl         time.Duration
	Registry    registry.ServiceRegistry
	DialOptions []grpc.DialOption
	Scheme      string
}

// Option to pass as arg while creating new service
type Option func(*Pool)

func NewGrpcClientPool(opts ...Option) *Pool {
	return &Pool{
		ttl: DefaultTTL,
	}
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
func (p *Pool) Get(address string, opts ...interface{}) (pool.Connection, error) {
	if !p.isInitialized(address) {
		err := p.Init(address)
		if err != nil {
			return nil, err
		}
	}

	conn, ok := p.connections[address]
	if !ok {
		return nil, errors.New("Could not find the requested connection")
	}

	return conn, nil
}

// Init gRPC client
func (p *Pool) Init(opts ...interface{}) error {
	address := opts[0].(string)
	if address == "" {
		return errors.New("Invalid service address provided.")
	}

	if p.isInitialized(address) {
		// Refresh TTL timer if the connection is called before TTL deadline
		p.connections[address].timer.Reset(DefaultTTL)
		return nil
	}

	gc, err := NewGrpcClient(address, p.ttl, p.DialOptions...)
	if err != nil {
		return err
	}

	p.connections[address] = gc

	return nil
}

// Return whether client has been intialized
func (p *Pool) isInitialized(name string) bool {
	_, ok := p.connections[name]
	return ok
}

// Release - releases and existing connection
func (p *Pool) Release(name string) error {
	if !p.isInitialized(name) {
		return fmt.Errorf("Connection does not exist in the pool")
	}

	defer delete(p.connections, name)
	return p.connections[name].Close()
}
