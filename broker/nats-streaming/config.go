package nats_streaming

import (
	"github.com/Netflix/go-env"
	goconfig "github.com/easeq/go-config"
)

type Config struct {
	// ClusterID is the name of your nats streaming cluster.
	ClusterID string `env:"BROKER_NATS_STREAMING_CLUSTER_ID,default=test-cluster"`
	// ClientID is the name given to your client
	ClientID string `env:"BROKER_NATS_STREAMING_CLIENT_ID,default=go-service-ns-client"`
}

// GetConfig returns the DB config
func GetConfig() *Config {
	return goconfig.NewEnvConfig(new(Config)).(*Config)
}

// UnmarshalEnv env.EnvSet to Config
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}
