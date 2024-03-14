package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ทรัพย์
// ข้อมูล user จะต้องอยู่หลังจาก user ถูกลบไปแล้ว
type Asset struct {
	BaseModel
	// เลขที่ทรัพย์
	No          *string `json:"no,omitempty" gorm:"type:varchar(255);" validate:"omitempty" filter:"="`
	ProjectName *string `json:"project_name,omitempty" gorm:"-" filter:"projects.name.like"`

	// FK to Project
	Project   *Project  `json:"project,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" `
	ProjectID uuid.UUID `json:"project_id" validate:"uuid" gorm:"type:uuid;unique,composite:idx_asset_project_id_user_id" filter:"="`

	// FK to User
	UserID uuid.UUID `json:"user_id" validate:"uuid" gorm:"type:uuid;unique,composite:idx_asset_project_id_user_id" filter:"="`
	User   *User     `json:"user,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" `

	UserFirstName *string `json:"user_first_name,omitempty" gorm:"-" filter:"users.first_name.like" `
	UserLastName  *string `json:"user_last_name,omitempty" gorm:"-" filter:"users.last_name.like"`
	// Meta data
	// รายละเอียด
	Description *string `json:"description,omitempty" gorm:"type:text;" validate:"omitempty" filter:"like"`
	// แผนที่ (Google Map)
	Map *string `json:"map,omitempty" gorm:"type:text;" validate:"omitempty"`
	// ขนาด (ตร.ม.)
	Size *float64 `json:"size,omitempty" gorm:"type:numeric;" validate:"omitempty"`
	// โซน/เขต
	Zone *string `json:"zone,omitempty" gorm:"type:varchar(255);" validate:"omitempty" filter:"="`
	// ประเภท
	Type *string `json:"type,omitempty" gorm:"type:varchar(255);" validate:"omitempty" filter:"assets.type.="`
	// ราคา (ซื้อ/ขาย)
	Price *float64 `json:"price,omitempty" gorm:"type:numeric;" validate:"omitempty"`
}

type AssetCreate struct {
	ID uuid.UUID `json:"id" form:"-" query:"-"`
	// เลขที่ทรัพย์
	No *string `json:"no,omitempty" validate:"omitempty" form:"no" query:"no"`

	// FK to Project
	ProjectID uuid.UUID `json:"project_id" validate:"uuid" form:"project_id" query:"project_id"`

	// FK to User
	UserID *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,uuid" form:"user_id" query:"user_id"`

	// Meta data
	// รายละเอียด
	Description *string  `json:"description,omitempty" validate:"omitempty" form:"description" query:"description"`
	Map         *string  `json:"map,omitempty" validate:"omitempty" form:"map" query:"map"`
	Size        *float64 `json:"size,omitempty" validate:"omitempty" form:"size" query:"size"`
	Zone        *string  `json:"zone,omitempty" validate:"omitempty" form:"zone" query:"zone"`
	Type        *string  `json:"type,omitempty" validate:"omitempty" form:"type" query:"type"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty" form:"price" query:"price"`
}

func (AssetCreate) TableName() string {
	return "assets"
}
func (s *AssetCreate) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

type AssetUpdate struct {
	ID          uuid.UUID  `json:"id" validate:"required,uuid" form:"-" query:"-"`
	No          *string    `json:"no,omitempty" validate:"omitempty" form:"no" query:"no"`
	ProjectID   *uuid.UUID `json:"project_id,omitempty" validate:"omitempty,uuid" form:"project_id" query:"project_id"`
	UserID      *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,uuid" form:"user_id" query:"user_id"`
	Description *string    `json:"description,omitempty" validate:"omitempty" form:"description" query:"description"`
	Map         *string    `json:"map,omitempty" validate:"omitempty" form:"map" query:"map"`
	Size        *float64   `json:"size,omitempty" validate:"omitempty" form:"size" query:"size"`
	Zone        *string    `json:"zone,omitempty" validate:"omitempty" form:"zone" query:"zone"`
	Type        *string    `json:"type,omitempty" validate:"omitempty" form:"type" query:"type"`
	Price       *float64   `json:"price,omitempty" validate:"omitempty" form:"price" query:"price"`
}

func (AssetUpdate) TableName() string {
	return "assets"
}
