package grpc

import (
	"fmt"

	"github.com/Netflix/go-env"
	// Manages env config
	_ "github.com/easeq/go-config"
)

// Config manages the HTTP server config
type Config struct {
	Host     string `env:"GRPC_HOST,default=localhost"`
	Port     int    `env:"GRPC_PORT,default=9090"`
	HTTPHost string `env:"HTTP_HOST,default=localhost"`
	HTTPPort int    `env:"HTTP_PORT,default=8080"`
}

// UnmarshalEnv env.EnvSet to GatewayConfig
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

// Address returns the full formatted http address
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// HTTPAddress returns the full formatted http address
func (c *Config) HTTPAddress() string {
	return fmt.Sprintf("%s:%d", c.HTTPHost, c.HTTPPort)
}
