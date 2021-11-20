package nsq

import (
	"github.com/easeq/go-service/broker"
)

// Subscriber holds additional options for nsq subscription
type subscriber struct {
	channel string
}

// NewNsqSubscriber returns a new subscriber instance for NSQ subscription
func NewNsqSubscriber(n *Nsq, topic string, opts ...broker.SubscribeOption) *subscriber {
	subscriber := &subscriber{
		channel: n.Channel(topic),
	}

	for _, opt := range opts {
		opt(subscriber)
	}

	return subscriber
}

// WithChannelName defines a channel name for the subscriber
func WithChannelName(name string) broker.SubscribeOption {
	return func(s broker.Subscriber) {
		s.(*subscriber).channel = name
	}
}
