package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// AddAttribute adds a custom attribute to the current span if it is recording.
func AddAttribute(ctx context.Context, key string, value interface{}) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	switch typedValue := value.(type) {
	case string:
		span.SetAttributes(attribute.String(key, typedValue))
	case int:
		span.SetAttributes(attribute.Int(key, typedValue))
	case bool:
		span.SetAttributes(attribute.Bool(key, typedValue))
	}
}

// StartSpan creates a new span with the given name and returns the updated context and span.
// This is a convenience wrapper around the OpenTelemetry tracer Start method.
func StartSpan(ctx context.Context, spanName string) (context.Context, trace.Span) {
	tracer := otel.Tracer("application")

	return tracer.Start(ctx, spanName)
}
