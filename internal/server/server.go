package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Amin-MAG/md2azw3/config"
	"github.com/Amin-MAG/md2azw3/internal/handler"
	ravandlog "github.com/Amin-MAG/md2azw3/pkg/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// New creates and configures a new Echo server.
func New(cfg config.Config, logger *ravandlog.Logger) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(requestLogger(logger))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Conversion endpoint
	convertHandler := handler.NewConvertHandler(logger)
	e.POST("/convert", convertHandler.Convert)

	return e
}

// Start starts the Echo server on the configured port.
func Start(e *echo.Echo, cfg config.Config, logger *ravandlog.Logger) error {
	addr := fmt.Sprintf(":%d", cfg.MD2AZW3.Port)
	logger.Infof(context.Background(), "starting HTTP server on %s", addr)
	return e.Start(addr)
}

func requestLogger(logger *ravandlog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			ctx := context.WithValue(c.Request().Context(), ravandlog.ContextKeyRequestUUID, reqID)
			c.SetRequest(c.Request().WithContext(ctx))

			err := next(c)

			logger.With("method", c.Request().Method).
				With("uri", c.Request().RequestURI).
				With("status", c.Response().Status).
				With("latency_ms", time.Since(start).Milliseconds()).
				With("request_id", reqID).
				Info(ctx, "request handled")

			return err
		}
	}
}
