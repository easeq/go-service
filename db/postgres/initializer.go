package postgres

import (
	"context"

	"github.com/easeq/go-service/logger"
)

type Initializer struct {
	pg *Postgres
}

// NewInitializer returns a new JetStream Initialiazer
func NewInitializer(pg *Postgres) *Initializer {
	return &Initializer{
		pg: pg,
	}
}

// AddDependency adds necessary service components as dependencies
func (i *Initializer) AddDependency(dep interface{}) error {
	switch v := dep.(type) {
	case logger.Logger:
		i.pg.logger = v
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
	// Run migrations
	if err := i.pg.Migrate(); err != nil {
		i.pg.logger.Debugw(
			"Database migration failed",
			"error", err,
		)
	}

	return nil
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanStop() bool {
	return false
}

// Stop - stops the running
func (i *Initializer) Stop(ctx context.Context) error {
	i.pg.logger.Infow("Closing db connection", "method", "Stop")
	return i.pg.Handle.Close()
}
