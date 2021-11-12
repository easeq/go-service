package simple

import (
	"fmt"
	"strings"

	"github.com/Netflix/go-env"
	// Manages env config
	_ "github.com/easeq/go-config"
	"github.com/easeq/go-service/registry"
)

// Config manages the HTTP server config
type Config struct {
	Host string `env:"SERVER_HOST,defaut="`
	Port int    `env:"SERVER_PORT,default=8080"`
	Tags string `env:"SERVER_CONSUL_TAGS,default="`
}

// UnmarshalEnv env.EnvSet to GatewayConfig
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
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
