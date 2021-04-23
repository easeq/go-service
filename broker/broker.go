package broker

type Message interface{}

type Handler interface{}

type Broker interface {
	// Publish a message
	Publish(subject string, message []byte, opts ...interface{}) error
	// Subscribe to a subject
	Subscribe(subject string, handler Handler, opts ...interface{}) error
}
