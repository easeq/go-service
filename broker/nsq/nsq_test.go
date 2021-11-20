package nsq

import (
	"context"
	"testing"
)

func TestPublish(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		topic   string
		message []byte
	}{
		{
			name:    "simpleString",
			ctx:     context.TODO(),
			topic:   "test-topic",
			message: []byte("test-message"),
		},
		// {
		// 	name:  "mapOfStringKV",
		// 	ctx:   context.TODO(),
		// 	topic: "test-interface-topic",
		// 	message: map[string]interface{}{
		// 		"name":  "test-name",
		// 		"value": "test-value",
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
