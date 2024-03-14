package services

import (
	"errors"
	"go_base/database"
	"go_base/domain"
	helper "go_base/domain/helper"
	"go_base/storage"
	"go_base/xerror"

	"github.com/labstack/echo/v4"
)

type IBaseService[T domain.Asset, U domain.AssetUpdate, C domain.AssetCreate] struct {
	store     *database.Store
	services  *domain.AllServices
	cache     *storage.Cache
	baseStore *database.BaseStore[T, U, C]
}

func NewAssetService[T domain.Asset, U domain.AssetUpdate, C domain.AssetCreate](store *database.Store, base *database.BaseStore[T, U, C], services *domain.AllServices, cache *storage.Cache) *IBaseService[T, U, C] {
	return &IBaseService[T, U, C]{store: store, services: services, baseStore: base, cache: cache}
}

func (s *IBaseService[T, U, C]) CreateScope(ctx echo.Context, c *C) error {

	user := domain.UserFromContext(ctx)
	cTmp := helper.Copy[domain.AssetCreate](c)

	if user == nil {
		if err := s.store.DB.First(&user, cTmp.UserID).Error; err != nil {
			return xerror.ENotFound()
		}
	}

	if err := s.store.DB.Where("user_id = ? AND project_id = ?", cTmp.UserID, cTmp.ProjectID).First(&domain.Asset{}).Error; err == nil {
		return xerror.EConflict(errors.New("you have already had this asset"))
	}

	return s.baseStore.CreateC(ctx, c)
}
