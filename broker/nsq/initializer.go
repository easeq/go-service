package nsq

import (
	"context"

	"github.com/easeq/go-service/logger"
	"github.com/easeq/go-service/tracer"
)

type Initializer struct {
	n *Nsq
}

// NewInitializer returns a new JetStream Initialiazer
func NewInitializer(n *Nsq) *Initializer {
	return &Initializer{n}
}

// AddDependency adds necessary service components as dependencies
func (i *Initializer) AddDependency(dep interface{}) error {
	switch v := dep.(type) {
	case logger.Logger:
		i.n.logger = v
	case tracer.Tracer:
		i.n.tracer = v
	}

	return nil
}

// Dependencies returns the string names of service components
// that are required as dependencies for this component
func (i *Initializer) Dependencies() []string {
	return []string{"logger", "tracer"}
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanRun() bool {
	return false
}

// Run start the service component
func (i *Initializer) Run(ctx context.Context) error {
	i.n.logger.Infow("Unimplemented", "method", "Run", "package", "goservice.broker.nsq")
	return nil
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanStop() bool {
	return true
}

// Stop - stops the running
func (i *Initializer) Stop(ctx context.Context) error {
	i.n.Producer.Stop()

	for _, consumer := range i.n.Consumers {
		consumer.Stop()
	}

	return nil
}
