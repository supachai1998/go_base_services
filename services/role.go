package services

import (
	"encoding/json"
	"go_base/database"
	"go_base/domain"
	"go_base/logger"
	"go_base/storage"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
)

type RoleService struct {
	store     *database.Store
	services  *domain.AllServices
	cache     *storage.Cache
	roleStore *database.RoleStore
}

func NewRoleService(store *database.Store, role *database.RoleStore, services *domain.AllServices, cache *storage.Cache) *RoleService {
	return &RoleService{store: store, services: services, roleStore: role, cache: cache}
}

func (s *RoleService) HasPermission(ctx echo.Context, roleID *uuid.UUID, requiredPermissions ...string) bool {
	role, err := s.roleStore.GetByID(ctx, roleID.String())
	if err != nil {
		logger.L().Errorf("error while getting role: %v\n", err)
	}
	// if permission is empty, then return false
	if role.Permissions == nil {
		return false
	}

	for _, requiredPermission := range requiredPermissions {
		if hasPermission(role.Permissions, requiredPermission) {
			return true
		}
	}

	return false
}

/*
hasPermission checks if the user has the required permission to access the resource

	permission from the role of the user
	allowedPermission is the permission required to access the resource in controllers
*/
func hasPermission(permission datatypes.JSON, allowedPermission string) bool {
	args := strings.Split(allowedPermission, ".")

	if len(args) < 4 {
		logger.L().Infof("invalid permission format: %s\n", allowedPermission)

		return false
	}

	system, resource, action, scope := args[0], args[1], args[2], args[3]

	var pt domain.PermissionTree
	if err := json.Unmarshal(permission, &pt); err != nil {
		logger.L().Errorf("error while unmarshalling permission: %v\n", err)
	}
	// if found map[*:map[*:map[*:]]] in the permission tree, then return true
	// map[*:map[*:map[*:]]]
	if pt["*"]["*"]["*"] == "*" {
		return true
	}

	// case permission view is disabled return false
	if pt[system][resource]["view"] == "false" {
		return false
	}

	return pt[system][resource][action] == scope
}

func (s *RoleService) Create(ctx echo.Context, role *domain.Role) error {
	return s.roleStore.Create(ctx, role)
}

func (s *RoleService) Update(ctx echo.Context, role *domain.RoleUpdate) error {
	return s.roleStore.UpdateU(ctx, role)
}

func (s *RoleService) GetByID(ctx echo.Context, id string) (*domain.Role, error) {
	return s.roleStore.GetByID(ctx, id)
}

func (s *RoleService) UpdateUserPermissions(ctx echo.Context, userID uuid.UUID, permissionIDs []uuid.UUID) error {
	return s.roleStore.UpdateStaffPermissions(ctx, userID, permissionIDs)
}

// GET /roles
func (s *RoleService) Find(ctx echo.Context, pagination domain.Pagination[domain.Role]) (*domain.Pagination[domain.Role], error) {
	return s.roleStore.Find(ctx, pagination)
}
