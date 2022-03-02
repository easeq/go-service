package tracer

import "github.com/easeq/go-service/component"

const (
	TRACER              = "tracer"
	DEFAULT_TRACER_NAME = "github.com/easeq/go-service"
)

type Tracer interface {
	component.Component
}
