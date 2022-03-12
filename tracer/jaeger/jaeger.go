package jaeger

import (
	"context"
	"log"

	"github.com/Netflix/go-env"
	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	Endpoint string `env:"JAEGER_ENDPOINT"`
}

// NewConfig returns the parsed config for jaeger from env
func NewConfig() *Config {
	c := new(Config)
	component.NewConfig(c)

	return c
}

// UnmarshalEnv env.EnvSet to Config
func (c *Config) UnmarshalEnv(es env.EnvSet) error {
	return env.Unmarshal(es, c)
}

type Jaeger struct {
	i      component.Initializer
	logger logger.Logger
	tracer *sdktrace.TracerProvider
}

func NewJaeger() *Jaeger {
	cfg := NewConfig()

	jaegerExporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Endpoint)),
	)
	if err != nil {
		log.Fatalln("Couldn't initialize jaeger", err)
	}

	resources, err := resource.New(
		context.Background(),
		resource.WithFromEnv(),
		resource.WithProcess(),
	)
	if err != nil {
		log.Fatalln("Couldn't initialize jaeger opentel resources")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(jaegerExporter),
		sdktrace.WithResource(resources),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	j := &Jaeger{tracer: tp}
	j.i = NewInitializer(j)
	return j
}

func (j *Jaeger) HasInitializer() bool {
	return true
}

func (j *Jaeger) Initializer() component.Initializer {
	return j.i
}
