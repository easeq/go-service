package grpc

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/easeq/go-service/client"
	"github.com/easeq/go-service/pool"
	"github.com/easeq/go-service/registry"
	"google.golang.org/grpc"
)

const (
	// Default registry connection scheme
	defaultScheme = "http"
)

var (
	// ErrInvalidAddress returned when the address provided is invalid
	ErrInvalidAddress = errors.New("invalid service address")
	// ErrInvalidRegistry returned when the registry provided does not implement regsitry.ServiceRegsitry
	ErrInvalidRegistry = errors.New("registry provided is invalid")
	// ErrTooFewArgs returned when the args provided is less the args required
	ErrTooFewArgs = errors.New("too few arguments. Required 3")
	// ErrTooFewFactoryArgs returned when the args provided is less the args required
	ErrTooFewFactoryArgs = errors.New("too few arguments for the factory. Required 1 address")
	// ErrInvalidDialOptions returned when the dial options provided are not valid
	ErrInvalidDialOptions = errors.New("dial options provided are invalid")
	// ErrInvalidGrpcClient returned when type assertion to GrpcClient fails
	ErrInvalidGrpcClient = errors.New("invalid GrpcClient")
	// ErrInvalidStreamDescription returned when the variable passed is not of grpc.StreamDesc type
	ErrInvalidStreamDescription = errors.New("invalid stream description")
)

// ClientOption to pass as arg while creating new service
type ClientOption func(*Grpc)

// Grpc client that holds the reference to the pool,
// along with other configuration required to create the pool
// It holds a reference to the service registry used by the service.
type Grpc struct {
	pool      *pool.ConnectionPool
	factory   pool.Factory
	closeFunc pool.CloseFunc
	Registry  registry.ServiceRegistry
	sync.RWMutex
}

// NewGrpc creates a new gRPC client
func NewGrpc(opts ...ClientOption) *Grpc {
	c := new(Grpc)

	for _, opt := range opts {
		opt(c)
	}

	c.pool = pool.NewPool(
		pool.WithFactory(c.factory),
		pool.WithSize(10),
		pool.WithCloseFunc(c.closeFunc),
	)

	return c
}

// WithRegistry passes services registry externally
func WithRegistry(registry registry.ServiceRegistry) ClientOption {
	return func(c *Grpc) {
		c.Registry = registry
	}
}

// WithFactory defines the client connection creation factory
func WithFactory(factory pool.Factory) ClientOption {
	return func(c *Grpc) {
		c.factory = factory
	}
}

// WithCloseFunc passes the callback function to close the gRPC connection in the pool
func WithCloseFunc(closeFunc pool.CloseFunc) ClientOption {
	return func(c *Grpc) {
		c.closeFunc = closeFunc
	}
}

// Dial creates/gets a connection from the pool using the address from the service registry
func (c *Grpc) Dial(name string, opts ...client.DialOption) (pool.Connection, error) {
	address := c.Registry.ConnectionString(name, defaultScheme)
	return c.pool.Get(address)
}

// Call gRPC method
func (c *Grpc) Call(
	ctx context.Context,
	sc client.ServiceClient,
	method string,
	req interface{},
	res interface{},
	opts ...client.CallOption,
) error {
	pcc, err := c.Dial(sc.GetServiceName())
	if err != nil {
		return err
	}

	defer pcc.Close()

	cc, ok := pcc.Conn().(*grpc.ClientConn)
	if !ok {
		return fmt.Errorf("invalid factory client connection")
	}

	callOpts := make([]grpc.CallOption, len(opts))
	for i, opt := range opts {
		callOpts[i] = opt.(grpc.CallOption)
	}

	return cc.Invoke(ctx, method, req, res, callOpts...)
}

// Stream gRPC method
func (c *Grpc) Stream(
	ctx context.Context,
	sc client.ServiceClient,
	desc interface{},
	method string,
	req interface{},
	opts ...client.CallOption,
) (client.StreamClient, error) {
	pcc, err := c.Dial(sc.GetServiceName())
	if err != nil {
		return nil, err
	}

	cc, ok := pcc.Conn().(*grpc.ClientConn)
	if !ok {
		return nil, fmt.Errorf("invalid connection")
	}

	callOpts := make([]grpc.CallOption, len(opts))
	for i, opt := range opts {
		callOpts[i] = opt.(grpc.CallOption)
	}

	serviceDesc, ok := desc.(*grpc.StreamDesc)
	if !ok {
		return nil, ErrInvalidStreamDescription
	}

	stream, err := cc.NewStream(ctx, serviceDesc, method, callOpts...)
	if err != nil {
		return nil, err
	}

	gs := &GrpcStreamClient{stream, pcc}
	if req == nil {
		return gs, nil
	}

	if err := gs.stream.SendMsg(req); err != nil {
		return nil, err
	}

	return gs, nil
}

// GrpcStreamClient is the gRPC client that allows streaming. It holds the stream and the connection to the gRPC server.
type GrpcStreamClient struct {
	stream grpc.ClientStream
	conn   pool.Connection
}

// Recv receive a message from the stream
func (sc *GrpcStreamClient) Recv(res interface{}) error {
	if err := sc.stream.RecvMsg(res); err != nil {
		return err
	}

	return nil
}

// Send sends a message
func (sc *GrpcStreamClient) Send(req interface{}) error {
	return sc.stream.SendMsg(req)
}

// CloseAndRecv first close the server stream and receives messages on the client stream
func (sc *GrpcStreamClient) CloseAndRecv(res interface{}) error {
	if err := sc.stream.CloseSend(); err != nil {
		return err
	}

	if err := sc.stream.RecvMsg(res); err != nil {
		return err
	}

	return nil
}

// CloseConn closes the connection in the pool if the pool is full or adds it back to the pool
func (sc *GrpcStreamClient) CloseConn() error {
	return sc.conn.Close()
}
