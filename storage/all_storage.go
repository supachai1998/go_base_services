package storage

import "gorm.io/gorm"

type AllStorage struct {
	Cache *Cache
	DB    *gorm.DB
}
