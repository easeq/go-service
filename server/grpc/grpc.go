package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"

	goconfig "github.com/easeq/go-config"
	grpc_client "github.com/easeq/go-service/client/grpc"
	"github.com/easeq/go-service/db"
	goservice_db "github.com/easeq/go-service/db"
	"github.com/easeq/go-service/pool"
	"github.com/easeq/go-service/registry"
	goservice_registry "github.com/easeq/go-service/registry"
	"github.com/easeq/go-service/server"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

var (
	// ErrRequiredGRPCRegistrar returned when gRPC registration callback is not provided
	ErrRequiredGRPCRegistrar = errors.New("gRPC registration callback is required")
	// ErrGRPCConfigLoad returned when env config for GRPC results in an error
	ErrGRPCConfigLoad = errors.New("error loading grpc config")
)

// ServiceRegistrar gRPC service registration func
type ServiceRegistrar func(*grpc.Server, *Grpc)

// Grpc holds gRPC config
type Grpc struct {
	ServerOptions    []grpc.ServerOption
	ServiceRegistrar ServiceRegistrar
	DialOptions      []grpc.DialOption
	Database         goservice_db.ServiceDatabase
	Registry         goservice_registry.ServiceRegistry
	ClientPool       pool.Pool
	exit             chan os.Signal
	*Gateway
	*Config
}

// Option to pass as arg while creating new service
type Option func(*Grpc)

// NewGrpc creates a new gRPC
func NewGrpc(opts ...Option) server.Server {
	defaultRegistry := goservice_registry.NewRegistry()
	g := &Grpc{
		DialOptions:   []grpc.DialOption{grpc.WithInsecure()},
		ServerOptions: []grpc.ServerOption{},
		Config:        goconfig.NewEnvConfig(new(Config)).(*Config),
		Database:      goservice_db.NewPostgres(),
		exit:          make(chan os.Signal),
		Gateway:       NewGateway(),
		Registry:      defaultRegistry,
	}

	for _, opt := range opts {
		opt(g)
	}

	if len(g.MuxOptions) > 0 {
		g.Mux = runtime.NewServeMux(g.MuxOptions...)
	} else {
		g.Mux = runtime.NewServeMux()
	}

	g.ClientPool = grpc_client.NewGrpcClientPool(
		grpc_client.WithRegistry(g.Registry),
		grpc_client.WithDialOptions(g.DialOptions...),
		grpc_client.WithScheme("http"),
		grpc_client.WithTTL(grpc_client.DefaultTTL),
	)

	return g
}

// WithGrpcServerOptions adds gRPC options
func WithGrpcServerOptions(opts ...grpc.ServerOption) Option {
	return func(g *Grpc) {
		g.ServerOptions = opts
	}
}

// WithGrpcServiceRegistrar adds gRPC service registration callback
func WithGrpcServiceRegistrar(registrar ServiceRegistrar) Option {
	return func(g *Grpc) {
		g.ServiceRegistrar = registrar
	}
}

// WithGRPCDialOptions overrides custom dial options
func WithGRPCDialOptions(opts ...grpc.DialOption) Option {
	return func(g *Grpc) {
		g.DialOptions = opts
	}
}

// WithDatabase passes databases externally
func WithDatabase(database db.ServiceDatabase) Option {
	return func(g *Grpc) {
		if g.Database != nil {
			g.Database.Close()
		}

		g.Database = database
	}
}

// WithRegistry passes services registry externally
func WithRegistry(registry goservice_registry.ServiceRegistry) Option {
	return func(g *Grpc) {
		g.Registry = registry
	}
}

// Client creates if not exists and returns the client to call the service
func (g *Grpc) GetClient(name string) (pool.Connection, error) {
	return g.ClientPool.Get(name)
}

// Register registers the grpc server with the service registry
func (g *Grpc) Register(ctx context.Context, name string) *registry.ErrRegistryRegFailed {
	return g.Registry.Register(ctx, name, g.Host, g.Port)
}

// Run runs gRPC service
func (g *Grpc) Run(ctx context.Context) error {
	if err := g.Database.Setup(); err != nil {
		log.Printf("Database setup err -> %s", err)
	}

	if err := g.Database.UpdateHandle(); err != nil {
		return err
	}

	defer g.Database.Close()

	// Run migrations
	if err := g.Database.Migrate(); err != nil {
		log.Println(err)
	}

	// Register service
	server := grpc.NewServer(g.ServerOptions...)
	if g.ServiceRegistrar == nil {
		return ErrRequiredGRPCRegistrar
	}

	g.ServiceRegistrar(server, g)
	// reflection.Register(server) // for EVANS CLI

	listener, err := net.Listen("tcp", g.Config.Address())
	if err != nil {
		return err
	}

	// graceful shutdown
	signal.Notify(g.exit, os.Interrupt)
	go func() {
		for range g.exit {
			// sig is a ^C, handle it
			log.Println("Closing DB connection")
			g.Database.Close()

			log.Println("Shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	// run HTTP gateway
	go func() {
		_ = g.Gateway.Run(ctx, g)
	}()

	// start gRPC server
	log.Println("Starting gRPC server...")
	return server.Serve(listener)
}
