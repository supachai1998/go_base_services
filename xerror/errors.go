package xerror

import "errors"

var ErrUnauthorized = errors.New("unauthorized")
var ErrPermissionDenied = errors.New("permission denied")
var ErrForbidden = errors.New("forbidden")
var ErrInvalidIDFormat = errors.New("invalid id format")
var ErrEndpointNotFound = errors.New("endpoint not found")
var ErrItemNotFound = errors.New("record not found")
var ErrInputValidationFailed = errors.New("input validation failed")
var ErrInvalidInput = errors.New("invalid input")
var ErrRequired = errors.New("required")

var ErrBadInput = errors.New("bad input")
