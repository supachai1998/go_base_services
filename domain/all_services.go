package domain

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// T = Model DB, U = UpdateModel, C = CreateModel
type IBaseService[T, U, C any] interface {
	GET(ctx echo.Context, id string) (*T, error)
	Create(ctx echo.Context, m *T) error
	CreateC(ctx echo.Context, m *C) error
	Update(ctx echo.Context, m *T) error
	UpdateU(ctx echo.Context, m *U) error
	Delete(ctx echo.Context, id uuid.UUID) error
	Find(ctx echo.Context, pagination Pagination[T]) (*Pagination[T], error)
	FindWithUserID(ctx echo.Context, pagination Pagination[T], ignoreRelations ...string) (*Pagination[T], error)
	GetWithUserID(ctx echo.Context, idStr string) (*T, error)
	UpdateWithUserID(ctx echo.Context, model *U, typeLog ...string) error
	DeleteWithUserID(ctx echo.Context, id uuid.UUID) error
}

type AllServices struct {
	Staff      StaffService
	AuthAdmin  AdminAuthService
	AuthUser   UserAuthService
	Role       RoleService
	User       UserService
	IDeveloper IBaseService[Developer, DeveloperUpdate, DeveloperCreate]
	IProject   IBaseService[Project, ProjectUpdate, ProjectCreate]
	IAsset     IBaseService[Asset, AssetUpdate, AssetCreate]
	Asset      IAssetService[Asset, AssetUpdate, AssetCreate]
}
