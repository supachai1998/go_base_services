package domain

import "github.com/labstack/echo/v4"

type UserService interface {
	Get(ctx echo.Context, id string) (*User, error)
	Find(ctx echo.Context, pagination Pagination[User]) (*Pagination[User], error)
	Create(ctx echo.Context, user UserCreate) (*User, error)
	LoginWithEmailPassword(ctx echo.Context, login UserLogin) (*AuthResult, error)
	Unlock(ctx echo.Context, email UserUnlock) error
	Delete(ctx echo.Context) error
	GetByEmail(ctx echo.Context, email SensitiveString) (*User, error)
	Update(ctx echo.Context, userUpdate UserUpdate) (*UserUpdate, error)
	Verify(ctx echo.Context, user UserVerifyToken) (*UserVerifyTokenResponse, error)
	GetToken(ctx echo.Context, user UserGetToken) (*UserGetTokenResponse, error)
	UpdatePassword(ctx echo.Context, user UserUpdatePassword) error
	GetLog(ctx echo.Context, user UserGetLog) (*Pagination[*Logs[User]], error)

	GetMe(ctx echo.Context) (*UserMe, error)
	UpdateMe(ctx echo.Context, user UserUpdate) error

	// log
	GetLogMe(ctx echo.Context) (*Pagination[*Logs[User]], error)

	DeleteByIds(ctx echo.Context, ids Ids) error
}
