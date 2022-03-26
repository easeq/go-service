package jetstream

import (
	// Manages env config
	"fmt"

	"github.com/Netflix/go-env"
	"github.com/easeq/go-service/component"
)

// Config holds the etcd configuration
type Config struct {
	Host   string `env:"NATS_HOST,default=127.0.0.1"`
	Port   string `env:"NATS_PORT,default=4222"`
	Bucket string `env:"JETSTREAM_BUCKET,default="`
}

// NewConfig returns the parsed config for jetstream from env
func NewConfig() *Config {
	c := new(Config)
	component.NewConfig(c)

	return c
}

// UnmarshalEnv env.EnvSet to Config
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

// Address returns the formatted address for the producer
func (c *Config) Address() string {
	return fmt.Sprintf("nats://%s:%s", c.Host, c.Port)
}
