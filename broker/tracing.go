package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/easeq/go-service/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type Trace struct {
	b          Broker
	tracer     trace.Tracer
	attrs      []attribute.KeyValue
	propagator propagation.TextMapPropagator
}

func NewTrace(b Broker) *Trace {
	return &Trace{
		b: b,
		tracer: otel.GetTracerProvider().Tracer(
			tracer.DEFAULT_TRACER_NAME,
		),
		attrs: []attribute.KeyValue{
			semconv.MessagingSystemKey.String(b.String()),
		},
		propagator: otel.GetTextMapPropagator(),
	}
}

func (t *Trace) Publish(ctx context.Context, topic string, payload []byte, publish func(*TraceMsgCarrier) error) error {
	tm := NewTraceMsgCarrier(topic, payload)
	if !trace.SpanFromContext(ctx).IsRecording() {
		return publish(tm)
	}

	ctx = t.propagator.Extract(ctx, tm)

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.MessageTypeSent,
			semconv.MessagingDestinationKindTopic,
			semconv.MessagingDestinationKey.String(topic),
		),
	}

	opName := fmt.Sprintf("%s.publish %s", t.b.String(), topic)
	ctx, span := t.tracer.Start(ctx, opName, opts...)
	t.propagator.Inject(ctx, tm)
	defer span.End()

	err := publish(tm)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}

// TraceSubscribe starts a trace on message receive
func (t *Trace) Subscribe(ctx context.Context, topic string, tmBytes []byte, subscribe func(*TraceMsgCarrier) error) error {
	tm := NewTraceMsgCarrierFromBytes(tmBytes)
	if tm == nil {
		return errors.New("payload empty")
	}

	ctx = t.propagator.Extract(ctx, tm)

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.MessageTypeReceived,
			semconv.MessagingDestinationKindTopic,
			semconv.MessageTypeKey.String(topic),
		),
	}

	opName := fmt.Sprintf("%s.subscribe %s", t.b.String(), topic)
	_, span := t.tracer.Start(ctx, opName, opts...)
	defer span.End()

	err := subscribe(tm)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}
