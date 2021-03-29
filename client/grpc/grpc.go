package grpc

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/easeq/go-service/pool"
	"github.com/easeq/go-service/registry"
	"github.com/easeq/go-service/server"
	"google.golang.org/grpc"
)

const (
	// Default registry connection scheme
	defaultScheme = "http"
)

var (
	// ErrInvalidAddress returned when the address provided is invalid
	ErrInvalidAddress = errors.New("Invalid service address")
	// ErrInvalidRegistry returned when the registry provided does not implement regsitry.ServiceRegsitry
	ErrInvalidRegistry = errors.New("Registry provided is invalid")
	// ErrTooFewArgs returned when the args provided is less the args required
	ErrTooFewArgs = errors.New("Too few arguments. Required 3")
	// ErrTooFewArgs returned when the args provided is less the args required
	ErrTooFewFactoryArgs = errors.New("Too few arguments for the factory. Required 1 address")
	// ErrInvalidDialOptions returned when the dial options provided are not valid
	ErrInvalidDialOptions = errors.New("Dial options provided are invalid")
	// ErrInvalidGrpcClient returned when type assertion to GrpcClient fails
	ErrInvalidGrpcClient = errors.New("Invalid GrpcClient")
)

type GrpcClient struct {
	pool.Connection
	address string
	p       *Pool
	sync.RWMutex
}

// NewGrpcClient creates and rea new gRPC client connection
func NewGrpcClientConn(args ...interface{}) (pool.Factory, error) {
	if len(args) < 2 {
		return nil, ErrTooFewArgs
	}

	registry, ok := args[0].(registry.ServiceRegistry)
	if !ok {
		return nil, ErrInvalidRegistry
	}

	scheme, ok := args[1].(string)
	if !ok {
		scheme = defaultScheme
	}

	opts := make([]grpc.DialOption, 0)
	if len(args) == 3 {
		dialOpts, ok := args[2].([]grpc.DialOption)
		if !ok {
			return nil, ErrInvalidDialOptions
		}

		opts = dialOpts
	}

	return func(args ...interface{}) (pool.Connection, error) {
		if len(args) != 1 {
			return nil, ErrTooFewFactoryArgs
		}

		address, ok := args[0].(string)
		if !ok {
			return nil, ErrInvalidAddress
		}

		cc, err := grpc.Dial(registry.ConnectionString(address, scheme), opts...)
		if err != nil {
			return nil, fmt.Errorf("gRPC connection failed: %v", err)
		}

		return cc, nil
	}, nil
}

// Helper method to get GrpcClient
func Get(server server.Server, name string) (*GrpcClient, error) {
	conn, err := server.GetClient(name)
	if err != nil {
		return nil, err
	}

	client, ok := conn.(*GrpcClient)
	if !ok {
		return nil, ErrInvalidGrpcClient
	}

	return client, nil
}

// Call - calls method on the connected gRPC client
func (gc *GrpcClient) Call(
	ctx context.Context,
	method string,
	req interface{},
	res interface{},
	opts ...interface{},
) error {
	cc, ok := gc.Connection.(*grpc.ClientConn)
	if !ok {
		return fmt.Errorf("Invalid connection")
	}

	callOpts := make([]grpc.CallOption, len(opts))
	for i, opt := range opts {
		callOpts[i] = opt.(grpc.CallOption)
	}

	return cc.Invoke(ctx, method, req, res, callOpts...)
}

// Close - closes the connection to the gRPC service server
func (gc *GrpcClient) Close() error {
	gc.RLock()
	defer gc.RUnlock()

	if gc.Connection != nil {
		return gc.Connection.Close()
	}

	return gc.p.add(gc.address, gc.Connection)
}
