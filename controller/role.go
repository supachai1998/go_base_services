package controller

import (
	"fmt"
	"go_base/domain"
	"go_base/validate"
	"go_base/xerror"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RoleHandler struct {
	Services *domain.AllServices
}

// Create new role /roles
func (h RoleHandler) Create(ctx echo.Context) error {
	var role domain.RoleSwaggerCreate
	if err := ctx.Bind(&role); err != nil {
		return xerror.EInvalidInput(err)
	}
	if err := validate.Struct(role); err != nil {
		return xerror.EInvalidInput(err)
	}
	if err := h.Services.Role.Create(ctx, &domain.Role{
		Type:        role.Type,
		Name:        role.Name,
		Description: role.Description,
		Permissions: role.Permissions,
	}); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusCreated)
}

// Find roles /roles
func (h RoleHandler) Find(ctx echo.Context) error {
	roles, err := h.Services.Role.Find(ctx, domain.PaginationFromCtx[domain.Role](ctx))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, roles)
}

// Update role /roles/:id
func (h RoleHandler) Update(ctx echo.Context) error {
	id, uid := domain.GetUUIDFromParam(ctx, "id")
	if uid == uuid.Nil {
		return xerror.EInvalidInput(fmt.Errorf("invalid id: %s", id))
	}
	var role domain.RoleUpdate
	role.ID = uid
	if err := ctx.Bind(&role); err != nil {
		return xerror.EInvalidInput(err)
	}
	if err := validate.Struct(role); err != nil {
		return xerror.EInvalidInput(err)
	}

	if err := h.Services.Role.Update(ctx, &role); err != nil {
		return err
	}
	return ctx.NoContent(http.StatusOK)
}
