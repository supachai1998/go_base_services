package database

import (
	"go_base/domain"
	"go_base/storage"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserStore struct {
	*BaseStore[domain.User, domain.UserUpdate, domain.UserCreate]
}

func NewUserStore(db *gorm.DB, allStorage *storage.AllStorage) *UserStore {
	config := &BaseStoreConfig{WriteChangelog: true, CacheExpire: time.Minute}
	return &UserStore{
		BaseStore: NewBaseStore[domain.User, domain.UserUpdate, domain.UserCreate](db, config, allStorage),
	}
}

func (s *UserStore) DeleteExistData(ctx echo.Context, user *domain.User) error {
	return s.DeleteIfExist(ctx, user, user.Email, user.Email)
}

func (s *UserStore) GetMe(ctx echo.Context) (*domain.UserMe, error) {
	staff := domain.UserFromContext(ctx)
	var staffMe domain.UserMe
	if err := s.DB.WithContext(ctx.Request().Context()).Preload(clause.Associations).First(&staffMe, "id = ?", staff.ID).Error; err != nil {
		return nil, err
	}
	return &staffMe, nil
}

func (s *UserStore) UpdateTime(ctx echo.Context, id uuid.UUID) error {
	return s.DB.WithContext(ctx.Request().Context()).Table("staffs").Where("id", lo.ToPtr(id)).Update("last_login", time.Now()).Error
}

func (s *UserStore) UpdateTokenVerify(ctx echo.Context, id uuid.UUID) error {
	var staff domain.User
	if err := s.DB.WithContext(ctx.Request().Context()).First(&staff, "id = ?", id).Error; err != nil {
		return err
	}

	staff.IsVerified = true
	staff.VerifyToken = ""
	if err := s.DB.WithContext(ctx.Request().Context()).Save(&staff).Error; err != nil {
		return err
	}
	return nil
}
