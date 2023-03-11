package object_store

import (

	// Manages env config
	"github.com/Netflix/go-env"
	"github.com/easeq/go-service/component"
)

// Config holds the etcd configuration
type Config struct {
	// Endpoint to your minio instance
	Endpoint string `env:"OBJECT_STORE_MINIO_ENDPOINT,omitempty"`
	// AccessKeyID
	AccessKeyID string `env:"OBJECT_STORE_MINIO_ACCESS_KEY_ID,omitempty"`
	// SecretAccessKey
	SecretAccessKey string `env:"OBJECT_STORE_MINIO_SECRET_ACCESS_KEY,omitempty"`
	// UseSSL
	UseSSL bool `env:"OBJECT_STORE_MINIO_USE_SSL,omitempty"`
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
