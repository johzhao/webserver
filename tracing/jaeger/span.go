package jaeger

import "github.com/opentracing/opentracing-go"

type span struct {
	span     opentracing.Span
	finished func()
}

func (s *span) Finish() {
	s.span.Finish()
	if s.finished != nil {
		s.finished()
		s.finished = nil
	}
}
