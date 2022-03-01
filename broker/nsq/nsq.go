package nsq

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/easeq/go-service/broker"
	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/logger"
	"github.com/easeq/go-service/tracer"
	"github.com/nsqio/go-nsq"
)

var (
	// ErrInvalidMessageHandler returned when the message handler doesn't implement the underlying interface
	ErrInvalidMessageHandler = errors.New("invalid message handler provided")
)

// Nsq holds our broker instance
type Nsq struct {
	i         component.Initializer
	t         *broker.Trace
	logger    logger.Logger
	tracer    tracer.Tracer
	Producer  *nsq.Producer
	Consumers map[string]*nsq.Consumer
	*Config
}

// NewNsq returns a new instance of NSQ
func NewNsq() *Nsq {
	config := NewConfig()

	producer, err := nsq.NewProducer(config.Producer.Address(), config.NSQConfig())
	if err != nil {
		panic("error starting nsq producer")
	}

	n := &Nsq{
		Producer:  producer,
		Consumers: make(map[string]*nsq.Consumer),
		Config:    config,
	}

	n.t = broker.NewTrace(n)
	n.i = NewInitializer(n)
	return n
}

// Logger returns the initialized logger instance
func (n *Nsq) Logger() logger.Logger {
	return n.logger
}

// Publish publishes the topic message
func (n *Nsq) Publish(ctx context.Context, topic string, message interface{}, opts ...broker.PublishOption) error {
	payload, err := json.Marshal(message)
	if err != nil {
		broker.LogError(n.logger, "NSQ unmarshalling message error", topic, err)
		return err
	}

	return n.t.Publish(ctx, topic, func(t *broker.TraceMsg) error {
		// Add the payload/original message
		t.Write(payload)

		// Send the message with span over NSQ
		err := n.Producer.Publish(topic, payload)
		broker.LogError(n.logger, "NSQ publish error", topic, err)

		return err
	})
}

// Subscribe subcribes for the given topic
func (n *Nsq) Subscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) error {
	subscriber := NewNsqSubscriber(n, topic, opts...)
	consumer, err := nsq.NewConsumer(topic, subscriber.channel, n.NSQConfig())
	if err != nil {
		broker.LogError(n.logger, "NSQ new consumer error", topic, err)
		return err
	}

	// TODO: add trace
	nsqHandler := NewNsqHandler(n, topic, handler)
	consumer.AddHandler(nsqHandler)
	if err := consumer.ConnectToNSQD(n.Config.Producer.Address()); err != nil {
		broker.LogError(n.logger, "NSQ consumer connect to NSQD error", topic, err)
		return err
	}

	if err := consumer.ConnectToNSQLookupd(n.Config.Lookupd.Address()); err != nil {
		broker.LogError(n.logger, "NSQ consumer connect to NSQLookupd error", topic, err)
		return err
	}

	n.Consumers[topic] = consumer

	return nil
}

// Ubsubscribe method is not applicable
func (n *Nsq) Unsubscribe(topic string) error {
	return nil
}

func (n *Nsq) HasInitializer() bool {
	return true
}

func (n *Nsq) Initializer() component.Initializer {
	return n.i
}
