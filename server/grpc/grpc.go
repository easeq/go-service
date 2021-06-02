package grpc

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"strings"

	goconfig "github.com/easeq/go-config"
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
	// SERVER_TYPE is the type of the server.
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

// GetLogger returns a new zap.Logger
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

// Address returns the server address
func (g *Grpc) Address() string {
	return g.Config.Address()
}

// GetMetadata returns the metadata by key
func (g *Grpc) GetMetadata(key string) interface{} {
	return nil
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
	listener, err := net.Listen("tcp", g.Config.Address())
	if err != nil {
		return err
	}

	// start gRPC server
	log.Println("Starting gRPC server...")
	return g.Server.Serve(listener)
}

// ShutDown - gracefully stops the server
func (g *Grpc) ShutDown(ctx context.Context) error {
	g.Server.GracefulStop()
	return nil
}

// AddRegistryTags - sets the registry tags for the server
func (g *Grpc) AddRegistryTags(tags ...string) {
	g.Config.Tags = strings.Join(
		append(g.Config.GetTags(), tags...),
		registry.TAGS_SEPARATOR,
	)
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
