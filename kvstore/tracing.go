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

func (t *Trace) start(ctx context.Context, op string, key string) (context.Context, trace.Span) {
	if !trace.SpanFromContext(ctx).IsRecording() {
		return ctx, nil
	}

	opts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(t.attrs...),
		trace.WithAttributes(
			semconv.DBOperationKey.String(op),
		),
	}

	opName := fmt.Sprintf("%s.%s", t.s.String(), op)
	if key == "" {
		opName = fmt.Sprintf("%s.%s", opName, key)
	}

	return t.tracer.Start(ctx, opName, opts...)
}

func (t *Trace) Put(ctx context.Context, record *Record, put PutCallback, opts ...SetOpt) (*Record, error) {
	ctx, span := t.start(ctx, "PUT", "")
	if span == nil {
		return put(ctx, record, opts...)
	}
	defer span.End()

	record, err := put(ctx, record, opts...)
	t.setSpanError(span, err)

	return record, err
}

func (t *Trace) Get(ctx context.Context, key string, get GetCallback, opts ...GetOpt) ([]*Record, error) {
	ctx, span := t.start(ctx, "GET", key)
	if span == nil {
		return get(ctx, key, opts...)
	}
	defer span.End()

	records, err := get(ctx, key, opts...)
	t.setSpanError(span, err)

	return records, err
}

func (t *Trace) Delete(ctx context.Context, key string, delete DeleteCallback) error {
	ctx, span := t.start(ctx, "DELETE", key)
	if span == nil {
		return delete(ctx, key)
	}
	defer span.End()

	err := delete(ctx, key)
	t.setSpanError(span, err)

	return err
}

func (t *Trace) Txn(ctx context.Context, handler TxnHandler, txn TxnCallback) error {
	ctx, span := t.start(ctx, "TXN", "")
	if span == nil {
		return txn(ctx, handler)
	}
	defer span.End()

	err := txn(ctx, handler)
	t.setSpanError(span, err)

	return err
}

func (t *Trace) Subscribe(ctx context.Context, key string, handler SubscribeHandler, subscribe SubscribeCallback) error {
	ctx, span := t.start(ctx, "SUBSCRIBE", "")
	if span == nil {
		return subscribe(ctx, key, handler)
	}
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
