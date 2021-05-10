package client

import (
	"context"
	"errors"

	"github.com/easeq/go-service/pool"
)

var (
	ErrInvalidStream = errors.New("invalid client stream")
)

// Callback function for Call method
type CallFn func(conn pool.Connection) error

type CallOption interface{}

type DialOption interface{}

type Client interface {
	// Call client method
	Call(ctx context.Context, sc ServiceClient, method string, req interface{}, res interface{}, opts ...CallOption) error
	// Stream client
	Stream(ctx context.Context, sc ServiceClient, desc interface{}, method string, req interface{}, opts ...CallOption) (StreamClient, error)
	// Get client connection
	Dial(address string, opts ...DialOption) (pool.Connection, error)
}

type ServiceClient interface {
	// Get service name of the client
	GetServiceName() string
	// Get dial options for the client
	GetDialOptions() []DialOption
}

type StreamClient interface {
	// Send request
	Send(req interface{}) error
	// Close send and receive response
	CloseAndRecv(res interface{}) error
	// Receive response
	Recv(res interface{}) error
}
