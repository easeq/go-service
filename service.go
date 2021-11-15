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
	"github.com/easeq/go-service/kvstore"
	"github.com/easeq/go-service/registry"
	"github.com/easeq/go-service/server"
	"github.com/easeq/go-service/tracer"
	"github.com/easeq/go-service/utils"
)

// Service handles config required by the service
type Service struct {
	Server   server.Server
	Broker   broker.Broker
	Database db.ServiceDatabase
	Registry registry.ServiceRegistry
	Client   client.Client
	KVStore  kvstore.KVStore
	Tracer   tracer.Tracer
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

// WithKVStore passes the kvstore used by the service
func WithKVStore(kvStore kvstore.KVStore) ServiceOption {
	return func(s *Service) {
		s.KVStore = kvStore
	}
}

// WithTracer assigns the tracer to be used by the service
func WithTracer(tracer tracer.Tracer) ServiceOption {
	return func(s *Service) {
		s.Tracer = tracer
	}
}

// // WithLogger sets the logger used by the service
// func WithLogger(logger logger.Logger) ServiceOption {
// 	return func(s *Service) {
// 		s.Logger = logger
// 	}
// }

// ShutDown shuts down the service and all its associated connections.
func (s *Service) ShutDown(ctx context.Context) {
	// graceful shutdown
	signal.Notify(s.exit, os.Interrupt)
	go func() {
		select {
		case <-s.exit:
			goto exit
		case <-ctx.Done():
			goto exit
		}

	exit:
		log.Println("Shutting down service")

		// sig is a ^C, handle it
		if s.Database != nil {
			log.Println("Closing DB connection")
			s.Database.Close()
		}

		if s.Broker != nil {
			log.Println("Closing broker")
			if err := s.Broker.Close(); err != nil {
				log.Println("Error closing broker", err)
			}
		}

		if s.KVStore != nil {
			log.Println("Closing KVStore")
			if err := s.KVStore.Close(); err != nil {
				log.Println("Error closing KVStore client", err)
			}
		}

		if s.Server != nil {
			log.Println("Shutting down server...")
			if err := s.Server.ShutDown(ctx); err != nil {
				log.Println("Error shutting down server", err)
			}
		}

		if s.Tracer != nil {
			log.Println("Closing tracer connection")
			s.Tracer.Stop()
		}
	}()
}

func (s *Service) RunResource(ctx context.Context, r interface{}) <-chan error {
	if r == nil {
		return nil
	}

	cErr := make(chan error, 1)

	go func() {
		defer close(cErr)

		switch v := r.(type) {
		case db.ServiceDatabase:
			log.Println("Initialize database...")
			if err := v.Init(); err != nil {
				cErr <- err
			}
		case registry.ServiceRegistry:
			log.Println("Register services...")
			if err := v.Register(ctx, s.Name, s.Server); err != nil {
				cErr <- err
			}
		case server.Server:
			log.Println("Run server...")
			if err := v.Run(ctx); err != nil {
				cErr <- err
			}
		}
	}()

	return cErr
}

// Run runs both the HTTP and gRPC server
func (s *Service) Run(ctx context.Context) error {
	return utils.WaitForError(
		s.RunResource(ctx, s.Database),
		s.RunResource(ctx, s.Registry),
		s.RunResource(ctx, s.Server),
	)
}
