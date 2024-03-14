package domain

import (
	"fmt"
	"go_base/logger"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	StaffCtx                       = "staff"
	StaffAuthCache                 = "backlist:email:%s"
	StaffVerifyTokenType TokenType = "verify_token"

	// Log
	UnlockLog = "unlock"

	LoginLog  = "login"
	LoginFail = "login_failed"

	ChangePasswordLog       = "change_password"
	ChangePasswordFailedLog = "change_password_failed"
)

// ----
type Name string
type Email string
type Status string

// ใช้งาน,ระงับใช้งาน,รอพิจารณา (อย่าลืมไปเพิ่มตรง ValidateStaffStatus ด้วยน่ะจ้าา)
const (
	StaffActive   Status = "active"
	StaffInactive Status = "inactive"
	StaffPending  Status = "pending"
)

type Staff struct {
	BaseModel
	Email       SensitiveString `json:"email" validate:"required,email" query:"email" swagger:"desc(email)" form:"email" gorm:"index:,option:CONCURRENTLY,unique" filter:"like" `
	FirstName   string          `json:"first_name" validate:"required" query:"first_name" swagger:"desc(first_name)" form:"first_name" gorm:"varchar(255);not null" filter:"="`
	LastName    string          `json:"last_name" validate:"required" query:"last_name" swagger:"desc(last_name)" form:"last_name" gorm:"varchar(255);not null" filter:"="`
	Password    Password        `json:"-" query:"password" swagger:"desc(password)" form:"password" gorm:"not null" password:"true"`
	TmpPassword string          `json:"tmp_password,omitempty" query:"-" swagger:"desc(tmp_password)" form:"-" gorm:"-"`
	LastLogin   *time.Time      `json:"last_login,omitempty" gorm:"index"`
	IsVerified  bool            `json:"is_verified" gorm:"default:false" validate:"bool"`
	VerifyToken string          `json:"-" gorm:"default:''" validate:"lowercase"`
	Status      Status          `json:"status" gorm:"default:pending" validate:"staff_status" filter:"="`
	Phone       *string         `json:"phone,omitempty" gorm:"varchar(255);" validate:"omitempty,phone" filter:"="`

	// fk role nullable
	RoleID *uuid.UUID `json:"-" gorm:"type:uuid;index:,option:CONCURRENTLY;" validate:"omitempty,uuid" filter:"="`
	Role   *Role      `json:"role,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (s *Staff) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Password != "" {
		s.Password = s.Password.Hash()
	}
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

type StaffVerifyToken struct {
	Token string `json:"token" validate:"required" query:"token" swagger:"desc(token),required" form:"token"`
}

type StaffVerifyTokenResponse struct {
	Email string `json:"email"`
}

type StaffLogin struct {
	Email    SensitiveString `json:"email" validate:"required" query:"email" swagger:"desc(email),required" form:"email" `
	Password string          `json:",omitempty" validate:"required" query:"password" swagger:"desc(password),required" form:"password"`
}

type StaffUnlock struct {
	Email string `json:"email" validate:"required,email" query:"email" swagger:"desc(email),required" filter:"email" form:"email"`
}

type StaffCreate struct {
	Email     SensitiveString `json:"email" validate:"required,email" query:"email" swagger:"desc(email),required" form:"email" `
	FirstName string          `json:"first_name" validate:"required" query:"first_name" swagger:"desc(first_name),required" form:"first_name"`
	LastName  string          `json:"last_name" validate:"required" query:"last_name" swagger:"desc(last_name),required" form:"last_name"`
	Phone     *string         `json:"phone,omitempty" validate:"omitempty,phone" query:"phone" swagger:"desc(phone)" form:"phone"`

	RoleID *string `json:"role_id" validate:"omitempty,uuid" query:"role_id" swagger:"desc(role_id)" form:"role_id"`
}
type StaffUpdate struct {
	ID        uuid.UUID `json:"id" query:"-" form:"-"`
	Email     *string   `json:"email,omitempty"  query:"email" swagger:"desc(email)" form:"email"`
	FirstName *string   `json:"first_name,omitempty" query:"first_name" swagger:"desc(first_name)" form:"first_name"`
	LastName  *string   `json:"last_name,omitempty" query:"last_name" swagger:"desc(last_name)" form:"last_name"`
	Password  *string   `json:",omitempty" query:"password" swagger:"desc(password)" form:"password"`
	Status    *Status   `json:"status" query:"status" swagger:"desc(status)" form:"status" validate:"omitempty,staff_status"`
	Phone     *string   `json:"phone,omitempty" gorm:"varchar(255);default:''" validate:"omitempty,phone"`

	RoleID *string `json:"role_id,omitempty" query:"role_id" swagger:"desc(role_id)" form:"role_id" validate:"omitempty,uuid"`
}

func (StaffUpdate) TableName() string {
	return "staffs"
}

type StaffGetToken struct {
	Email       string `json:"email" validate:"required,email" query:"email" swagger:"desc(email),required" filter:"email" form:"email"`
	TmpPassword string `json:"tmp_password" validate:"required" query:"tmp_password" swagger:"desc(tmp_password),required" form:"tmp_password"`
}
type StaffGetTokenResponse struct {
	Token string `json:"token"`
}

type StaffUpdatePassword struct {
	Email       SensitiveString `json:"email" validate:"required,email" query:"-" swagger:"-" form:"-"`
	Password    string          `json:"password" validate:"required" query:"password" swagger:"desc(password),required" form:"password"`
	OldPassword string          `json:"old_password" validate:"required" query:"old_password" swagger:"desc(old_password),required" form:"old_password"`
}

type StaffGetLog struct {
	ID uuid.UUID `json:"id" validate:"required" query:"id" form:"id"`
}

type StaffMe struct {
	Email     SensitiveString `json:"email"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	LastLogin *time.Time      `json:"last_login,omitempty"`
	Status    Status          `json:"status"`
	Phone     *string         `json:"phone,omitempty"`
	// fk role
	RoleID *uuid.UUID `json:"role_id,omitempty"`
	Role   *Role      `json:"role,omitempty" `
}

func (StaffMe) TableName() string {
	return "staffs"
}

type StaffFK struct {
	BaseModel
	Email     SensitiveString `json:"email" validate:"required,email" query:"email" swagger:"desc(email)" form:"email" gorm:"index:,option:CONCURRENTLY,unique" `
	FirstName string          `json:"first_name" validate:"required" query:"first_name" swagger:"desc(first_name)" form:"first_name" gorm:"varchar(255);not null"`
	LastName  string          `json:"last_name" validate:"required" query:"last_name" swagger:"desc(last_name)" form:"last_name" gorm:"varchar(255);not null"`
}

func (StaffFK) TableName() string {
	return "staffs"
}
func StaffGetName(staff *Staff) string {
	if staff == nil {
		return ""
	}
	if staff.FirstName == "" && staff.LastName == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", staff.FirstName, staff.LastName)
}

type VerifyTokenClaims struct {
	jwt.RegisteredClaims
	Token     string    `json:"token"`
	TokenType TokenType `json:"token_type"`
}

func GenerateVerifyToken(token, secret string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := VerifyTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		TokenType: StaffVerifyTokenType,
		Token:     token,
	}

	t, err := GenerateAccessToken(claims, string(StaffVerifyTokenType))
	if err != nil {
		logger.L().Error(err)
		return "", err
	}
	return t, nil
}

func ParseVerifyToken(token string) error {
	var claims VerifyTokenClaims
	if _, err := ParseToken(token, &claims, string(StaffVerifyTokenType)); err != nil {
		logger.L().Error(err)
		return err
	}
	return nil
}
