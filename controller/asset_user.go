package controller

import (
	"go_base/domain"
	"go_base/validate"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type AssetUserHandler struct {
	Services *domain.AllServices
}

// GET /assets/user
func (h AssetHandler) FindUser(ctx echo.Context) error {
	m, err := h.Services.IAsset.FindWithUserID(ctx, domain.PaginationFromCtx[domain.Asset](ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, m)
}

// GET /assets/user/:id
func (h AssetHandler) GetUser(ctx echo.Context) error {
	idStr, _ := domain.GetUUIDFromParam(ctx, "id")
	m, err := h.Services.IAsset.GetWithUserID(ctx, idStr)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, m)
}

// POST /assets/user
func (h AssetHandler) CreateUser(ctx echo.Context) error {
	user := domain.UserFromContext(ctx)
	var m domain.AssetCreate
	if err := ctx.Bind(&m); err != nil {
		return err
	}
	if err := validate.Struct(m); err != nil {
		return err
	}
	m.UserID = lo.ToPtr(user.ID)
	if err := h.Services.Asset.CreateScope(ctx, &m); err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, m)
}

// Update /assets/user/:id
func (h AssetHandler) UpdateUser(ctx echo.Context) error {
	_, id := domain.GetUUIDFromParam(ctx, "id")
	var m domain.AssetUpdate
	m.ID = id
	if err := ctx.Bind(&m); err != nil {
		return err
	}
	if err := validate.Struct(m); err != nil {
		return err
	}
	if err := h.Services.IAsset.UpdateWithUserID(ctx, &m); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, m)
}

// DELETE /assets/user/:id
func (h AssetHandler) DeleteUser(ctx echo.Context) error {
	_, id := domain.GetUUIDFromParam(ctx, "id")
	if err := h.Services.IAsset.DeleteWithUserID(ctx, id); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
