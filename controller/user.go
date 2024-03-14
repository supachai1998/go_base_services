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

type UserHandler struct {
	Services *domain.AllServices
}

// GET /user
func (h UserHandler) Find(ctx echo.Context) error {
	user, err := h.Services.User.Find(ctx, domain.PaginationFromCtx[domain.User](ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

// (POST /user/login)
func (h UserHandler) LoginWithEmailPassword(ctx echo.Context) error {
	var login domain.UserLogin
	if err := ctx.Bind(&login); err != nil {
		return err
	}
	if err := validate.Struct(login); err != nil {
		return err
	}
	jwt, err := h.Services.User.LoginWithEmailPassword(ctx, login)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	ctx.Response().Header().Set(domain.AuthHeaderKeyUser, domain.BearerKey+jwt.AccessToken)
	return ctx.JSON(http.StatusOK, jwt)
}

// POST /user/unlock
func (h UserHandler) Unlock(ctx echo.Context) error {
	var unlock domain.UserUnlock
	if err := ctx.Bind(&unlock); err != nil {
		return err
	}
	if err := validate.Struct(unlock); err != nil {
		return err
	}
	return h.Services.User.Unlock(ctx, unlock)
}

// POST /user
func (h UserHandler) Create(ctx echo.Context) error {
	var user domain.UserCreate
	if err := ctx.Bind(&user); err != nil {
		return err
	}
	if err := validate.Struct(user); err != nil {
		return err
	}
	s, err := h.Services.User.Create(ctx, user)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, s)
}

// DELETE /user/:id
func (h UserHandler) Delete(ctx echo.Context) error {
	if err := h.Services.User.Delete(ctx); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}

// UPDATE /user/:id
func (h UserHandler) Update(ctx echo.Context) error {
	var user domain.UserUpdate
	if err := ctx.Bind(&user); err != nil {
		return err
	}
	if err := validate.Struct(user); err != nil {
		return err
	}
	_user, err := h.Services.User.Update(ctx, user)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, _user)
}

// MOCK Verify Token /user/verify
func (h UserHandler) VerifyToken(ctx echo.Context) error {
	var user domain.UserVerifyToken
	if err := ctx.Bind(&user); err != nil {
		return xerror.EInvalidInput(nil)
	}
	if err := validate.Struct(user); err != nil {
		return xerror.EInvalidInput(nil)
	}
	_user, err := h.Services.User.Verify(ctx, user)
	if err != nil {
		return xerror.EInvalidInput(nil)
	}
	return ctx.JSON(http.StatusOK, _user)
}

// Mock Get token /user/token
func (h UserHandler) GetToken(ctx echo.Context) error {
	var user domain.UserGetToken
	if err := ctx.Bind(&user); err != nil {
		return xerror.EInvalidInput(nil)
	}
	if err := validate.Struct(user); err != nil {
		return xerror.EInvalidInput(nil)
	}
	_user, err := h.Services.User.GetToken(ctx, user)
	if err != nil {
		return xerror.EInvalidInput(nil)
	}
	return ctx.JSON(http.StatusOK, _user)
}

// Get Log /user/log/:id
func (h UserHandler) GetLog(ctx echo.Context) error {
	_, id := domain.GetUUIDFromParam(ctx, "id")
	user := domain.UserGetLog{ID: id}
	userLog, err := h.Services.User.GetLog(ctx, user)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, userLog)
}

// Update Password /user/me/password
func (h UserHandler) UpdatePassword(ctx echo.Context) error {
	var user domain.UserUpdatePassword
	userCtx := domain.UserFromContext(ctx)
	user.Email = userCtx.Email
	if err := ctx.Bind(&user); err != nil {
		return xerror.EInvalidInput(nil)
	}
	if err := validate.Struct(user); err != nil {
		return xerror.EInvalidInput(nil)
	}

	if err := h.Services.User.UpdatePassword(ctx, user); err != nil {
		return xerror.EInvalidInput(nil)
	}
	return ctx.NoContent(http.StatusOK)
}

// Update Me /user/me
func (h UserHandler) UpdateMe(ctx echo.Context) error {
	var user domain.UserUpdate
	userCtx := domain.UserFromContext(ctx)
	if userCtx == nil {
		return xerror.EUnAuthorized().SetDebugInfo("dev", "user not found in context")
	}
	if err := ctx.Bind(&user); err != nil {
		return xerror.EInvalidInput(nil).SetDebugInfo("dev", "bind error")
	}
	if err := validate.Struct(user); err != nil {
		return xerror.EInvalidInput(nil).SetDebugInfo("dev", "validate error")
	}
	user.ID = userCtx.ID
	if err := h.Services.User.UpdateMe(ctx, user); err != nil {
		return xerror.EInvalidInput(nil).SetDebugInfo("dev", "update error")
	}
	return ctx.NoContent(http.StatusOK)
}

// Get Me /users/me
func (h UserHandler) GetMe(ctx echo.Context) error {
	user, err := h.Services.User.GetMe(ctx)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, user)
}

// Delete /users/ids
func (h UserHandler) DeleteIds(ctx echo.Context) error {
	var ids domain.Ids
	if err := ctx.Bind(&ids); err != nil {
		return xerror.EInvalidInput(nil)
	}
	if err := validate.Struct(ids); err != nil {
		return xerror.EInvalidInput(nil)
	}
	if err := h.Services.User.DeleteByIds(ctx, ids); err != nil {
		return xerror.EInvalidInput(nil)
	}
	return ctx.NoContent(http.StatusOK)
}

// GetMeLogs /users/me/logs
func (h UserHandler) GetLogMe(ctx echo.Context) error {
	user, err := h.Services.User.GetLogMe(ctx)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, user)
}
