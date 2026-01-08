package metrics

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		ctx := r.Context()
		spanContext := trace.SpanContextFromContext(ctx)

		ctx = trace.ContextWithSpanContext(ctx,
			spanContext.WithTraceFlags(spanContext.TraceFlags()))

		newCtx, span := NewTracer("route").StartSpan(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.String()))
		defer span.Finish()

		c.Request = r.WithContext(newCtx)

		c.Next()
	}
}
