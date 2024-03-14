package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// staff migration ไม่ได้ ignore ฟิลด์ password และ tmp_password
type StaffMigration struct {
	BaseModel
	Email       SensitiveString `json:"email" validate:"required,email" `
	FirstName   string          `json:"first_name" validate:"required" `
	LastName    string          `json:"last_name" validate:"required" `
	Password    Password        `json:"password" `
	LastLogin   *time.Time      `json:"last_login,omitempty" `
	IsVerified  bool            `json:"is_verified"`
	VerifyToken string          `json:"-"`
	Status      Status          `json:"status"`
	Phone       *string         `json:"phone,omitempty"`
	RoleID      *uuid.UUID      `json:"role_id" validate:"required,uuid" `
}

func (StaffMigration) TableName() string {
	return "staffs"
}

func (s *StaffMigration) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Password != "" {
		s.Password = s.Password.Hash()
	}
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

// ไม่ได้ bind password
type StaffMock struct {
	BaseModel
	Email       SensitiveString `json:"email" validate:"required,email" `
	FirstName   string          `json:"first_name" validate:"required" `
	LastName    string          `json:"last_name" validate:"required" `
	Password    string          `json:"password" `
	LastLogin   *time.Time      `json:"last_login,omitempty" `
	IsVerified  bool            `json:"is_verified"`
	Status      Status          `json:"status"`
	Phone       *string         `json:"phone,omitempty"`
	RoleID      *uuid.UUID      `json:"role_id" validate:"required,uuid" `
	VerifyToken string          `json:"-"`
}

func (StaffMock) TableName() string {
	return "staffs"
}
