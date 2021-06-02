package gateway

import "fmt"

// Metadata holds the metadata used by the gateway server
type Metadata struct {
	GrpcHost string `env:"GRPC_HOST,default="`
	GrpcPort int    `env:"GRPC_PORT,default=9090"`
}

// Get returns the gateway metadata by key
func (m *Metadata) Get(key string) interface{} {
	switch key {
	case "GRPC_HOST":
		return m.GrpcHost
	case "GRPC_PORT":
		return m.GrpcPort
	case "GRPC_ADDRESS":
		return fmt.Sprintf("%s:%d", m.GrpcHost, m.GrpcPort)
	}

	return nil
}
