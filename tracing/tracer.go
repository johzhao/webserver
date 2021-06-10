package tracing

type Tracer interface {
	Close()

	StartSpan(operationName string) Span
}
