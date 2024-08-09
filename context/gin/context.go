package context

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	datadogSpanID  ctxKey = "dd.span_id"
	datadogTraceID ctxKey = "dd.trace_id"
)

// WithDatadogSpanId sets the Datadog span ID in the context.
func WithDatadogSpanId(ctx context.Context, spanID uint64) context.Context {
	return context.WithValue(ctx, datadogSpanID, spanID)
}

// SpanID returns the Datadog span ID from the context.
func SpanID(ctx context.Context) (uint64, bool) {
	spanID, ok := ctx.Value(datadogTraceID).(uint64)
	return spanID, ok
}

// WithDatadogTraceId sets the Datadog trace ID in the context.
func WithDatadogTraceId(ctx context.Context, traceID uint64) context.Context {
	return context.WithValue(ctx, datadogTraceID, traceID)
}

// TraceID returns the Datadog trace ID from the context.
func TraceID(ctx context.Context) (uint64, bool) {
	traceID, ok := ctx.Value(datadogTraceID).(uint64)
	return traceID, ok
}

// GinNewContextWithObservability parses the gin context and return a new context with observability.
// Includes Datadog trace and span IDs to the context.
func GinNewContextWithObservability(c *gin.Context) (nctx context.Context) {
	ctx := c.Request.Context()

	ddSpanID := c.GetUint64("dd.span_id")
	ddTraceID := c.GetUint64("dd.trace_id")

	if ddSpanID != 0 {
		ctx = WithDatadogSpanId(ctx, ddSpanID)
	}

	if ddTraceID != 0 {
		ctx = WithDatadogTraceId(ctx, ddTraceID)
	}

	return ctx
}
