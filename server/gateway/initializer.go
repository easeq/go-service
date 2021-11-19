package gateway

import (
	"context"

	"github.com/easeq/go-service/logger"
)

type Initializer struct {
	g *Gateway
}

// NewInitializer returns a new JetStream Initialiazer
func NewInitializer(g *Gateway) *Initializer {
	return &Initializer{g}
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
	return []string{logger.LOGGER}
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanRun() bool {
	return true
}

// Run start the service component
func (i *Initializer) Run(ctx context.Context) error {
	i.g.logger.Infow(
		"starting HTTP/REST gRPC gateway...",
		"method", "goservice.server.gateway.Run",
	)
	return i.g.Server.ListenAndServe()
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanStop() bool {
	return true
}

// Run start the service component
func (i *Initializer) Stop(ctx context.Context) error {
	i.g.logger.Infow(
		"shutting down HTTP/REST gRPC gateway...",
		"method", "goservice.server.gateway.Stop",
	)
	return i.g.Server.Shutdown(ctx)
}
