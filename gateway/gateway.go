package gateway

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	goconfig "github.com/easeq/go-config"
	grac_gateway "github.com/easeq/go-redis-access-control/gateway"

	// "github.com/easeq/pign/service-shopify/internal/config"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/protobuf/runtime/protoiface"
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

// Middleware type
type Middleware func(http.Handler) http.Handler

// HTTPServiceHandlerRegistrar HTTP service handler registration func
type HTTPServiceHandlerRegistrar func(context.Context) error

// Gateway handle gRPC gateway
type Gateway struct {
	Mux                         *runtime.ServeMux
	Middleware                  Middleware
	HTTPServiceHandlerRegistrar HTTPServiceHandlerRegistrar
	MuxOptions                  []runtime.ServeMuxOption
	exit                        chan os.Signal
	*Config
}

// Option to pass as arg while creating new service
type Option func(*Gateway)

// NewGateway creates a new gRPC gateway
func NewGateway(opts ...Option) *Gateway {
	// modifier := gateway.NewModifier(s.store, config)
	g := &Gateway{
		Mux:        runtime.NewServeMux(),
		Middleware: grac_gateway.Middleware,
		Config:     goconfig.NewEnvConfig(new(Config)).(*Config),
	}

	for _, opt := range opts {
		opt(g)
	}

	if g.HTTPServiceHandlerRegistrar == nil {
		log.Println(ErrNotDefinedHTTPServiceHandlerRegistrar)
	}

	if len(g.MuxOptions) > 0 {
		g.Mux = runtime.NewServeMux(g.MuxOptions...)
	}

	return g
}

// WithMux creats and assigns a new serve mux
func WithMux(opts ...runtime.ServeMuxOption) Option {
	return func(g *Gateway) {
		g.Mux = runtime.NewServeMux(opts...)
	}
}

// WithMuxOptions adds gRPC options
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

// RedirectHandler gRPC metadata and redirect if redirection headers are set
func RedirectHandler(ctx context.Context, w http.ResponseWriter, resp protoiface.MessageV1) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])

		if code, ok := headers["Grpc-Metadata-Code"]; ok {
			codeInt, err := strconv.Atoi(code[0])
			if err != nil {
				return err
			}

			w.WriteHeader(codeInt)
		} else {
			w.WriteHeader(http.StatusFound)
		}
	}

	return nil
}

// Run runs the HTTP server
func (g *Gateway) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := g.HTTPServiceHandlerRegistrar(ctx)
	if err != nil {
		return ErrHTTPServiceHandlerRegFailed
	}

	srv := &http.Server{
		Addr:    g.Address(),
		Handler: g.Middleware(g.Mux),
	}

	// Graceful shutdown
	signal.Notify(g.exit, os.Interrupt)
	go func() {
		for range g.exit {
			// sig is a ^C, handle it
			log.Println("Shutting down HTTP server...")
			srv.Shutdown(ctx)
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		srv.Shutdown(ctx)
	}()

	log.Println("starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}
