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
	"github.com/easeq/go-service/logger"
	"github.com/easeq/go-service/registry"
	"github.com/easeq/go-service/server"
	"github.com/easeq/go-service/tracer"
	"google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

type ServiceComponent interface {
	// AddDependency adds the service component dependency
	// requested by a service component
	AddDependency(dep interface{}) error
	// CanRun returns whether the service component has a Run function defined
	CanRun() bool
	// Run - runs/starts the service components
	Run(ctx context.Context) error
	// Dependencies returns the list of string service component dependency names
	Dependencies() []string
	// Stop - stops/closes the service components
	Stop(ctx context.Context) error
}

// Service handles config required by the service
type Service struct {
	// Server server.Server
	// Broker     broker.Broker
	// Database   db.ServiceDatabase
	// Registry registry.ServiceRegistry
	// Client   client.Client
	// KVStore    kvstore.KVStore
	// Tracer     tracer.Tracer
	// Logger     logger.Logger
	exit       chan os.Signal
	components map[string]interface{}
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
		// Logger:     zap.NewZap(),
		components: make(map[string]interface{}),
		exit:       make(chan os.Signal),
	}

	svc.components["logger"] = svc.Logger

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

// WithServer passes the server
func WithServer(server server.Server) ServiceOption {
	return func(s *Service) {
		s.components["server"] = server
	}
}

// WithBroker sets a broker for the service
func WithBroker(broker broker.Broker) ServiceOption {
	return func(s *Service) {
		s.components["broker"] = broker
	}
}

// WithDatabase sets the database used by the service
func WithDatabase(database db.ServiceDatabase) ServiceOption {
	return func(s *Service) {
		s.components["database"] = database
	}
}

// WithRegistry passes services registry externally
func WithRegistry(registry registry.ServiceRegistry) ServiceOption {
	return func(s *Service) {
		s.components["registry"] = registry
	}
}

// WithClient registers the server's client
func WithClient(client client.Client) ServiceOption {
	return func(s *Service) {
		s.components["client"] = client
	}
}

// WithKVStore passes the kvstore used by the service
func WithKVStore(kvStore kvstore.KVStore) ServiceOption {
	return func(s *Service) {
		s.components["kv-store"] = kvStore
	}
}

// WithTracer assigns the tracer to be used by the service
func WithTracer(tracer tracer.Tracer) ServiceOption {
	return func(s *Service) {
		s.components["tracer"] = tracer
	}
}

// WithLogger sets the logger used by the service
func WithLogger(logger logger.Logger) ServiceOption {
	return func(s *Service) {
		s.components["logger"] = logger
	}
}

// Broker returns the instance as broker.Broker
func (s *Service) Broker() broker.Broker {
	return s.components["broker"].(broker.Broker)
}

// Server returns the instance as server.Server
func (s *Service) Server() server.Server {
	return s.components["server"].(server.Server)
}

// Tracer returns the instance as tracer.Tracer
func (s *Service) Tracer() tracer.Tracer {
	return s.components["tracer"].(tracer.Tracer)
}

// Database returns the instance as database.Database
func (s *Service) Database() database.Database {
	return s.components["database"].(database.Database)
}

// KVStore returns the instance as kvstore.KVStore
func (s *Service) KVStore() kvstore.KVStore {
	return s.components["kvstore"].(kvstore.KVStore)
}

// Client returns the instance as client.Client
func (s *Service) Client() client.Client {
	return s.components["client"].(client.Client)
}

// Registry returns the instance as registry.ServiceRegistry
func (s *Service) Registry() registry.ServiceRegistry {
	return s.components["registry"].(registry.ServiceRegistry)
}

// Logger returns the instance as logger.Logger
func (s *Service) Logger() logger.Logger {
	return s.components["logger"].(logger.Logger)
}

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
			if err := s.Broker().Close(); err != nil {
				log.Println("Error closing broker", err)
			}
		}

		if s.KVStore != nil {
			log.Println("Closing KVStore")
			if err := s.KVStore().Close(); err != nil {
				log.Println("Error closing KVStore client", err)
			}
		}

		if s.Server != nil {
			log.Println("Shutting down server...")
			if err := s.Server().ShutDown(ctx); err != nil {
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

func (s *Service) Init() error {
	s.configure()
	return nil
}

func (s *Service) IterateComponents(cb func(comp ServiceComponent) error) error {
	for _, comp := range s.components {
		svcComponent, ok := comp.(ServiceComponent)
		if !ok {
			log.Println("Not a valid service component")
		}

		go cb(svcComponent)
	}

	return nil
}

// Run runs both the HTTP and gRPC server
func (s *Service) Run(ctx context.Context) error {
	s.IterateComponents(func(comp ServiceComponent) error {
		if !comp.CanRun() {
			return nil
		}

		return comp.Run(ctx)
	})

	// return utils.WaitForError(
	// 	s.RunResource(ctx, s.Database),
	// 	s.RunResource(ctx, s.Registry),
	// 	s.RunResource(ctx, s.Server),
	// )
}

func (s *Service) configure() error {
	s.IterateComponents(func(comp ServiceComponent) error {
		deps := comp.Dependencies()
		if len(deps) == 0 {
			return nil
		}

		s.configureComponentDeps(comp, deps)
		return nil
	})

	return nil
}

func (s *Service) configureComponentDeps(comp ServiceComponent, deps []string) {
	for _, dep := range deps {
		comp.AddDependency(s.components[dep])
	}
}
