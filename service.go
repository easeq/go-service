package goservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/Netflix/go-env"
	goconfig "github.com/easeq/go-config"
	"github.com/easeq/go-redis-access-control/manager"
	"github.com/easeq/go-service/db"
	goservice_grpc "github.com/easeq/go-service/grpc"
	goservice_registry "github.com/easeq/go-service/registry"
)

var (
	// ErrDatabaseNameNotProvided returned when database name is not provided
	ErrDatabaseNameNotProvided = errors.New("database name not provided")
)

// ServiceOption to pass as arg while creating new service
type ServiceOption func(*ServiceConfig)

// Config - Service configuration
type Config struct {
	Name string `env:"SERVICE_NAME"`
	Grac manager.Config
}

// UnmarshalEnv env.EnvSet to GatewayConfig
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

// ServiceConfig handles config required by the service
type ServiceConfig struct {
	Database db.ServiceDatabase
	Grpc     *goservice_grpc.Grpc
	Registry goservice_registry.ServiceRegistry
	*Config
}

// NewService creates a new service
func NewService(opts ...ServiceOption) *ServiceConfig {
	cfg := new(Config)
	cfg.UnmarshalEnv(goconfig.EnvSet())

	sc := &ServiceConfig{
		Grpc:     goservice_grpc.NewGrpc(),
		Registry: goservice_registry.NewRegistry(),
		Config:   cfg,
	}

	for _, opt := range opts {
		opt(sc)
	}

	return sc
}

// WithGrpc passes gRPC as option to servic
func WithGrpc(g *goservice_grpc.Grpc) ServiceOption {
	return func(s *ServiceConfig) {
		s.Grpc = g
	}
}

// Run runs both the HTTP and gRPC server
func (s *ServiceConfig) Run(ctx context.Context) error {
	if err := s.Registry.Register(ctx, s.Config.Name, s.Grpc.Host, s.Grpc.Port); err != nil {
		return fmt.Errorf("Consul registration failed: %s", err)
	}

	return s.Grpc.Run(ctx)
}
