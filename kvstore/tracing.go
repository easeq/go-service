package kvstore

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
	s          KVStore
	tracer     trace.Tracer
	attrs      []attribute.KeyValue
	propagator propagation.TextMapPropagator
}

func NewTrace(s KVStore) *Trace {
	return &Trace{
		s: s,
		tracer: otel.GetTracerProvider().Tracer(
			tracer.DEFAULT_TRACER_NAME,
		),
		attrs: []attribute.KeyValue{
			semconv.DBSystemKey.String(s.String()),
		},
		propagator: otel.GetTextMapPropagator(),
	}
}

func (t *Trace) Put(
	ctx context.Context,
	record *Record,
	put PutCallback,
	putOpts ...SetOpt,
) (*Record, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return put(ctx, record, putOpts...)
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.DBOperationKey.String("PUT"),
		),
	}

	opName := fmt.Sprintf("%s.PUT", t.s.String())
	ctx, span := t.tracer.Start(ctx, opName, opts...)
	defer span.End()

	record, err := put(ctx, record, putOpts...)
	t.setSpanError(span, err)

	return record, err
}

func (t *Trace) Get(
	ctx context.Context,
	key string,
	get GetCallback,
	getOpts ...GetOpt,
) ([]*Record, error) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return get(ctx, key, getOpts...)
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.DBOperationKey.String("GET"),
		),
	}

	opName := fmt.Sprintf("%s.GET.%s", t.s.String(), key)
	ctx, span := t.tracer.Start(ctx, opName, opts...)
	defer span.End()

	records, err := get(ctx, key, getOpts...)
	t.setSpanError(span, err)

	return records, err
}

func (t *Trace) Delete(
	ctx context.Context,
	key string,
	delete DeleteCallback,
) error {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return delete(ctx, key)
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.DBOperationKey.String("DELETE"),
		),
	}

	opName := fmt.Sprintf("%s.DELETE.%s", t.s.String(), key)
	ctx, span := t.tracer.Start(ctx, opName, opts...)
	defer span.End()

	err := delete(ctx, key)
	t.setSpanError(span, err)

	return err
}

func (t *Trace) Txn(
	ctx context.Context,
	handler TxnHandler,
	txn TxnCallback,
) error {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return txn(ctx, handler)
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.DBOperationKey.String("TXN"),
		),
	}

	opName := fmt.Sprintf("%s.TXN", t.s.String())
	ctx, span := t.tracer.Start(ctx, opName, opts...)
	defer span.End()

	err := txn(ctx, handler)
	t.setSpanError(span, err)

	return err
}

func (t *Trace) Subscribe(
	ctx context.Context,
	key string,
	handler SubscribeHandler,
	subscribe SubscribeCallback,
) error {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return subscribe(ctx, key, handler)
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.DBOperationKey.String("SUBSCRIBE"),
		),
	}

	opName := fmt.Sprintf("%s.SUBSCRIBE.%s", t.s.String(), key)
	ctx, span := t.tracer.Start(ctx, opName, opts...)
	defer span.End()

	err := subscribe(ctx, key, handler)
	t.setSpanError(span, err)

	return err
}

func (t *Trace) setSpanError(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
}
