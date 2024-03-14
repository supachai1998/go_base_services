package services

import (
	"go_base/database"
	"go_base/domain"

	"github.com/labstack/echo/v4"
)

type AuthUserService struct {
	store    *database.AuthStore
	services *domain.AllServices
}

func NewAuthUserService(store *database.AuthStore, services *domain.AllServices) *AuthUserService {
	return &AuthUserService{store: store, services: services}
}

// CreateAuth
func (s *AuthUserService) CreateAuth(ctx echo.Context, auth *domain.Auth) error {
	if err := s.store.CreateAuth(ctx, auth); err != nil {
		return err
	}
	return nil
}

// UpdateAuth
func (s *AuthUserService) UpdateAuth(ctx echo.Context, userID string, update domain.Auth) error {
	if err := s.store.UpdateAuth(ctx, userID, update); err != nil {
		return err
	}
	return nil
}

// FindAuth
func (s *AuthUserService) FindAuth(ctx echo.Context, userID string) (*domain.Auth, error) {
	auth, err := s.store.FindAuth(ctx, userID)
	if err != nil {
		return nil, err
	}
	return auth, nil
}
