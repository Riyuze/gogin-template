package identifier

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func GetTraceId(c *gin.Context) string {
	return fmt.Sprint(trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID())
}
