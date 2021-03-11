package registry

import "context"

// ServiceRegistry - service registry
type ServiceRegistry interface {
	Register(ctx context.Context, name string, host string, port int) error
}
