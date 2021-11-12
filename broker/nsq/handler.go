package nsq

import (
	"github.com/easeq/go-service/broker"
	"github.com/nsqio/go-nsq"
)

type nsqHandler struct {
	handler broker.Handler
}

// NewNsqHandler creates a new nsq message Handler
func NewNsqHandler(handler broker.Handler) *nsqHandler {
	return &nsqHandler{handler}
}

// HandleMessage handles the nsq Message as a standard go-service broker Message.
func (n *nsqHandler) HandleMessage(message *nsq.Message) error {
	return n.handler.Handle(&broker.Message{
		Body:   message.Body,
		Extras: message,
	})
}
