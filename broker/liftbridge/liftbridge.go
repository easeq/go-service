package liftbridge

import (
	"context"
	"fmt"
	"log"

	"github.com/easeq/go-service/broker"
	lift "github.com/liftbridge-io/go-liftbridge/v2"
)

type Liftbridge struct {
	client lift.Client
	*Config
}

// NewPostgres returns new connection to the postgres db
func NewLiftbridge(opts ...lift.ClientOption) broker.Broker {
	cfg := GetConfig()

	client, err := lift.Connect(cfg.Addrs, opts...)
	if err != nil {
		panic(fmt.Errorf("Error creating to liftbrige client -> %s", err))
	}

	return &Liftbridge{
		client: client,
		Config: cfg,
	}
}

// Initialize stream
func (l *Liftbridge) Init(ctx context.Context, args map[string]interface{}, opts ...interface{}) error {
	subject, ok := args["subject"].(string)
	if !ok {
		return fmt.Errorf("Liftbridge init error -> invalid subject")
	}

	stream, ok := args["stream"].(string)
	if !ok {
		return fmt.Errorf("Liftbridge init error -> invalid stream name")
	}

	liftOpts := []lift.StreamOption{}
	for _, opt := range opts {
		liftOpt, ok := opt.(lift.StreamOption)
		if !ok {
			return fmt.Errorf("Invalid liftbridge stream option")
		}

		liftOpts = append(liftOpts, liftOpt)
	}

	if err := l.client.CreateStream(ctx, subject, stream, liftOpts...); err != nil {
		return err
	}

	return nil
}

// Publish - synchronously publishes a message for a given subject
func (l *Liftbridge) Publish(ctx context.Context, stream string, message []byte) error {
	if _, err := l.client.Publish(ctx, stream, message); err != nil {
		return err
	}

	return nil
}

// Subscribe - subscribes to a subject using a handler
func (l *Liftbridge) Subscribe(
	ctx context.Context,
	stream string,
	handler broker.Handler,
	opts ...interface{},
) error {
	liftHandler, ok := handler.(lift.Handler)
	if !ok {
		return fmt.Errorf("Invalid liftbridge handler")
	}

	liftOpts := []lift.SubscriptionOption{}
	for _, opt := range opts {
		liftOpt, ok := opt.(lift.SubscriptionOption)
		if !ok {
			return fmt.Errorf("Invalid liftbridge subscription option")
		}

		liftOpts = append(liftOpts, liftOpt)
	}

	err := l.client.Subscribe(ctx, stream, liftHandler, liftOpts...)
	if err != nil {
		return fmt.Errorf("liftbridge subscription failed")
	}

	return nil
}

// Unsubscribe - unsubscribes from a given subject
func (l *Liftbridge) Unsubscribe(subject string) error {
	log.Println("Unsupport liftbridge method: Unsubscribe")
	return nil
}

// Close - closes an existing nats-streaming connection
func (l *Liftbridge) Close() error {
	return l.client.Close()
}
