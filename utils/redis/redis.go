package redis

import (
	goredis "github.com/go-redis/redis/v8"
)

// Redis hold the redis config and redis client
type Redis struct {
	*Config
	Client *goredis.Client
}

// NewRedisClient creates a new redis client using the env config
func NewRedisClient(config *Config) *Redis {
	client := goredis.NewClient(&goredis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return &Redis{config, client}
}
