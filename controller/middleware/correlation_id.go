package middleware

import (
	"context"

	"go_base/domain"
	"go_base/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CorrelationID sets correlation id to the X-Correlation-Id header and attaches it to the request context.
func CorrelationID() echo.MiddlewareFunc {
	return middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: func(c echo.Context, id string) {
			l := logger.L().With("trace.id", id)

			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, domain.CorrelationIDKey, id)
			ctx = logger.ContextWithLogger(ctx, l)

			c.SetRequest(c.Request().WithContext(ctx))
		},
		TargetHeader: echo.HeaderXCorrelationID,
		Generator:    uuid.NewString,
	})
}
