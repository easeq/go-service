package grpc

import (
	"fmt"
	"strings"

	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/registry"
)

// Config manages the HTTP server config
type Config struct {
	Host string `env:"GRPC_HOST,defaut="`
	Port int    `env:"GRPC_PORT,default=9090"`
	Tags string `env:"GRPC_CONSUL_TAGS,default="`
}

// NewConfig returns the parsed config for jetstream from env
func NewConfig() *Config {
	c := new(Config)
	component.NewConfig(c)

	return c
}

// Address returns the full formatted http address
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetTags return the service registry tags slice
func (c *Config) GetTags() []string {
	if c.Tags == "" {
		return []string{}
	}

	return strings.Split(c.Tags, registry.TAGS_SEPARATOR)
}
