package goservice

import (
	"context"

	"github.com/Netflix/go-env"
	goconfig "github.com/easeq/go-config"
	_ "github.com/easeq/go-consul-registry/v2/consul"
	"github.com/easeq/go-service/server"
	goservice_grpc "github.com/easeq/go-service/server/grpc"
)

// ServiceOption to pass as arg while creating new service
type ServiceOption func(*ServiceConfig)

// Config - Service configuration
type Config struct {
	Name string `env:"SERVICE_NAME"`
}

// UnmarshalEnv env.EnvSet to GatewayConfig
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

// ServiceConfig handles config required by the service
type ServiceConfig struct {
	Server server.Server
	*Config
}

// NewService creates a new service
func NewService(opts ...ServiceOption) *ServiceConfig {
	cfg := new(Config)
	cfg.UnmarshalEnv(goconfig.EnvSet())

	sc := &ServiceConfig{
		Config: cfg,
	}

	for _, opt := range opts {
		opt(sc)
	}

	if sc.Server == nil {
		sc.Server = goservice_grpc.NewGrpc()
	}

	return sc
}

// WithServer passes the server
func WithServer(server server.Server) ServiceOption {
	return func(s *ServiceConfig) {
		s.Server = server
	}
}

// Run runs both the HTTP and gRPC server
func (s *ServiceConfig) Run(ctx context.Context) error {
	if err := s.Server.Register(ctx, s.Config.Name); err != nil {
		return err
	}

	return s.Server.Run(ctx)
}
