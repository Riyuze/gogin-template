package identifier

import (
	"fmt"

	"github.com/google/uuid"

	"context"

	"github.com/gin-gonic/gin"
)

type LogReffCtxKey struct{}

func GetLogReff(c *gin.Context) string {
	logReff, _ := c.Request.Context().Value(LogReffCtxKey{}).(string)

	if logReff == "" {
		logReff = fmt.Sprint(uuid.New())
	}

	return logReff
}

func SetLogReff(c *gin.Context, logReff string) {
	nctx := context.WithValue(c.Request.Context(), LogReffCtxKey{}, logReff)
	c.Request = c.Request.WithContext(nctx)
}

func SetLogReffWithRandom(c *gin.Context) {
	SetLogReff(c, fmt.Sprint(uuid.New()))
}
