//go:generate mockery --name=TracerProvider --with-expecter
package metrics

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type (
	TraceRecorder struct {
		tracer trace.Tracer
	}

	SpanRecorder struct {
		span trace.Span
	}
)

const (
	traceAppPrefix = "alpha:auth-api"
)

func NewTracer(tracerName string) TraceRecorder {
	return TraceRecorder{
		tracer: otel.GetTracerProvider().Tracer(fmt.Sprintf("%s/%s", traceAppPrefix, tracerName)),
	}
}

func (s TraceRecorder) StartSpan(ctx context.Context, spanName string) (context.Context, SpanRecorder) {
	return s.StartSpanWithLabels(ctx, spanName)
}

func (s TraceRecorder) StartSpanWithLabels(ctx context.Context, spanName string, attributes ...attribute.KeyValue) (context.Context, SpanRecorder) {
	newCtx, span := s.tracer.Start(ctx, spanName, trace.WithAttributes(attributes...))

	return newCtx, SpanRecorder{
		span: span,
	}
}

func (s SpanRecorder) SetLabels(attributes ...attribute.KeyValue) {
	s.span.SetAttributes(attributes...)
}

func (s SpanRecorder) Finish() {
	s.span.End()
}
