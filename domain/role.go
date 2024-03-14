package domain

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	RoleTypeSuperAdmin RoleType = "SUPER_ADMIN"
	RoleTypeAdmin      RoleType = "ADMIN"
	RoleTypeDoctor     RoleType = "DOCTOR"

	PermissionActionMenu   PermissionAction = "menu"
	PermissionActionFind   PermissionAction = "find"
	PermissionActionCreate PermissionAction = "create"
	PermissionActionUpdate PermissionAction = "update"
	PermissionActionDelete PermissionAction = "delete"
	PermissionActionExport PermissionAction = "export"

	PermissionScopeAll PermissionScope = "all"
	PermissionScopeOrg PermissionScope = "organization"
)

type RoleType string

func NewRoleType(str string) (RoleType, error) {
	switch str {
	case RoleTypeSuperAdmin.String():
		return RoleTypeSuperAdmin, nil
	case RoleTypeAdmin.String():
		return RoleTypeAdmin, nil
	case RoleTypeDoctor.String():
		return RoleTypeDoctor, nil
	default:
		return "", errors.New("type not supported")
	}

}

func (r RoleType) String() string {
	return string(r)
}

// system.resource.action.scope
// e.g.
//
//	{
//		"admin": {
//			"booking": {
//				"menu": "true",
//				"create": "true",
//				"update": "true",
//				"delete": "true",
//				"export": "false"
//			}
//		}
//	}
type PermissionTree map[string]map[string]map[string]string
type Permission string

type PermissionAction string
type PermissionScope string

func (s PermissionAction) String() string { return string(s) }

func (s PermissionScope) String() string { return string(s) }

type Role struct {
	BaseModel
	Type        RoleType       ` json:"type" gorm:"index:,unique,composite:idx_type_name_tier_level"`
	Name        string         ` json:"name" gorm:"index:,unique,composite:idx_type_name_tier_level"`
	Description string         ` json:"description"`
	Permissions datatypes.JSON ` json:"permissions" gorm:"type:jsonb"`
	CountStaff  *int64         ` json:"count_staff,omitempty" gorm:"-"`
}
type RoleUpdate struct {
	ID          uuid.UUID       `json:"id" form:"-" query:"-" validate:"required,uuid4"`
	Type        *RoleType       ` json:"type" form:"type" query:"type" validate:"omitempty"`
	Name        *string         ` json:"name" form:"name" query:"name" validate:"omitempty,max=20"`
	Description *string         ` json:"description" form:"description" query:"description" validate:"omitempty,max=100"`
	Permissions *datatypes.JSON ` json:"permissions" form:"permissions" query:"permissions" validate:"omitempty,valid_permissions"`
}

func (r *RoleUpdate) TableName() string {
	return "roles"
}

type RoleSwaggerCreate struct {
	Type        RoleType       ` json:"type" validate:"required"`
	Name        string         ` json:"name" validate:"required,max=20"`
	Description string         ` json:"description" validate:"max=100"`
	Permissions datatypes.JSON ` json:"permissions" validate:"required,valid_permissions"`
}

type RoleMetadata struct {
	ID   uuid.UUID `json:"id" `
	Type RoleType  `json:"type" `
	Name string    `json:"name" `
}

type RoleWithStaffCount struct {
	Type        RoleType ` json:"type"`
	Name        string   ` json:"name"`
	Description string   ` json:"description"`
	StaffCount  int      ` json:"staff_count"`
}

// inject reflect name of struct to name
func (RoleWithStaffCount) TableName() string {
	return "roles"
}

type RoleWithPermissionInfoResponse struct {
	Type        RoleType       `json:"type"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Permissions PermissionTree `json:"permissions"`
}

type RoleStore interface {
	// FindByIDs(ctx context.Context, IDs []uuid.UUID) ([]*Model[*Role], error)
	// FindByTypeName(ctx context.Context, roleType RoleType, name string) (*Model[*Role], error)
}

func NewRole(roleType RoleType, name, description string, permissions datatypes.JSON) *Role {
	return &Role{
		Type:        roleType,
		Name:        name,
		Description: description,
		Permissions: permissions,
	}
}
