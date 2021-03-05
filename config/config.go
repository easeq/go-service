package config

import (
	"log"

	env "github.com/Netflix/go-env"
	"github.com/easeq/go-redis-access-control/manager"
)

// Config provides app configuration
type Config struct {
	ServiceName string `env:"SERVICE_NAME"`
	ServiceHost string `env:"SERVICE_HOST,default=localhost"`
	GRPCPort    int    `env:"GRPC_PORT,default=9090"`
	HTTPPort    int    `env:"HTTP_PORT,default=8080"`
	DB          Database
	Consul      Consul
	Grac        manager.Config
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() *Config {
	var config Config
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}
