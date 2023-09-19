// Package ginx provides an extended gin.Engine.
package ginx

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alukart32/effective-mobile-test-task/internal/pkg/zerologx"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

var (
	once sync.Once
	g    *gin.Engine
)

// Get returns extended gin.Engine instance.
func Get() (*gin.Engine, error) {
	var err error

	once.Do(func() {
		mode := os.Getenv("GIN_MODE")
		if len(mode) == 0 {
			mode = gin.DebugMode // default to debug mode
		}

		gin.SetMode(mode)

		g = gin.New()

		// gin Logger: os.Stdout (default)
		logger := &zerologHandler{}
		g.Use(logger.Handle)

		// gin Recovery
		g.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
			if err, ok := recovered.(string); ok {
				c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
			}
			c.AbortWithStatus(http.StatusInternalServerError)
		}))
	})

	return g, err
}

// zerologHandler represents a logger for gin router.
type zerologHandler struct{}

// CorrID represents the corresponding ID for the request.
type CorrID string

// Handle adds zerlog context to the request context.
func (h *zerologHandler) Handle(c *gin.Context) {
	t := time.Now()

	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery

	correlationID := xid.New().String()
	ctx := context.WithValue(c.Request.Context(), CorrID("correlation_id"), correlationID)
	c.Request = c.Request.WithContext(ctx)

	l := zerologx.Get().
		With().
		Str("correlation_id", correlationID).
		Logger()

	l.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("correlation_id", correlationID)
	})

	c.Request = c.Request.WithContext(l.WithContext(c.Request.Context()))

	c.Next()

	if raw != "" {
		path = path + "?" + raw
	}

	l.Info().
		Str("method", c.Request.Method).
		Str("path", path).
		Int("status", c.Writer.Status()).
		Dur("elapsed_ms", time.Since(t)).
		Msg("incoming request")
}
