package jetstream

import (
	"github.com/easeq/go-service/broker"
	"github.com/nats-io/nats.go"
)

// Subscriber holds additional options for jetstream subscription
type publisher struct {
	opts []nats.PubOpt
}

// NewPublisher returns a new publisher instance for jetstream publisher
func NewPublisher(j *JetStream, topic string, opts ...broker.PublishOption) *publisher {
	p := &publisher{
		opts: []nats.PubOpt{},
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// WithNatsPubOpts defines a additional jetstream publish options
func WithNatsPubOpts(opts ...nats.PubOpt) broker.PublishOption {
	return func(p broker.Publisher) {
		pubOpts := p.(*publisher).opts
		pubOpts = append(pubOpts, opts...)
		p.(*publisher).opts = pubOpts
	}
}
