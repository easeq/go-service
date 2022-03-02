package broker

import (
	"bytes"
	"context"
	"fmt"

	"github.com/easeq/go-service/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type Trace struct {
	b      Broker
	tracer trace.Tracer
	attrs  []attribute.KeyValue
}

func NewTrace(b Broker) *Trace {
	return &Trace{
		b: b,
		tracer: otel.GetTracerProvider().Tracer(
			tracer.DEFAULT_TRACER_NAME,
		),
		attrs: []attribute.KeyValue{
			{
				Key:   semconv.MessagingSystemKey,
				Value: attribute.StringValue(b.String()),
			},
		},
	}
}

// TraceMsg is used to trace the message using an opentracing
type TraceMsg struct {
	bytes.Buffer
}

// Prepare a trace message with the broker sent bytes
func NewTraceMsg(data []byte) *TraceMsg {
	b := bytes.NewBuffer(data)
	return &TraceMsg{*b}
}

// TracePublish starts a trace on message publish
func (t *Trace) Publish(ctx context.Context, topic string, publish func(*TraceMsg) error) error {
	var tm TraceMsg

	if !trace.SpanFromContext(ctx).IsRecording() {
		return publish(&tm)
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.MessageTypeKey.String(topic),
		),
	}

	opName := fmt.Sprintf("Publish message (%s)", topic)
	_, span := t.tracer.Start(ctx, opName, opts...)
	defer span.End()

	return publish(&tm)
}

// TraceSubscribe starts a trace on message receive
func (t *Trace) Subscribe(ctx context.Context, topic string, dataWithSpanCtx []byte, subscribe func([]byte) error) error {
	tm := NewTraceMsg(dataWithSpanCtx)

	if !trace.SpanFromContext(ctx).IsRecording() {
		return subscribe(tm.Bytes())
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.MessageTypeKey.String(topic),
		),
	}

	opName := fmt.Sprintf("Receive message (%s)", topic)
	_, span := t.tracer.Start(ctx, opName, opts...)
	defer span.End()

	return subscribe(tm.Bytes())
}
