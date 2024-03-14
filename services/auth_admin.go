package services

import (
	"go_base/database"
	"go_base/domain"

	"github.com/labstack/echo/v4"
)

type AuthAdminService struct {
	store    *database.AuthStore
	services *domain.AllServices
}

func NewAuthAdminService(store *database.AuthStore, services *domain.AllServices) *AuthAdminService {
	return &AuthAdminService{store: store, services: services}
}

// CreateAuth
func (s *AuthAdminService) CreateAuth(ctx echo.Context, auth *domain.Auth) error {
	if err := s.store.CreateAuth(ctx, auth); err != nil {
		return err
	}
	return nil
}

// UpdateAuth
func (s *AuthAdminService) UpdateAuth(ctx echo.Context, userID string, update domain.Auth) error {
	if err := s.store.UpdateAuth(ctx, userID, update); err != nil {
		return err
	}
	return nil
}

// FindAuth
func (s *AuthAdminService) FindAuth(ctx echo.Context, userID string) (*domain.Auth, error) {
	auth, err := s.store.FindAuth(ctx, userID)
	if err != nil {
		return nil, err
	}
	return auth, nil
}
