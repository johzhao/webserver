package jaeger

import (
	_ "container/list"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
	"io"
	"sync"
	"webserver/tracing"
)

func NewTracer(serviceName string, address string, logger *zap.Logger) (tracing.Tracer, error) {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: address,
		},
	}

	jLogger := jaegerlog.NullLogger
	jMetricsFactory := metrics.NullFactory

	t, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		logger.Error("failed to create jaeger tracer",
			zap.Error(err))
		return nil, err
	}

	return &tracer{
		logger: logger,
		tracer: t,
		closer: closer,
		spans:  nil,
	}, nil
}

type tracer struct {
	logger *zap.Logger
	mutex  sync.Mutex

	tracer opentracing.Tracer
	closer io.Closer
	spans  []span
}

func (t *tracer) Close() {
	_ = t.closer.Close()
}

func (t *tracer) StartSpan(operationName string, tags map[string]interface{}) tracing.Span {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	opentracingSpan := t.tracer.StartSpan(operationName, opentracing.ChildOf(t.lastSpanContext()))
	for key, value := range tags {
		opentracingSpan.SetTag(key, value)
	}

	jaegerSpan := span{
		span:     opentracingSpan,
		finished: t.spanFinished,
	}
	t.spans = append(t.spans, jaegerSpan)

	return &jaegerSpan
}

func (t *tracer) lastSpanContext() opentracing.SpanContext {
	if len(t.spans) == 0 {
		return nil
	}
	return t.spans[len(t.spans)-1].span.Context()
}

func (t *tracer) spanFinished() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.spans = t.spans[:len(t.spans)-1]
}
