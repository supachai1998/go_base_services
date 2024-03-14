package database

import "go_base/domain"

type AllStores struct {
	Base      *Store
	Staff     *StaffStore
	Auth      *AuthStore
	Role      *RoleStore
	User      *UserStore
	Developer *BaseStore[domain.Developer, domain.DeveloperUpdate, domain.DeveloperCreate]
	Project   *BaseStore[domain.Project, domain.ProjectUpdate, domain.ProjectCreate]
	Asset     *BaseStore[domain.Asset, domain.AssetUpdate, domain.AssetCreate]
}
