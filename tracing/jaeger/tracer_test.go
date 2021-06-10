package jaeger

import (
	"testing"
	"time"
)

const (
	JaegerAddress = "" //xxx: fill the jaeger address, for example: http://127.0.0.1:14268/api/traces
)

func TestJaegerTracer(t *testing.T) {
	jaegerTracer, err := NewTracer("jaeger test", JaegerAddress, nil)
	if err != nil {
		t.Fatalf("create jaeger tracer failed with error: (%v)\n", err)
	}
	defer jaegerTracer.Close()

	s1 := jaegerTracer.StartSpan("span 1")
	defer s1.Finish()
	time.Sleep(100 * time.Millisecond)

	s2 := jaegerTracer.StartSpan("span 2")
	defer s2.Finish()
	time.Sleep(200 * time.Millisecond)

	s3 := jaegerTracer.StartSpan("span 3")
	defer s3.Finish()
	time.Sleep(300 * time.Millisecond)
}
