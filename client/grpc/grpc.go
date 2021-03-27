package grpc

import (
	"context"
	"errors"
	"fmt"
	"sync"

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
	ErrInvalidAddress = errors.New("Invalid service address")
	// ErrInvalidRegistry returned when the registry provided does not implement regsitry.ServiceRegsitry
	ErrInvalidRegistry = errors.New("Registry provided is invalid")
	// ErrTooFewArgs returned when the args provided is less the args required
	ErrTooFewArgs = errors.New("Too few arguments. Required 3")
	// ErrInvalidDialOptions returned when the dial options provided are not valid
	ErrInvalidDialOptions = errors.New("Dial options provided are invalid")
)

type GrpcClient struct {
	pool.Connection
	address string
	p       *Pool
	*sync.RWMutex
	// cc      *grpc.ClientConn
	// timer   *time.Timer
}

// NewGrpcClient creates and rea new gRPC client connection
func NewGrpcClientConn(args ...interface{}) (pool.Connection, error) {
	if len(args) < 3 {
		return nil, ErrTooFewArgs
	}

	address, ok := args[0].(string)
	if !ok {
		return nil, ErrInvalidAddress
	}

	registry, ok := args[1].(registry.ServiceRegistry)
	if !ok {
		return nil, ErrInvalidRegistry
	}

	scheme, ok := args[2].(string)
	if !ok {
		scheme = defaultScheme
	}

	opts := make([]grpc.DialOption, 0)
	if len(args) == 4 {
		dialOpts, ok := args[3].([]grpc.DialOption)
		if !ok {
			return nil, ErrInvalidDialOptions
		}

		opts = dialOpts
	}

	cc, err := grpc.Dial(registry.ConnectionString(address, scheme), opts...)
	if err != nil {
		return nil, fmt.Errorf("gRPC connection failed: %v", err)
	}

	return cc, nil
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

	// if gc.cc != nil && gc.cc.GetState().String() != "SHUTDOWN" {
	// 	// Release connection from pool
	// 	gc.p.Release(gc.address)
	// 	return gc.cc.Close()
	// }

	return gc.p.add(gc.address, gc.Connection)
}
