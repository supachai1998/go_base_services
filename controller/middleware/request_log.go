package middleware

import (
	"strings"

	"go_base/domain"
	"go_base/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func BodyDumpWithConfig(env string) echo.MiddlewareFunc {
	return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Skipper: func(c echo.Context) bool {
			return env == "production" || strings.HasPrefix(c.Request().URL.String(), "/v1/docs")
		},
		Handler: func(c echo.Context, req []byte, res []byte) {
			c.Set("req_body", req)
			c.Set("res_body", res)
		},
	})
}

// RequestLogger returns an Echo middleware that logs request information.
// The middleware logs the following information:
// - Latency
// - Remote IP
// - Host
// - HTTP Method
// - URI
// - User Agent
// - Response Status
// - Request Body (if available)
// - Response Body (if available)
// - Correlation ID
//
// The middleware can be configured to log different combinations of request and response information.
// The `env` parameter specifies the environment (e.g., "development", "production") in which the middleware is used.
// The `req` parameter specifies whether to log request information.
// The `rep` parameter specifies whether to log response information.
// The middleware returns an Echo MiddlewareFunc that can be used with the Echo framework.
func RequestLogger(env string, req, rep bool) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:   true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURI:       true,
		LogUserAgent: true,
		LogStatus:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if rep && req {
				reqBody, _ := c.Get("req_body").([]byte)
				resBody, _ := c.Get("res_body").([]byte)

				logger.L().Infow("request",
					"remote_ip", v.RemoteIP,
					"user_agent", v.UserAgent,
					"host", v.Host,
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					zap.ByteString("request_body", reqBody),
					zap.ByteString("response_body", resBody),
					"correlation_id", c.Request().Context().Value(domain.CorrelationIDKey),
				)
				return nil
			}
			if req && !rep {
				reqBody, _ := c.Get("req_body").([]byte)
				logger.L().Infow("request",
					"remote_ip", v.RemoteIP,
					"user_agent", v.UserAgent,
					"host", v.Host,
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					zap.ByteString("request_body", reqBody),
					"correlation_id", c.Request().Context().Value(domain.CorrelationIDKey),
				)
				return nil
			}
			if !req && rep {
				resBody, _ := c.Get("res_body").([]byte)
				logger.L().Infow("request",
					"remote_ip", v.RemoteIP,
					"user_agent", v.UserAgent,
					"host", v.Host,
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"latency", v.Latency,
					zap.ByteString("response_body", resBody),
					"correlation_id", c.Request().Context().Value(domain.CorrelationIDKey),
				)
				return nil
			}
			return nil
		},
	})
}
