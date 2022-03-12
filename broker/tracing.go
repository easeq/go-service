package broker

import (
	"context"
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

// Publish adds tracer details to the message
func (t *Trace) Publish(ctx context.Context, tm *TraceMsgCarrier, publish PublishCallback) error {
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
			semconv.MessagingDestinationKey.String(tm.Topic),
		),
	}

	opName := fmt.Sprintf("%s.publish %s", t.b.String(), tm.Topic)
	ctx, span := t.tracer.Start(ctx, opName, opts...)
	defer span.End()

	t.propagator.Inject(ctx, tm)

	err := publish(tm)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}

// Subscribe adds the tracer details to the context
func (t *Trace) Subscribe(ctx context.Context, tm *TraceMsgCarrier, subscribe SubscribeCallback) error {
	ctx = t.propagator.Extract(ctx, tm)

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.MessageTypeReceived,
			semconv.MessagingDestinationKindTopic,
			semconv.MessageTypeKey.String(tm.Topic),
		),
	}

	opName := fmt.Sprintf("%s.subscribe %s", t.b.String(), tm.Topic)
	ctx, span := t.tracer.Start(ctx, opName, opts...)
	defer span.End()

	t.propagator.Inject(ctx, tm)

	err := subscribe(ctx, tm)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}
