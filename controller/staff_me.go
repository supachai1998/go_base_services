package controller

import (
	"go_base/domain"
	"go_base/validate"
	"go_base/xerror"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StaffMeHandler struct {
	Services *domain.AllServices
}

// Get log me delete /staff/me/log
func (h StaffMeHandler) GetLogMeDelete(ctx echo.Context) error {
	stafflog, err := h.Services.Staff.GetLogMe(ctx)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, stafflog)
}

// Update Password /staff/password
func (h StaffMeHandler) UpdatePassword(ctx echo.Context) error {
	var staff domain.StaffUpdatePassword
	staffCtx := domain.StaffFromContext(ctx)
	staff.Email = staffCtx.Email
	if err := ctx.Bind(&staff); err != nil {
		return xerror.EInvalidInput(nil)
	}
	if err := validate.Struct(staff); err != nil {
		return xerror.EInvalidInput(nil)
	}

	if err := h.Services.Staff.UpdatePassword(ctx, staff); err != nil {
		return xerror.EInvalidInput(nil)
	}
	return ctx.NoContent(http.StatusOK)
}

// Get me /staff/me
func (h StaffMeHandler) GetMe(ctx echo.Context) error {
	staff, err := h.Services.Staff.GetMe(ctx)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, staff)
}
