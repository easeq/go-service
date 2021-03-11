package grpc

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	// "github.com/easeq/pign/service-shopify/internal/config"
	"github.com/easeq/go-redis-access-control/gateway"
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
type HTTPServiceHandlerRegistrar func(context.Context, *Grpc) error

// Gateway handle gRPC gateway
type Gateway struct {
	Mux                         *runtime.ServeMux
	Middleware                  Middleware
	HTTPServiceHandlerRegistrar HTTPServiceHandlerRegistrar
	MuxOptions                  []runtime.ServeMuxOption
	exit                        chan os.Signal
}

// NewGateway creates and returns gRPC-gateway
func NewGateway() *Gateway {
	return &Gateway{
		Mux:        runtime.NewServeMux(),
		Middleware: gateway.Middleware,
		MuxOptions: []runtime.ServeMuxOption{},
		exit:       make(chan os.Signal),
	}
}

// WithMuxOptions adds mux options
func WithMuxOptions(opts ...runtime.ServeMuxOption) Option {
	return func(g *Grpc) {
		g.MuxOptions = opts
	}
}

// WithMiddleware adds middleware to the rest handler
func WithMiddleware(middleware Middleware) Option {
	return func(g *Grpc) {
		g.Middleware = middleware
	}
}

// WithHTTPServiceHandlerRegistrar add HTTP service handle registration callback
func WithHTTPServiceHandlerRegistrar(registrar HTTPServiceHandlerRegistrar) Option {
	return func(g *Grpc) {
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
func (g *Gateway) Run(ctx context.Context, grpc *Grpc) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if g.HTTPServiceHandlerRegistrar == nil {
		log.Println("gRPC running without gateway")
		return nil
	}

	err := g.HTTPServiceHandlerRegistrar(ctx, grpc)
	if err != nil {
		return ErrHTTPServiceHandlerRegFailed
	}

	srv := &http.Server{
		Addr:    grpc.HTTPAddress(),
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
