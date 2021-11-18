package etcd

import (
	"context"

	"github.com/easeq/go-service/logger"
)

type Initializer struct {
	e *Etcd
}

// NewInitializer returns a new JetStream Initialiazer
func NewInitializer(e *Etcd) *Initializer {
	return &Initializer{
		e: e,
	}
}

// AddDependency adds necessary service components as dependencies
func (i *Initializer) AddDependency(dep interface{}) error {
	switch v := dep.(type) {
	case logger.Logger:
		i.e.logger = v
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
	return true
}

// Run start the service component
func (i *Initializer) Run(ctx context.Context) error {
	i.e.logger.Infow("Unimplemented", "method", "Run", "package", "goservice.kvstore.etcd")
	return nil
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanStop() bool {
	return false
}

// Stop - stops the running
func (i *Initializer) Stop(ctx context.Context) error {
	i.e.logger.Infow("Closing etcd watcher", "method", "Stop")
	i.e.Client.Watcher.Close()

	i.e.logger.Infow("Closing etcd client", "method", "Stop")
	return i.e.Client.Close()
}
