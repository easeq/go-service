package broker

import "context"

// Message structure
type Message struct {
	Body   []byte
	Extras interface{}
}

// Handler used by the subscriber
type Handler interface {
	// Handles the subscribed message
	Handle(m *Message) error
}

// Broker interface for adding new brokers
type Broker interface {
	// Publish a message
	Publish(ctx context.Context, topic string, message interface{}, opts ...PublishOption) error
	// Subscribe to a subject
	Subscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) error
	// Unsubscribe from a subject
	Unsubscribe(topic string) error
	// Close a connection
	Close() error
}

// Interface for a message publisher
type Publisher interface{}

// PublishOption to pass as arg while publishing a message
type PublishOption func(Publisher)

// Interface for a subscription subscriber
type Subscriber interface{}

// SubscribeOption to pass as arg while creating subscription
type SubscribeOption func(Subscriber)

// Interface for a broker runner instance
type Runner interface{}

// RunOption to pass as arg while calling the run method for the broker
type RunOption func(Runner)
