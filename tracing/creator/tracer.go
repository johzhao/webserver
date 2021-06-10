package creator

import (
	"go.uber.org/zap"
	"webserver/tracing"
	"webserver/tracing/jaeger"
)

func NewTracer(serviceName string, serverAddress string, logger *zap.Logger) (tracing.Tracer, error) {
	return jaeger.NewTracer(serviceName, serverAddress, logger)
}
