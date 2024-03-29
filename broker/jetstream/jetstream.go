package jetstream

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/easeq/go-service/broker"
	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/logger"
	"github.com/easeq/go-service/tracer"
	"github.com/easeq/go-service/utils"
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
	w             *broker.Wrapper
	nc            *nats.Conn
	logger        logger.Logger
	tracer        tracer.Tracer
	jsCtx         nats.JetStreamContext
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
		Config:        config,
		Subscriptions: make(map[string]*nats.Subscription),
	}

	for _, opt := range opts {
		opt(j)
	}

	j.w = broker.NewWrapper(j)
	j.i = NewInitializer(j)

	return j
}

// AddStream defines a the stream in which to publish the message
func AddStream(name string, subjects ...string) broker.Option {
	return func(b broker.Broker) {
		if len(subjects) == 0 {
			subjectAll := fmt.Sprintf("%s.>", name)
			subjects = []string{subjectAll}
		}

		b.(*JetStream).createStream(name, subjects...)
	}
}

// Logger returns the initialized logger instance
func (j *JetStream) Logger() logger.Logger {
	return j.logger
}

// StreamExists returns whether a stream by theErrorw given name exists
func (j *JetStream) streamExists(name string) *nats.StreamInfo {
	info, err := j.jsCtx.StreamInfo(name)
	if err != nil {
		return nil
	}

	return info
}

// createStream creates a new JS stream if it doens't exist and
// attaches the provided subjects to the stream.
// If the stream, the new subjects are appended to the stream
func (j *JetStream) createStream(name string, subjects ...string) error {
	if streamInfo := j.streamExists(name); streamInfo != nil {
		newStreamCfg := streamInfo.Config
		newStreamCfg.Retention = nats.InterestPolicy
		newStreamCfg.Subjects = utils.Unique(
			append(newStreamCfg.Subjects, subjects...),
		)

		_, err := j.jsCtx.UpdateStream(&newStreamCfg)
		return err
	}

	_, err := j.jsCtx.AddStream(&nats.StreamConfig{
		Name:       name,
		Subjects:   subjects,
		Retention:  nats.InterestPolicy,
		Duplicates: 0 * time.Second,
		Discard:    nats.DiscardOld,
		// Retention: nats.InterestPolicy,
		// Replicas:  1,
	})

	return err
}

// Publish publishes the topic message
func (j *JetStream) Publish(ctx context.Context, topic string, message interface{}, opts ...broker.PublishOption) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshalling error: %v", err)
	}

	publisher := NewPublisher(j, topic, opts...)
	return j.w.Publish(ctx, topic, payload, func(t *broker.TraceMsgCarrier) error {
		data, err := t.Bytes()
		if err != nil {
			j.logger.Errorw("publish[data error]", "topic", topic, "err", err)
			return err
		}

		// Send the message with span over NATS
		_, err = j.jsCtx.Publish(topic, data, publisher.opts...)
		if err != nil {
			j.logger.Errorw("publish error", "topic", topic, "err", err)
			return err
		}

		return nil
	})
}

// Subscribe subcribes for the given topic.
func (j *JetStream) Subscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) error {
	subscriber := NewSubscriber(j, topic, opts...)
	natsHandler := func(m *nats.Msg) {
		// Create new TraceMsg from normal NATS message.
		j.w.Subscribe(ctx, m.Subject, m.Data, func(
			ctx context.Context,
			t *broker.TraceMsgCarrier,
		) error {
			if err := handler.Handle(ctx, &broker.Message{
				Body: t.Message,
				Extras: map[string]interface{}{
					broker.KEY_TRACE_MSG_CARRIER: t,
					broker.KEY_BROKER_MSG:        m,
				},
			}); err != nil {
				m.Nak()
				j.logger.Errorw("subscribe handle error", "topic", topic, "err", err)
				return fmt.Errorf("subscribe handle error: %v", err)
			}

			m.Ack()
			return nil
		})
	}

	subscription, err := subscriber.Subscribe(natsHandler)
	j.logger.Infow("subscription", "s", subscription)
	if err != nil {
		j.logger.Errorw("subscribe error", "topic", topic, "err", err)
		return err
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

func (j *JetStream) String() string {
	return "nats.jetstream"
}
