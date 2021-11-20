package consul

import "github.com/easeq/go-service/component"

// Config - consul configuration
type Config struct {
	ServiceName string `env:"SERVICE_NAME"`
	Host        string `env:"CONSUL_HOST,default=localhost"`
	Port        int    `env:"CONSUL_PORT,default=8500"`
	TTL         int    `env:"CONSUL_TTL,default=15"`
}

// NewConfig returns the parsed config for jetstream from env
func NewConfig() *Config {
	c := new(Config)
	component.NewConfig(c)

	return c
}
