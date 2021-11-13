package simple

import (
	"context"
	"errors"
	"os"
	"runtime"
	"strings"

	goconfig "github.com/easeq/go-config"
	"github.com/easeq/go-service/registry"

	"go.uber.org/zap"
)

var (
	// ErrRequiredGRPCRegistrar returned when gRPC registration callback is not provided
	ErrRequiredGRPCRegistrar = errors.New("gRPC registration callback is required")
	// ErrGRPCConfigLoad returned when env config for GRPC results in an error
	ErrGRPCConfigLoad = errors.New("error loading grpc config")
)

const (
	// SERVER_TYPE is the type of the server.
	SERVER_TYPE = "simple"
)

// Grpc holds gRPC config
type Simple struct {
	Logger *zap.Logger
	exit   chan os.Signal
	*Config
}

// Option to pass as arg while creating new service
type Option func(*Simple)

// NewGrpc creates a new gRPC server
func NewSimple(opts ...Option) *Simple {
	g := &Simple{
		Config: goconfig.NewEnvConfig(new(Config)).(*Config),
		exit:   make(chan os.Signal),
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// Address returns the server address
func (s *Simple) Address() string {
	return s.Config.Address()
}

// Host returns gRPC server hostname
func (s *Simple) Host() string {
	return s.Config.Host
}

// Port returns gRPC server port
func (s *Simple) Port() int {
	return s.Config.Port
}

// RegistryTags returns gRPC server registry tags
func (s *Simple) RegistryTags() []string {
	return s.Config.GetTags()
}

// GetMetadata returns the metadata by key
func (s *Simple) GetMetadata(key string) interface{} {
	return nil
}

// Run runs gRPC service
func (s *Simple) Run(ctx context.Context) error {
	runtime.Goexit()
	return nil
}

// ShutDown - gracefully stops the server
func (s *Simple) ShutDown(ctx context.Context) error {
	return nil
}

// AddRegistryTags - sets the registry tags for the server
func (s *Simple) AddRegistryTags(tags ...string) {
	s.Config.Tags = strings.Join(
		append(s.Config.GetTags(), tags...),
		registry.TAGS_SEPARATOR,
	)
}

// String - Returns the type of the server
func (s *Simple) String() string {
	return SERVER_TYPE
}
