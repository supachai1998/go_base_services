package controller

import (
	"go_base/domain"
	"go_base/validate"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProjectHandler struct {
	Services *domain.AllServices
}

// GET /developers
func (h ProjectHandler) Find(ctx echo.Context) error {
	m, err := h.Services.IProject.Find(ctx, domain.PaginationFromCtx[domain.Project](ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, m)
}

// GET /developers/:id
func (h ProjectHandler) Get(ctx echo.Context) error {
	idStr, _ := domain.GetUUIDFromParam(ctx, "id")
	m, err := h.Services.IProject.GET(ctx, idStr)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, m)
}

// POST /developers
func (h ProjectHandler) Create(ctx echo.Context) error {
	var m domain.ProjectCreate
	if err := ctx.Bind(&m); err != nil {
		return err
	}
	if err := validate.Struct(m); err != nil {
		return err
	}
	if err := h.Services.IProject.CreateC(ctx, &m); err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, m)
}

// Update /developers/:id
func (h ProjectHandler) Update(ctx echo.Context) error {
	_, id := domain.GetUUIDFromParam(ctx, "id")
	var m domain.ProjectUpdate
	m.ID = id
	if err := ctx.Bind(&m); err != nil {
		return err
	}
	if err := validate.Struct(m); err != nil {
		return err
	}
	if err := h.Services.IProject.UpdateU(ctx, &m); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, m)
}

// DELETE /developers/:id
func (h ProjectHandler) Delete(ctx echo.Context) error {
	_, id := domain.GetUUIDFromParam(ctx, "id")
	if err := h.Services.IProject.Delete(ctx, id); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
