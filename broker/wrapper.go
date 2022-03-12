package broker

import (
	"context"
	"errors"
)

type Wrapper struct {
	b     Broker
	trace *Trace
}

type PublishCallback func(*TraceMsgCarrier) error
type SubscribeCallback func(context.Context, *TraceMsgCarrier) error

func NewWrapper(b Broker) *Wrapper {
	return &Wrapper{b, NewTrace(b)}
}

// Publish - publishes a message with the traceparent if tracer is defined
func (w *Wrapper) Publish(ctx context.Context, topic string, payload []byte, publish PublishCallback) error {
	tm := NewTraceMsgCarrier(topic, payload)
	if w.trace == nil {
		return publish(tm)
	}

	return w.trace.Publish(ctx, tm, publish)
}

// Subscribe - subscribes to a message and adds traceparent to the ctx if tracer is defined
func (w *Wrapper) Subscribe(ctx context.Context, topic string, tmBytes []byte, subscribe SubscribeCallback) error {
	tm := NewTraceMsgCarrierFromBytes(tmBytes)
	if tm == nil {
		return errors.New("payload empty")
	}

	if w.trace == nil {
		return subscribe(ctx, tm)
	}

	return w.trace.Subscribe(ctx, tm, subscribe)
}
