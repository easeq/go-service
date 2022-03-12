package jetstream

import (
	"github.com/easeq/go-service/broker"
	"github.com/nats-io/nats.go"
)

// Subscriber holds additional options for jetstream subscription
type subscriber struct {
	opts []nats.SubOpt
}

// NewSubscriber returns a new subscriber instance for jetstream subscription
func NewSubscriber(j *JetStream, topic string, opts ...broker.SubscribeOption) *subscriber {
	s := &subscriber{
		opts: []nats.SubOpt{
			nats.MaxDeliver(3),
			nats.ManualAck(),
			nats.DeliverNew(),
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithNatsSubOpts defines a additional jetstream subscribe options
func WithNatsSubOpts(opts ...nats.SubOpt) broker.SubscribeOption {
	return func(s broker.Subscriber) {
		subOpts := s.(*subscriber).opts
		subOpts = append(subOpts, opts...)
		s.(*subscriber).opts = subOpts
	}
}
