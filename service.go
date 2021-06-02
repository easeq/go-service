package goservice

import (
	"context"
	"log"
	"os"
	"os/signal"

	goconfig "github.com/easeq/go-config"
	"github.com/easeq/go-service/broker"
	"github.com/easeq/go-service/client"
	"github.com/easeq/go-service/db"
	"github.com/easeq/go-service/registry"
	"github.com/easeq/go-service/server"
)

// ServiceConfig handles config required by the service
type Service struct {
	Server   server.Server
	Broker   broker.Broker
	Database db.ServiceDatabase
	Registry registry.ServiceRegistry
	Client   client.Client
	exit     chan os.Signal
	*Config
}

// ServiceOption to pass as arg while creating new service
type ServiceOption func(*Service)

// NewService creates a new service
func NewService(opts ...ServiceOption) *Service {
	cfg := new(Config)
	if err := cfg.UnmarshalEnv(goconfig.EnvSet()); err != nil {
		panic("Error loading env vars")
	}

	svc := &Service{
		Config: cfg,
		exit:   make(chan os.Signal),
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

// WithServer passes the server
func WithServer(server server.Server) ServiceOption {
	return func(s *Service) {
		s.Server = server
	}
}

// WithBroker sets a broker for the service
func WithBroker(broker broker.Broker) ServiceOption {
	return func(s *Service) {
		s.Broker = broker
	}
}

// WithDatabase sets the database used by the service
func WithDatabase(database db.ServiceDatabase) ServiceOption {
	return func(s *Service) {
		s.Database = database
	}
}

// WithRegistry passes services registry externally
func WithRegistry(registry registry.ServiceRegistry) ServiceOption {
	return func(s *Service) {
		s.Registry = registry
	}
}

// WithClient registers the server's client
func WithClient(client client.Client) ServiceOption {
	return func(s *Service) {
		s.Client = client
	}
}

// // WithLogger sets the logger used by the service
// func WithLogger(logger logger.Logger) ServiceOption {
// 	return func(s *Service) {
// 		s.Logger = logger
// 	}
// }

// Shutdown service and all the connections
func (s *Service) ShutDown(ctx context.Context) {
	// graceful shutdown
	signal.Notify(s.exit, os.Interrupt)
	go func() {
		for range s.exit {
			log.Println("Shutting down service")

			// sig is a ^C, handle it
			if s.Database != nil {
				log.Println("Closing DB connection")
				s.Database.Close()
			}

			if s.Server != nil {
				log.Println("Shutting down server...")
				if err := s.Server.ShutDown(ctx); err != nil {
					log.Println("Error shutting down server", err)
				}
			}

			<-ctx.Done()
		}
	}()
}

// Run runs both the HTTP and gRPC server
func (s *Service) Run(ctx context.Context) error {
	if s.Database != nil {
		if err := s.Database.Init(); err != nil {
			return err
		}
	}

	if s.Server != nil && s.Registry != nil {
		if err := s.Server.Register(ctx, s.Name, s.Registry); err != nil {
			return err
		}
	}

	err := s.Server.Run(ctx)
	if err != nil {
		return err
	}

	s.ShutDown(ctx)

	return nil
}
