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
	opName := fmt.Sprintf("%s.publish %s", t.b.String(), tm.Topic)
	ctx, span := t.start(ctx, tm, opName, trace.SpanKindProducer, []attribute.KeyValue{
		semconv.MessageTypeSent,
	})
	defer span.End()

	err := publish(tm)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}

// Subscribe adds the tracer details to the context
func (t *Trace) Subscribe(ctx context.Context, tm *TraceMsgCarrier, subscribe SubscribeCallback) error {
	opName := fmt.Sprintf("%s.subscribe %s", t.b.String(), tm.Topic)
	ctx, span := t.start(ctx, tm, opName, trace.SpanKindConsumer, []attribute.KeyValue{
		semconv.MessageTypeReceived,
	})
	defer span.End()

	err := subscribe(ctx, tm)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}

func (t *Trace) start(
	ctx context.Context,
	tm *TraceMsgCarrier,
	opName string,
	kind trace.SpanKind,
	attrs []attribute.KeyValue,
) (context.Context, trace.Span) {
	ctx = t.propagator.Extract(ctx, tm)

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(kind),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(attrs...),
		trace.WithAttributes(
			semconv.MessagingDestinationKindTopic,
			semconv.MessageTypeKey.String(tm.Topic),
		),
	}

	ctx, span := t.tracer.Start(ctx, opName, opts...)
	t.propagator.Inject(ctx, tm)

	return ctx, span
}
