package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ผู้พัฒนา
type Developer struct {
	BaseModel
	Name string `json:"name" gorm:"type:varchar(255);not null;uniqueIndex:,option:CONCURRENTLY;" validate:"required" filter:"like"`
}
type DeveloperCreate struct {
	ID   uuid.UUID `json:"id" form:"-" query:"-"`
	Name string    `json:"name" validate:"required" form:"name" query:"name"`
}

func (DeveloperCreate) TableName() string {
	return "developers"
}
func (s *DeveloperCreate) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

type DeveloperUpdate struct {
	ID   uuid.UUID `json:"id" validate:"required,uuid" form:"-" query:"-"`
	Name *string   `json:"name,omitempty" validate:"omitempty" form:"name" query:"name"`
}

func (DeveloperUpdate) TableName() string {
	return "developers"
}
