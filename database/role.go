package database

import (
	"go_base/domain"
	"go_base/storage"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type RoleStore struct {
	*BaseStore[domain.Role, domain.RoleUpdate, domain.Role]
}

func NewRoleStore(db *gorm.DB, allStorage *storage.AllStorage) *RoleStore {
	return &RoleStore{
		BaseStore: NewBaseStore[domain.Role, domain.RoleUpdate, domain.Role](db, &BaseStoreConfig{WriteChangelog: true}, allStorage),
	}
}

func (s *RoleStore) GetRolesByRoleIDs(ctx echo.Context, roleIDs []uuid.UUID) ([]*domain.Role, error) {
	var roles []*domain.Role
	if err := s.DB.WithContext(ctx.Request().Context()).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RoleStore) UpdateStaffPermissions(ctx echo.Context, staffID uuid.UUID, roleIDs []uuid.UUID) error {
	return s.DB.WithContext(ctx.Request().Context()).Model(&domain.Staff{
		BaseModel: domain.BaseModel{ID: staffID},
	}).Association("Roles").Replace(roleIDs)
}

func (s *RoleStore) Find(ctx echo.Context, pagination domain.Pagination[domain.Role]) (*domain.Pagination[domain.Role], error) {
	var _db = s.DB.WithContext(ctx.Request().Context())
	roles, err := pagination.Paginate(ctx, _db)
	if err != nil {
		return nil, err
	}
	// find roles with CountStaff
	for i := range roles.Items {
		var count int64
		_db.Model(&domain.Staff{}).Where("role_id = ?", roles.Items[i].ID).Count(&count)
		roles.Items[i].CountStaff = lo.ToPtr(count)
	}
	return roles, nil

}
