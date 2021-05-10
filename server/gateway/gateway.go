package gateway

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/easeq/go-redis-access-control/gateway"
	"github.com/easeq/go-service/pool"
	"github.com/easeq/go-service/registry"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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

// Option to pass as arg while creating new service
type Option func(*Gateway)

// Middleware type
type Middleware func(http.Handler) http.Handler

// HTTPServiceHandlerRegistrar HTTP service handler registration func
type HTTPServiceHandlerRegistrar func(context.Context, *Gateway) error

// Gateway handle gRPC gateway
type Gateway struct {
	Mux                         *runtime.ServeMux
	Middleware                  Middleware
	HTTPServiceHandlerRegistrar HTTPServiceHandlerRegistrar
	MuxOptions                  []runtime.ServeMuxOption
	Server                      *http.Server
	exit                        chan os.Signal
	*Config
}

// NewGateway creates and returns gRPC-gateway
func NewGateway(opts ...Option) *Gateway {
	g := &Gateway{
		Mux:        runtime.NewServeMux(),
		Middleware: gateway.Middleware,
		MuxOptions: []runtime.ServeMuxOption{},
		exit:       make(chan os.Signal),
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
}

// WithMuxOptions adds mux options
func WithMuxOptions(opts ...runtime.ServeMuxOption) Option {
	return func(g *Gateway) {
		g.MuxOptions = opts
	}
}

// WithMiddleware adds middleware to the rest handler
func WithMiddleware(middleware Middleware) Option {
	return func(g *Gateway) {
		g.Middleware = middleware
	}
}

// WithHTTPServiceHandlerRegistrar add HTTP service handle registration callback
func WithHTTPServiceHandlerRegistrar(registrar HTTPServiceHandlerRegistrar) Option {
	return func(g *Gateway) {
		g.HTTPServiceHandlerRegistrar = registrar
	}
}

// Client creates if not exists and returns the client to call the service
func (g *Gateway) GetClient(address string) (pool.Connection, error) {
	// var err error
	// for i := 0; i < maxBadClientConnRetries; i++ {
	// 	conn, err := g.getClientConn(address)
	// 	if err == nil {
	// 		return conn, nil
	// 	}
	// }

	// return nil, err
	return nil, fmt.Errorf("No avaialable client")
}

// Register registers the grpc server with the service registry
func (g *Gateway) Register(
	ctx context.Context,
	name string,
	registry registry.ServiceRegistry,
) *registry.ErrRegistryRegFailed {
	return registry.Register(ctx, name, g.Host, g.Port, g.GetTags()...)
}

// Run runs the HTTP server
func (g *Gateway) Run(ctx context.Context) error {
	g.Server = &http.Server{
		Addr:    g.Address(),
		Handler: g.Middleware(g.Mux),
	}

	log.Println("starting HTTP/REST gateway...")
	return g.Server.ListenAndServe()
}

// Shutdown HTTP server
func (g *Gateway) ShutDown(ctx context.Context) error {
	return g.Server.Shutdown(ctx)
}
