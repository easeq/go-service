package nsq

import (
	"context"
	"encoding/json"
	"errors"

	goconfig "github.com/easeq/go-config"
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
	config := goconfig.NewEnvConfig(new(Config)).(*Config)

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
		return err
	}

	return n.t.Publish(topic, func(t *broker.TraceMsg) error {
		// Add the payload/original message
		t.Write(payload)

		// Send the message with span over NSQ
		return n.Producer.Publish(topic, payload)
	})
}

// Subscribe subcribes for the given topic
func (n *Nsq) Subscribe(ctx context.Context, topic string, handler broker.Handler, opts ...broker.SubscribeOption) error {
	subscriber := NewNsqSubscriber(n, topic, opts...)
	consumer, err := nsq.NewConsumer(topic, subscriber.channel, n.NSQConfig())
	if err != nil {
		return err
	}

	nsqHandler := NewNsqHandler(n, topic, handler)
	consumer.AddHandler(nsqHandler)
	if err := consumer.ConnectToNSQD(n.Config.Producer.Address()); err != nil {
		return err
	}

	if err := consumer.ConnectToNSQLookupd(n.Config.Lookupd.Address()); err != nil {
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
