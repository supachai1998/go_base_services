package domain

import (
	"errors"
	"strings"

	"go_base/xerror"
)

// App error codes.
// Deprecated: use code from xerror package instead.
const (
	ErrorCodeConflict       = xerror.ErrCodeConflict
	ErrorCodeInternalError  = xerror.ErrCodeInternalServerError
	ErrorCodeInvalidInput   = xerror.ErrCodeInvalidInput
	ErrorCodeNotFound       = xerror.ErrCodeNotFound
	ErrorCodeNotImplemented = xerror.ErrCodeNotImplemented
	ErrorCodeUnauthorized   = xerror.ErrCodeUnauthorized
	ErrorCodeForbidden      = xerror.ErrCodeForbidden
)

type BaseError struct {
	code    string
	message string
}

func NewError(code, message string) *BaseError {
	return &BaseError{
		code:    code,
		message: message,
	}
}

func (e *BaseError) Code() string {
	return e.code
}

func (e *BaseError) Message() string {
	return e.message
}

func (e *BaseError) Error() string {
	return e.message
}

// FieldError contains error on single field validation.
type FieldError struct {
	Field string
	Err   string
}

// FieldErrors contains errors on multiple fields (struct) validation.
type FieldErrors []*FieldError

func (errs FieldErrors) Code() string {
	return ErrorCodeInvalidInput
}

func (errs FieldErrors) Message() string {
	return "invalid input"
}

func (errs FieldErrors) Error() string {
	var s strings.Builder

	for i, err := range errs {
		if i > 0 {
			s.WriteString("; ")
		}
		s.WriteString(err.Err)
	}

	return s.String()
}

// InvalidIDFormatError represents error when id format is invalid.
type InvalidIDFormatError struct {
	*BaseError
	ID string
}

// NewInvalidIDFormatError returns new *InvalidIDFormatError.
func NewInvalidIDFormatError(id string) *InvalidIDFormatError {
	return &InvalidIDFormatError{
		BaseError: NewError(ErrorCodeInvalidInput, "invalid id format: "+id),
		ID:        id,
	}
}

// NotFoundError represents resource not found error.
type NotFoundError struct {
	*BaseError
	Resource string
}

func (e *NotFoundError) Is(err error) bool {
	var target *NotFoundError
	if !errors.As(err, &target) {
		return false
	}

	return e.Resource == target.Resource
}

// NewInvalidInputError returns *BaseError with the given error message.
func NewInvalidInputError(message string) *BaseError {
	return &BaseError{
		code:    ErrorCodeInvalidInput,
		message: message,
	}
}

// NewNotFoundError returns *NotFoundError with the given resource name.
func NewNotFoundError(resource string) *NotFoundError {
	return &NotFoundError{
		BaseError: NewError(ErrorCodeNotFound, resource+" not found"),
		Resource:  resource,
	}
}

// PermissionDeniedError represents error when current staff/user have no permission to service
type PermissionDeniedError struct {
	*BaseError
}

func NewPermissionDeniedError() *PermissionDeniedError {
	return &PermissionDeniedError{
		BaseError: NewError(ErrorCodeForbidden, "permission denied"),
	}
}

var ErrPermissionDenied = NewPermissionDeniedError()
