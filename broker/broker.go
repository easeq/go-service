package broker

import "context"

// Message sent by the Publisher
type Message interface{}

// Handler used by the subscriber
type Handler interface{}

// Broker interface for adding new brokers
type Broker interface {
	// Run broker
	Run(ctx context.Context, opts ...RunOption) error
	// Publish a message
	Publish(ctx context.Context, topic string, message Message) error
	// Subscribe to a subject
	Subscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) error
	// Unsubscribe from a subject
	Unsubscribe(topic string) error
	// Close a connection
	Close() error
}

// Interface for a subscription subscriber
type Subscriber interface{}

// SubscribeOption to pass as arg while creating subscription
type SubscribeOption func(Subscriber)

// Interface for a broker runner instance
type Runner interface{}

// RunOption to pass as arg while calling the run method for the broker
type RunOption func(Runner)
