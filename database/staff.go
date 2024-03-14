package database

import (
	"fmt"
	"go_base/domain"
	"go_base/storage"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stoewer/go-strcase"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	groupStaffCache = "group_cache_staff_%s"
)

type StaffStore struct {
	*BaseStore[domain.Staff, domain.StaffUpdate, domain.StaffCreate]
}

func NewStaffStore(db *gorm.DB, allStorage *storage.AllStorage) *StaffStore {
	config := &BaseStoreConfig{WriteChangelog: true, CacheExpire: time.Minute}
	return &StaffStore{
		BaseStore: NewBaseStore[domain.Staff, domain.StaffUpdate, domain.StaffCreate](db, config, allStorage),
	}
}

func (s *StaffStore) GetByEmail(ctx echo.Context, email domain.SensitiveString) (*domain.Staff, error) {
	var staff domain.Staff
	staff.Email = email

	if err := s.DB.WithContext(ctx.Request().Context()).First(&staff, "email = ?", staff.Email).Error; err != nil {
		return nil, err
	}
	return &staff, nil
}

func (s *StaffStore) DeleteExistData(ctx echo.Context, staff *domain.Staff) error {
	return s.DeleteIfExist(ctx, staff, staff.Email, staff.Email)
}

func (s *StaffStore) GetMe(ctx echo.Context) (*domain.StaffMe, error) {
	staff := domain.StaffFromContext(ctx)
	var staffMe domain.StaffMe
	if err := s.DB.WithContext(ctx.Request().Context()).Preload(clause.Associations).First(&staffMe, "id = ?", staff.ID).Error; err != nil {
		return nil, err
	}
	return &staffMe, nil
}

func (s *StaffStore) UpdateTime(ctx echo.Context, staffID uuid.UUID) error {
	return s.DB.WithContext(ctx.Request().Context()).Table("staffs").Where("id", lo.ToPtr(staffID)).Update("last_login", time.Now()).Error
}

func (s *StaffStore) UpdateTokenVerify(ctx echo.Context, staffID uuid.UUID) error {
	var staff domain.Staff
	if err := s.DB.WithContext(ctx.Request().Context()).First(&staff, "id = ?", staffID).Error; err != nil {
		return err
	}

	staff.IsVerified = true
	staff.VerifyToken = ""
	if err := s.DB.WithContext(ctx.Request().Context()).Save(&staff).Error; err != nil {
		return err
	}
	return nil
}

func (s *StaffStore) CountJsonGroupRole(ctx echo.Context) (*map[string]int64, error) {
	// redis cache
	if cache, err := s.getCache(ctx, fmt.Sprintf(groupStaffCache, "role")); err == nil {
		return cache, nil
	}
	var roles []domain.Role
	if err := s.DB.WithContext(ctx.Request().Context()).Find(&roles).Error; err != nil {
		return nil, err
	}

	countRoles := map[string]int64{}
	countRolesAll := int64(0)
	for _, role := range roles {
		var count int64
		if err := s.DB.WithContext(ctx.Request().Context()).Model(&domain.Staff{}).Where("role_id = ?", role.ID.String()).Count(&count).Error; err != nil {
			return nil, err
		}
		name := strcase.SnakeCase(role.Name)
		countRoles[name] = count
		countRolesAll += count
	}
	countRoles["all"] = countRolesAll

	if err := s.setCache(ctx, fmt.Sprintf(groupStaffCache, "role"), countRoles); err != nil {
		return nil, err
	}
	return &countRoles, nil
}
