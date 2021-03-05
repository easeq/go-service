package goservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/easeq/go-consul-registry/v2/consul"
	"github.com/easeq/go-redis-access-control/gateway"
	"github.com/easeq/go-service/config"
	"github.com/easeq/go-service/db"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/runtime/protoiface"

	grac_grpc "github.com/easeq/go-redis-access-control/grpc"

	// migration source file
	_ "github.com/golang-migrate/migrate/v4/source/file"
	// consul resolver
	_ "github.com/easeq/go-consul-registry/v2/consul"
	// postgres sql driver
	_ "github.com/lib/pq"
)

var (
	// ErrRequiredGRPCRegistrar returned when gRPC registration callback is not provided
	ErrRequiredGRPCRegistrar = errors.New("gRPC registration callback is required")
	// ErrNotDefinedHTTPServiceHandlerRegistrar thrown when http service registration handler is not provided
	ErrNotDefinedHTTPServiceHandlerRegistrar = errors.New("HTTP service handler registration callback not provided")
	// ErrDatabaseNameNotProvided returned when database name is not provided
	ErrDatabaseNameNotProvided = errors.New("database name not provided")
	// ErrHTTPServiceHandlerRegFailed returned when any HTTP service handler registration fails
	ErrHTTPServiceHandlerRegFailed = errors.New("HTTP service handler registration failed")
)

// Middleware type
type Middleware func(http.Handler) http.Handler

// GRPCServiceRegistrar gRPC service registration func
type GRPCServiceRegistrar func(*grpc.Server, *ServiceConfig)

// HTTPServiceHandlerRegistrar HTTP service handler registration func
type HTTPServiceHandlerRegistrar func(context.Context, *ServiceConfig) error

// ServiceOption to pass as arg while creating new service
type ServiceOption func(*ServiceConfig)

// Service handles a new service
type Service interface {
	Run()
}

// ServiceConfig handles config required by the service
type ServiceConfig struct {
	Config                      *config.Config
	Mux                         *runtime.ServeMux
	middleware                  Middleware
	gRPCServerOption            []grpc.ServerOption
	gRPCServiceRegistrar        GRPCServiceRegistrar
	httpServiceHandlerRegistrar HTTPServiceHandlerRegistrar
	GRPCDialOptions             []grpc.DialOption
	exit                        chan os.Signal
}

// NewService creates a new service
func NewService(opts ...ServiceOption) *ServiceConfig {
	serviceConfig := &ServiceConfig{
		Config:          config.LoadConfig(),
		GRPCDialOptions: []grpc.DialOption{grpc.WithInsecure()},
		exit:            make(chan os.Signal, 1),
	}

	for _, opt := range opts {
		opt(serviceConfig)
	}

	if serviceConfig.gRPCServiceRegistrar == nil {
		panic(ErrRequiredGRPCRegistrar)
	}

	if serviceConfig.httpServiceHandlerRegistrar == nil {
		log.Println(ErrNotDefinedHTTPServiceHandlerRegistrar)
	}

	if serviceConfig.gRPCServerOption == nil {
		serviceConfig.gRPCServerOption = serviceConfig.GetDefaultGRPCServerOption()
	}

	if serviceConfig.Mux == nil {
		serviceConfig.Mux = runtime.NewServeMux(serviceConfig.GetDefaultServeMuxOptions()...)
	}

	if serviceConfig.middleware == nil {
		serviceConfig.middleware = gateway.Middleware
	}

	return serviceConfig
}

// WithMux creats and assigns a new serve mux
func WithMux(opts ...runtime.ServeMuxOption) ServiceOption {
	return func(s *ServiceConfig) {
		s.Mux = runtime.NewServeMux(opts...)
	}
}

// WithMiddleware adds middleware to the rest handler
func WithMiddleware(middleware Middleware) ServiceOption {
	return func(s *ServiceConfig) {
		s.middleware = middleware
	}
}

// WithGrpcServerOption adds gRPC options
func WithGrpcServerOption(opts ...grpc.ServerOption) ServiceOption {
	return func(s *ServiceConfig) {
		s.gRPCServerOption = opts
	}
}

// WithGrpcServiceRegistrar adds gRPC service registration callback
func WithGrpcServiceRegistrar(registrar GRPCServiceRegistrar) ServiceOption {
	return func(s *ServiceConfig) {
		s.gRPCServiceRegistrar = registrar
	}
}

// WithHTTPServiceHandlerRegistrar add HTTP service handle registration callback
func WithHTTPServiceHandlerRegistrar(registrar HTTPServiceHandlerRegistrar) ServiceOption {
	return func(s *ServiceConfig) {
		s.httpServiceHandlerRegistrar = registrar
	}
}

// WithGRPCDialOptions overrides custom dial options
func WithGRPCDialOptions(opts ...grpc.DialOption) ServiceOption {
	return func(s *ServiceConfig) {
		s.GRPCDialOptions = opts
	}
}

// GetGRPCAddress returns gRPC address string
func (s *ServiceConfig) GetGRPCAddress() string {
	return fmt.Sprintf(":%d", s.Config.GRPCPort)
}

// GetHTTPAddress return HTTP address string
func (s *ServiceConfig) GetHTTPAddress() string {
	return fmt.Sprintf(":%d", s.Config.HTTPPort)
}

// GetDefaultServeMuxOptions returns the default mux options
func (s *ServiceConfig) GetDefaultServeMuxOptions() []runtime.ServeMuxOption {
	modifier := gateway.NewModifier(nil, &s.Config.Grac)

	return []runtime.ServeMuxOption{
		runtime.WithForwardResponseOption(modifier.ResponseModifier),
		runtime.WithForwardResponseOption(redirectHandler),
		runtime.WithMetadata(modifier.MetadataAnnotator),
	}
}

// GetDefaultGRPCServerOption returns the default gRPC server options
func (s *ServiceConfig) GetDefaultGRPCServerOption() []grpc.ServerOption {
	// Register gRPC interceptors to handle authentication
	authInterceptor := grac_grpc.NewAuthInterceptor(&s.Config.Grac.JWT)
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(authInterceptor.Unary()),
		grpc.StreamInterceptor(authInterceptor.Stream()),
	}
}

// SetupDatabase creates a database if not exists
func (s *ServiceConfig) SetupDatabase() error {
	if s.Config.DB.Name == "" {
		return ErrDatabaseNameNotProvided
	}

	// Setup database if not exists
	conn := db.NewConnection(s.Config.DB.Driver, s.Config.DB.GetURI(false, true))
	defer conn.DB.Close()

	if err := conn.SetupDatabase(s.Config.DB.Name, s.Config.DB.User, s.Config.DB.Password); err != nil {
		return err
	}

	return nil
}

// GetDBConnection returns a new DB connection
func (s *ServiceConfig) GetDBConnection() *db.Connection {
	// Create new connection to the database
	conn := db.NewConnection(s.Config.DB.Driver, s.Config.DB.GetURI(true, false))

	// Run migrations
	if err := conn.RunMigrations("file://configs/sql/migrations", s.Config.DB.Driver); err != nil {
		log.Println(err)
	}

	return conn
}

// GetClient returns the gRPC client for the given service name
func (s *ServiceConfig) GetClient(name string) (*grpc.ClientConn, error) {
	return grpc.Dial(s.Config.Consul.GetConnectionString(name))
}

// RegisterWithConsul registers this service with consul
func (s *ServiceConfig) RegisterWithConsul(ctx context.Context) error {
	if err := consul.Register(
		ctx,
		s.Config.ServiceName,
		s.Config.ServiceHost,
		s.Config.GRPCPort,
		s.Config.Consul.Address(),
		s.Config.Consul.TTL,
	); err != nil {
		return fmt.Errorf("Consul register error: %v", err)
	}

	return nil
}

// RunHTTPServer runs the HTTP server
func (s *ServiceConfig) RunHTTPServer(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := s.httpServiceHandlerRegistrar(ctx, s)
	if err != nil {
		return ErrHTTPServiceHandlerRegFailed
	}

	srv := &http.Server{
		Addr:    s.GetHTTPAddress(),
		Handler: s.middleware(s.Mux),
	}

	// Graceful shutdown
	signal.Notify(s.exit, os.Interrupt)
	go func() {
		for range s.exit {
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

// RunGRPCServer runs gRPC service
func (s *ServiceConfig) RunGRPCServer(ctx context.Context) error {
	// Register service
	server := grpc.NewServer(s.gRPCServerOption...)
	s.gRPCServiceRegistrar(server, s)
	// reflection.Register(server) // for EVANS CLI

	listener, err := net.Listen("tcp", s.GetGRPCAddress())
	if err != nil {
		return err
	}

	// graceful shutdown
	signal.Notify(s.exit, os.Interrupt)
	go func() {
		for range s.exit {
			// sig is a ^C, handle it
			log.Println("Shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("Starting gRPC server...")
	return server.Serve(listener)
}

// Run runs both the HTTP and gRPC servers
func (s *ServiceConfig) Run(ctx context.Context) error {
	if err := s.SetupDatabase(); err != nil {
		log.Printf("Database setup was unsuccesful: %s", err)
	}

	// conn := s.GetDBConnection()
	if err := s.RegisterWithConsul(ctx); err != nil {
		return fmt.Errorf("Consul registration failed: %s", err)
	}

	// run HTTP gateway
	go func() {
		_ = s.RunHTTPServer(ctx)
	}()

	return s.RunGRPCServer(ctx)
}

// Receive gRPC metadata and redirect if redirection headers are set
func redirectHandler(ctx context.Context, w http.ResponseWriter, resp protoiface.MessageV1) error {
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
