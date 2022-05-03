package rest

import (
	"context"

	"github.com/easeq/go-service/logger"
)

type Initializer struct {
	r *Rest
}

// NewInitializer returns a new REST server initialiazer
func NewInitializer(r *Rest) *Initializer {
	return &Initializer{r}
}

// AddDependency adds necessary service components as dependencies
func (i *Initializer) AddDependency(dep interface{}) error {
	switch v := dep.(type) {
	case logger.Logger:
		i.r.logger = v
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
	i.r.logger.Infow(
		"Starting REST gateway...",
	)
	return i.r.App.Listen(i.r.Config.Address())
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanStop() bool {
	return true
}

// Run start the service component
func (i *Initializer) Stop(ctx context.Context) error {
	i.r.logger.Infow(
		"Shutting down HTTP/REST gRPC gateway...",
	)
	return i.r.App.Shutdown()
}
