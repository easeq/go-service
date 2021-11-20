package client

import (
	"context"
	"errors"

	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/pool"
)

const (
	CLIENT = "client"
)

var (
	// ErrInvalidStream returned when the stream is invalid
	ErrInvalidStream = errors.New("invalid client stream")
)

// CallFn function for Call method
type CallFn func(conn pool.Connection) error

// CallOption is used to pass client call options
type CallOption interface{}

// DialOption is used to pass clients' dial options
type DialOption interface{}

// Client interface to implement custom clients
type Client interface {
	component.Component
	// Call client method
	Call(ctx context.Context, sc ServiceClient, method string, req interface{}, res interface{}, opts ...CallOption) error
	// Stream client
	Stream(ctx context.Context, sc ServiceClient, desc interface{}, method string, req interface{}, opts ...CallOption) (StreamClient, error)
	// Get client connection
	Dial(address string, opts ...DialOption) (pool.Connection, error)
}

// ServiceClient is used by the generated code to return service config
type ServiceClient interface {
	// Get service name of the client
	GetServiceName() string
	// Get dial options for the client
	GetDialOptions() []DialOption
}

// StreamClient interface is used by client implementation for streaming
type StreamClient interface {
	// Send request
	Send(req interface{}) error
	// Close send and receive response
	CloseAndRecv(res interface{}) error
	// Receive response
	Recv(res interface{}) error
	// Close client connection
	CloseConn() error
}
