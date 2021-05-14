package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"strings"

	goconfig "github.com/easeq/go-config"
	"github.com/easeq/go-service/pool"
	"github.com/easeq/go-service/registry"

	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	// ErrRequiredGRPCRegistrar returned when gRPC registration callback is not provided
	ErrRequiredGRPCRegistrar = errors.New("gRPC registration callback is required")
	// ErrGRPCConfigLoad returned when env config for GRPC results in an error
	ErrGRPCConfigLoad = errors.New("error loading grpc config")
)

const (
	SERVER_TYPE = "grpc"
)

// Grpc holds gRPC config
type Grpc struct {
	ServerOptions []grpc.ServerOption
	DialOptions   []grpc.DialOption
	// Broker     broker.Broker
	Logger *zap.Logger
	Server *grpc.Server
	exit   chan os.Signal
	*Config
}

// Option to pass as arg while creating new service
type Option func(*Grpc)

// NewGrpc creates a new gRPC
func NewGrpc(opts ...Option) *Grpc {
	g := &Grpc{
		DialOptions:   []grpc.DialOption{grpc.WithInsecure()},
		ServerOptions: []grpc.ServerOption{},
		Config:        goconfig.NewEnvConfig(new(Config)).(*Config),
		exit:          make(chan os.Signal),
	}

	for _, opt := range opts {
		opt(g)
	}

	g.Server = grpc.NewServer(g.ServerOptions...)

	return g
}

// Get zap logger
func GetLogger() *zap.Logger {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Println("ZapLogger failed!")
	}
	defer zapLogger.Sync()

	return zapLogger
}

// WithGrpcServerOptions adds gRPC options
func WithGrpcServerOptions(opts ...grpc.ServerOption) Option {
	return func(g *Grpc) {
		g.ServerOptions = opts
	}
}

// WithGRPCDialOptions overrides custom dial options
func WithGRPCDialOptions(opts ...grpc.DialOption) Option {
	return func(g *Grpc) {
		g.DialOptions = opts
	}
}

// WithBroker passes the message broker externally
// func WithBroker(opts broker.Broker) Option {
// 	return func(g *Grpc) {
// 		g.Broker = opts
// 	}
// }

// func (g *Grpc) getClientConn(name string) (pool.Connection, error) {
// 	conn, err := g.ClientPool.Get(name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return conn, nil
// }

// Client creates if not exists and returns the client to call the service
func (g *Grpc) GetClient(address string) (pool.Connection, error) {
	var err error
	// for i := 0; i < maxBadClientConnRetries; i++ {
	// 	conn, err := g.getClientConn(address)
	// 	if err == nil {
	// 		return conn, nil
	// 	}
	// }

	return nil, err
}

// Register registers the grpc server with the service registry
func (g *Grpc) Register(
	ctx context.Context,
	name string,
	registry registry.ServiceRegistry,
) *registry.ErrRegistryRegFailed {
	return registry.Register(ctx, name, g.Host, g.Port, g.GetTags()...)
}

// Run runs gRPC service
func (g *Grpc) Run(ctx context.Context) error {
	// Register service
	server := grpc.NewServer(g.ServerOptions...)

	listener, err := net.Listen("tcp", g.Config.Address())
	if err != nil {
		return err
	}

	// start gRPC server
	log.Println("Starting gRPC server...")
	return server.Serve(listener)
}

// Shutdown - gracefully stops the server
func (g *Grpc) ShutDown(ctx context.Context) error {
	g.Server.GracefulStop()
	return nil
}

// SetRegistryTags - sets the registry tags for the server
func (g *Grpc) SetRegistryTags(tags ...string) {
	g.Config.Tags = strings.Join(tags, ",")
}

// String - Returns the type of the server
func (g *Grpc) String() string {
	return SERVER_TYPE
}

func init() {
	otCfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Println("Error setting up opentracing: ", err)
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, _, err := otCfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)

	// TODO: find a place to add closer.Close() to avoid premature closing
	opentracing.SetGlobalTracer(tracer)
}
