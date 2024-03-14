package domain

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var TableNames []string = []string{
	"Role",
	"Staff",
	"Token_expire",
	"Auth",
	"User",
	"Asset",
	"Developer",
	"Project",
}

type BaseModel struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;primaryKey;index:,option:CONCURRENTLY"`
	CreatedAt time.Time       `json:"created_at" gorm:"default:now();autoCreateTime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"default:now();autoUpdateTime"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type BaseModelIDOnly struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;index:,option:CONCURRENTLY"`
}
type Ids struct {
	IDs []string `json:"ids" validate:"required,unique" form:"ids" query:"ids"`
}

func (m *BaseModel) IsZeroID() bool {
	return m.ID == uuid.Nil
}

func (m *BaseModel) IsDeleted() bool {
	return !m.DeletedAt.Time.IsZero()
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.IsZeroID() {
		m.ID = uuid.New()
	}
	if !m.IsZeroID() {
		if err := uuid.Validate(m.ID.String()); err != nil {
			return err
		}
	}
	return nil
}

func (m *BaseModelIDOnly) IsZeroID() bool {
	return m.ID == uuid.Nil
}

func (m *BaseModelIDOnly) BeforeCreate(tx *gorm.DB) error {
	if m.IsZeroID() {
		m.ID = uuid.New()
	}
	if !m.IsZeroID() {
		if err := uuid.Validate(m.ID.String()); err != nil {
			return err
		}
	}
	return nil
}

func (m *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = TimeNow()
	return nil
}

func ConvertAnyIntoBaseModel(modelAny any) BaseModel {
	// ptr
	if reflect.TypeOf(modelAny).Kind() == reflect.Ptr {
		modelAny = reflect.ValueOf(modelAny).Elem().Interface()
	}
	jsonMarshal, _ := json.Marshal(modelAny)
	var model BaseModel
	json.Unmarshal(jsonMarshal, &model)
	return model

}

// Scope
func WithUserID(userID uuid.UUID) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userID)
	}
}
