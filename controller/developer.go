package controller

import (
	"go_base/domain"
	"go_base/validate"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DeveloperHandler struct {
	Services *domain.AllServices
}

// GET /developers
func (h DeveloperHandler) Find(ctx echo.Context) error {
	m, err := h.Services.IDeveloper.Find(ctx, domain.PaginationFromCtx[domain.Developer](ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, m)
}

// GET /developers/:id
func (h DeveloperHandler) Get(ctx echo.Context) error {
	idStr, _ := domain.GetUUIDFromParam(ctx, "id")
	m, err := h.Services.IDeveloper.GET(ctx, idStr)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, m)
}

// POST /developers
func (h DeveloperHandler) Create(ctx echo.Context) error {
	var m domain.DeveloperCreate
	if err := ctx.Bind(&m); err != nil {
		return err
	}
	if err := validate.Struct(m); err != nil {
		return err
	}
	if err := h.Services.IDeveloper.CreateC(ctx, &m); err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, m)
}

// Update /developers/:id
func (h DeveloperHandler) Update(ctx echo.Context) error {
	_, id := domain.GetUUIDFromParam(ctx, "id")
	var m domain.DeveloperUpdate
	m.ID = id
	if err := ctx.Bind(&m); err != nil {
		return err
	}
	if err := validate.Struct(m); err != nil {
		return err
	}
	if err := h.Services.IDeveloper.UpdateU(ctx, &m); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, m)
}

// DELETE /developers/:id
func (h DeveloperHandler) Delete(ctx echo.Context) error {
	_, id := domain.GetUUIDFromParam(ctx, "id")
	if err := h.Services.IDeveloper.Delete(ctx, id); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
