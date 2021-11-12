package jetstream

import (
	"github.com/easeq/go-service/broker"
	nats "github.com/nats-io/nats.go"
)

// JetStreamPublisher holds additional options for jetstream publishing
type JetStreamPublisher struct {
	js       *JetStream
	stream   string
	subjects []string
}

// NewJetStreamPublisher returns a new publisher instance for Nats publishing
func NewJetStreamPublisher(js *JetStream, opts ...broker.PublishOption) *JetStreamPublisher {
	publisher := &JetStreamPublisher{
		js: js,
	}

	for _, opt := range opts {
		opt(publisher)
	}

	return publisher
}

// WithStream defines a the stream in which to publish the message
func WithStream(name string) broker.PublishOption {
	return func(p broker.Publisher) {
		p.(*JetStreamPublisher).stream = name
	}
}

// WithStreamSubjects defines a the stream in which to publish the message
func WithStreamSubjects(subjects ...string) broker.PublishOption {
	return func(p broker.Publisher) {
		p.(*JetStreamPublisher).subjects = subjects
	}
}

// createStream creates a new JS stream if it doens't exist and
// attaches the pre-defined subjects to the stream
func (jp *JetStreamPublisher) createStream() error {
	// Check if the ORDERS stream already exists; if not, create it.
	js := jp.js.jsCtx
	stream, err := js.StreamInfo(jp.stream)
	if err != nil {
		return err
	}

	if stream == nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     jp.stream,
			Subjects: jp.subjects,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
