package registry

import (
	"context"
	"fmt"
)

// ErrRegistryRegFailed returned when service registry with registry fails
type ErrRegistryRegFailed struct {
	value error
}

func (e *ErrRegistryRegFailed) Error() string {
	return fmt.Sprintf("service registration failed: [%s]", e.value)
}

// ServiceRegistry - service registry
type ServiceRegistry interface {
	// Registers the service
	Register(ctx context.Context, name string, host string, port int) *ErrRegistryRegFailed
	// Address returns the address of the registry
	Address() string
	// ConnectionString returns the full formatted connection string
	ConnectionString(...interface{}) string
	// Returns the string name of the registry
	ToString() string
}
