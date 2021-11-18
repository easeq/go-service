package grpc

import (
	"errors"
	"os"
	"strings"

	goconfig "github.com/easeq/go-config"
	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/logger"
	"github.com/easeq/go-service/registry"

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
	i             component.Initializer
	logger        logger.Logger
	ServerOptions []grpc.ServerOption
	DialOptions   []grpc.DialOption
	Server        *grpc.Server
	exit          chan os.Signal
	*Config
}

// Option to pass as arg while creating new service
type Option func(*Grpc)

// NewGrpc creates a new gRPC server
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
	g.i = NewInitializer(g)

	return g
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

// Address returns the server address
func (g *Grpc) Address() string {
	return g.Config.Address()
}

// Host returns gRPC server hostname
func (g *Grpc) Host() string {
	return g.Config.Host
}

// Port returns gRPC server port
func (g *Grpc) Port() int {
	return g.Config.Port
}

// RegistryTags returns gRPC server registry tags
func (g *Grpc) RegistryTags() []string {
	return g.Config.GetTags()
}

// GetMetadata returns the metadata by key
func (g *Grpc) GetMetadata(key string) interface{} {
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

func (g *Grpc) HasInitializer() bool {
	return true
}

func (g *Grpc) Initializer() component.Initializer {
	return g.i
}
