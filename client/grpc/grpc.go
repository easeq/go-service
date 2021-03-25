package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

type GrpcClient struct {
	cc    *grpc.ClientConn
	timer *time.Timer
}

// NewGrpcClient creates a new gRPC client connection
func NewGrpcClient(address string, ttl time.Duration, opts ...grpc.DialOption) (*GrpcClient, error) {
	// Get gRPC client connection
	cc, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, fmt.Errorf("gRPC connection failed: %v", err)
	}

	// Create new gRPC client
	gc := &GrpcClient{cc, time.NewTimer(ttl)}
	go func() {
		// Close connection after TTL
		<-gc.timer.C
		gc.cc.Close()
	}()

	return gc, nil
}

// isInitialized returns whether the client connection has already been initialzed
func (gc *GrpcClient) isInitialized() bool {
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
	if gc.cc != nil && gc.cc.GetState().String() != "SHUTDOWN" {
		return gc.cc.Close()
	}

	return nil
}
