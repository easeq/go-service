package redis

import (
	goconfig "github.com/easeq/go-config"
	goredis "github.com/go-redis/redis/v8"
)

// Redis hold the redis config and redis client
type Redis struct {
	*Config
	Client *goredis.Client
}

// NewRedisClient creates a new redis client using the env config
func NewRedisClient() *Redis {
	cfg := GetConfig()

	return &Redis{Config: cfg}
}

// GetConfig returns the DB config
func GetConfig() *Config {
	return goconfig.NewEnvConfig(new(Config)).(*Config)
}

// Get returns redis client connection
func (r *Redis) Get() *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})
}
