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

type AdminAuthService struct {
	staffService domain.StaffService
	store        domain.AdminAuthStore
	cfg          *AuthConfig
	cache        *storage.Cache
	email        email.Email
}

func (s *AdminAuthService) validateStaffAccount(ctx echo.Context, creds domain.LoginCredentials) (*domain.Staff, int, error) {
	if err := validate.Struct(creds); err != nil {
		return nil, 0, xerror.EInvalidInput(err)
	}

	staff, err := s.staffService.GetByEmail(ctx, creds.Email)
	if err != nil {
		if xerror.IsNotFoundError(err) {
			return nil, 0, xerror.E(xerror.ErrUnauthorized).SetErrorCode(xerror.ErrInvalidCredentials).
				SetStatusCode(xerror.ErrCodeUnauthorized)
		}
		return nil, 0, xerror.E(err)
	}

	if staff.Password == "" {
		return nil, 0, xerror.EInvalidInput(errors.New("please setup new password before login")).SetErrorCode(xerror.ErrAuthAdminLoginMustSetPassword)
	}

	loginAttempt, err := s.cache.GetStrikes(ctx.Request().Context(), getAdminLoginAttemptCountKey(staff.ID.String()))
	if err != nil {
		return nil, 0, xerror.E(err)
	}
	if loginAttempt >= s.cfg.AccountLockoutMaxAttempts {
		return nil, 0, xerror.EInvalidInput(errors.New("login attempt reach limit")).SetErrorCode(xerror.ErrAuthAdminLoginReachLimit)
	}

	return staff, loginAttempt, nil
}

func getAdminLoginAttemptCountKey(userID string) string {
	return "admin:login_attempt_count:" + userID
}

func getAdminLoginKey(email string) string {
	return "admin:login:" + email
}

func getAdminResendOTPKey(email string) string {
	return "admin:resend:" + email
}

func getAdminResendOTPLockedKey(email string) string {
	return "admin:resend_otp_locked:" + email
}

func getAdminVerifyOTPKey(ref string) string {
	return "admin:verify:" + ref
}

func getAdminVerifyOTPLockedKey(email string) string {
	return "admin:verify_otp_locked:" + email
}

func getAdminLoggedInKey(email string) string {
	return "admin:logged_in:" + email
}
