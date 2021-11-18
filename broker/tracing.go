package broker

import (
	"bytes"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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
func (t *Trace) Publish(topic string, publish func(*TraceMsg) error) error {
	var tm TraceMsg

	operationName := fmt.Sprintf("Publish message (%s)", topic)
	span := opentracing.StartSpan(operationName, ext.SpanKindProducer)
	if span != nil {
		defer span.Finish()

		ext.MessageBusDestination.Set(span, topic)
		if err := opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.Binary,
			&tm,
		); err != nil {
			t.b.Logger().Debugw(
				"ERROR: Injecting tracer",
				"method", "goservice.broker.Publish",
				"error", err,
			)
			return err
		}
	}

	return publish(&tm)
}

// TraceSubscribe starts a trace on message receive
func (t *Trace) Subscribe(topic string, dataWithSpanCtx []byte, subscribe func([]byte) error) error {
	tm := NewTraceMsg(dataWithSpanCtx)

	// Extract the span context.
	sc, err := opentracing.GlobalTracer().Extract(opentracing.Binary, tm)
	if err != nil {
		t.b.Logger().Debugw(
			"ERROR: Extracting span from tracer",
			"method", "goservice.broker.Subscribe",
			"error", err,
		)
		return err
	}

	operationName := fmt.Sprintf("Receive message (%s)", topic)
	span := opentracing.StartSpan(
		operationName,
		ext.SpanKindConsumer,
		opentracing.FollowsFrom(sc),
	)
	if span != nil {
		defer span.Finish()
		ext.MessageBusDestination.Set(span, topic)
	}

	return subscribe(tm.Bytes())
}
