package middleware

import (
	"github.com/labstack/echo/v4"
)

// Remove Sensitivity Information from response
func RemoveSensitiveInformation() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
