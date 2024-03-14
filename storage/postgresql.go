package storage

import (
	"encoding/json"
	"go_base/configs"
	"go_base/domain"
	helper "go_base/domain/helper"
	"go_base/logger"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Client struct {
	Client *gorm.DB
	// cfg        *configs.Config
}

func NewPostgresClient(dsn string, cfg *gorm.Config, serverCfg *configs.Config) (*Client, error) {
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}

	if ping, err := db.DB(); err != nil {
		return nil, err
	} else {
		if err := ping.Ping(); err != nil {
			return nil, err
		}
	}
	db.Callback().Create().Before("gorm:create").Register("hash_password", hashPassword)
	db.Callback().Update().Before("gorm:update").Register("hash_password", hashPassword)
	logger.L().Infof("Postgres initialized: %v", dsn)

	if err := migration(db); err != nil {
		return nil, err
	}
	logger.L().Info("Postgres migration completed")

	return &Client{
		Client: db,
	}, nil
}

func migration(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.Role{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&domain.Staff{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&domain.TokenExpires{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&domain.Auth{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		return err
	}
	return nil
}

func Seed[T any](db *gorm.DB, file string) error {
	dir := "storage/seed"

	filePath := filepath.Join(configs.Root, dir, file)
	_, ext := strings.Split(file, ".")[0], strings.Split(file, ".")[1]
	var t T
	tableName, err := helper.GetTableNameByStruct(t)
	if err != nil {
		return err
	}
	tableName = strings.ReplaceAll(tableName, "_migrations", "s")

	logger.L().Infof("Seeding %s", tableName)
	switch ext {
	case "csv":
		logger.L().Infof("Seeding %s with %s", tableName, filePath)
		if err := db.Transaction(func(tx *gorm.DB) error {
			// clear records
			if err := tx.Exec("TRUNCATE TABLE " + tableName + " CASCADE").Error; err != nil {
				return err
			}
			if err := tx.Exec("COPY " + tableName + " FROM '" + filePath + "' WITH (FORMAT CSV, DELIMITER ',', HEADER true)").Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	case "json":
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		var items []T
		if err := json.Unmarshal(content, &items); err != nil {
			return err
		}
		logger.L().Infof("Seeding %s with %d items", tableName, len(items))
		if err := db.Transaction(func(tx *gorm.DB) error {
			// clear records
			if err := tx.Exec("TRUNCATE TABLE " + tableName + " CASCADE").Error; err != nil {
				return err
			}
			// insert bulk
			logger.L().Infof("Inserting %s items %d", tableName, len(items))
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error; err != nil {
				logger.L().Errorf("Inserting %s items %d failed: %v", tableName, len(items), err)
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

// Migration : Migrate the database data don't drop the table
func Migration[T any](db *gorm.DB, file string) error {
	dir := "storage/seed"

	filePath := filepath.Join(configs.Root, dir, file)
	_, ext := strings.Split(file, ".")[0], strings.Split(file, ".")[1]
	var t T
	tableName, err := helper.GetTableNameByStruct(t)
	if err != nil {
		return err
	}
	tableName = strings.ReplaceAll(tableName, "_migrations", "s")

	logger.L().Infof("Migrating %s", tableName)
	if ext != "json" { // only support json
		return nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var items []T
	if err := json.Unmarshal(content, &items); err != nil {
		return err
	}
	logger.L().Infof("Migrating %s with %d items", tableName, len(items))
	if err := db.Transaction(func(tx *gorm.DB) error {
		// insert bulk
		logger.L().Infof("Inserting %s items %d", tableName, len(items))
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&items).Error; err != nil {
			logger.L().Errorf("Inserting %s items %d failed: %v", tableName, len(items), err)
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func GetAllIntoRedis[T any](db *gorm.DB) error {
	var items []T
	if err := db.Find(&items).Error; err != nil {
		return err
	}
	return nil
}

func Mock[T any](total int, db *gorm.DB, mockData func() T) error {
	var items []T

	tableName, err := helper.GetTableNameByStruct(items)
	if err != nil {
		return err
	}
	logger.L().Infof("Mock %s total :%d", tableName, total)
	maxGoroutines := 20
	if total < maxGoroutines {
		maxGoroutines = total
	}
	guard := make(chan struct{}, maxGoroutines)
	totalInsert := 0
	var wg sync.WaitGroup
	for i := 0; i < total; i += maxGoroutines {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			guard <- struct{}{} // would block if guard channel is already filled
			defer func() {
				<-guard
			}()

			data := make([]T, maxGoroutines)
			bulkMax := 100
			if end-start < bulkMax {
				bulkMax = end - start
			}
			for j := 0; j < bulkMax; j++ {
				data[j] = mockData()
			}
			if len(data) == 0 {
				return
			}
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&data).Error; err != nil {
				return
			}
			totalInsert += len(data)
			logger.L().Infof("Inserted %s items : %d ( %f%% )", tableName, totalInsert, float64(totalInsert)/float64(total)*100)
		}(i, i+maxGoroutines)
	}

	wg.Wait()
	logger.L().Infof("Mock %s finished", tableName)

	return nil
}

func Drop[T any](db *gorm.DB) error {
	var t T
	tableName, err := helper.GetTableNameByStruct(t)
	if err != nil {
		return err
	}
	logger.L().Infof("Drop %s", tableName)
	if err := db.Exec("DROP TABLE " + tableName + " CASCADE").Error; err != nil {
		return err
	}
	db.AutoMigrate(t)
	return nil
}

func hashPassword(db *gorm.DB) {

	if db.Statement.Schema != nil {
		maxGoroutines := 5
		guard := make(chan struct{}, maxGoroutines)
		var wg sync.WaitGroup

		for _, field := range db.Statement.Schema.Fields {
			wg.Add(1)
			go func(field *schema.Field) {
				defer wg.Done()
				guard <- struct{}{} // would block if guard channel is already filled
				defer func() {
					<-guard
				}()

				tagPassword := field.Tag.Get("password")
				if tagPassword != "true" {
					return
				}
				switch db.Statement.ReflectValue.Kind() {
				case reflect.Slice:
					if db.Statement.ReflectValue.Len() < maxGoroutines {
						maxGoroutines = db.Statement.ReflectValue.Len()
					}
					for i := 0; i < db.Statement.ReflectValue.Len(); i += maxGoroutines {

						valArr := db.Statement.ReflectValue.Index(i)
						password := valArr.FieldByName(field.Name).String()
						if password != "" && !domain.Password(password).IsHashed() {
							db.Statement.ReflectValue.Index(i).FieldByName(field.Name).SetString(domain.Password(password).Hash().String())
						}
					}

				case reflect.String, reflect.Struct:
					password := db.Statement.ReflectValue.FieldByName(field.Name).String()
					if password != "" && !domain.Password(password).IsHashed() {
						db.Statement.ReflectValue.FieldByName(field.Name).SetString(domain.Password(password).Hash().String())
					}
				}
			}(field)
		}
		wg.Wait()
	}
}
