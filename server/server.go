package server

import (
	"context"
)

// Metadata interface for getting the server related metadata
type Metadata interface {
	// Get the value for the given key
	Get(key string) string
}

// Server interface for implementing custom servers
type Server interface {
	// // Registers the server with the service registry
	// Register(context.Context, string, registry.ServiceRegistry) *registry.ErrRegistryRegFailed
	// Runs the server
	Run(context.Context) error
	// Host returns server host
	Host() string
	// Port returns server port
	Port() int
	// Address returns the server address
	Address() string
	// RegistryTags returns server registry tags
	RegistryTags() []string
	// GetMetdata returns the server metadata
	GetMetadata(key string) interface{}
	// AddRegistryTags appends new tags to the existing tags slice
	AddRegistryTags(tags ...string)
	// Get string identifier of the server
	String() string
	// Method to shut down server
	ShutDown(ctx context.Context) error
}
