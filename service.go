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
	"github.com/easeq/go-service/utils"
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
	cfg.UnmarshalEnv(goconfig.EnvSet())

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

// Register service at service registry
func (s *Service) Register(ctx context.Context) <-chan error {
	if s.Registry == nil {
		return nil
	}

	cErr := make(chan error, 1)
	go func() {
		defer close(cErr)
		cErr <- s.Server.Register(ctx, s.Name, s.Registry)
	}()

	return cErr
}

// Initialize database
func (s *Service) InitDatabase() <-chan error {
	if s.Database == nil {
		return nil
	}

	cErr := make(chan error, 1)
	go func() {
		defer close(cErr)
		cErr <- s.Database.Init()
	}()

	return cErr
}

// Start server
func (s *Service) RunServer(ctx context.Context) <-chan error {
	if s.Server == nil {
		return nil
	}

	cErr := make(chan error, 1)
	go func() {
		defer close(cErr)
		cErr <- s.Server.Run(ctx)
	}()

	return cErr
}

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
	// Register service with service registry
	cRegisterErr := s.Register(ctx)
	// Initialize database
	cInitDBErr := s.InitDatabase()
	// Run server
	cRunServerErr := s.RunServer(ctx)
	// Shutdown service and all its connections
	s.ShutDown(ctx)

	// Exit on error
	return utils.WaitForError(cRegisterErr, cInitDBErr, cRunServerErr)
}
