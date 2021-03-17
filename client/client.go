package client

import (
	"context"
)

type Client interface {
	Init(name string) error
	// Return whether client has been initialized
	IsInitialized() bool
	// Calls a method on a service
	Call(ctx context.Context, method string, req interface{}, res interface{}, opts interface{}) error
	// Closes the connection to the service
	Close() error
}
