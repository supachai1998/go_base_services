package middleware

import (
	"go_base/domain"
	"go_base/xerror"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RestrictPermissions(
	hasPermissionFn func(ctx echo.Context, roleID *uuid.UUID, requiredPermissions ...string) bool,
) func(...string) echo.MiddlewareFunc {
	return func(permissionNames ...string) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				// ctx := c.Request().Context()
				currentStaff := domain.StaffFromContext(c)
				if currentStaff == nil {
					return xerror.E(xerror.ErrForbidden).SetStatusCode(xerror.ErrCodeForbidden).SetDebugInfo("msg", "currentStaff is nil")
				}
				if len(permissionNames) == 0 {
					return next(c)
				}
				if currentStaff.RoleID == nil {
					return xerror.E(xerror.ErrForbidden).SetStatusCode(xerror.ErrCodeForbidden).SetDebugInfo("msg", "currentStaff.Role is nil")
				}
				if hasPermissionFn(c, currentStaff.RoleID, permissionNames...) {
					return next(c)
				}

				return xerror.E(xerror.ErrForbidden).SetStatusCode(xerror.ErrCodeForbidden)
			}
		}
	}
}
