package config

import "fmt"

// Consul configuration
type Consul struct {
	Host string `env:"CONSUL_HOST,default=localhost"`
	Port int    `env:"CONSUL_PORT,default=8500"`
	TTL  int    `env:"CONSUL_TTL,default=15"`
}

// Address returns the prepared consul address
func (c *Consul) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// GetConnectionString returns consul connections string the given service
func (c *Consul) GetConnectionString(service string) string {
	return fmt.Sprintf(
		"consul://%s/%s?scheme=https",
		c.Address(),
		service,
	)
}
