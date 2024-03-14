package auth

import (
	"context"
	"fmt"
	"go_base/domain"
	"go_base/xerror"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthConfig struct {
	JWTSecret                 string
	AccessTokenDuration       time.Duration
	RefreshTokenDuration      time.Duration
	VerifyTokenDuration       time.Duration
	AccountLockoutMaxAttempts int
	ResendOTPMaxAttempts      int
	VerifyOTPMaxAttempts      int
	ReturnOTP                 bool
	DemoUser                  struct {
		Email string
		Tel   string
		Pin   string
	}
	LenTempPwd int
}

func issueToken(tokenType domain.TokenType, secret, userID string, duration time.Duration, now time.Time) (*domain.TokenExpires, error) {
	claims := domain.AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		TokenType: tokenType,
		UserID:    userID,
	}

	token, err := domain.GenerateAccessToken(claims, secret)
	if err != nil {
		return nil, xerror.E(err)
	}

	return &domain.TokenExpires{Token: token, ExpireAt: now.Add(duration)}, nil
}

func IssueAccessRefreshToken(ctx echo.Context, userID uuid.UUID, cfg *AuthConfig,
	findFunc func(ctx echo.Context, userID string) (*domain.Auth, error),
	updateFunc func(ctx echo.Context, userID string, update domain.Auth) error,
	createFunc func(ctx echo.Context, auth *domain.Auth) error,
	cacheFunc func(ctx context.Context, key string, value any, exp time.Duration) error,
) (*domain.AuthResult, error) {
	id := userID.String()
	exists := true
	auth, err := findFunc(ctx, id)
	if err != nil {
		if xerror.IsNotFoundError(err) {
			exists = false
		} else {
			return nil, xerror.E(err)
		}
	}

	now := time.Now()
	at, err := issueToken(domain.TokenTypeAccess, cfg.JWTSecret, id, cfg.AccessTokenDuration, now)
	if err != nil {
		return nil, xerror.E(err)
	}

	if err := cacheFunc(ctx.Request().Context(), getWhitelistKey(id), at.Token, at.ExpireAt.Sub(time.Now())); err != nil {
		return nil, xerror.E(err)
	}

	rt, err := issueToken(domain.TokenTypeRefresh, cfg.JWTSecret, id, cfg.RefreshTokenDuration, now)
	if err != nil {
		return nil, xerror.E(err)
	}

	if exists {
		auth.RefreshToken = rt
		if err := updateFunc(ctx, id, *auth); err != nil {
			return nil, xerror.E(err)
		}
	} else {
		auth := &domain.Auth{
			UserID:       userID,
			RefreshToken: rt,
			BaseModel: domain.BaseModel{
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		if err := createFunc(ctx, auth); err != nil {
			return nil, xerror.E(err)
		}
	}

	result := &domain.AuthResult{AccessToken: at.Token, RefreshToken: rt.Token}

	return result, nil
}

func getWhitelistKey(userID string) string {
	return fmt.Sprintf(domain.WhitelistAccessTokenCacheKey, userID)
}
