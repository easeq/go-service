package liftbridge

import (
	"strings"

	"github.com/Netflix/go-env"
	goconfig "github.com/easeq/go-config"
)

type Config struct {
	// Addrs contains the client addresses
	Addrs string `env:"BROKER_LIFTBRIDGE_ADDRESSES,default=localhost:9292"`
}

// GetConfig returns the DB config
func GetConfig() *Config {
	return goconfig.NewEnvConfig(new(Config)).(*Config)
}

// UnmarshalEnv env.EnvSet to Config
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

// GetAddresses returns an array of addresses
func (c *Config) GetAddresses() []string {
	return strings.Split(c.Addrs, ",")
}
