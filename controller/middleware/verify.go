package middleware

import (
	"fmt"
	"go_base/domain"
	"go_base/xerror"

	"github.com/labstack/echo/v4"
)

// Verify custom logic to the context using id in the context.
// IMPORTANT: id must be set before this middleware.
func Verify(
	getUser func(ctx echo.Context, id string) (*domain.User, error),
	getStaff func(ctx echo.Context, id string) (*domain.Staff, error),
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := domain.UserID(c)
			if id == "" {
				return xerror.EForbidden().SetDebugInfo("Verify User id", fmt.Sprintf("%v", id))
			}
			isUserKey, _ := c.Get(string(domain.IsUserKey)).(bool)
			// Custom logic for user
			if isUserKey {
				user, err := getUser(c, id)
				if err != nil {
					return xerror.EForbidden().SetDebugInfo("Verify User err", fmt.Sprintf("%v", err))
				}
				if !user.IsVerified {
					return xerror.EForbidden().SetDebugInfo("Verify User IsVerified", fmt.Sprintf("%v", id))
				}

				// ต้องทำให้ส่ง id แล้ว get role แทน
				c.Set(string(domain.UserKey), user)
				return next(c)
			}
			// Custom logic for staff
			if staff, err := getStaff(c, id); err == nil {
				if !staff.IsVerified {
					return xerror.EForbidden().SetDebugInfo("Verify User IsVerified", fmt.Sprintf("%v", id))
				}
				if staff.Status != domain.StaffActive {
					return xerror.EForbidden().SetDebugInfo("Verify User Status", fmt.Sprintf("%v", staff.Status))
				}
				c.Set(string(domain.StaffKey), staff)
				return next(c)
			}
			return xerror.EForbidden().SetDebugInfo("invalid user id", fmt.Sprintf("%v", id))
		}
	}
}
