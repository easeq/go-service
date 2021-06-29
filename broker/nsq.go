package broker

import (
	"context"
	"encoding/json"
	"errors"

	goconfig "github.com/easeq/go-config"
	"github.com/nsqio/go-nsq"
)

var (
	// ErrInvalidMessageHandler returned when the message handler doesn't implement the underlying interface
	ErrInvalidMessageHandler = errors.New("invalid message handler provided")
)

// Nsq holds our broker instance
type Nsq struct {
	Producer  *nsq.Producer
	Consumers map[string]*nsq.Consumer
	*Config
}

// NewNsq returns a new instance of NSQ
func NewNsq() *Nsq {
	return &Nsq{
		Consumers: make(map[string]*nsq.Consumer),
		Config:    goconfig.NewEnvConfig(new(Config)).(*Config),
	}
}

// Run the broker until it's stopped
func (n *Nsq) Run(ctx context.Context, opts ...RunOption) error {
	producer, err := nsq.NewProducer(n.Config.Producer.Address(), n.NSQConfig())
	if err != nil {
		return err
	}

	n.Producer = producer

	<-ctx.Done()

	return nil
}

// Publish publishes the topic message
func (n *Nsq) Publish(ctx context.Context, topic string, message Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return n.Producer.Publish(topic, payload)
}

// Subscribe subcribes for the given topic
func (n *Nsq) Subscribe(ctx context.Context, topic string, handler Handler, opts ...SubscribeOption) error {
	subscriber := NewNsqSubscriber(n, topic, opts...)
	consumer, err := nsq.NewConsumer(topic, subscriber.channel, n.NSQConfig())
	if err != nil {
		return err
	}

	nsqHandler, ok := handler.(nsq.Handler)
	if !ok {
		return ErrInvalidMessageHandler
	}

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

// Close shuts down the broker
func (n *Nsq) Close() error {
	n.Producer.Stop()

	for _, consumer := range n.Consumers {
		consumer.Stop()
	}

	return nil
}

// NsqSubscriber holds additional options for nsq subscription
type NsqSubscriber struct {
	channel string
}

// NewNsqSubscriber returns a new subscriber instance for NSQ subscription
func NewNsqSubscriber(n *Nsq, topic string, opts ...SubscribeOption) *NsqSubscriber {
	subscriber := &NsqSubscriber{
		channel: n.Channel(topic),
	}

	for _, opt := range opts {
		opt(subscriber)
	}

	return subscriber
}

// WithChannelName defines a channel name for the subscriber
func WithChannelName(name string) SubscribeOption {
	return func(s Subscriber) {
		s.(*NsqSubscriber).channel = name
	}
}
