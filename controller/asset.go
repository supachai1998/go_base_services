package controller

import (
	"go_base/domain"
	"go_base/validate"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AssetHandler struct {
	Services *domain.AllServices
}

// GET /assets
func (h AssetHandler) Find(ctx echo.Context) error {
	m, err := h.Services.IAsset.Find(ctx, domain.PaginationFromCtx[domain.Asset](ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, m)
}

// GET /assets/:id
func (h AssetHandler) Get(ctx echo.Context) error {
	idStr, _ := domain.GetUUIDFromParam(ctx, "id")
	m, err := h.Services.IAsset.GET(ctx, idStr)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, m)
}

// POST /assets
func (h AssetHandler) Create(ctx echo.Context) error {
	var m domain.AssetCreate
	if err := ctx.Bind(&m); err != nil {
		return err
	}
	if err := validate.Struct(m); err != nil {
		return err
	}
	if err := h.Services.IAsset.CreateC(ctx, &m); err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, m)
}

// Update /assets/:id
func (h AssetHandler) Update(ctx echo.Context) error {
	_, id := domain.GetUUIDFromParam(ctx, "id")
	var m domain.AssetUpdate
	m.ID = id
	if err := ctx.Bind(&m); err != nil {
		return err
	}
	if err := validate.Struct(m); err != nil {
		return err
	}
	if err := h.Services.IAsset.UpdateU(ctx, &m); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, m)
}

// DELETE /assets/:id
func (h AssetHandler) Delete(ctx echo.Context) error {
	_, id := domain.GetUUIDFromParam(ctx, "id")
	if err := h.Services.IAsset.Delete(ctx, id); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
