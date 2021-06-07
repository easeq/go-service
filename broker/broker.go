package broker

import "context"

// Message sent by the Publisher
type Message interface{}

// Handler used by the subscriber
type Handler interface{}

// Broker interface for adding new brokers
type Broker interface {
	// Initialize
	Init(ctx context.Context, opts ...interface{}) error
	// Publish a message
	Publish(ctx context.Context, topic string, message Message) error
	// Subscribe to a subject
	Subscribe(ctx context.Context, topic string, handler Handler, opts ...interface{}) error
	// Unsubscribe from a subject
	Unsubscribe(topic string) error
	// Close a connection
	Close() error
}
