package auth

import (
	"errors"
	"go_base/domain"
	"go_base/storage"
	"go_base/validate"
	"go_base/xerror"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/email"
)

type UserAuthService struct {
	userService domain.UserService
	store       domain.UserAuthStore
	cfg         *AuthConfig
	cache       *storage.Cache
	email       email.Email
}

func (s *UserAuthService) validateUserAccount(ctx echo.Context, creds domain.LoginCredentials) (*domain.User, int, error) {
	if err := validate.Struct(creds); err != nil {
		return nil, 0, xerror.EInvalidInput(err)
	}

	user, err := s.userService.GetByEmail(ctx, creds.Email)
	if err != nil {
		if xerror.IsNotFoundError(err) {
			return nil, 0, xerror.E(xerror.ErrUnauthorized).SetErrorCode(xerror.ErrInvalidCredentials).
				SetStatusCode(xerror.ErrCodeUnauthorized)
		}
		return nil, 0, xerror.E(err)
	}

	if user.Password == "" {
		return nil, 0, xerror.EInvalidInput(errors.New("please setup new password before login")).SetErrorCode(xerror.ErrAuthAdminLoginMustSetPassword)
	}

	loginAttempt, err := s.cache.GetStrikes(ctx.Request().Context(), getUserLoginAttemptCountKey(user.ID.String()))
	if err != nil {
		return nil, 0, xerror.E(err)
	}
	if loginAttempt >= s.cfg.AccountLockoutMaxAttempts {
		return nil, 0, xerror.EInvalidInput(errors.New("login attempt reach limit")).SetErrorCode(xerror.ErrAuthAdminLoginReachLimit)
	}

	return user, loginAttempt, nil
}

func getUserLoginAttemptCountKey(userID string) string {
	return "user:login_attempt_count:" + userID
}

func getUserLoginKey(email string) string {
	return "user:login:" + email
}

func getUserResendOTPKey(email string) string {
	return "user:resend:" + email
}

func getUserResendOTPLockedKey(email string) string {
	return "user:resend_otp_locked:" + email
}

func getUserVerifyOTPKey(ref string) string {
	return "user:verify:" + ref
}

func getUserVerifyOTPLockedKey(email string) string {
	return "user:verify_otp_locked:" + email
}

func getUserLoggedInKey(email string) string {
	return "user:logged_in:" + email
}
