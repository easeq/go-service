package server

import (
	"github.com/easeq/go-service/component"
)

const (
	SERVER = "server"
)

// Metadata interface for getting the server related metadata
type Metadata interface {
	// Get the value for the given key
	Get(key string) string
}

// Server interface for implementing custom servers
type Server interface {
	component.Component
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
}
