package services

import (
	"go_base/database"
	"go_base/domain"
	"go_base/storage"
	"reflect"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type BaseService[T, U, C any] struct {
	store     *database.Store
	services  *domain.AllServices
	cache     *storage.Cache
	baseStore *database.BaseStore[T, U, C]
}

func NewBaseService[T, U, C any](store *database.Store, base *database.BaseStore[T, U, C], services *domain.AllServices, cache *storage.Cache) *BaseService[T, U, C] {
	var m T
	if reflect.TypeOf(m) != nil {
		store.DB.AutoMigrate(m)
	}
	return &BaseService[T, U, C]{store: store, services: services, baseStore: base, cache: cache}
}

func (s *BaseService[T, U, C]) GET(ctx echo.Context, id string) (*T, error) {
	return s.baseStore.GetByID(ctx, id)
}

func (s *BaseService[T, U, C]) Create(ctx echo.Context, m *T) error {
	return s.baseStore.Create(ctx, m)
}

func (s *BaseService[T, U, C]) CreateC(ctx echo.Context, m *C) error {
	return s.baseStore.CreateC(ctx, m)
}

func (s *BaseService[T, U, C]) Update(ctx echo.Context, m *T) error {
	return s.baseStore.Update(ctx, m)
}
func (s *BaseService[T, U, C]) UpdateU(ctx echo.Context, m *U) error {
	return s.baseStore.UpdateU(ctx, m)
}

func (s *BaseService[T, U, C]) Delete(ctx echo.Context, id uuid.UUID) error {
	return s.baseStore.Delete(ctx, id)
}

func (s *BaseService[T, U, C]) Find(ctx echo.Context, pagination domain.Pagination[T]) (*domain.Pagination[T], error) {
	return s.baseStore.Find(ctx, pagination)
}

// for user_id
func (s *BaseService[T, U, C]) FindWithUserID(ctx echo.Context, pagination domain.Pagination[T], ignoreRelations ...string) (*domain.Pagination[T], error) {
	return s.baseStore.FindWithUserID(ctx, pagination, ignoreRelations...)
}

func (s *BaseService[T, U, C]) GetWithUserID(ctx echo.Context, idStr string) (*T, error) {
	return s.baseStore.GetWithUserID(ctx, idStr)
}

func (s *BaseService[T, U, C]) UpdateWithUserID(ctx echo.Context, model *U, typeLog ...string) error {
	return s.baseStore.UpdateWithUserID(ctx, model, typeLog...)
}

func (s *BaseService[T, U, C]) DeleteWithUserID(ctx echo.Context, id uuid.UUID) error {
	return s.baseStore.DeleteWithUserID(ctx, id)
}
