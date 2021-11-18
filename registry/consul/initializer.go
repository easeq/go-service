package consul

import (
	"context"

	"github.com/easeq/go-service/logger"
)

type Initializer struct {
	c *Consul
}

// NewInitializer returns a new JetStream Initialiazer
func NewInitializer(c *Consul) *Initializer {
	return &Initializer{c}
}

// AddDependency adds necessary service components as dependencies
func (i *Initializer) AddDependency(dep interface{}) error {
	switch v := dep.(type) {
	case logger.Logger:
		i.c.logger = v
	}

	return nil
}

// Dependencies returns the string names of service components
// that are required as dependencies for this component
func (i *Initializer) Dependencies() []string {
	return []string{"logger", "server"}
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanRun() bool {
	return true
}

// Run start the service component
func (i *Initializer) Run(ctx context.Context) error {
	return i.c.Register(ctx, "<service-name>", i.c.server)
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanStop() bool {
	return false
}

// Stop - stops the running
func (i *Initializer) Stop(ctx context.Context) error {
	i.c.logger.Infow("Unimplemented", "method", "Stop", "package", "goservice.client.grpc")
	return nil
}
