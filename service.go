package goservice

import (
	"context"
	"log"
	"os"
	"os/signal"

	goconfig "github.com/easeq/go-config"
	"github.com/easeq/go-service/broker"
	"github.com/easeq/go-service/client"
	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/db"
	"github.com/easeq/go-service/kvstore"
	"github.com/easeq/go-service/logger"
	"github.com/easeq/go-service/logger/zap"
	"github.com/easeq/go-service/registry"
	"github.com/easeq/go-service/server"
	"github.com/easeq/go-service/tracer"
	"github.com/easeq/go-service/utils"
)

// Service handles config required by the service
type Service struct {
	exit       chan os.Signal
	components map[string]component.Component
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
		components: make(map[string]component.Component),
		exit:       make(chan os.Signal),
	}

	svc.components["logger"] = zap.NewZap()

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

// WithServer passes the server
func WithServer(srv server.Server) ServiceOption {
	return func(s *Service) {
		s.components[server.SERVER] = srv
	}
}

// WithBroker sets a broker for the service
func WithBroker(b broker.Broker) ServiceOption {
	return func(s *Service) {
		s.components[broker.BROKER] = b
	}
}

// WithDatabase sets the database used by the service
func WithDatabase(database db.ServiceDatabase) ServiceOption {
	return func(s *Service) {
		s.components[db.DATABASE] = database
	}
}

// WithRegistry passes services registry externally
func WithRegistry(r registry.ServiceRegistry) ServiceOption {
	return func(s *Service) {
		s.components[registry.REGISTRY] = r
	}
}

// WithClient registers the server's client
func WithClient(c client.Client) ServiceOption {
	return func(s *Service) {
		s.components[client.CLIENT] = c
	}
}

// WithKVStore passes the kvstore used by the service
func WithKVStore(kvStore kvstore.KVStore) ServiceOption {
	return func(s *Service) {
		s.components[kvstore.KV_STORE] = kvStore
	}
}

// WithTracer assigns the tracer to be used by the service
func WithTracer(t tracer.Tracer) ServiceOption {
	return func(s *Service) {
		s.components[tracer.TRACER] = t
	}
}

// WithLogger sets the logger used by the service
func WithLogger(l logger.Logger) ServiceOption {
	return func(s *Service) {
		s.components[logger.LOGGER] = l
	}
}

// Broker returns the instance as broker.Broker
func (s *Service) Broker() broker.Broker {
	return s.components[broker.BROKER].(broker.Broker)
}

// Server returns the instance as server.Server
func (s *Service) Server() server.Server {
	return s.components[server.SERVER].(server.Server)
}

// Tracer returns the instance as tracer.Tracer
func (s *Service) Tracer() tracer.Tracer {
	return s.components[tracer.TRACER].(tracer.Tracer)
}

// Database returns the instance as database.Database
func (s *Service) Database() db.ServiceDatabase {
	return s.components[db.DATABASE].(db.ServiceDatabase)
}

// KVStore returns the instance as kvstore.KVStore
func (s *Service) KVStore() kvstore.KVStore {
	return s.components[kvstore.KV_STORE].(kvstore.KVStore)
}

// Client returns the instance as client.Client
func (s *Service) Client() client.Client {
	return s.components[client.CLIENT].(client.Client)
}

// Registry returns the instance as registry.ServiceRegistry
func (s *Service) Registry() registry.ServiceRegistry {
	return s.components[registry.REGISTRY].(registry.ServiceRegistry)
}

// Logger returns the instance as logger.Logger
func (s *Service) Logger() logger.Logger {
	return s.components[logger.LOGGER].(logger.Logger)
}

// Init initializes the service
// Configures dependencies
func (s *Service) Init() error {
	s.configure()
	return nil
}

func (s *Service) configure() error {
	return s.IterateComponents(func(comp component.Component) error {
		if !comp.HasInitializer() {
			return nil
		}

		initializer := comp.Initializer()
		deps := initializer.Dependencies()
		if len(deps) == 0 {
			return nil
		}

		for _, dep := range deps {
			initializer.AddDependency(s.components[dep])
		}

		return nil
	})
}

// IterateComponents - iterates over all the service components and invokes the callback
func (s *Service) IterateComponents(cb func(comp component.Component) error) error {
	var errcList []<-chan error
	for _, comp := range s.components {
		svcComponent, ok := comp.(component.Component)
		if !ok {
			log.Println("Not a valid service component")
		}

		cErr := make(chan error, 1)
		go func() {
			if err := cb(svcComponent); err != nil {
				cErr <- err
			}
		}()
		errcList = append(errcList, cErr)
	}

	return utils.WaitForError(errcList...)
}

// Run runs both the HTTP and gRPC server
func (s *Service) Run(ctx context.Context) error {
	return s.IterateComponents(func(comp component.Component) error {
		if !comp.HasInitializer() {
			return nil
		}

		initializer := comp.Initializer()
		if !initializer.CanRun() {
			return nil
		}

		return initializer.Run(ctx)
	})
}

// Shutdown - shuts down the service by stopping all the components
func (s *Service) ShutDown(ctx context.Context) {
	signal.Notify(s.exit, os.Interrupt)
	go func() {
		select {
		case <-s.exit:
			goto exit
		case <-ctx.Done():
			goto exit
		}

	exit:
		s.Logger().Info("Shutting down service and it's components")
		s.IterateComponents(func(comp component.Component) error {
			if !comp.HasInitializer() {
				return nil
			}

			initializer := comp.Initializer()
			if !initializer.CanStop() {
				return nil
			}

			return initializer.Stop(ctx)
		})
	}()
}
