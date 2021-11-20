package tracer

import "github.com/easeq/go-service/component"

const (
	TRACER = "tracer"
)

type Tracer interface {
	component.Component
}
