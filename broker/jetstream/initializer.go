package jetstream

import (
	"context"

	"github.com/easeq/go-service/logger"
	"github.com/easeq/go-service/tracer"
)

type Initializer struct {
	j *JetStream
}

// NewInitializer returns a new JetStream Initialiazer
func NewInitializer(j *JetStream) *Initializer {
	return &Initializer{
		j: j,
	}
}

// AddDependency adds necessary service components as dependencies
func (i *Initializer) AddDependency(dep interface{}) error {
	switch v := dep.(type) {
	case logger.Logger:
		i.j.logger = v
	case tracer.Tracer:
		i.j.tracer = v
	}

	return nil
}

// Dependencies returns the string names of service components
// that are required as dependencies for this component
func (i *Initializer) Dependencies() []string {
	return []string{logger.LOGGER, tracer.TRACER}
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanRun() bool {
	return false
}

// Run start the service component
func (i *Initializer) Run(ctx context.Context) error {
	i.j.logger.Infow("Unimplemented", "method", "goservice.broker.jetstream.Run")
	return nil
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanStop() bool {
	return true
}

// Stop - stops the running
func (i *Initializer) Stop(ctx context.Context) error {
	for _, sub := range i.j.Subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			i.j.logger.Errorw(
				"JetStream close connection error: %s",
				"error", err,
				"method", "goservice.broker.jetstream.Stop",
			)
			return err
		}
	}

	i.j.nc.Close()
	return nil
}
