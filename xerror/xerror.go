package xerror

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
)

const (
	minimumCallerDepth = 2
	maximumCallerDepth = 15
)

var excludePkgName = map[string]bool{
	"logger": true,
	"xerror": true,
}

type Xerror struct {
	Message    string         `json:"message"`
	Err        error          `json:"error_raw,omitempty"`
	ErrCode    string         `json:"error_code"`
	StatusCode string         `json:"-"`
	SourceFunc string         `json:"source_func,omitempty"`
	SourceFile string         `json:"source_file,omitempty"`
	DebugInfo  map[string]any `json:"debug_info,omitempty"`
	ExtraInfo  map[string]any `json:"extra_info,omitempty"`
}

func (e Xerror) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("message", e.Message)
	enc.AddString("status", e.StatusCode)
	enc.AddString("code", e.ErrCode)
	enc.AddString("function", e.SourceFunc)
	enc.AddString("file", e.SourceFile)

	switch v := e.Err.(type) {
	case nil:
	case *Xerror:
		if err := enc.AddObject("error", v); err != nil {
			return err
		}
	default:
		enc.AddString("error", v.Error())
	}

	enc.OpenNamespace("extra_info")
	for k, v := range e.ExtraInfo {
		if err := enc.AddReflected(k, v); err != nil {
			return err
		}
	}

	return nil
}

// Error implements error interface.
func (e Xerror) Error() string {
	if e.Err != nil && e.Err.Error() == "" {
		return e.ErrCode
	}

	if e.Message == "" {
		return e.ErrCode
	}
	return e.Message
}

// SetMessage overrides error message with a given message
func (e *Xerror) SetMessage(format string, args ...any) *Xerror {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

// SetErr sets inner error to a given error
// If the inner error is *xerror.Xerror, it will copy error contexts (eg. message, error code) to the current one.
func (e *Xerror) SetErr(err error) *Xerror {
	if err != nil {
		if e.Message == "" {
			e.Message = err.Error()
		}

		// handle case not found
		if innerError, ok := err.(*Xerror); ok {
			e.ErrCode = innerError.ErrCode
			e.StatusCode = innerError.StatusCode
			e.DebugInfo = innerError.DebugInfo
			e.ExtraInfo = innerError.ExtraInfo
		} else if innerError, ok := err.(DomainError); ok {
			e.Message = innerError.Message()
			e.ErrCode = innerError.Code()
			e.StatusCode = innerError.Code()
		} else if _, ok := err.(*echo.HTTPError); ok {
			e.Message = ErrCodeInvalidInput
			e.ErrCode = ErrCodeInvalidInput
			e.StatusCode = ErrCodeInvalidInput
		}
	}
	e.Err = err
	return e
}

// SetDebugInfo sets debugging information to an error.
// This error will be logged, but not show to user.
func (e *Xerror) SetDebugInfo(key string, value any) *Xerror {
	e.DebugInfo[key] = value
	return e
}

// SetExtraInfo sets extra informations to an error.
// This will show to user.
func (e *Xerror) SetExtraInfo(key string, value any) *Xerror {
	e.ExtraInfo[key] = value
	return e
}

// SetErrorCode sets a specific error code that will show to user.
func (e *Xerror) SetErrorCode(code string) *Xerror {
	e.ErrCode = code
	return e
}

// SetStatusCode sets status code and also error code if it is empty.
func (e *Xerror) SetStatusCode(code string) *Xerror {
	if e.ErrCode == "" {
		e.ErrCode = code
	}
	e.StatusCode = code

	return e
}

func (e *Xerror) Unwrap() error {
	return e.Err
}

func buildError(err error) *Xerror {
	fn, file := getCaller()

	xerr := &Xerror{
		SourceFunc: fn,
		SourceFile: file,
		StatusCode: ErrCodeInternalServerError,
		ExtraInfo:  make(map[string]any),
		DebugInfo:  make(map[string]any),
	}

	return xerr.SetErr(err)
}

func New() *Xerror {
	return buildError(nil)
}

func E(err error) *Xerror {
	return buildError(err)
}

func EErrorCode(code string) *Xerror {
	return New().SetErrorCode(code)
}

func EStatusCode(code string) *Xerror {
	return New().SetStatusCode(code)
}

func EInvalidInput(err ...error) *Xerror {
	var e error = ErrInvalidInput
	if len(err) > 0 {
		e = err[0]
	}
	return E(e).SetStatusCode(ErrCodeInvalidInput)
}

func EInvalidInputOk() *Xerror {
	return E(nil).SetStatusCode(fmt.Sprint(http.StatusOK))
}

func ENotFound() *Xerror {
	return E(ErrItemNotFound).SetStatusCode(ErrCodeNotFound)
}

func ENotFoundResource(resource string) *Xerror {
	return EStatusCode(ErrCodeNotFound).SetMessage(resource + " not found")
}

func EInternal() *Xerror {
	return New().SetStatusCode(ErrCodeInternalServerError)
}

func EForbidden() *Xerror {
	return New().SetStatusCode(ErrCodeForbidden)
}
func EUnAuthorized() *Xerror {
	return New().SetStatusCode(ErrCodeUnauthorized)
}

func EInternalError() *Xerror {
	return EStatusCode(ErrCodeInternalServerError)
}

func ErrInvalidField(field string) *Xerror {
	return EStatusCode(ErrCodeInvalidInput).SetMessage("invalid field: %s", field)
}

func ErrInvalidOperator(operators map[string]bool) *Xerror {
	return EStatusCode(ErrCodeInvalidInput).SetMessage("invalid operator").SetExtraInfo("operator", operators)
}
func EInvalidInputField(field string) *Xerror {
	return EStatusCode(ErrCodeInvalidInput).SetMessage("invalid field: %s", field)
}

func EConflict(err error) *Xerror {
	return EStatusCode(ErrCodeConflict).SetErr(err)
}

func IsNotFoundError(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	if err, ok := err.(*Xerror); ok {
		return err.StatusCode == ErrCodeNotFound || err.Err == ErrItemNotFound
	}
	return false
}

func EInvalidParameter(err error) *Xerror {
	return E(err).SetStatusCode(ErrCodeInvalidInput)
}

type DomainError interface {
	Code() string
	Message() string
	Error() string
}

// Deprecated: Wrap wraps the given error with context message.
func Wrap(err error, message string) error {
	return E(err)
}

// Deprecated: Wrap wraps the given error with context message format.
func Wrapf(err error, format string, args ...any) error {
	return Wrap(err, fmt.Sprintf(format, args...))
}

func getCaller() (string, string) {
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if !excludePkgName[pkg] {
			return f.Function, fmt.Sprintf("%s:%d", f.File, f.Line)
		}
	}

	return "", ""
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
