package redis

import (
	"time"

	"github.com/Netflix/go-env"
)

type Config struct {
	Network  string `env:"REDIS_NETWORK,default=tcp"`
	Addr     string `env:"REDIS_ADDRESS,default=localhost:6379"`
	Username string `env:"REDIS_USERNAME,default="`
	Password string `env:"REDIS_PASSWORD,default="`

	DB int `env:"REDIS_DB,default="`

	MaxRetries      int           `env:"REDIS_MAX_RETRIES,default=3"`
	MinRetryBackoff time.Duration `env:"REDIS_MIN_RETRY_BACKOFF,default=8ms"`
	MaxRetryBackoff time.Duration `env:"REDIS_MAX_RETRY_BACKOFF,default=512ms"`

	DialTimeout  time.Duration `env:"REDIS_DIAL_TIMEOUT,default=5s"`
	ReadTimeout  time.Duration `env:"REDIS_READ_TIMEOUT,default=3s"`
	WriteTimeout time.Duration `env:"REDIS_WRITE_TIMEOUT,default="`

	PoolSize           int           `env:"REDIS_POOL_SIZE,default=10"`
	MinIdleConns       int           `env:"REDIS_MIN_IDLE_CONNS,default="`
	MaxConnAge         time.Duration `env:"REDIS_MAX_CONN_AGE,default="`
	PoolTimeout        time.Duration `env:"REDIS_POOL_TIMEOUT,default="`
	IdleTimeout        time.Duration `env:"REDIS_IDLE_TIMEOUT,default="`
	IdleCheckFrequency time.Duration `env:"REDIS_IDLE_CHEKC_FREQUENCY,default=60s"`
}

// UnmarshalEnv env.EnvSet to Config
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}
