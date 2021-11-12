package nsq

import (
	"fmt"

	"github.com/Netflix/go-env"
	"github.com/nsqio/go-nsq"

	// Manages env config
	_ "github.com/easeq/go-config"
)

// Producer holds config for NSQ producer
type Producer struct {
	Host string `env:"BROKER_PRODUCER_HOST,default=127.0.0.1"`
	Port string `env:"BROKER_PRODUCER_PORT,default=4150"`
}

// Address returns the formatted address for the producer
func (p *Producer) Address() string {
	return fmt.Sprintf("%s:%s", p.Host, p.Port)
}

// Lookup holds config for NSQLookupd
type Lookupd struct {
	Host string `env:"BROKER_NSQ_LOOKUPD_HOST,default=127.0.0.1"`
	Port string `env:"BROKER_NSQ_LOOKUPD_PORT,default=4161"`
}

// Address returns the formatted address for the nsq lookupd
func (l *Lookupd) Address() string {
	return fmt.Sprintf("%s:%s", l.Host, l.Port)
}

// Config holds database configuration
type Config struct {
	Producer Producer
	Lookupd  Lookupd
}

// UnmarshalEnv env.EnvSet to Config
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

// Channel returns the channel name for the topic
func (c *Config) Channel(topic string) string {
	return fmt.Sprintf("Channel_%s", topic)
}

// NSQConfig returns the new config
func (c *Config) NSQConfig() *nsq.Config {
	return nsq.NewConfig()
}
