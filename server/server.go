package server

import (
	"context"

	"github.com/easeq/go-service/registry"
)

type Server interface {
	// Registers the server with the service registry
	Register(context.Context, string) *registry.ErrRegistryRegFailed
	// Runs the server
	Run(context.Context) error
}
