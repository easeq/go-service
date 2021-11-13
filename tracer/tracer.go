package tracer

import "context"

type Tracer interface {
	// Start tracer
	Start(ctx context.Context) error
	// Stop tracer
	Stop() error
}
