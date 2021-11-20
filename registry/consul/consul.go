package consul

import (
	"context"
	"errors"
	"fmt"

	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/logger"
	"github.com/easeq/go-service/server"

	"github.com/Netflix/go-env"
	"github.com/easeq/go-consul-registry/v2/consul"
)

var (
	// ErrConsulConfigLoad returned when env config for consul results in an error
	ErrConsulConfigLoad = errors.New("error loading consul config")
)

// Consul registry
type Consul struct {
	i      component.Initializer
	logger logger.Logger
	server server.Server
	*Config
}

// UnmarshalEnv env.EnvSet to GatewayConfig
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

// NewConsul returns a new consul registry
func NewConsul() *Consul {
	c := &Consul{Config: NewConfig()}
	c.i = NewInitializer(c)
	return c
}

// Register registers service with the registry.
func (c *Consul) Register(
	ctx context.Context,
	server server.Server,
) error {
	if err := consul.Register(
		ctx,
		c.ServiceName,
		server.Host(),
		server.Port(),
		c.Address(),
		c.TTL,
		server.RegistryTags()...,
	); err != nil {
		c.logger.Errorw(
			"Service registration failed",
			"error", err.Error(),
		)
		return err
	}

	c.logger.Infof("Successfully registered service: %s", c.ServiceName)
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

func (c *Consul) HasInitializer() bool {
	return true
}

func (c *Consul) Initializer() component.Initializer {
	return c.i
}
