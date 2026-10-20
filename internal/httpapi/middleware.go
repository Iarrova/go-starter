package httpapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Iarrova/go-starter/internal/logging"
)

const requestIDHeader = "X-Request-ID"

// requestID reuses the caller's X-Request-ID when present, otherwise generates one, and echoes it back in the response.
// The ID is stored in the request context, so every context-aware log call for this request includes it.
func requestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" || len(id) > 64 {
			id = uuid.NewString()
		}
		c.Request = c.Request.WithContext(logging.WithRequestID(c.Request.Context(), id))
		c.Header(requestIDHeader, id)
		c.Next()
	}
}

// timeout sets a deadline on the request context, so database queries and other context-aware work stop once d has passed.
func timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// logRequests logs one line per request. Successful requests to skipPaths (e.g. health probes) are not logged; failing ones still are.
func logRequests(logger *slog.Logger, skipPaths ...string) gin.HandlerFunc {
	skip := make(map[string]bool, len(skipPaths))
	for _, path := range skipPaths {
		skip[path] = true
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		if skip[c.Request.URL.Path] && status < http.StatusBadRequest {
			return
		}

		attrs := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"route", c.FullPath(),
			"status", status,
			"duration", time.Since(start),
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.Errors())
		}

		level := slog.LevelInfo
		if status >= http.StatusInternalServerError {
			level = slog.LevelError
		}
		logger.Log(c.Request.Context(), level, "request", attrs...)
	}
}

// recoverPanics turns a panic in a handler into a logged error and a 500.
func recoverPanics(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			p := recover()
			if p == nil {
				return
			}
			if p == http.ErrAbortHandler {
				panic(p) // Let net/http abort the connection as intended
			}

			logger.ErrorContext(c.Request.Context(), "panic recovered",
				"panic", p,
				"stack", string(debug.Stack()),
			)
			c.Header("Content-Type", "application/problem+json")
			c.AbortWithStatus(http.StatusInternalServerError)
			_ = json.NewEncoder(c.Writer).Encode(huma.NewError(http.StatusInternalServerError, "unexpected error occurred"))
		}()
		c.Next()
	}
}
