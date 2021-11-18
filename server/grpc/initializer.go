package grpc

import (
	"context"
	"log"
	"net"

	"github.com/easeq/go-service/logger"
)

type Initializer struct {
	g *Grpc
}

// NewInitializer returns a new JetStream Initialiazer
func NewInitializer(g *Grpc) *Initializer {
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
	return []string{"logger"}
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanRun() bool {
	return true
}

// Run start the service component
func (i *Initializer) Run(ctx context.Context) error {
	i.g.logger.Infow(
		"Starting gRPC server",
		"method", "Run",
		"package", "goservice.server.grpc",
	)
	listener, err := net.Listen("tcp", i.g.Config.Address())
	if err != nil {
		return err
	}

	// start gRPC server
	log.Println("Starting gRPC server...")
	return i.g.Server.Serve(listener)
}

// CanRun returns true if the component has anything to Run
func (i *Initializer) CanStop() bool {
	return true
}

// Run start the service component
func (i *Initializer) Stop(ctx context.Context) error {
	i.g.logger.Infow(
		"gracefully stop gRPC server",
		"method", "Stop",
		"package", "goservice.server.gateway",
	)

	i.g.Server.GracefulStop()
	return nil
}
