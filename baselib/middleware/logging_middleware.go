package middleware

import (
	"bufio"
	"bytes"
	"gogin-template/bootstrap"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type responseBodyLogger struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyLogger) Write(b []byte) (int, error) {
	if w.body == nil {
		w.body = &bytes.Buffer{}
	}
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w responseBodyLogger) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w responseBodyLogger) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w responseBodyLogger) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func LoggingMiddleware(cfg *bootstrap.Container) gin.HandlerFunc {
	return func(c *gin.Context) {

		ignoredPaths := cfg.GetConfig().GetStringSlice("log.ignore")
		skipLog := false

		if len(ignoredPaths) > 0 {
			for _, path := range ignoredPaths {
				if strings.HasPrefix(c.Request.URL.Path, path) {
					skipLog = true
					break
				}
			}
		}

		if !skipLog {
			// Read the request body
			var requestBody []byte
			contentType := c.Request.Header.Get("Content-Type")
			if strings.Contains(contentType, "application/json") && c.Request.Body != nil {
				requestBody, _ = io.ReadAll(c.Request.Body)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			}

			rw := &responseBodyLogger{body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
			c.Writer = rw
			// Process the request
			c.Next()

			responseBody := rw.body.String()

			// Log the request and response
			cfg.Logger().Printf("METHOD : %s", c.Request.Method)
			cfg.Logger().Printf("ENDPOINT : %s", c.Request.URL.Path)
			cfg.Logger().Printf("REQUEST BODY : %s", string(requestBody))
			cfg.Logger().Printf("RESPONSE STATUS : %d", c.Writer.Status())
			cfg.Logger().Printf("RESPONSE BODY : %s", responseBody)
		}
	}
}
