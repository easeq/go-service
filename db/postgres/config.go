package postgres

import (
	"fmt"

	// Manages env config
	"github.com/easeq/go-service/component"
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

// NewConfig returns the parsed config for postgres from env
func NewConfig() *Config {
	c := new(Config)
	component.NewConfig(c)

	return c
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
