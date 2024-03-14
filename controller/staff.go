package controller

import (
	"errors"
	"go_base/domain"
	"go_base/validate"
	"go_base/xerror"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type StaffHandler struct {
	Services *domain.AllServices
}

// GET /staffs
func (h StaffHandler) Find(ctx echo.Context) error {
	staffs, err := h.Services.Staff.Find(ctx, domain.PaginationFromCtx[domain.Staff](ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, staffs)
}

// (POST /staff/login)
func (h StaffHandler) LoginWithEmailPassword(ctx echo.Context) error {
	var login domain.StaffLogin
	if err := ctx.Bind(&login); err != nil {
		return err
	}
	if err := validate.Struct(login); err != nil {
		return err
	}
	jwt, err := h.Services.Staff.LoginWithEmailPassword(ctx, login)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	ctx.Response().Header().Set(domain.AuthHeaderKeyStaff, domain.BearerKey+jwt.AccessToken)
	return ctx.JSON(http.StatusOK, jwt)
}

// POST /staff/unlock
func (h StaffHandler) Unlock(ctx echo.Context) error {
	var unlock domain.StaffUnlock
	if err := ctx.Bind(&unlock); err != nil {
		return err
	}
	if err := validate.Struct(unlock); err != nil {
		return err
	}
	return h.Services.Staff.Unlock(ctx, unlock)
}

// POST /staffs
func (h StaffHandler) Create(ctx echo.Context) error {
	var staff domain.StaffCreate
	if err := ctx.Bind(&staff); err != nil {
		return err
	}
	if err := validate.Struct(staff); err != nil {
		return err
	}
	s, err := h.Services.Staff.Create(ctx, staff)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, s)
}

// DELETE /staff/:id
func (h StaffHandler) Delete(ctx echo.Context) error {
	if err := h.Services.Staff.Delete(ctx); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}

// UPDATE /staff/:id
func (h StaffHandler) Update(ctx echo.Context) error {
	var staff domain.StaffUpdate
	_, id := domain.GetUUIDFromParam(ctx, "id")
	staff.ID = id
	if err := ctx.Bind(&staff); err != nil {
		return err
	}
	if err := validate.Struct(staff); err != nil {
		return err
	}
	_staff, err := h.Services.Staff.Update(ctx, staff)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, _staff)
}

// MOCK Verify Token /staff/verify
func (h StaffHandler) VerifyToken(ctx echo.Context) error {
	var staff domain.StaffVerifyToken
	if err := ctx.Bind(&staff); err != nil {
		return xerror.EInvalidInput(nil)
	}
	if err := validate.Struct(staff); err != nil {
		return xerror.EInvalidInput(nil)
	}
	_staff, err := h.Services.Staff.Verify(ctx, staff)
	if err != nil {
		return xerror.EInvalidInput(nil)
	}
	return ctx.JSON(http.StatusOK, _staff)
}

// Mock Get token /staff/token
func (h StaffHandler) GetToken(ctx echo.Context) error {
	var staff domain.StaffGetToken
	if err := ctx.Bind(&staff); err != nil {
		return xerror.EInvalidInput(nil)
	}
	if err := validate.Struct(staff); err != nil {
		return xerror.EInvalidInput(nil)
	}
	_staff, err := h.Services.Staff.GetToken(ctx, staff)
	if err != nil {
		return xerror.EInvalidInput(nil)
	}
	return ctx.JSON(http.StatusOK, _staff)
}

// Get Log /staff/log/:id
func (h StaffHandler) GetLog(ctx echo.Context) error {
	_, id := domain.GetUUIDFromParam(ctx, "id")
	staff := domain.StaffGetLog{ID: id}
	staffLog, err := h.Services.Staff.GetLog(ctx, staff)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, staffLog)
}
