package rest

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/logger"
	"github.com/easeq/go-service/registry"
	"github.com/gofiber/fiber/v2"
)

var (
	// ErrGatewayConfigLoad returned when env config for gRPC-gateway results in an error
	ErrGatewayConfigLoad = errors.New("error loading gateway config")
	// ErrNotDefinedHTTPServiceHandlerRegistrar thrown when http service registration handler is not provided
	ErrNotDefinedHTTPServiceHandlerRegistrar = errors.New("http service handler registration callback not provided")
	// ErrHTTPServiceHandlerRegFailed returned when any HTTP service handler registration fails
	ErrHTTPServiceHandlerRegFailed = errors.New("http service handler registration failed")
	// ErrCannotAddMuxOptionAtPos returned when adding new mux option at the specified position is not possible
	ErrCannotAddMuxOptionAtPos = errors.New("cannot add mux option at the position specified")
)

const (
	// SERVER_TYPE is the type of the server.
	SERVER_TYPE = "rest"
)

// Option to pass as arg while creating new service
type Option func(*Rest)

// Rest server using gofiber
type Rest struct {
	i       component.Initializer
	logger  logger.Logger
	App     *fiber.App
	Options []fiber.Config
	Server  *http.Server
	exit    chan os.Signal
	*Config
}

// NewRest creates and returns rest server
func NewRest(opts ...Option) *Rest {
	r := &Rest{
		Options: []fiber.Config{},
		Config:  NewConfig(),
		exit:    make(chan os.Signal),
	}

	for _, opt := range opts {
		opt(r)
	}

	r.App = fiber.New(r.Options...)

	r.i = NewInitializer(r)
	return r
}

// WithMuxOptions adds mux options
func WithOptions(opts ...fiber.Config) Option {
	return func(r *Rest) {
		r.Options = opts
	}
}

// Address returns the server address
func (r *Rest) Address() string {
	return r.Config.Address()
}

// Host returns gateway server hostname
func (r *Rest) Host() string {
	return r.Config.Host
}

// Port returns gateway server port
func (r *Rest) Port() int {
	return r.Config.Port
}

// RegistryTags returns gateway server registry tags
func (r *Rest) RegistryTags() []string {
	return r.Config.GetTags()
}

// String - Returns the type of the server
func (r *Rest) String() string {
	return SERVER_TYPE
}

// AddRegistryTags - sets the registry tags for the server
func (r *Rest) AddRegistryTags(tags ...string) {
	r.Config.Tags = strings.Join(
		append(r.Config.GetTags(), tags...),
		registry.TAGS_SEPARATOR,
	)
}

func (r *Rest) HasInitializer() bool {
	return true
}

func (r *Rest) Initializer() component.Initializer {
	return r.i
}
