package broker

import (
	"context"
	"testing"
)

func TestPublish(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		topic   string
		message Message
	}{
		{
			name:    "simpleString",
			ctx:     context.TODO(),
			topic:   "test-topic",
			message: "test-message",
		},
		{
			name:  "mapOfStringKV",
			ctx:   context.TODO(),
			topic: "test-interface-topic",
			message: map[string]string{
				"name":  "test-name",
				"value": "test-value",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
