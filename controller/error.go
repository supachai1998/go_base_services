package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go_base/domain"
	"go_base/logger"
	"go_base/xerror"

	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ErrorResponse represents an API error response.
type ErrorResponse struct {
	Error         string         `json:"error"`
	Code          string         `json:"code"`
	CorrelationID string         `json:"correlation_id"`
	Data          map[string]any `json:"data"`
}

// ErrorHandler is a centralized error handler.
// It converts an error returned from all handlers to error response format and logs unhandled error.
func ErrorHandler(isDebug bool) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		log := logger.Ctx(c.Request().Context())
		var (
			xerr *xerror.Xerror
			ok   bool
		)

		statusCode := http.StatusInternalServerError
		if xerr, ok = lo.ErrorsAs[*xerror.Xerror](err); ok {
			statusCode = xErrorCodes[xerr.StatusCode]
		} else if domainErr, ok := lo.ErrorsAs[xerror.DomainError](err); ok {
			// For backward compatibility
			statusCode = errorCodes[domainErr.Code()]
			xerr = xerror.E(err).SetErrorCode(domainErr.Code())
		} else if echoErr, ok := lo.ErrorsAs[*echo.HTTPError](err); ok {
			// Bind error (eg. invalid json input format)
			// Always return invalid input error
			xerr = xerror.E(echoErr).
				SetStatusCode(xerror.ErrCodeInvalidInput).
				SetMessage("invalid input format")

			statusCode = errorCodes[domain.ErrorCodeInvalidInput]
		} else {
			// Unexpected error
			xerr = xerror.E(err)
			for _, invalidInputCase := range xerror.ErrInvalidInputs {
				if strings.Contains(err.Error(), invalidInputCase) {
					xerr = xerror.E(fmt.Errorf(domain.ErrorCodeInvalidInput)).SetErrorCode(domain.ErrorCodeInvalidInput).SetDebugInfo("err", err.Error())
					statusCode = errorCodes[domain.ErrorCodeInvalidInput]
					break
				}
			}
			for _, invalidInputCase := range xerror.ErrNotFound {
				if strings.Contains(err.Error(), invalidInputCase) {
					xerr = xerror.E(fmt.Errorf(domain.ErrorCodeNotFound)).SetErrorCode(domain.ErrorCodeNotFound).SetDebugInfo("err", err.Error())
					statusCode = errorCodes[domain.ErrorCodeNotFound]
					break
				}
			}
			for _, invalidConflictsCase := range xerror.ErrConflicts {
				if strings.Contains(err.Error(), invalidConflictsCase) {
					xerr = xerror.E(fmt.Errorf(domain.ErrorCodeConflict)).SetErrorCode(domain.ErrorCodeConflict).SetDebugInfo("err", err.Error())
					statusCode = errorCodes[domain.ErrorCodeConflict]
					break
				}
			}
			// record not found
			if errors.Is(err, gorm.ErrRecordNotFound) {
				xerr = xerror.E(fmt.Errorf(domain.ErrorCodeNotFound)).SetErrorCode(domain.ErrorCodeNotFound).SetDebugInfo("err", err.Error())
				statusCode = errorCodes[domain.ErrorCodeNotFound]
			}

		}

		if isDebug {
			log = log.With("debug_info", xerr.DebugInfo)
		}

		log.Desugar().With(
			zap.Int("status_code", statusCode),
			zap.Inline(xerr),
		).Error(xerr.Error())

		var errRes ErrorResponse
		errRes.Code = xerr.ErrCode
		errRes.CorrelationID, _ = c.Request().Context().Value(domain.CorrelationIDKey).(string)
		errRes.Data = xerr.ExtraInfo

		if xerr.ErrCode == xerror.ErrCodeInternalServerError {
			// hide server related information
			errRes.Error = "internal server error"
		} else {
			errRes.Error = xerr.Error()
		}

		c.JSON(statusCode, errRes)
	}
}

var errorCodes = map[string]int{
	domain.ErrorCodeConflict:       http.StatusConflict,
	domain.ErrorCodeInternalError:  http.StatusInternalServerError,
	domain.ErrorCodeInvalidInput:   http.StatusBadRequest,
	domain.ErrorCodeNotFound:       http.StatusNotFound,
	domain.ErrorCodeNotImplemented: http.StatusNotImplemented,
	domain.ErrorCodeUnauthorized:   http.StatusUnauthorized,
	domain.ErrorCodeForbidden:      http.StatusForbidden,
}

var xErrorCodes = map[string]int{
	xerror.ErrCodeUnauthorized:        http.StatusUnauthorized,
	xerror.ErrCodeForbidden:           http.StatusForbidden,
	xerror.ErrCodeInvalidInput:        http.StatusBadRequest,
	xerror.ErrCodeNotFound:            http.StatusNotFound,
	xerror.ErrCodeConflict:            http.StatusConflict,
	xerror.ErrCodeInternalServerError: http.StatusInternalServerError,
	xerror.ErrCodeNotImplemented:      http.StatusNotImplemented,
}
