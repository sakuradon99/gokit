package trace

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const keyGinTraceID = "trace_id"

func GinMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		traceID, ok := c.Value("trace_id").(string)
		if !ok || traceID == "" {
			traceID = uuid.NewString()
			c.Set(keyGinTraceID, traceID)
		}
		c.Header("X-Trace-ID", traceID)
		c.Next()
	}
}

func WithTraceID(ctx context.Context) context.Context {
	traceID := uuid.NewString()
	return context.WithValue(ctx, keyGinTraceID, traceID)
}

func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(keyGinTraceID).(string)
	if !ok {
		return ""
	}
	return traceID
}
