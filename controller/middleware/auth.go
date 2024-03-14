package middleware

import (
	"context"
	"fmt"
	"go_base/xerror"
	"strings"

	"go_base/domain"

	"github.com/labstack/echo/v4"
)

// Check token in header and verify it. If token is valid, set user id to context.
// For auth middleware, and/or verify middleware.
func Auth(secretAdmin string, secretUser string, cacheFunc func(context.Context, string) (string, error)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get(domain.AuthHeaderKeyUser)
			if header == "" {
				header = c.Request().Header.Get(domain.AuthHeaderKeyStaff)
			}

			var at string
			header = strings.TrimPrefix(header, domain.BearerKey)
			at = strings.TrimSpace(header)
			var claims domain.AuthClaims
			var isUser bool
			if _, err := domain.ParseToken(at, &claims, secretAdmin); err != nil {
				_, userErr := domain.ParseToken(at, &claims, secretUser)
				if userErr != nil {
					return xerror.E(xerror.ErrUnauthorized).SetStatusCode(xerror.ErrCodeUnauthorized)
				}
				isUser = true
			}
			if claims.TokenType != domain.TokenTypeAccess {
				return xerror.E(xerror.ErrUnauthorized).SetStatusCode(xerror.ErrCodeUnauthorized)
			}
			token, err := cacheFunc(c.Request().Context(), fmt.Sprintf(domain.WhitelistAccessTokenCacheKey, claims.UserID))
			if err != nil {
				if xerror.IsNotFoundError(err) {
					return xerror.E(xerror.ErrUnauthorized).SetStatusCode(xerror.ErrCodeUnauthorized)
				}
				return xerror.E(err)
			}
			if token != at {
				return xerror.E(xerror.ErrUnauthorized).SetStatusCode(xerror.ErrCodeUnauthorized)
			}
			// token is expired ?
			if claims.ExpiresAt.Unix() < domain.NowUnix() {
				return xerror.E(xerror.ErrUnauthorized).SetStatusCode(xerror.ErrCodeUnauthorized)
			}
			c.Set(string(domain.UserIDKey), claims.UserID)
			c.Set(string(domain.IsUserKey), isUser)
			return next(c)
		}
	}
}
