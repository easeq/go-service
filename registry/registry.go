package registry

import (
	"context"
	"fmt"
)

const (
	TAGS_SEPARATOR = ","
)

// ErrRegistryRegFailed returned when service registry with registry fails
type ErrRegistryRegFailed struct {
	Value error
}

func (e *ErrRegistryRegFailed) Error() string {
	return fmt.Sprintf("service registration failed: [%s]", e.Value)
}

// ServiceRegistry - service registry
type ServiceRegistry interface {
	// Registers the service
	Register(ctx context.Context, name string, host string, port int, tags ...string) *ErrRegistryRegFailed
	// Address returns the address of the registry
	Address() string
	// ConnectionString returns the full formatted connection string
	ConnectionString(...interface{}) string
	// Returns the string name of the registry
	ToString() string
}
