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
	RootUser       string `env:"DB_ROOT_USER"`
	RootPassword   string `env:"DB_ROOT_PASS"`
	MigrationsPath string `env:"DB_MIGRATIONS_PATH"`
}

// UnmarshalEnv env.EnvSet to Config
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

// GetUsername returns the required db username
func (c *Config) GetUsername(root bool) string {
	if root == true {
		return c.RootUser
	}

	return c.User
}

// GetPassword returns the required db password
func (c *Config) GetPassword(root bool) string {
	if root == true {
		return c.RootPassword
	}

	return c.Password
}

// GetURI generates and returns the database URI from the provided config
func (c *Config) GetURI(asRoot bool) string {
	dbName := c.Name
	if asRoot {
		dbName = ""
	}

	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		c.Driver,
		c.GetUsername(asRoot),
		c.GetPassword(asRoot),
		c.Host,
		c.Port,
		dbName,
		c.SSLMode,
	)
}
