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
	"gorm.io/gorm"
)

type UserService struct {
	store     *database.Store
	services  *domain.AllServices
	userStore database.UserStore
	cache     *storage.Cache
	cfg       *auth.AuthConfig
}

func NewUserService(store *database.Store, user *database.UserStore, services *domain.AllServices, cache *storage.Cache, cfg *auth.AuthConfig) *UserService {
	return &UserService{store: store, services: services, userStore: *user, cache: cache, cfg: cfg}
}

// GET /users/:id
func (s *UserService) Get(ctx echo.Context, id string) (*domain.User, error) {
	return s.userStore.GetByID(ctx, id)
}

// GET /users
func (s *UserService) Find(ctx echo.Context, pagination domain.Pagination[domain.User]) (*domain.Pagination[domain.User], error) {
	pg, err := s.userStore.Find(ctx, pagination)
	if err != nil {
		return nil, err
	}
	metaCountType, err := s.userStore.CountJsonGroup(ctx, "type")
	if err != nil {
		return nil, err
	}
	metaCountInterest, err := s.userStore.CountJsonGroup(ctx, "interest")
	if err != nil {
		return nil, err
	}
	metaCountTag, err := s.userStore.CountJsonGroup(ctx, "tag")
	if err != nil {
		return nil, err
	}

	/*
		pg.MetaCount  --> { "type" : {"admin": 1, "user": 2} }
	*/
	pg.MetaCount = map[string]interface{}{
		"type":     metaCountType,
		"interest": metaCountInterest,
		"tag":      metaCountTag,
	}

	return pg, nil
}

// POST /users
func (s *UserService) Create(ctx echo.Context, userCreate domain.UserCreate) (*domain.User, error) {
	verifyToken := hash.GenerateToken()
	claimsVerifyToken, err := domain.GenerateVerifyToken(verifyToken, s.cfg.JWTSecret, s.cfg.VerifyTokenDuration)
	if err != nil {
		return nil, err
	}
	ran := hash.GenerateRandomString(s.cfg.LenTempPwd)

	user := domain.User{
		Email:       domain.SensitiveString(userCreate.Email),
		FirstName:   userCreate.FirstName,
		LastName:    userCreate.LastName,
		VerifyToken: claimsVerifyToken,
		Password:    domain.Password(ran),
		TmpPassword: ran,
	}

	if err := s.userStore.Create(ctx, &user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if err := s.userStore.DeleteExistData(ctx, &user); err != nil {
				return nil, err
			} else {
				if err := s.userStore.Create(ctx, &user); err != nil {
					return nil, err
				}
				return nil, nil
			}
		}
		return nil, err
	}

	// update role default

	return &user, nil
}

// POST /users/login
func (s *UserService) LoginWithEmailPassword(ctx echo.Context, login domain.UserLogin) (*domain.AuthResult, error) {
	s.cache.DeleteStrikes(ctx.Request().Context(), fmt.Sprintf(domain.UserAuthCache, login.Email))
	if i, e := s.cache.GetStrikes(ctx.Request().Context(), fmt.Sprintf(domain.UserAuthCache, login.Email)); e == nil {
		if i >= s.cfg.AccountLockoutMaxAttempts {
			s.userStore.WriteLog(ctx, &domain.User{Email: domain.SensitiveString(login.Email)}, string(xerror.ErrCodeTooManyLoginAttempts))
			return nil, xerror.EForbidden().SetExtraInfo("reason", xerror.ErrCodeTooManyLoginAttempts)
		}
	}
	user, err := s.userStore.GetByKey(ctx, "email", login.Email.String())
	if err != nil {
		return s.userIncrementStrike(ctx, login)
	}
	password := string(user.Password)
	if !user.IsVerified {
		return s.userIncrementStrike(ctx, login)
	}
	if ok := hash.CompareBcrypt(password, login.Password); !ok {
		return s.userIncrementStrike(ctx, login)
	}
	// login success
	if err := s.userStore.Update(ctx, &domain.User{
		BaseModel: domain.BaseModel{ID: user.ID},
		LastLogin: domain.TimeNowPtr(),
	}, domain.LoginLog); err != nil {
		return nil, err
	}
	if err := s.cache.DeleteStrikes(ctx.Request().Context(), fmt.Sprintf(domain.UserAuthCache, login.Email)); err != nil {
		return nil, err
	}
	auth, err := auth.IssueAccessRefreshToken(ctx, user.ID, s.cfg, s.services.AuthUser.FindAuth, s.services.AuthUser.UpdateAuth, s.services.AuthUser.CreateAuth, s.cache.SetCache)
	if err != nil {
		return nil, err
	}
	return auth, nil
}

// POST /users/unlock
func (s *UserService) Unlock(ctx echo.Context, unlock domain.UserUnlock) error {

	err := s.cache.DeleteStrikes(ctx.Request().Context(), fmt.Sprintf(domain.UserAuthCache, unlock.Email))
	if err != nil {
		return err
	}
	// write log
	if user, err := s.GetByEmail(ctx, domain.SensitiveString(unlock.Email)); err == nil {
		if err := s.userStore.WriteLog(ctx, user, domain.GetActionFromContext(ctx)); err != nil {
			return err
		}
	} else {
		if err := s.userStore.WriteLog(ctx, &domain.User{
			Email: domain.SensitiveString(unlock.Email),
		}, domain.UnlockLog); err != nil {
			return err
		}
	}
	return nil
}

// DELETE /users/:id
func (s *UserService) Delete(ctx echo.Context) error {
	_, idUid := domain.GetUUIDFromParam(ctx, "id")
	if idUid == uuid.Nil {
		return xerror.EInvalidParameter(nil)
	}

	if err := s.userStore.Delete(ctx, idUid); err != nil {
		return err
	}

	return nil
}

// GET By Email
func (s *UserService) GetByEmail(ctx echo.Context, email domain.SensitiveString) (*domain.User, error) {
	return s.userStore.GetByKey(ctx, "email", email.String())
}

// UPDATE /users/:id
func (s *UserService) Update(ctx echo.Context, userUpdate domain.UserUpdate) (*domain.UserUpdate, error) {
	_, uid := domain.GetUUIDFromParam(ctx, "id")
	if uid == uuid.Nil {
		return nil, xerror.EInvalidParameter(nil)
	}
	userUpdate.ID = uid

	if err := s.userStore.UpdateU(ctx, &userUpdate); err != nil {
		return nil, err
	}
	return &userUpdate, nil
}

// Get log me delete /users/log
func (s *UserService) GetLogMe(ctx echo.Context) (*domain.Pagination[*domain.Logs[domain.User]], error) {
	user := domain.UserFromContext(ctx)
	if user == nil {
		return nil, xerror.EUnAuthorized()
	}

	return s.GetLog(ctx, domain.UserGetLog{ID: user.ID})
}

// Get log delete /users/log/:id
func (s *UserService) GetLog(ctx echo.Context, user domain.UserGetLog) (*domain.Pagination[*domain.Logs[domain.User]], error) {
	var userlog *domain.Logs[domain.User]
	result, err := s.userStore.LogsQueryModel(ctx, userlog, "doer.id", user.ID.String())
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Mock IsVerified /users/verify
func (s *UserService) Verify(ctx echo.Context, user domain.UserVerifyToken) (*domain.UserVerifyTokenResponse, error) {
	if _user, err := s.userStore.GetByKey(ctx, string(domain.UserVerifyTokenType), user.Token); err == nil {
		if err := domain.ParseVerifyToken(user.Token); err == nil {

			if err := s.userStore.UpdateTokenVerify(ctx, _user.ID); err != nil {
				return nil, err
			}
			// if err := s.userStore.Update(ctx, &domain.User{
			// 	BaseModel:   domain.BaseModel{ID: _user.ID},
			// 	IsVerified:  true,
			// 	VerifyToken: "",
			// }); err != nil {
			// 	return nil, err
			// }
			return &domain.UserVerifyTokenResponse{
				Email: _user.Email.String(),
			}, nil
		}
	}
	return nil, xerror.EInvalidInputOk()

}

// Mock Get token /users/token
func (s *UserService) GetToken(ctx echo.Context, user domain.UserGetToken) (*domain.UserGetTokenResponse, error) {
	if _user, err := s.userStore.GetByKey(ctx, "email", user.Email); err == nil {
		if _user.IsVerified {
			return nil, xerror.EInvalidInputOk()
		}
		if ok := hash.CompareBcrypt(string(_user.Password), user.TmpPassword); ok {
			return &domain.UserGetTokenResponse{
				Token: _user.VerifyToken,
			}, nil
		}
	}
	return nil, xerror.EInvalidInputOk()
}

// Update Password /users/password
func (s *UserService) UpdatePassword(ctx echo.Context, user domain.UserUpdatePassword) error {
	if _user, err := s.userStore.GetByKey(ctx, "email", user.Email.String()); err == nil {
		if _user.IsVerified {
			if ok := hash.CompareBcrypt(string(_user.Password), user.OldPassword); ok {
				if err := s.userStore.Update(ctx, &domain.User{
					BaseModel: domain.BaseModel{ID: _user.ID},
					Password:  domain.Password(user.Password).Hash(),
				}, domain.ChangePasswordLog); err != nil {
					return err
				}
				return nil
			}
		}
	}
	s.userStore.WriteLog(ctx, &domain.User{Email: domain.SensitiveString(user.Email)}, domain.ChangePasswordFailedLog)
	return xerror.EInvalidInputOk()
}

func (s *UserService) userIncrementStrike(ctx echo.Context, login domain.UserLogin) (*domain.AuthResult, error) {
	s.userStore.WriteLog(ctx, &domain.User{Email: domain.SensitiveString(login.Email)}, domain.LoginFail)
	increaseStrike, e := s.cache.IncreaseStrike(ctx.Request().Context(), fmt.Sprintf(domain.UserAuthCache, login.Email))
	if e != nil {
		return nil, e
	}
	return nil, xerror.EInvalidInput(nil).SetExtraInfo("counter", increaseStrike)
}

// GetMe /users/me
func (s *UserService) GetMe(ctx echo.Context) (*domain.UserMe, error) {
	user, err := s.userStore.GetMe(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update Me /users/me
func (s *UserService) UpdateMe(ctx echo.Context, user domain.UserUpdate) error {
	if err := s.userStore.UpdateU(ctx, &user); err != nil {
		return err
	}
	return nil
}

// Delete /users/ids
func (s *UserService) DeleteByIds(ctx echo.Context, ids domain.Ids) error {
	if err := s.userStore.DeleteIds(ctx, ids); err != nil {
		return err
	}
	return nil
}
