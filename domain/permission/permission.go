package permission

var (
	Resources = []string{
		"staff",
	}
)

const (
	SystemAdmin = "admin"
)

const (
	// if view is true, then the user can view the resource and menu, if false can't view and can't access resource (CURD)
	ActionFind   = "view"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionUnlock = "unlock"
)

const (
	ScopeAll = "true"
)

// Permission name format: {system}.{resource}.{action}.{scope}
const (
	ROLE_VIEW_ALL   = "admin.role.view.true"
	ROLE_CREATE_ALL = "admin.role.create.true"
	ROLE_UPDATE_ALL = "admin.role.update.true"
	ROLE_EXPORT_ALL = "admin.role.export.true"

	STAFF_VIEW_ALL   = "admin.staff.view.true"
	STAFF_CREATE_ALL = "admin.staff.create.true"
	STAFF_UPDATE_ALL = "admin.staff.update.true"
	STAFF_DELETE_ALL = "admin.staff.delete.true"
	STAFF_UNLOCK_ALL = "admin.staff.unlock.true"

	STAFF_ME_FIND_SELF = "admin.staff_me.view.true"
	STAFF_ME_LOG_SELF  = "admin.staff_me.log.true"

	USER_VIEW_ALL   = "admin.user.view.true"
	USER_UPDATE_ALL = "admin.user.update.true"
	USER_DELETE_ALL = "admin.user.delete.true"
	USER_UNLOCK_ALL = "admin.user.unlock.true"

	ROLE_FIND   = "admin.role.view.true"
	ROLE_CREATE = "admin.role.create.true"
	ROLE_UPDATE = "admin.role.update.true"

	DEVELOPER_VIEW_ALL   = "admin.developer.view.true"
	DEVELOPER_CREATE_ALL = "admin.developer.create.true"
	DEVELOPER_UPDATE_ALL = "admin.developer.update.true"
	DEVELOPER_EXPORT_ALL = "admin.developer.export.true"
	DEVELOPER_DELETE_ALL = "admin.developer.delete.true"

	PROJECT_VIEW_ALL   = "admin.project.view.true"
	PROJECT_CREATE_ALL = "admin.project.create.true"
	PROJECT_UPDATE_ALL = "admin.project.update.true"
	PROJECT_EXPORT_ALL = "admin.project.export.true"
	PROJECT_DELETE_ALL = "admin.project.delete.true"

	ASSET_VIEW_ALL   = "admin.asset.view.true"
	ASSET_CREATE_ALL = "admin.asset.create.true"
	ASSET_UPDATE_ALL = "admin.asset.update.true"
	ASSET_EXPORT_ALL = "admin.asset.export.true"
	ASSET_DELETE_ALL = "admin.asset.delete.true"
)
