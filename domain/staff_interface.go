package domain

import "github.com/labstack/echo/v4"

type StaffService interface {
	Get(ctx echo.Context, id string) (*Staff, error)
	Find(ctx echo.Context, pagination Pagination[Staff]) (*Pagination[Staff], error)
	Create(ctx echo.Context, staff StaffCreate) (*Staff, error)
	LoginWithEmailPassword(ctx echo.Context, login StaffLogin) (*AuthResult, error)
	Unlock(ctx echo.Context, email StaffUnlock) error
	Delete(ctx echo.Context) error
	GetByEmail(ctx echo.Context, email SensitiveString) (*Staff, error)
	Update(ctx echo.Context, staffUpdate StaffUpdate) (*StaffUpdate, error)
	Verify(ctx echo.Context, staff StaffVerifyToken) (*StaffVerifyTokenResponse, error)
	GetToken(ctx echo.Context, staff StaffGetToken) (*StaffGetTokenResponse, error)
	UpdatePassword(ctx echo.Context, staff StaffUpdatePassword) error
	GetLog(ctx echo.Context, staff StaffGetLog) (*Pagination[*Logs[Staff]], error)

	GetMe(ctx echo.Context) (*StaffMe, error)

	// log
	GetLogMe(ctx echo.Context) (*Pagination[*Logs[Staff]], error)
}
