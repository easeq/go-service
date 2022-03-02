package broker

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/easeq/go-service/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

type TraceMsgCarrier struct {
	Topic   string
	Message []byte
	Headers map[string]string
}

func NewTraceMsgCarrier(topic string, data []byte) *TraceMsgCarrier {
	return &TraceMsgCarrier{
		Topic:   topic,
		Message: data,
		Headers: make(map[string]string),
	}
}

func NewTraceMsgCarrierFromBytes(tmBytes []byte) *TraceMsgCarrier {
	var data bytes.Buffer
	dec := gob.NewDecoder(&data)

	var tm *TraceMsgCarrier
	if err := dec.Decode(&tm); err != nil {
		return nil
	}

	return tm
}

func (tm *TraceMsgCarrier) Get(key string) string {
	return tm.Headers[key]
}

func (tm *TraceMsgCarrier) Set(key string, value string) {
	tm.Headers[key] = value
}

func (tm *TraceMsgCarrier) Keys() []string {
	keys := make([]string, 0, len(tm.Headers))
	for k := range tm.Headers {
		keys = append(keys, k)
	}

	return keys
}

func (tm *TraceMsgCarrier) Bytes() ([]byte, error) {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	if err := enc.Encode(tm); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
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
			semconv.MessagingDestinationKindTopic,
			semconv.MessagingDestinationKey.String(topic),
		),
	}

	opName := fmt.Sprintf("%s.publish %s", t.b.String(), topic)
	ctx, span := t.tracer.Start(ctx, opName, opts...)
	t.propagator.Inject(ctx, tm)
	defer span.End()

	return publish(tm)
}

// TraceSubscribe starts a trace on message receive
func (t *Trace) Subscribe(ctx context.Context, topic string, tmBytes []byte, subscribe func(*TraceMsgCarrier) error) error {
	tm := NewTraceMsgCarrierFromBytes(tmBytes)

	if !trace.SpanFromContext(ctx).IsRecording() {
		return subscribe(tm)
	}

	ctx = t.propagator.Extract(ctx, tm)

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.MessageTypeKey.String(topic),
		),
	}

	opName := fmt.Sprintf("%s.subscribe %s", t.b.String(), topic)
	_, span := t.tracer.Start(ctx, opName, opts...)
	defer span.End()

	return subscribe(tm)
}
