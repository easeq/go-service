package etcd

import (
	"strings"
	"time"

	// Manages env config
	"github.com/Netflix/go-env"
	"github.com/easeq/go-service/component"
)

// Config holds the etcd configuration
type Config struct {
	// Endpoints is a comma separated list of etcd URLs
	Endpoints string `env:"KVSTORE_ETCD_ENDPOINTS"`
	// DialTimeout is the timeout for failing to establish a connection
	DialTimeout time.Duration `env:"KVSTORE_ETCD_DIALTIMEOUT,omitempty"`
	// Username is a username for authentication
	Username string `env:"KVSTORE_ETCD_USERNAME,omitempty"`
	// Password is the password for authentication
	Password string `env:"KVSTORE_ETCD_PASSWORD,omitempty"`
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

// GetEndpoints return the etcd server endpoints
func (c *Config) GetEndpoints() []string {
	if c.Endpoints == "" {
		return []string{}
	}

	return strings.Split(c.Endpoints, ",")
}
