package jetstream

import (
	"context"
	"encoding/json"
	"errors"

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
func NewJetStream() *JetStream {
	config := goconfig.NewEnvConfig(new(Config)).(*Config)
	nc, err := nats.Connect(config.Address())
	if err != nil {
		panic("error connecting to nats server")
	}

	js, err := nc.JetStream()
	if err != nil {
		panic("error creating JetStreamContext")
	}

	return &JetStream{
		nc:            nc,
		jsCtx:         js,
		Js:            js,
		Config:        config,
		Subscriptions: make(map[string]*nats.Subscription),
	}
}

// Publish publishes the topic message
func (j *JetStream) Publish(ctx context.Context, topic string, message interface{}, opts ...broker.PublishOption) error {
	publisher := NewJetStreamPublisher(j, opts...)
	publisher.createStream()

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = j.jsCtx.Publish(topic, payload)
	return err
}

// Subscribe subcribes for the given topic.
func (j *JetStream) Subscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) error {

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

	subscription, err := j.nc.Subscribe(topic, natsHandler)
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
