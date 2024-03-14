package domain

import (
	"fmt"

	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	UserCtx                       = "user"
	UserAuthCache                 = "backlist:email:%s"
	UserVerifyTokenType TokenType = "verify_token"
)

// ----

type User struct {
	BaseModel
	Email       SensitiveString `json:"email" validate:"required,email" query:"email" swagger:"desc(email)" form:"email" gorm:"index:,option:CONCURRENTLY,unique" filter:"like"`
	FirstName   string          `json:"first_name" validate:"required" query:"first_name" swagger:"desc(first_name)" form:"first_name" gorm:"varchar(255);not null" filter:"="`
	LastName    string          `json:"last_name" validate:"required" query:"last_name" swagger:"desc(last_name)" form:"last_name" gorm:"varchar(255);not null" filter:"="`
	Password    Password        `json:"-" query:"password" swagger:"desc(password)" form:"password" gorm:"not null" password:"true"`
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

func (s *User) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Password != "" {
		s.Password = s.Password.Hash()
	}
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

type UserVerifyToken struct {
	Token string `json:"token" validate:"required" query:"token" swagger:"desc(token),required" form:"token"`
}

type UserVerifyTokenResponse struct {
	Email string `json:"email"`
}

type UserLogin struct {
	Email    SensitiveString `json:"email" validate:"required" query:"email" swagger:"desc(email),required" form:"email" `
	Password string          `json:",omitempty" validate:"required" query:"password" swagger:"desc(password),required" form:"password"`
}

type UserUnlock struct {
	Email string `json:"email" validate:"required,email" query:"email" swagger:"desc(email),required" filter:"email" form:"email"`
}

type UserCreate struct {
	Email     SensitiveString `json:"email" validate:"required,email" query:"email" swagger:"desc(email),required" form:"email" `
	FirstName string          `json:"first_name" validate:"required" query:"first_name" swagger:"desc(first_name),required" form:"first_name"`
	LastName  string          `json:"last_name" validate:"required" query:"last_name" swagger:"desc(last_name),required" form:"last_name"`
}
type UserUpdate struct {
	ID        uuid.UUID        `json:"id" query:"-" form:"-"`
	Email     *SensitiveString `json:"email,omitempty" validate:"omitempty,email" query:"email" swagger:"desc(email)" form:"email"`
	FirstName *string          `json:"first_name,omitempty" query:"first_name" swagger:"desc(first_name)" form:"first_name"`
	LastName  *string          `json:"last_name,omitempty" query:"last_name" swagger:"desc(last_name)" form:"last_name"`

	// Meta data
	// งบประมาณ
	BudgetBuy      *float64 `json:"budget_buy,omitempty" form:"budget_buy" query:"budget_buy"`
	BudgetSell     *float64 `json:"budget_sell,omitempty" form:"budget_sell" query:"budget_sell"`
	BudgetPerMonth *float64 `json:"budget_per_month,omitempty" form:"budget_per_month" query:"budget_per_month"`

	// แหล่งที่มา
	Source *string `json:"source,omitempty" form:"source" query:"source"`
	//  พนักงานที่รับผิดชอบ
	StaffID *string `json:"staff_id,omitempty" form:"staff_id" query:"staff_id" validate:"omitempty,uuid"`

	// สิ่งที่ต้องทำ
	Todo   *string    `json:"todo,omitempty" form:"todo" query:"todo"`
	TodoAt *time.Time `json:"todo_at,omitempty" form:"todo_at" query:"todo_at"`

	// ประเภท [1: ผู้ซื้อ, 2: ผู้ขาย] can be all or null datatypes.JSON
	Type *datatypes.JSON `json:"type,omitempty" form:"type" query:"type" validate:"omitempty,valid_jsonb,enum=buyer seller"`
	// ความสนใจ [1: ขาย, 2: ซื้อ, 3: บริหาร] can be all or null datatypes.JSON
	Interest *datatypes.JSON `json:"interest,omitempty" form:"interest" query:"interest" validate:"omitempty,valid_jsonb,enum=sell buy manage"`

	// สถานะ
	Status *string `json:"status,omitempty" form:"status" query:"status" validate:"omitempty,max=255"`

	// แท็ก [1: คอนโด, 2:สุขุมวิท]
	Tag *datatypes.JSON `json:"tag,omitempty" form:"tag" query:"tag" validate:"omitempty,valid_jsonb"`

	// กิจกรรมล่าสุด
	LastActivityAt *time.Time `json:"last_activity_at,omitempty" form:"last_activity_at" query:"last_activity_at" validate:"omitempty"`
	LastActivity   *string    `json:"last_activity,omitempty" form:"last_activity" query:"last_activity" validate:"omitempty,max=255"`

	// Contact (Full Name,Display Name,DOB,Full Address)
	FullName    *string    `json:"full_name,omitempty" form:"full_name" query:"full_name" validate:"omitempty,max=255"`
	DisplayName *string    `json:"display_name,omitempty" form:"display_name" query:"display_name" validate:"omitempty,max=255"`
	DOB         *time.Time `json:"dob,omitempty" form:"dob" query:"dob" validate:"omitempty"`

	// Preferences (Language,Timezone,Date Format)
	Language   *string `json:"language,omitempty" form:"language" query:"language" validate:"omitempty,max=255"`
	Timezone   *string `json:"timezone,omitempty" form:"timezone" query:"timezone" validate:"omitempty,max=255"`
	DateFormat *string `json:"date_format,omitempty" form:"date_format" query:"date_format" validate:"omitempty,max=255"`

	Gender *string `json:"gender,omitempty" form:"gender" query:"gender" validate:"omitempty,max=255"`
}

type UserGetToken struct {
	Email       string `json:"email" validate:"required,email" query:"email" swagger:"desc(email),required" filter:"email" form:"email"`
	TmpPassword string `json:"tmp_password" validate:"required" query:"tmp_password" swagger:"desc(tmp_password),required" form:"tmp_password"`
}
type UserGetTokenResponse struct {
	Token string `json:"token"`
}

type UserUpdatePassword struct {
	Email       SensitiveString `json:"email" validate:"required,email" query:"-" swagger:"-" form:"-"`
	Password    string          `json:"password" validate:"required" query:"password" swagger:"desc(password),required" form:"password"`
	OldPassword string          `json:"old_password" validate:"required" query:"old_password" swagger:"desc(old_password),required" form:"old_password"`
}

type UserGetLog struct {
	ID uuid.UUID `json:"id" validate:"required" query:"id" form:"id"`
}

type UserMe struct {
	ID        uuid.UUID       `json:"id"`
	Email     SensitiveString `json:"email"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	LastLogin *time.Time      `json:"last_login,omitempty"`

	BudgetBuy      *float64 `json:"budget_buy,omitempty"`
	BudgetSell     *float64 `json:"budget_sell,omitempty"`
	BudgetPerMonth *float64 `json:"budget_per_month,omitempty"`

	Source *string `json:"source,omitempty"`

	StaffID *uuid.UUID `json:"staff_id,omitempty"`
	Staff   *StaffFK   `json:"staff,omitempty"`

	Todo   *string    `json:"todo,omitempty"`
	TodoAt *time.Time `json:"todo_at,omitempty"`

	Type     *datatypes.JSON `json:"type,omitempty"`
	Interest *datatypes.JSON `json:"interest,omitempty"`

	Status *string `json:"status,omitempty"`

	Tag *datatypes.JSON `json:"tag,omitempty"`

	LastActivityAt *time.Time `json:"last_activity_at,omitempty"`
	LastActivity   *string    `json:"last_activity,omitempty"`
}

func (UserMe) TableName() string {
	return "users"
}
func (UserUpdate) TableName() string {
	return "users"
}

func UserGetName(user *User) string {
	if user == nil {
		return ""
	}
	if user.FirstName == "" && user.LastName == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}
