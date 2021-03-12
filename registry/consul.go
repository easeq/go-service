package registry

import (
	"context"
	"errors"
	"fmt"

	goconfig "github.com/easeq/go-config"

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

// NewRegistry returns a new consul registry
func NewRegistry() ServiceRegistry {
	return &Consul{
		Config: goconfig.NewEnvConfig(new(Config)).(*Config),
	}
}

// Register registers service with the registry
func (c *Consul) Register(ctx context.Context, name string, host string, port int) *ErrRegistryRegFailed {
	if err := consul.Register(
		ctx,
		name,
		host,
		port,
		c.Address(),
		c.TTL,
	); err != nil {
		return &ErrRegistryRegFailed{err}
	}

	return nil
}

// Address returns the prepared consul address
func (c *Consul) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetConnectionString returns consul connections string the given service
func (c *Consul) GetConnectionString(service string) string {
	return fmt.Sprintf(
		"consul://%s/%s?scheme=https",
		c.Address(),
		service,
	)
}

func (c *Consul) ToString() string {
	return "consul"
}
