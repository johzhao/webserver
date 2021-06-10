package tracing

type Tracer interface {
	Close()

	StartSpan(operationName string, tags map[string]interface{}) Span
}
