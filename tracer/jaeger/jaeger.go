package jaeger

import (
	"context"
	"errors"
	"io"

	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

var (
	// ErrJaegerConnection returned when connecting to jaeger fails
	ErrJaegerConnection = errors.New("error connecting to jaeger")
	// ErrJaegerClose returned when closing connection with jaeger fails
	ErrJaegerClose = errors.New("error closing jaeger connection")
)

type jaeger struct {
	tracer opentracing.Tracer
	closer io.Closer
	config *jaegercfg.Configuration
}

func NewJaeger() *jaeger {
	config, err := jaegercfg.FromEnv()
	if err != nil {
		panic("error fetching jaeger config")
	}

	tracer, closer, err := config.NewTracer()
	if err != nil {
		panic(ErrJaegerConnection)
	}

	return &jaeger{tracer, closer, config}
}

// Start a connection with jaeger
func (j *jaeger) Start(ctx context.Context) error {
	opentracing.SetGlobalTracer(j.tracer)
	return nil
}

// Stop jaeger connection
func (j *jaeger) Stop() error {
	return j.closer.Close()
}
