package nats_streaming

import (
	"fmt"

	"github.com/easeq/go-service/broker"
	stan "github.com/nats-io/stan.go"
)

type NatsStreaming struct {
	conn          stan.Conn
	subscriptions map[string]stan.Subscription
	*Config
}

// NewPostgres returns new connection to the postgres db
func NewNatsStreaming() broker.Broker {
	cfg := GetConfig()

	sc, err := stan.Connect(cfg.ClusterID, cfg.ClientID)
	if err != nil {
		panic(fmt.Errorf("Error connecting to nats-streaming server -> %s", err))
	}

	return &NatsStreaming{
		conn:          sc,
		subscriptions: make(map[string]stan.Subscription),
		Config:        cfg,
	}
}

// Publish - synchronously publishes a message for a given subject
func (ns *NatsStreaming) Publish(subject string, message []byte, opts ...interface{}) error {
	return ns.conn.Publish(subject, message)
}

// Subscribe - subscribes to a subject using a handler
func (ns *NatsStreaming) Subscribe(subject string, handler broker.Handler, opts ...interface{}) error {
	natsHandler, ok := handler.(stan.MsgHandler)
	if !ok {
		return fmt.Errorf("Invalid nats-streaming handler")
	}

	natsOpts := []stan.SubscriptionOption{}
	for _, opt := range opts {
		natOpt, ok := opt.(stan.SubscriptionOption)
		if !ok {
			return fmt.Errorf("Invalid nats-steaming subscription option")
		}

		natsOpts = append(natsOpts, natOpt)
	}

	sub, err := ns.conn.Subscribe(subject, natsHandler, natsOpts...)
	if err != nil {
		return fmt.Errorf("nats-streaming subscription failed")
	}

	ns.subscriptions[subject] = sub

	return nil
}

// Unsubscribe - unsubscribes from a given subject
func (ns *NatsStreaming) Unsubscribe(subject string) error {
	sub, ok := ns.subscriptions[subject]
	if !ok {
		return fmt.Errorf("Invalid subscription subject: %s", subject)
	}

	return sub.Unsubscribe()
}

// Close - closes an existing nats-streaming connection
func (ns *NatsStreaming) Close() error {
	return ns.conn.Close()
}
