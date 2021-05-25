package db

import (
	"fmt"

	"github.com/Netflix/go-env"
	// Manages env config
	_ "github.com/easeq/go-config"
)

// Config holds database configuration
type Config struct {
	Name           string `env:"DB_NAME"`
	User           string `env:"DB_USER"`
	Password       string `env:"DB_PASS"`
	Driver         string `env:"DB_DRIVER,default=postgres"`
	Host           string `env:"DB_HOST,default=localhost"`
	Port           int    `env:"DB_PORT,default=5432"`
	SSLMode        string `env:"DB_SSL_MODE,default=disable"`
	MigrationsPath string `env:"DB_MIGRATIONS_PATH"`
}

// UnmarshalEnv env.EnvSet to Config
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

// GetURI generates and returns the database URI from the provided config
func (c *Config) GetURI() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		c.Driver,
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}
