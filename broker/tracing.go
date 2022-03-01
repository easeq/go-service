package broker

import (
	"bytes"
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
)

type Trace struct {
	b Broker
}

func NewTrace(b Broker) *Trace {
	return &Trace{b}
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

	tracer := otel.Tracer("tracer")
	operationName := fmt.Sprintf("Publish message (%s)", topic)
	ctx, span := tracer.Start(ctx, operationName)
	defer span.End()

	return publish(&tm)
}

// TraceSubscribe starts a trace on message receive
func (t *Trace) Subscribe(ctx context.Context, topic string, dataWithSpanCtx []byte, subscribe func([]byte) error) error {
	tracer := otel.Tracer("tracer")
	operationName := fmt.Sprintf("Receive message (%s)", topic)
	ctx, span := tracer.Start(ctx, operationName)
	defer span.End()

	tm := NewTraceMsg(dataWithSpanCtx)

	return subscribe(tm.Bytes())
}
