package jaeger

import (
	"fmt"
	"io"

	"github.com/easeq/go-service/component"
	"github.com/easeq/go-service/logger"
	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

type Jaeger struct {
	i      component.Initializer
	logger logger.Logger
	tracer opentracing.Tracer
	closer io.Closer
	config *jaegercfg.Configuration
}

func NewJaeger() *Jaeger {
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

	j := &Jaeger{tracer: tracer, closer: closer, config: config}
	j.i = NewInitializer(j)
	return j
}

func (j *Jaeger) HasInitializer() bool {
	return true
}

func (j *Jaeger) Initializer() component.Initializer {
	return j.i
}
