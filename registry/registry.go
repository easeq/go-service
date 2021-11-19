package registry

import (
	"context"

	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/server"
)

const (
	REGISTRY = "registry"
	// TAGS_SEPARATOR is the separator used to split the tags passed in the tag env var for the specific service registry.
	TAGS_SEPARATOR = ","
)

// ServiceRegistry - service registry
type ServiceRegistry interface {
	component.Component
	// Registers the service server
	Register(ctx context.Context, name string, server server.Server) error
	// Address returns the address of the registry
	Address() string
	// ConnectionString returns the full formatted connection string
	ConnectionString(...interface{}) string
	// Returns the string name of the registry
	ToString() string
}
