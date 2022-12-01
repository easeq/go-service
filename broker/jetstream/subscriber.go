package jetstream

import (
	"errors"

	"github.com/easeq/go-service/broker"
	"github.com/nats-io/nats.go"
)

const (
	DEFAULT = iota
	SYNC
	PULL
	QUEUE
)

var (
	// ErrReqDurableName returned when durable name is a required field
	ErrReqDurableName = errors.New("durable name is required - use withDurableName(name) subopt")
	// ErrReqQueueName returned when queue name is a required field
	ErrReqQueueName = errors.New("queue name is required - use withQueueName(name) subopt")
)

// Subscriber holds additional options for jetstream subscription
type subscriber struct {
	j           *JetStream
	topic       string
	durableName string
	queueName   string
	sType       uint
	opts        []nats.SubOpt
}

// NewSubscriber returns a new subscriber instance for jetstream subscription
func NewSubscriber(j *JetStream, topic string, opts ...broker.SubscribeOption) *subscriber {
	s := &subscriber{
		j:     j,
		topic: topic,
		sType: DEFAULT,
		opts: []nats.SubOpt{
			nats.MaxDeliver(3),
			nats.ManualAck(),
			nats.AckExplicit(),
			nats.DeliverAll(),
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

// WithQueueSubscription - used to create a queue subscriber
func WithQueueSubscription() broker.SubscribeOption {
	return func(s broker.Subscriber) {
		s.(*subscriber).sType = QUEUE
	}
}

// WithDurableName - used to provied a durable name for sync and pull subscription
func WithDurableName(name string) broker.SubscribeOption {
	return func(s broker.Subscriber) {
		s.(*subscriber).durableName = name
	}
}

// WithQueueName - used to provied a queue name for queue subscriptions
func WithQueueName(name string) broker.SubscribeOption {
	return func(s broker.Subscriber) {
		s.(*subscriber).queueName = name
	}
}

func (s *subscriber) Subscribe(handler func(m *nats.Msg)) (*nats.Subscription, error) {
	switch s.sType {
	case QUEUE:
		if s.queueName == "" {
			return nil, ErrReqQueueName
		}

		return s.j.jsCtx.QueueSubscribe(s.topic, s.queueName, handler, s.opts...)
	case SYNC:
		if s.durableName != "" {
			s.opts = append(s.opts, nats.Durable(s.durableName))
		}
		// TODO: make this possible
		return s.j.jsCtx.SubscribeSync(s.topic, s.opts...)
	case PULL:
		if s.durableName == "" {
			return nil, ErrReqDurableName
		}
		// TODO: make this possible
		return s.j.jsCtx.PullSubscribe(s.topic, s.durableName, s.opts...)
	default:
		return s.j.jsCtx.Subscribe(s.topic, handler, s.opts...)
	}
}
