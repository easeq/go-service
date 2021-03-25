package client

import "context"

type Client interface {
	// Calls a method on a service
	Call(ctx context.Context, method string, req interface{}, res interface{}, opts ...interface{}) error
}
