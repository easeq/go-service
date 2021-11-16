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
	t             *broker.Trace
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
		t:             new(broker.Trace),
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

	return j.t.Publish(topic, func(t *broker.TraceMsg) error {
		// Add the payload/original message
		t.Write(payload)

		// Send the message with span over NATS
		_, err = j.jsCtx.Publish(topic, t.Bytes())
		return err
	})
}

// Subscribe subcribes for the given topic.
func (j *JetStream) Subscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) error {
	subscriber := NewSubscriber(j, topic, opts...)
	natsHandler := func(m *nats.Msg) {
		// Create new TraceMsg from normal NATS message.
		j.t.Subscribe(m.Subject, m.Data, func(body []byte) error {
			if err := handler.Handle(&broker.Message{
				Body:   body,
				Extras: m,
			}); err != nil {
				m.Nak()
				return err
			}

			m.Ack()
			return nil
		})
	}

	subscription, err := j.jsCtx.Subscribe(topic, natsHandler, subscriber.opts...)
	if err != nil {
		return fmt.Errorf("JetStream subscription failed: %v", err)
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
