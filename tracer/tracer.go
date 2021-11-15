package tracer

type Tracer interface {
	// Stop tracer
	Stop() error
}
