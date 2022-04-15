package broker

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/logger"
)

const (
	BROKER = "broker"
)

// Message structure
type Message struct {
	Body   []byte
	Extras interface{}
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
	Subscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) error
	// Unsubscribe from a subject
	Unsubscribe(topic string) error
	// String returns the string name of the broker
	String() string
}

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

// EncodeMessage encodes a broker message
func EncodeMessage(msg interface{}) ([]byte, error) {
	var payload bytes.Buffer
	enc := gob.NewEncoder(&payload)

	if err := enc.Encode(msg); err != nil {
		return nil, fmt.Errorf("encoding error: %v", err)
	}

	return payload.Bytes(), nil
}

// DecodeMessage decodes a message in bytes to the given interface
func DecodeMessage(encMsg []byte, v interface{}) error {
	var payload bytes.Buffer
	dec := gob.NewDecoder(&payload)
	payload.Write(encMsg)

	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("decode message error: %v", err)
	}

	return nil
}
