package domain

import (
	"go_base/logger"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ContextKey represents a context key.
type ContextKey string

// Context keys
var (
	CorrelationIDKey = ContextKey("correlation_id")
	UserIDKey        = ContextKey("user_id")
	StaffKey         = ContextKey("staff")
	UserKey          = ContextKey("user")

	IsUserKey = ContextKey("is_user")
)

func UserID(ctx echo.Context) string {
	userID, ok := ctx.Get(string(UserIDKey)).(string)
	if !ok {
		return ""
	}
	_, err := uuid.Parse(userID)
	if err != nil {
		return ""
	}
	return userID
}

func StaffFromContext(ctx echo.Context) *Staff {
	staffCtx := ctx.Get(StaffCtx)
	if staffCtx == nil {
		return nil
	}
	staff, ok := staffCtx.(*Staff)
	if !ok {
		return nil
	}
	return staff
}

func UserFromContext(ctx echo.Context) *User {
	userCtx := ctx.Get(UserCtx)
	if userCtx == nil {
		return nil
	}
	user, ok := userCtx.(*User)
	if !ok {
		return nil
	}
	return user
}

func GetActionFromContext(ctx echo.Context) string {
	path := ctx.Path()
	if path == "" {
		return ""
	}
	paths := strings.Split(path, "/")
	// remove 2 first element
	paths = paths[3:]
	action := ctx.Request().Method + "|" + strings.Join(paths, "|")
	return action
}

func ErrLogGlsGo(ctx echo.Context, err error) {
	if err != nil {
		logger.L().Error(err)
	}
}

func GetUUIDFromParam(ctx echo.Context, key string) (string, uuid.UUID) {
	id := ctx.Param(key)
	isUid, err := uuid.Parse(id)
	if err != nil {
		return id, uuid.Nil
	}
	return isUid.String(), isUid
}

func GetUUID(idStr string) (string, uuid.UUID) {
	id := idStr
	isUid, err := uuid.Parse(id)
	if err != nil {
		return "", uuid.Nil
	}
	return isUid.String(), isUid
}
