package domain

import (
	"context"
	"go_base/xerror"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	BearerKey                    = "Bearer "
	WhitelistAccessTokenCacheKey = "whitelist:access_token:%s"

	AuthHeaderKeyStaff = "Authorization_Staff"
	AuthHeaderKeyUser  = "Authorization"
)

type Auth struct {
	BaseModel
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;primaryKey"`

	// foreign key TokenExpires
	RefreshToken   *TokenExpires `json:"refresh_token" gorm:"foreignKey:TokenExpiresID"`
	TokenExpiresID *string       `json:"-" gorm:"type:varchar(255)"`
}
type TokenExpires struct {
	BaseModel
	Token    string    `json:"token" gorm:"type:text;uniqueIndex"`
	ExpireAt time.Time `json:"expire_at"`
}

type AuthUpdate struct {
	TokenExpiresID string        `json:"-" gorm:"type:varchar(255)"`
	RefreshToken   *TokenExpires `json:"refresh_token"`
	UpdateAt       time.Time     `json:"update_at"`
}

func (rt TokenExpires) Expired() bool {
	return TimeNow().After(rt.ExpireAt)
}

type LoginCredentials struct {
	Email    SensitiveString `json:"email" validate:"required,email" `
	Password string          `json:"password" validate:"required"`
}

type AuthResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	DeviceID     string `json:"device_id,omitempty"`
}

type TokenType string

const (
	TokenTypeAccess  = "access_token"
	TokenTypeRefresh = "refresh_token"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	UserID    string    `json:"user_id"`
	TokenType TokenType `json:"token_type"`
}

type UserDevice struct {
	DeviceID string `json:"device_id" validate:"required,uuid4"`
}

type UserDeviceLocked struct {
	Locked bool `json:"locked"`
}

type UserDeviceNonce struct {
	Nonce string `json:"nonce"`
}

// swagger param header
type JWTHeader struct {
	Header string `json:"Authorization" swagger:"required"`
}

type AdminAuthService interface {
	// Login(ctx context.Context, creds LoginCredentials) (*AuthResult, error)
	// Logout(ctx context.Context, userID string) error
	// Refresh(ctx context.Context, refreshToken string) (*AuthResult, error)
	// ClearAdminLoginAttemptCount(ctx context.Context, userID string) error
	// ClearAdminWhitelistAccessToken(ctx context.Context, userID string) error
	// SecureLogin(ctx context.Context, creds LoginCredentials) (*ReferenceCode, error)
	// VerifyOTP(ctx context.Context, req VerifyEmailOTP) (*AuthResult, error)
	// ResendOTP(ctx context.Context, email UsersEmail) (*ReferenceCode, error)
	CreateAuth(ctx echo.Context, auth *Auth) error
	UpdateAuth(ctx echo.Context, userID string, update Auth) error
	FindAuth(ctx echo.Context, userID string) (*Auth, error)
}

type AdminAuthStore interface {
	CreateAdminAuth(ctx context.Context, auth *Auth) error
	FindAdminAuthByUserID(ctx context.Context, userID string) (*Auth, error)
	FindAdminAuthByRefreshToken(ctx context.Context, token string) (*Auth, error)
	UpdateAdminAuth(ctx context.Context, userID string, update AuthUpdate) error
}

type UserAuthService interface {
	// Login(ctx context.Context, creds LoginCredentials) (*AuthResult, error)
	// Logout(ctx context.Context, userID string) error
	// Refresh(ctx context.Context, refreshToken string) (*AuthResult, error)
	// ClearUserLoginAttemptCount(ctx context.Context, userID string) error
	// ClearUserWhitelistAccessToken(ctx context.Context, userID string) error
	// SecureLogin(ctx context.Context, creds LoginCredentials) (*ReferenceCode, error)
	// VerifyOTP(ctx context.Context, req VerifyEmailOTP) (*AuthResult, error)
	// ResendOTP(ctx context.Context, email UsersEmail) (*ReferenceCode, error)
	CreateAuth(ctx echo.Context, auth *Auth) error
	UpdateAuth(ctx echo.Context, userID string, update Auth) error
	FindAuth(ctx echo.Context, userID string) (*Auth, error)
}

type UserAuthStore interface {
	CreateUserAuth(ctx context.Context, auth *Auth) error
	FindUserAuthByUserID(ctx context.Context, userID string) (*Auth, error)
	FindUserAuthByRefreshToken(ctx context.Context, token string) (*Auth, error)
	UpdateUserAuth(ctx context.Context, userID string, update AuthUpdate) error
}

func GenerateAccessToken(claims jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr string, claims jwt.Claims, secret string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, xerror.E(xerror.ErrUnauthorized).SetStatusCode(xerror.ErrCodeUnauthorized)
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, xerror.E(xerror.ErrUnauthorized).SetStatusCode(xerror.ErrCodeUnauthorized).SetDebugInfo("token", tokenStr).SetDebugInfo("invalid_token", err)
	}

	return token, nil
}
