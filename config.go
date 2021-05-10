package goservice

import "github.com/Netflix/go-env"

// Config - Service configuration
type Config struct {
	Name string `env:"SERVICE_NAME"`
}

// UnmarshalEnv env.EnvSet to GatewayConfig
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}
