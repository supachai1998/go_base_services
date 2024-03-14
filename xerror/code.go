package xerror

type Level string

const (
	// specific error code
	ErrInvalidCredentials            = "invalid_credentials"
	ErrAuthAdminLoginMustSetPassword = "auth_admin_login_must_set_password"
	ErrAuthAdminLoginReachLimit      = "auth_admin_login_reach_limit"
)

const (
	// common error code
	ErrCodeUnauthorized         = "unauthorized"
	ErrCodeForbidden            = "forbidden"
	ErrCodeInvalidInput         = "invalid_input"
	ErrCodeNotFound             = "not_found"
	ErrCodeConflict             = "conflict"
	ErrCodeInternalServerError  = "internal_server_error"
	ErrCodeNotImplemented       = "not_implemented"
	ErrCodeInternal             = "internal"
	ErrCodeTooManyLoginAttempts = "too_many_login_attempts"
)

const (
	ErrLevelWarn  Level = "warn"
	ErrLevelError Level = "error"
)

// GORMError is a struct that contains error message from gorm
var ErrInvalidInputs = []string{
	"SQLSTATE 42703",
	"SQLSTATE 42P01",
}
var ErrConflicts = []string{
	"SQLSTATE 23505",
}
var ErrNotFound = []string{
	"SQLSTATE 02000",
	"SQLSTATE 23503",
}
