package database

import (
	"go_base/domain"
	"go_base/storage"
	"go_base/xerror"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var (
	AuthAction = "auth"
)

type AuthStore struct {
	*BaseStore[domain.Auth, domain.AuthUpdate, domain.Auth]
}

func NewAuthStore(db *gorm.DB, allStorage *storage.AllStorage) *AuthStore {
	return &AuthStore{NewBaseStore[domain.Auth, domain.AuthUpdate, domain.Auth](db, &BaseStoreConfig{
		WriteChangelog: true,
		CacheExpire:    time.Minute,
	}, allStorage)}
}

func (s *AuthStore) CreateAuth(ctx echo.Context, auth *domain.Auth) error {
	if err := s.Create(ctx, auth); err != nil {
		return err
	}
	return nil
}

func (s *AuthStore) UpdateAuth(ctx echo.Context, userID string, update domain.Auth) error {
	id, idUUID := domain.GetUUID(userID)
	if idUUID == uuid.Nil {
		return xerror.EInvalidParameter(nil)
	}
	if err := s.UpdateWhereID(ctx, &update, id, domain.LoginLog); err != nil {
		return err
	}
	return nil
}

func (s *AuthStore) FindAuth(ctx echo.Context, userID string) (*domain.Auth, error) {
	id, idUUID := domain.GetUUID(userID)
	if idUUID == uuid.Nil {
		return nil, xerror.EInvalidParameter(nil)
	}
	var result domain.Auth
	if err := s.DB.WithContext(ctx.Request().Context()).Where("user_id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
