package logger

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLogMiddleware(cfg *RequestLogConfig, l Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.Enabled || l == nil {
			c.Next()
			return
		}

		// Check if url is white-listed or not
		urlString := c.Request.URL.String()
		if isSliceContainsPrefix(urlString, cfg.WhiteList) {
			c.Next()
			return
		}

		header := c.Request.Header

		// Setting up writer to get api response
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Getting request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			l.Error(c, fmt.Sprintf("[LOGGER] Error while reading the request body: '%s'", err))
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		c.Next()

		// Get api body response if enabled
		resBody := "Not Enabled"
		if cfg.LoggingResponse {
			resBody = blw.body.String()
		}

		// Setup Logging Body Options
		opts := []Field{
			zap.String("method", c.Request.Method),
			zap.String("url", urlString),
			zap.Int("status", c.Writer.Status()),
			zap.String("user_agent", header.Get("User-Agent")),
			zap.String("request_body", string(body)),
			zap.String("response_body", resBody),
		}

		l.Info(c, "[Request Log]", opts...)
	}
}

func isSliceContainsPrefix(str string, arr []string) bool {
	if str == "" {
		return false
	}

	for _, v := range arr {
		if strings.HasPrefix(v, str) {
			return true
		}
	}

	return false
}
