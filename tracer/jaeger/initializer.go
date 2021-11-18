package jaeger

import (
	"context"

	"github.com/easeq/go-service/logger"
)

type Initializer struct {
	j *Jaeger
}

// NewInitializer returns a new JetStream Initialiazer
func NewInitializer(j *Jaeger) *Initializer {
	return &Initializer{j}
}

// AddDependency adds necessary service components as dependencies
func (i *Initializer) AddDependency(dep interface{}) error {
	switch v := dep.(type) {
	case logger.Logger:
		i.j.logger = v
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
	i.j.logger.Infow("Unimplemented", "method", "Run", "package", "goservice.tracer.jaeger")
	return nil
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanStop() bool {
	return false
}

// Run start the service component
func (i *Initializer) Stop(ctx context.Context) error {
	i.j.logger.Infow("Stop Jaeger", "method", "Stop", "package", "goservice.tracer.jaeger")
	return i.j.closer.Close()
}
