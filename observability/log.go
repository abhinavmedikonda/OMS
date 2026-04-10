package observability

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/trace"
)

func TraceID(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return ""
	}
	return sc.TraceID().String()
}

func LogWithTrace(ctx context.Context, format string, args ...interface{}) {
	traceID := TraceID(ctx)
	if traceID == "" {
		log.Printf(format, args...)
		return
	}
	args = append(args, traceID)
	log.Printf(format+" trace_id=%s", args...)
}
