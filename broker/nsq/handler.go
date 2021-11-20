package nsq

import (
	"github.com/easeq/go-service/broker"
	"github.com/nsqio/go-nsq"
)

type nsqHandler struct {
	n       *Nsq
	topic   string
	handler broker.Handler
}

// NewNsqHandler creates a new nsq message Handler
func NewNsqHandler(n *Nsq, topic string, handler broker.Handler) *nsqHandler {
	return &nsqHandler{n, topic, handler}
}

// HandleMessage handles the nsq Message as a standard go-service broker Message.
func (h *nsqHandler) HandleMessage(message *nsq.Message) error {
	return h.n.t.Subscribe(h.topic, message.Body, func(body []byte) error {
		if err := h.handler.Handle(&broker.Message{
			Body:   body,
			Extras: message,
		}); err != nil {
			return err
		}

		return nil
	})
}
