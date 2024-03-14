package domain

import (
	"fmt"
	helper "go_base/domain/helper"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	DoerTypeStaff  = "staff"
	DoerTypeUser   = "user"
	DoerTypeSystem = "system"
)

var removeWordTableName = map[string]string{
	"Update": "",
	"Delete": "",
	"Create": "",
}

type Logs[T any] struct {
	BaseModel
	Model     datatypes.JSON `json:"model" gorm:"type:jsonb;not null"`
	Action    string         `json:"action" gorm:"type:varchar(255);not null"`
	FromTable *string        `json:"from_table" gorm:"type:varchar(255);"`
	// doer
	// ex: {"id":1,"name":"admin","email":"admin@localhost", type:"staff"}
	Doer datatypes.JSON `json:"doer" gorm:"type:jsonb;not null"`

	LogModel T `json:"-" gorm:"-"`
}

type Doer struct {
	ID     uuid.UUID  `json:"id"`
	Name   string     `json:"name"`
	Email  string     `json:"email"`
	Type   string     `json:"type"`
	RoleID *uuid.UUID `json:"role_id,omitempty"`
	Role   *Role      `json:"role,omitempty"`
}

func NewLogs[T any]() *Logs[T] {
	return &Logs[T]{}
}

func (l *Logs[T]) TableName() string {
	fieldName := reflect.TypeOf(l.LogModel).Name()
	for k, v := range removeWordTableName {
		fieldName = strings.Replace(fieldName, k, v, -1)
	}
	return fmt.Sprintf("%s_logs", helper.ToSnakeCase(fieldName))
}

// Find is pagination for logs
func (l *Logs[T]) Find(ctx echo.Context, db *gorm.DB) (*Pagination[Logs[T]], error) {

	db = db.Where("model @> ?", l.Model)
	pg := PaginationFromCtx[Logs[T]](ctx)

	result, err := pg.Paginate(ctx, db)
	if err != nil {
		return nil, err
	}
	return result, nil
}
