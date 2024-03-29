package nsq

import (
	"context"

	"github.com/easeq/go-service/broker"
	"github.com/nsqio/go-nsq"
)

type nsqHandler struct {
	ctx     context.Context
	n       *Nsq
	topic   string
	handler broker.Handler
}

// NewNsqHandler creates a new nsq message Handler
func NewNsqHandler(ctx context.Context, n *Nsq, topic string, handler broker.Handler) *nsqHandler {
	return &nsqHandler{ctx, n, topic, handler}
}

// HandleMessage handles the nsq Message as a standard go-service broker Message.
func (h *nsqHandler) HandleMessage(message *nsq.Message) error {
	return h.n.w.Subscribe(
		h.ctx,
		h.topic,
		message.Body,
		func(ctx context.Context, t *broker.TraceMsgCarrier) error {
			if err := h.handler.Handle(ctx, &broker.Message{
				Body: t.Message,
				Extras: map[string]interface{}{
					broker.KEY_TRACE_MSG_CARRIER: t,
					broker.KEY_BROKER_MSG:        message,
				},
			}); err != nil {
				return err
			}

			return nil
		},
	)
}
