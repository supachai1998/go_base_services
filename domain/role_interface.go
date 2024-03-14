package domain

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RoleService interface {
	Create(ctx echo.Context, role *Role) error
	Update(ctx echo.Context, role *RoleUpdate) error
	GetByID(ctx echo.Context, id string) (*Role, error)
	Find(ctx echo.Context, pagination Pagination[Role]) (*Pagination[Role], error)
	HasPermission(ctx echo.Context, roleID *uuid.UUID, requiredPermissions ...string) bool
	// FindList(ctx context.Context, filter *Filter[RoleFilter]) (*Pagination[*Model[*RoleWithStaffCount]], error)
	// GetByTypeName(ctx context.Context, roleType RoleType, name string) (*Model[*Role], error)
	// GetByIDs(ctx context.Context, IDs []uuid.UUID) ([]*Model[*Role], error)
	// CreateWithPermissions(ctx context.Context, payload Role) (*Model[*Role], error)
	// UpdateWithPermissions(ctx context.Context, id string, payload Role) (*Model[*Role], error)
	// G
	// GetNames(ctx context.Context, roleType RoleType) ([]string, error)
}
