package redis

import (
	"time"

	"github.com/Netflix/go-env"
)

type Config struct {
	Network  string `env:"REDIS_NETWORK"`
	Addr     string `env:"REDIS_ADDRESS"`
	Username string `env:"REDIS_USERNAME"`
	Password string `env:"REDIS_PASSWORD"`

	DB int `env:"REDIS_DB"`

	MaxRetries      int           `env:"REDIS_MAX_RETRIES"`
	MinRetryBackoff time.Duration `env:"REDIS_MIN_RETRY_BACKOFF"`
	MaxRetryBackoff time.Duration `env:"REDIS_MAX_RETRY_BACKOFF"`

	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT"`

	PoolSize           int           `env:"REDIS_POOL_SIZE"`
	MinIdleConns       int           `env:"REDIS_MIN_IDLE_CONNS"`
	MaxConnAge         time.Duration `env:"REDIS_MAX_CONN_AGE"`
	PoolTimeout        time.Duration `env:"REDIS_POOL_TIMEOUT"`
	IdleTimeout        time.Duration `env:"REDIS_IDLE_TIMEOUT"`
	IdleCheckFrequency time.Duration `env:"REDIS_IDLE_CHEKC_FREQUENCY"`
}

// UnmarshalEnv env.EnvSet to Config
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}
