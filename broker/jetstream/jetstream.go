package jetstream

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/easeq/go-service/broker"
	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/logger"
	"github.com/easeq/go-service/tracer"
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
	i             component.Initializer
	t             *broker.Trace
	nc            *nats.Conn
	logger        logger.Logger
	tracer        tracer.Tracer
	jsCtx         nats.JetStreamContext
	Js            nats.JetStreamContext
	Subscriptions map[string]*nats.Subscription
	*Config
}

// NewJetStream returns a new instance of nats jetstream
func NewJetStream(opts ...broker.Option) *JetStream {
	config := NewConfig()
	nc, err := nats.Connect(config.Address())
	if err != nil {
		panic("error connecting to nats server")
	}

	js, err := nc.JetStream()
	if err != nil {
		panic("error creating JetStreamContext")
	}

	j := &JetStream{
		i:             nil,
		nc:            nc,
		jsCtx:         js,
		Js:            js,
		Config:        config,
		Subscriptions: make(map[string]*nats.Subscription),
	}

	for _, opt := range opts {
		opt(j)
	}

	j.t = broker.NewTrace(j)
	j.i = NewInitializer(j)

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

// Logger returns the initialized logger instance
func (j *JetStream) Logger() logger.Logger {
	return j.logger
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
		j.logger.Errorf("JetStream publish payload error: %s", err)
		return err
	}

	return j.t.Publish(topic, func(t *broker.TraceMsg) error {
		// Add the payload/original message
		t.Write(payload)

		// Send the message with span over NATS
		_, err = j.jsCtx.Publish(topic, t.Bytes())

		j.logger.Errorf("JetStream publish error: %s", err)
		return err
	})
}

// Subscribe subcribes for the given topic.
func (j *JetStream) Subscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) error {
	j.logger.Info("Subscribe message", topic)
	subscriber := NewSubscriber(j, topic, opts...)
	natsHandler := func(m *nats.Msg) {
		// Create new TraceMsg from normal NATS message.
		j.t.Subscribe(m.Subject, m.Data, func(body []byte) error {
			if err := handler.Handle(&broker.Message{
				Body:   body,
				Extras: m,
			}); err != nil {
				j.logger.Errorf("JetStream subcribe handle error: %s", err)
				m.Nak()
				return err
			}

			m.Ack()
			return nil
		})
	}

	subscription, err := j.jsCtx.Subscribe(topic, natsHandler, subscriber.opts...)
	if err != nil {
		j.logger.Errorf("JetStream subcription error: %s", err)
		return fmt.Errorf("JetStream subscription failed: %v", err)
	}

	j.Subscriptions[topic] = subscription

	return nil
}

// Ubsubscribe method is not applicable
func (j *JetStream) Unsubscribe(topic string) error {
	return j.Subscriptions[topic].Unsubscribe()
}

func (j *JetStream) HasInitializer() bool {
	return true
}

func (j *JetStream) Initializer() component.Initializer {
	return j.i
}
