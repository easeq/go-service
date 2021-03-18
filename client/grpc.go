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
	wg   sync.WaitGroup
}

type GrpcClientOptions struct {
	Scheme      string
	Registry    registry.ServiceRegistry
	DialOptions []grpc.DialOption
}

// Init gRPC client
func (gc *GrpcClient) Init(name string) error {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	defer gc.wg.Add(1)

	if gc.IsInitialized() {
		return nil
	}

	cc, err := grpc.Dial(gc.Opts.Registry.ConnectionString(name, gc.Opts.Scheme), gc.Opts.DialOptions...)
	if err != nil {
		return fmt.Errorf("gRPC connection failed: %v", err)
	}

	gc.cc = cc

	return nil
}

// Return whether client has been intialized
func (gc *GrpcClient) IsInitialized() bool {
	state := ""
	if gc.cc != nil {
		state = gc.cc.GetState().String()
	}

	return gc.cc != nil && state != "SHUTDOWN"
}

// Call - calls method on the connected gRPC client
func (gc *GrpcClient) Call(
	ctx context.Context,
	method string,
	req interface{},
	res interface{},
	opts ...interface{},
) error {
	callOpts := make([]grpc.CallOption, len(opts))
	for i, opt := range opts {
		callOpts[i] = opt.(grpc.CallOption)
	}

	return gc.cc.Invoke(ctx, method, req, res, callOpts...)
}

// Close - closes the connection to the gRPC service server
func (gc *GrpcClient) Close() error {
	gc.wg.Done()
	gc.wg.Wait()

	if gc.cc != nil && gc.cc.GetState().String() != "SHUTDOWN" {
		return gc.cc.Close()
	}

	return nil
}
