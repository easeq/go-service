package jetstream

import (
	"fmt"

	"github.com/Netflix/go-env"

	// Manages env config
	_ "github.com/easeq/go-config"
)

// Config holds database configuration
type Config struct {
	Host string `env:"NATS_HOST,default=127.0.0.1"`
	Port string `env:"NATS_PORT,default=4222"`
}

// Address returns the formatted address for the producer
func (c *Config) Address() string {
	return fmt.Sprintf("nats://%s:%s", c.Host, c.Port)
}

// UnmarshalEnv env.EnvSet to Config
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}
