package broker

import (
	"bytes"
	"encoding/gob"
)

// TraceMsgCarrier implements TextMapPropagator
type TraceMsgCarrier struct {
	Topic   string
	Message []byte
	Headers map[string]string
}

// NewTraceMsgCarrier creates a new instance of opentel TextMapPropagator
func NewTraceMsgCarrier(topic string, data []byte) *TraceMsgCarrier {
	return &TraceMsgCarrier{
		Topic:   topic,
		Message: data,
		Headers: make(map[string]string),
	}
}

// NewTraceMsgCarrierFromBytes converts carrier bytes to TraceMsgCarrier
func NewTraceMsgCarrierFromBytes(tmBytes []byte) *TraceMsgCarrier {
	data := bytes.NewBuffer(tmBytes)
	dec := gob.NewDecoder(data)

	tm := new(TraceMsgCarrier)
	if err := dec.Decode(tm); err != nil {
		return nil
	}

	return tm
}

// Get returns the key value from the headers property
func (tm *TraceMsgCarrier) Get(key string) string {
	return tm.Headers[key]
}

// Set sets the key value of the headers property
func (tm *TraceMsgCarrier) Set(key string, value string) {
	tm.Headers[key] = value
}

// Keys returns the list of keys in the headers
func (tm *TraceMsgCarrier) Keys() []string {
	keys := make([]string, 0, len(tm.Headers))
	for k := range tm.Headers {
		keys = append(keys, k)
	}

	return keys
}

// Bytes converts the TraceMsgCarrier instance to bytes
func (tm *TraceMsgCarrier) Bytes() ([]byte, error) {
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	if err := enc.Encode(tm); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}
