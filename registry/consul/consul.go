package consul

import (
	"context"
	"errors"
	"fmt"

	goconfig "github.com/easeq/go-config"
	"github.com/easeq/go-service/registry"
	"github.com/easeq/go-service/server"

	"github.com/Netflix/go-env"
	"github.com/easeq/go-consul-registry/v2/consul"
)

var (
	// ErrConsulConfigLoad returned when env config for consul results in an error
	ErrConsulConfigLoad = errors.New("error loading consul config")
)

// Config - consul configuration
type Config struct {
	Host string `env:"CONSUL_HOST,default=localhost"`
	Port int    `env:"CONSUL_PORT,default=8500"`
	TTL  int    `env:"CONSUL_TTL,default=15"`
}

// Consul registry
type Consul struct {
	*Config
}

// UnmarshalEnv env.EnvSet to GatewayConfig
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

// NewConsul returns a new consul registry
func NewConsul() *Consul {
	return &Consul{
		Config: goconfig.NewEnvConfig(new(Config)).(*Config),
	}
}

// Register registers service with the registry.
func (c *Consul) Register(
	ctx context.Context,
	name string,
	server server.Server,
) *registry.ErrRegistryRegFailed {
	if err := consul.Register(
		ctx,
		name,
		server.Host(),
		server.Port(),
		c.Address(),
		c.TTL,
		server.RegistryTags()...,
	); err != nil {
		return &registry.ErrRegistryRegFailed{Value: err}
	}

	return nil
}

// ConnectionString returns the formatted connection string using the config loaded
func (c *Consul) ConnectionString(args ...interface{}) string {
	return fmt.Sprintf(
		"consul://%s/%s?scheme=%s",
		c.Address(),
		args[0],
		args[1],
	)
}

// Address returns the prepared consul address
func (c *Consul) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// ToString returns the string name of the service registry
func (c *Consul) ToString() string {
	return "consul"
}
