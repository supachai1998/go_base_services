package middleware

import (
	"fmt"
	"go_base/domain"
	"go_base/xerror"

	"github.com/labstack/echo/v4"
)

// Attach to the context using id in the context.
// IMPORTANT: id must be set before this middleware.
func Attach(
	getUser func(ctx echo.Context, id string) (*domain.User, error),
	getStaff func(ctx echo.Context, id string) (*domain.Staff, error),
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := domain.UserID(c)
			if id == "" {
				return xerror.EForbidden().SetDebugInfo("AttachUser id", fmt.Sprintf("%v", id))
			}
			isUserKey, _ := c.Get(string(domain.IsUserKey)).(bool)
			if isUserKey {
				user, err := getUser(c, id)
				if err != nil {
					return xerror.EForbidden().SetDebugInfo("AttachUser err", fmt.Sprintf("%v", err))
				}
				// ต้องทำให้ส่ง id แล้ว get role แทน
				c.Set(string(domain.UserKey), user)
				return next(c)
			}
			if staff, err := getStaff(c, id); err == nil {
				c.Set(string(domain.StaffKey), staff)
				return next(c)
			}
			return xerror.EForbidden().SetDebugInfo("invalid user id", fmt.Sprintf("%v", id))
		}
	}
}
