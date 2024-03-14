package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// โครงการ
type Project struct {
	BaseModel
	Name string `json:"name" gorm:"type:varchar(255);not null;uniqueIndex:,option:CONCURRENTLY;" validate:"required" filter:"like"`

	// FK to Developer
	DeveloperID *uuid.UUID `json:"developer_id,omitempty" gorm:"type:uuid;index:,option:CONCURRENTLY;" validate:"omitempty,uuid" filter:"="`
	Developer   *Developer `json:"developer,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// โครงการ
type ProjectCreate struct {
	ID   uuid.UUID `json:"id" form:"-" query:"-"`
	Name string    `json:"name" form:"name" query:"name" validate:"required"`

	// FK to Developer
	DeveloperID *string `json:"developer_id,omitempty" validate:"omitempty,uuid" form:"developer_id" query:"developer_id"`
}

func (ProjectCreate) TableName() string {
	return "projects"
}
func (s *ProjectCreate) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

type ProjectUpdate struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid" form:"-" query:"-"`
	Name        *string   `json:"name,omitempty" validate:"omitempty" form:"name" query:"name"`
	DeveloperID *string   `json:"developer_id,omitempty" validate:"omitempty,uuid" form:"developer_id" query:"developer_id"`
}

func (ProjectUpdate) TableName() string {
	return "projects"
}
