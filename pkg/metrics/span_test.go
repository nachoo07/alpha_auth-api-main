package metrics

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func TestNewTracer(t *testing.T) {
	tracer := NewTracer("test-span")
	if tracer.tracer == nil {
		t.Errorf("Expected tracer to be initialized")
	}
}

func TestTraceRecorder_StartSpan(t *testing.T) {
	tracer := NewTracer("test-span")
	ctx := context.Background()
	newCtx, spanRecorder := tracer.StartSpan(ctx, "test-span")

	if newCtx == nil {
		t.Errorf("Expected new context to be returned")
	}
	if spanRecorder.span == nil {
		t.Errorf("Expected span to be initialized")
	}
}

func TestTraceRecorder_StartSpanWithLabels(t *testing.T) {
	tracer := NewTracer("test-span")
	ctx := context.Background()
	attributes := []attribute.KeyValue{attribute.String("key", "value")}
	newCtx, spanRecorder := tracer.StartSpanWithLabels(ctx, "test-span", attributes...)

	if newCtx == nil {
		t.Errorf("Expected new context to be returned")
	}
	if spanRecorder.span == nil {
		t.Errorf("Expected span to be initialized")
	}
}

func TestSpanRecorder_SetLabels(t *testing.T) {
	tracer := NewTracer("test-span")
	ctx := context.Background()
	_, spanRecorder := tracer.StartSpan(ctx, "test-span")

	attributes := []attribute.KeyValue{attribute.String("key", "value")}
	spanRecorder.SetLabels(attributes...)
}

func TestSpanRecorder_Finish(t *testing.T) {
	tracer := NewTracer("test-span")
	ctx := context.Background()
	_, spanRecorder := tracer.StartSpan(ctx, "test-span")

	spanRecorder.Finish()
}
