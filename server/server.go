package server

import (
	"context"

	"github.com/easeq/go-service/pool"
	"github.com/easeq/go-service/registry"
)

type Server interface {
	// Registers the server with the service registry
	Register(context.Context, string, registry.ServiceRegistry) *registry.ErrRegistryRegFailed
	// Runs the server
	Run(context.Context) error
	// Client creates if not exists and returns the client to call the service
	GetClient(address string) (pool.Connection, error)
	// SetTags sets the registry tags for the server
	SetRegistryTags(tags ...string)
	// Get string identifier of the server
	String() string
	// Method to shut down server
	ShutDown(ctx context.Context) error
}
