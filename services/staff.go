package services

import (
	"errors"
	"fmt"
	"go_base/database"
	"go_base/domain"
	"go_base/hash"
	"go_base/services/auth"
	"go_base/storage"
	"go_base/xerror"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type StaffService struct {
	store      *database.Store
	services   *domain.AllServices
	staffStore database.StaffStore
	cache      *storage.Cache
	cfg        *auth.AuthConfig
}

func NewStaffService(store *database.Store, staff *database.StaffStore, services *domain.AllServices, cache *storage.Cache, cfg *auth.AuthConfig) *StaffService {
	return &StaffService{store: store, services: services, staffStore: *staff, cache: cache, cfg: cfg}
}

// GET /staff/:id
func (s *StaffService) Get(ctx echo.Context, id string) (*domain.Staff, error) {
	return s.staffStore.GetByID(ctx, id)
}

// GET /staffs
func (s *StaffService) Find(ctx echo.Context, pagination domain.Pagination[domain.Staff]) (*domain.Pagination[domain.Staff], error) {
	pg, err := s.staffStore.Find(ctx, pagination, "Roles")
	if err != nil {
		return nil, err
	}

	// attach meta data with total_all , group by role
	countRole, err := s.staffStore.CountJsonGroupRole(ctx)
	if err != nil {
		return nil, err
	}

	pg.MetaCount = map[string]interface{}{
		"role": countRole,
	}
	return pg, nil
}

// POST /staffs
func (s *StaffService) Create(ctx echo.Context, staffCreate domain.StaffCreate) (*domain.Staff, error) {
	verifyToken := hash.GenerateToken()
	claimsVerifyToken, err := domain.GenerateVerifyToken(verifyToken, s.cfg.JWTSecret, s.cfg.VerifyTokenDuration)
	if err != nil {
		return nil, err
	}
	ran := hash.GenerateRandomString(s.cfg.LenTempPwd)

	staff := domain.Staff{
		Email:       domain.SensitiveString(staffCreate.Email),
		FirstName:   staffCreate.FirstName,
		LastName:    staffCreate.LastName,
		Phone:       staffCreate.Phone,
		VerifyToken: claimsVerifyToken,
		Password:    domain.Password(ran),
		TmpPassword: ran,
	}

	if staffCreate.RoleID != nil {
		staff.RoleID = lo.ToPtr(uuid.MustParse(*staffCreate.RoleID))
	}

	if err := s.staffStore.Create(ctx, &staff); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if err := s.staffStore.DeleteExistData(ctx, &staff); err != nil {
				return nil, err
			} else {
				if err := s.staffStore.Create(ctx, &staff); err != nil {
					return nil, err
				}
				return nil, nil
			}
		}
		return nil, err
	}

	// update role default

	return &staff, nil
}

// POST /staff/login
func (s *StaffService) LoginWithEmailPassword(ctx echo.Context, login domain.StaffLogin) (*domain.AuthResult, error) {
	if i, e := s.cache.GetStrikes(ctx.Request().Context(), fmt.Sprintf(domain.StaffAuthCache, login.Email)); e == nil {
		if i >= s.cfg.AccountLockoutMaxAttempts {
			s.staffStore.WriteLog(ctx, &domain.Staff{Email: domain.SensitiveString(login.Email)}, string(xerror.ErrCodeTooManyLoginAttempts))
			return nil, xerror.EForbidden().SetExtraInfo("reason", xerror.ErrCodeTooManyLoginAttempts)
		}
	}
	staff, err := s.staffStore.GetByEmail(ctx, login.Email)
	if err != nil {
		return s.staffIncrementStrike(ctx, login)
	}
	password := string(staff.Password)
	if !staff.IsVerified {
		return s.staffIncrementStrike(ctx, login)
	}
	if ok := hash.CompareBcrypt(password, login.Password); !ok {
		return s.staffIncrementStrike(ctx, login)
	}
	// login success
	if err := s.staffStore.Update(ctx, &domain.Staff{
		BaseModel: domain.BaseModel{ID: staff.ID},
		LastLogin: domain.TimeNowPtr(),
	}, domain.LoginLog); err != nil {
		return nil, err
	}
	if err := s.cache.DeleteStrikes(ctx.Request().Context(), fmt.Sprintf(domain.StaffAuthCache, login.Email)); err != nil {
		return nil, err
	}

	auth, err := auth.IssueAccessRefreshToken(ctx, staff.ID, s.cfg, s.services.AuthAdmin.FindAuth, s.services.AuthAdmin.UpdateAuth, s.services.AuthAdmin.CreateAuth, s.cache.SetCache)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// POST /staff/unlock
func (s *StaffService) Unlock(ctx echo.Context, unlock domain.StaffUnlock) error {

	err := s.cache.DeleteStrikes(ctx.Request().Context(), fmt.Sprintf(domain.StaffAuthCache, unlock.Email))
	if err != nil {
		return err
	}
	// write log
	if staff, err := s.GetByEmail(ctx, domain.SensitiveString(unlock.Email)); err == nil {
		if err := s.staffStore.WriteLog(ctx, staff, domain.GetActionFromContext(ctx)); err != nil {
			return err
		}
	} else {
		if err := s.staffStore.WriteLog(ctx, &domain.Staff{
			Email: domain.SensitiveString(unlock.Email),
		}, domain.UnlockLog); err != nil {
			return err
		}
	}
	return nil
}

// DELETE /staff/:id
func (s *StaffService) Delete(ctx echo.Context) error {
	_, idUid := domain.GetUUIDFromParam(ctx, "id")
	if idUid == uuid.Nil {
		return xerror.EInvalidParameter(nil)
	}

	if err := s.staffStore.Delete(ctx, idUid); err != nil {
		return err
	}

	return nil
}

// GET By Email
func (s *StaffService) GetByEmail(ctx echo.Context, email domain.SensitiveString) (*domain.Staff, error) {
	return s.staffStore.GetByEmail(ctx, email)
}

// UPDATE /staff/:id
func (s *StaffService) Update(ctx echo.Context, staffUpdate domain.StaffUpdate) (*domain.StaffUpdate, error) {
	if err := s.staffStore.UpdateU(ctx, &staffUpdate); err != nil {
		return nil, err
	}
	return &staffUpdate, nil
}

// Get log me delete /staff/log
func (s *StaffService) GetLogMe(ctx echo.Context) (*domain.Pagination[*domain.Logs[domain.Staff]], error) {
	staff := domain.StaffFromContext(ctx)
	if staff == nil {
		return nil, xerror.EUnAuthorized()
	}

	return s.GetLog(ctx, domain.StaffGetLog{ID: staff.ID})
}

// Get log delete /staff/log/:id
func (s *StaffService) GetLog(ctx echo.Context, staff domain.StaffGetLog) (*domain.Pagination[*domain.Logs[domain.Staff]], error) {
	var stafflog *domain.Logs[domain.Staff]
	result, err := s.staffStore.LogsQueryModel(ctx, stafflog, "doer.id", staff.ID.String())
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Mock IsVerified /staff/verify
func (s *StaffService) Verify(ctx echo.Context, staff domain.StaffVerifyToken) (*domain.StaffVerifyTokenResponse, error) {
	if _staff, err := s.staffStore.GetByKey(ctx, string(domain.StaffVerifyTokenType), staff.Token); err == nil {
		if err := domain.ParseVerifyToken(staff.Token); err == nil {

			if err := s.staffStore.UpdateTokenVerify(ctx, _staff.ID); err != nil {
				return nil, err
			}
			return &domain.StaffVerifyTokenResponse{
				Email: _staff.Email.String(),
			}, nil
		}
	}
	return nil, xerror.EInvalidInputOk()

}

// Mock Get token /staff/token
func (s *StaffService) GetToken(ctx echo.Context, staff domain.StaffGetToken) (*domain.StaffGetTokenResponse, error) {
	if _staff, err := s.staffStore.GetByKey(ctx, "email", staff.Email); err == nil {
		if _staff.IsVerified {
			return nil, xerror.EInvalidInputOk()
		}
		if ok := hash.CompareBcrypt(string(_staff.Password), staff.TmpPassword); ok {
			return &domain.StaffGetTokenResponse{
				Token: _staff.VerifyToken,
			}, nil
		}
	}
	return nil, xerror.EInvalidInputOk()
}

// Update Password /staff/password
func (s *StaffService) UpdatePassword(ctx echo.Context, staff domain.StaffUpdatePassword) error {
	if _staff, err := s.staffStore.GetByEmail(ctx, staff.Email); err == nil {
		if _staff.IsVerified {
			if ok := hash.CompareBcrypt(string(_staff.Password), staff.OldPassword); ok {
				if err := s.staffStore.Update(ctx, &domain.Staff{
					BaseModel: domain.BaseModel{ID: _staff.ID},
					Password:  domain.Password(staff.Password).Hash(),
				}, domain.ChangePasswordLog); err != nil {
					return err
				}
				return nil
			}
		}
	}
	s.staffStore.WriteLog(ctx, &domain.Staff{Email: domain.SensitiveString(staff.Email)}, domain.ChangePasswordFailedLog)
	return xerror.EInvalidInputOk()
}

func (s *StaffService) staffIncrementStrike(ctx echo.Context, login domain.StaffLogin) (*domain.AuthResult, error) {
	s.staffStore.WriteLog(ctx, &domain.Staff{Email: domain.SensitiveString(login.Email)}, domain.LoginFail)
	increaseStrike, e := s.cache.IncreaseStrike(ctx.Request().Context(), fmt.Sprintf(domain.StaffAuthCache, login.Email))
	if e != nil {
		return nil, e
	}
	return nil, xerror.EInvalidInput(nil).SetExtraInfo("counter", increaseStrike)
}

// GetMe /staff/me
func (s *StaffService) GetMe(ctx echo.Context) (*domain.StaffMe, error) {
	staff, err := s.staffStore.GetMe(ctx)
	if err != nil {
		return nil, err
	}
	return staff, nil
}
