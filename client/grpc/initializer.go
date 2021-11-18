package grpc

import (
	"context"

	"github.com/easeq/go-service/logger"
)

type Initializer struct {
	g *Grpc
}

// NewInitializer returns a new JetStream Initialiazer
func NewInitializer(g *Grpc) *Initializer {
	return &Initializer{
		g: g,
	}
}

// AddDependency adds necessary service components as dependencies
func (i *Initializer) AddDependency(dep interface{}) error {
	switch v := dep.(type) {
	case logger.Logger:
		i.g.logger = v
	}

	return nil
}

// Dependencies returns the string names of service components
// that are required as dependencies for this component
func (i *Initializer) Dependencies() []string {
	return []string{"logger"}
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanRun() bool {
	return false
}

// Run start the service component
func (i *Initializer) Run(ctx context.Context) error {
	i.g.logger.Infow("Unimplemented", "method", "Run", "package", "goservice.client.grpc")
	return nil
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanStop() bool {
	return false
}

// Stop - stops the running
func (i *Initializer) Stop(ctx context.Context) error {
	i.g.logger.Infow("Unimplemented", "method", "Stop", "package", "goservice.client.grpc")
	return nil
}
