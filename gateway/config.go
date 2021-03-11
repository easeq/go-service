package gateway

import (
	"fmt"

	"github.com/Netflix/go-env"
	// Manages env config
	_ "github.com/easeq/go-config"
)

// Config manages the HTTP server config
type Config struct {
	Host string `env:"HTTP_HOST,default=localhost"`
	Port int    `env:"HTTP_PORT,default=8080"`
}

// UnmarshalEnv env.EnvSet to GatewayConfig
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

// Address returns the full formatted http address
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
