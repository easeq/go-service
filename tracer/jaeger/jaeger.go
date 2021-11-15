package jaeger

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
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

	tracer, closer, err := config.NewTracer(
		jaegercfg.Logger(jaegerlog.StdLogger),
		jaegercfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		panic(fmt.Errorf("jeaeger-err: %s", err))
	}

	opentracing.SetGlobalTracer(tracer)

	return &jaeger{tracer, closer, config}
}

// Stop jaeger connection
func (j *jaeger) Stop() error {
	return j.closer.Close()
}
