package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// user migration ไม่ได้ ignore ฟิลด์ password และ tmp_password
type UserMigration struct {
	BaseModel
	Email       SensitiveString `json:"email" validate:"required,email" query:"email" swagger:"desc(email)" form:"email" gorm:"index:,option:CONCURRENTLY,unique" `
	FirstName   string          `json:"first_name" validate:"required" query:"first_name" swagger:"desc(first_name)" form:"first_name" gorm:"varchar(255);not null"`
	LastName    string          `json:"last_name" validate:"required" query:"last_name" swagger:"desc(last_name)" form:"last_name" gorm:"varchar(255);not null"`
	Password    Password        `json:"password" query:"password" swagger:"desc(password)" form:"password" gorm:"not null" password:"true"`
	TmpPassword string          `json:"tmp_password,omitempty" query:"-" swagger:"desc(tmp_password)" form:"-" gorm:"-"`
	LastLogin   *time.Time      `json:"last_login,omitempty" gorm:"index"`
	IsVerified  bool            `json:"is_verified" gorm:"default:false" validate:"bool"`
	VerifyToken string          `json:"-" gorm:"default:''" validate:"lowercase"`

	// Meta data
	// งบประมาณ (ซื้อ)
	BudgetBuy *float64 `json:"budget_buy,omitempty" gorm:"type:numeric(17,2);default:0.00"`
	// งบประมาณ (ขาย)
	BudgetSell *float64 `json:"budget_sell,omitempty" gorm:"type:numeric(17,2);default:0.00"`
	// งบประมาณ (เช่า)
	BudgetPerMonth *float64 `json:"budget_per_month,omitempty" gorm:"type:numeric(17,2);default:0.00"`

	Phone *string `json:"phone,omitempty" gorm:"varchar(255);" validate:"omitempty,phone" filter:"="`

	// แหล่งที่มา
	Source *string `json:"source,omitempty" gorm:"varchar(255);"`
	//  พนักงานที่รับผิดชอบ
	StaffID *uuid.UUID `json:"staff_id,omitempty" gorm:"type:uuid;index:,option:CONCURRENTLY;" validate:"omitempty,uuid" filter:"="`
	Staff   *StaffFK   `json:"staff,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// สิ่งที่ต้องทำ
	Todo   *string    `json:"todo,omitempty" gorm:"varchar(255);"`
	TodoAt *time.Time `json:"todo_at,omitempty" gorm:"index"`

	// ประเภท [1: ผู้ซื้อ, 2: ผู้ขาย] can be all or null datatypes.JSON
	Type *datatypes.JSON `json:"type,omitempty" gorm:"type:jsonb;default:'[]'" filter:"in" validate:"omitempty,valid_jsonb,enum=buyer seller"`
	// ความสนใจ [1: ขาย, 2: ซื้อ, 3: บริหาร] can be all or null datatypes.JSON
	Interest *datatypes.JSON `json:"interest,omitempty" gorm:"type:jsonb;default:'[]'" filter:"in" validate:"omitempty,valid_jsonb,enum=sell buy manage"`

	// สถานะ
	Status *string `json:"status,omitempty" gorm:"varchar(255);" filter:"="`

	// แท็ก [1: คอนโด, 2:สุขุมวิท]
	Tag *datatypes.JSON `json:"tag,omitempty" gorm:"type:jsonb;default:'[]'" filter:"in"`

	// กิจกรรมล่าสุด
	LastActivityAt *time.Time `json:"last_activity_at,omitempty" gorm:"index"`
	LastActivity   *string    `json:"last_activity,omitempty" gorm:"varchar(255);"`

	// Contact (Full Name,Display Name,DOB,Full Address)
	FullName    *string    `json:"full_name,omitempty" gorm:"varchar(255);"`
	DisplayName *string    `json:"display_name,omitempty" gorm:"varchar(255);"`
	DOB         *time.Time `json:"dob,omitempty" gorm:"index"`
	FullAddress *string    `json:"full_address,omitempty" gorm:"varchar(255);"`

	// Preferences (Language,Timezone,Date Format)
	Language   *string `json:"language,omitempty" gorm:"varchar(255);"`
	Timezone   *string `json:"timezone,omitempty" gorm:"varchar(255);"`
	DateFormat *string `json:"date_format,omitempty" gorm:"varchar(255);"`

	Gender *string `json:"gender,omitempty" gorm:"varchar(255);"`
}

func (UserMigration) TableName() string {
	return "users"
}

func (s *UserMigration) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Password != "" {
		s.Password = s.Password.Hash()
	}
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
