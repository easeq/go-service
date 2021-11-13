package jetstream

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	goconfig "github.com/easeq/go-config"
	"github.com/easeq/go-service/broker"
	nats "github.com/nats-io/nats.go"
)

var (
	// ErrInvalidMessageHandler returned when the message handler doesn't implement the underlying interface
	ErrInvalidMessageHandler = errors.New("invalid message handler provided")
	// ErrSubscriptionFailed returned when subscription fails
	ErrSubscriptionFailed = errors.New("nats subscription failed")
)

// Nsq holds our broker instance
type JetStream struct {
	nc            *nats.Conn
	jsCtx         nats.JetStreamContext
	Js            nats.JetStreamContext
	Subscriptions map[string]*nats.Subscription
	*Config
}

// NewJetStream returns a new instance of nats jetstream
func NewJetStream(opts ...broker.Option) *JetStream {
	config := goconfig.NewEnvConfig(new(Config)).(*Config)
	nc, err := nats.Connect(config.Address())
	if err != nil {
		panic("error connecting to nats server")
	}

	js, err := nc.JetStream()
	if err != nil {
		panic("error creating JetStreamContext")
	}

	j := &JetStream{
		nc:            nc,
		jsCtx:         js,
		Js:            js,
		Config:        config,
		Subscriptions: make(map[string]*nats.Subscription),
	}

	for _, opt := range opts {
		opt(j)
	}

	return j
}

// AddStream defines a the stream in which to publish the message
func AddStream(name string, subjects ...string) broker.Option {
	return func(b broker.Broker) {
		if len(subjects) == 0 {
			subjectAll := fmt.Sprintf("%s.*", name)
			subjects = []string{subjectAll}
		}

		b.(*JetStream).createStream(name, subjects...)
	}
}

// StreamExists returns whether a stream by the given name exists
func (j *JetStream) streamExists(name string) bool {
	if _, err := j.jsCtx.StreamInfo(name); err != nil {
		return false
	}

	return true
}

// createStream creates a new JS stream if it doens't exist and
// attaches the pre-defined subjects to the stream
func (j *JetStream) createStream(name string, subjects ...string) error {
	if j.streamExists(name) {
		return nil
	}

	_, err := j.jsCtx.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: subjects,
	})

	return err
}

// Publish publishes the topic message
func (j *JetStream) Publish(ctx context.Context, topic string, message interface{}, opts ...broker.PublishOption) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = j.jsCtx.Publish(topic, payload)
	return err
}

// Subscribe subcribes for the given topic.
func (j *JetStream) Subscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) error {
	subscriber := NewSubscriber(j, topic, opts...)
	natsHandler := func(m *nats.Msg) {
		err := handler.Handle(&broker.Message{
			Body:   m.Data,
			Extras: m,
		})
		if err != nil {
			m.Nak()
		}
		m.Ack()
	}

	subscription, err := j.jsCtx.Subscribe(topic, natsHandler, subscriber.subOpts...)
	if err != nil {
		return ErrSubscriptionFailed
	}

	j.Subscriptions[topic] = subscription

	return nil
}

// Ubsubscribe method is not applicable
func (j *JetStream) Unsubscribe(topic string) error {
	return j.Subscriptions[topic].Unsubscribe()
}

// Close shuts down the broker
func (j *JetStream) Close() error {
	for _, sub := range j.Subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			return err
		}
	}

	j.nc.Close()
	return nil
}