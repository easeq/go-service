package broker

import "context"

type Message interface{}

type Handler interface{}

type Broker interface {
	// Initialize
	Init(ctx context.Context, args map[string]interface{}, opts ...interface{}) error
	// Publish a message
	Publish(ctx context.Context, subject string, message []byte) error
	// Subscribe to a subject
	Subscribe(ctx context.Context, subject string, handler Handler, opts ...interface{}) error
	// Unsubscribe from a subject
	Unsubscribe(subject string) error
	// Close a connection
	Close() error
}
