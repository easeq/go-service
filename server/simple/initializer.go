package simple

import (
	"context"

	"github.com/easeq/go-service/logger"
)

type Initializer struct {
	s *Simple
}

// NewInitializer returns a new JetStream Initialiazer
func NewInitializer(s *Simple) *Initializer {
	return &Initializer{s}
}

// AddDependency adds necessary service components as dependencies
func (i *Initializer) AddDependency(dep interface{}) error {
	switch v := dep.(type) {
	case logger.Logger:
		i.s.logger = v
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
	i.s.logger.Infow(
		"Running simple server...",
		"method", "goservice.server.simple.Run",
	)
	<-ctx.Done()
	return nil
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanStop() bool {
	return false
}

// Run start the service component
func (i *Initializer) Stop(ctx context.Context) error {
	i.s.logger.Infow("Unimplemented", "method", "goservice.server.simple.Stop")
	return nil
}
