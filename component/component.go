package component

import "context"

type Initializer interface {
	// AddDependency adds the service component dependency
	// requested by a service component
	AddDependency(dep interface{}) error
	// CanRun returns whether the service component has a Run function defined
	CanRun() bool
	// Run - runs/starts the service components
	Run(ctx context.Context) error
	// Dependencies returns the list of string service component dependency names
	Dependencies() []string
	// CanStop returns whether the component is stoppable
	CanStop() bool
	// Stop - stops/closes the service components
	Stop(ctx context.Context) error
}

type Component interface {
	// HasInitializer returns whether a service component has
	// an initializer defined. Every service component needs to define this method
	HasInitializer() bool
	// Initializer returns the initializer for the component
	Initializer() Initializer
}
