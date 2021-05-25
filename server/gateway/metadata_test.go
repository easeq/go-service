package gateway

import (
	"testing"
)

func TestGet(t *testing.T) {
	metadata := Metadata{
		GrpcHost: "localhost",
		GrpcPort: 9090,
	}

	tests := []struct {
		name string
		data Metadata
		want interface{}
	}{
		{
			name: "GRPC_HOST",
			data: metadata,
			want: "localhost",
		}, {
			name: "GRPC_PORT",
			data: metadata,
			want: 9090,
		}, {
			name: "GRPC_ADDRESS",
			data: metadata,
			want: "localhost:9090",
		}, {
			name: "NIL_RESULT",
			data: metadata,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.data.Get(tt.name)
			if result != tt.want {
				t.Errorf("Metadata Get(%s) was incorrect, got: %v, want: %v", tt.name, result, tt.want)
			}
		})
	}
}
