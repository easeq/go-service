package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/easeq/go-service/registry"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	cc   *grpc.ClientConn
	mu   sync.RWMutex
	Opts *GrpcClientOptions
}

type GrpcClientOptions struct {
	Scheme      string
	Registry    registry.ServiceRegistry
	DialOptions []grpc.DialOption
}

// Init gRPC client
func (gc *GrpcClient) Init(name string) error {
	cc, err := grpc.Dial(gc.Opts.Registry.ConnectionString(name, gc.Opts.Scheme), gc.Opts.DialOptions...)
	if err != nil {
		return fmt.Errorf("gRPC connection failed: %v", err)
	}

	gc.mu.Lock()
	gc.cc = cc
	gc.mu.Unlock()

	return nil
}

// Return whether client has been intialized
func (gc *GrpcClient) IsInitialized() bool {
	return gc.cc != nil
}

// Call - calls method on the connected gRPC client
func (gc *GrpcClient) Call(
	ctx context.Context,
	method string,
	req interface{},
	res interface{},
	opts interface{},
) error {
	callOpts, ok := opts.(*[]grpc.CallOption)
	if !ok {
		return fmt.Errorf("Invalid call option provided")
	}

	return gc.cc.Invoke(ctx, method, req, res, *callOpts...)
}

// Close - closes the connection to the gRPC service server
func (gc *GrpcClient) Close() error {
	return gc.cc.Close()
}
