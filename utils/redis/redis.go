package redis

import (
	goconfig "github.com/easeq/go-config"
	goredis "github.com/go-redis/redis/v8"
)

type Redis struct {
	*Config
	Client *goredis.Client
}

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
