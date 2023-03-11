package broker

import (
	"context"

	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/logger"
)

const (
	BROKER = "broker"

	KEY_TRACE_MSG_CARRIER = "trace_msg_carrier"
	KEY_BROKER_MSG        = "msg"
)

// Message structure
type Message struct {
	Body   []byte
	Extras map[string]interface{}
}

// Handler used by the subscriber
type Handler interface {
	// Handles the subscribed message
	Handle(ctx context.Context, m *Message) error
}

// Broker interface for adding new brokers
type Broker interface {
	component.Component
	// Return logger
	Logger() logger.Logger
	// Publish a message
	Publish(ctx context.Context, topic string, message interface{}, opts ...PublishOption) error
	// Subscribe to a subject
	Subscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) (Subscription, error)
	// Unsubscribe from a subject
	Unsubscribe(topic string) error
	// String returns the string name of the broker
	String() string
}

type Subscription interface{}

// Option to pass as arg while creating new broker instance
type Option func(Broker)

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

func LogError(l logger.Logger, msg string, topic string, err error) {
	l.Errorw(msg, "topic", topic, "error", err)
}
