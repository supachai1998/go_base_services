package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_base/domain"
	helper "go_base/domain/helper"
	"go_base/logger"
	"go_base/storage"
	"go_base/xerror"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jtolds/gls"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stoewer/go-strcase"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Store represents a data store.
type Store struct {
	DB         *gorm.DB
	cache      *storage.Cache
	allStorage *storage.AllStorage
}

// NewStore creates a new store.
func NewStore(DB *gorm.DB, cache *storage.Cache, allStorage *storage.AllStorage) *Store {
	return &Store{DB: DB, cache: cache, allStorage: allStorage}
}

/*
Base

	T = base in pg init db
	U = Update domain
	C = Create domain if we want different form T
*/
type BaseStore[T, U, C any] struct {
	DB         *gorm.DB
	cache      *storage.Cache
	cfg        *BaseStoreConfig
	allStorage *storage.AllStorage
}

var (
	CreateLog = "create"
	UpdateLog = "update"
	DeleteLog = "delete"
)

type BaseStoreConfig struct {
	WriteChangelog bool
	CacheExpire    time.Duration
}

var (
	groupCache = "group_cache_%s"
)

type GroupTypeCount struct {
	Field string
	Count int64
}

// new base on store
func NewBaseStore[T, U, C any](DB *gorm.DB, cfg *BaseStoreConfig, allStorage *storage.AllStorage) *BaseStore[T, U, C] {
	cache := allStorage.Cache
	if cfg.WriteChangelog {
		// check if table exists
		model := domain.NewLogs[T]()
		tableName := model.TableName()
		if tableName != "" {
			tx := DB.WithContext(context.Background())
			tx = tx.Table(tableName)
			if err := tx.AutoMigrate(model); err != nil {
				panic(err)
			}
		}

	}

	return &BaseStore[T, U, C]{DB: DB, cfg: cfg, cache: cache, allStorage: allStorage}
}

// find base on store
func (s *BaseStore[T, U, C]) Find(ctx echo.Context, pagination domain.Pagination[T], ignoreRelations ...string) (*domain.Pagination[T], error) {
	iDB := s.DB.WithContext(ctx.Request().Context())
	if len(ignoreRelations) > 0 {
		iDB = iDB.Omit(ignoreRelations...)
	}

	result, err := pagination.Paginate(ctx, iDB)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// get base on store
func (s *BaseStore[T, U, C]) GetBy(ctx echo.Context, field any, value any) (*T, error) {
	var result T
	filedName, err := helper.GetFieldNameByField(result, field)
	if err != nil {
		return nil, err
	}
	// reflect type if pointer
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}

	if err := s.DB.WithContext(ctx.Request().Context()).Where(filedName+" = ?", value).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

// get by id base on store
func (s *BaseStore[T, U, C]) GetByID(ctx echo.Context, idStr string) (*T, error) {
	id, idUUID := domain.GetUUID(idStr)
	if idUUID == uuid.Nil {
		return nil, xerror.EInvalidParameter(nil)
	}

	var result T

	if err := s.DB.WithContext(ctx.Request().Context()).Preload(clause.Associations).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// get by key base on store
func (s *BaseStore[T, U, C]) GetByKey(ctx echo.Context, key string, value string) (*T, error) {
	var result T

	if err := s.DB.WithContext(ctx.Request().Context()).Where(key+" = ?", value).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// create base on store
func (s *BaseStore[T, U, C]) Create(ctx echo.Context, model *T, typeLog ...string) error {
	err := s.DB.WithContext(ctx.Request().Context()).Create(model).Error
	if err != nil {
		return err
	}

	if s.cfg.WriteChangelog {
		log := CreateLog
		if len(typeLog) > 0 {
			log = typeLog[0]
		}
		if err := s.WriteLog(ctx, model, log); err != nil {
			return err
		}
	}

	return nil
}
func (s *BaseStore[T, U, C]) CreateC(ctx echo.Context, model *C, typeLog ...string) error {
	if err := s.DB.WithContext(ctx.Request().Context()).Create(model).Error; err != nil {
		return err
	}

	if s.cfg.WriteChangelog {
		log := CreateLog
		if len(typeLog) > 0 {
			log = typeLog[0]
		}
		if err := s.WriteLog(ctx, model, log); err != nil {
			return err
		}
	}

	return nil
}

// update base on store
func (s *BaseStore[T, U, C]) Update(ctx echo.Context, model *T, typeLog ...string) error {
	err := s.DB.WithContext(ctx.Request().Context()).Updates(model).Error
	if err != nil {
		return err
	}

	if s.cfg.WriteChangelog {
		log := UpdateLog
		if len(typeLog) > 0 {
			log = typeLog[0]
		}
		if err := s.updateLog(ctx, model, log); err != nil {
			return err
		}
	}
	return nil
}

func (s *BaseStore[T, U, C]) UpdateU(ctx echo.Context, model *U, typeLog ...string) error {
	err := s.DB.WithContext(ctx.Request().Context()).Updates(model).Error
	if err != nil {
		return err
	}

	if s.cfg.WriteChangelog {
		log := UpdateLog
		if len(typeLog) > 0 {
			log = typeLog[0]
		}
		if err := s.updateULog(ctx, model, log); err != nil {
			return err
		}
	}
	return nil
}

// update one field base on store
func (s *BaseStore[T, U, C]) UpdateOne(ctx echo.Context, model *T, field any, typeLog ...string) error {
	filedName, err := helper.GetFieldNameByField(model, field)
	if err != nil {
		return err
	}

	value, err := helper.GetValueFromStructByFieldName(model, filedName)
	if err != nil {
		return err
	}
	if value == nil {
		logger.L().Panic("dev error: value is nil")
		return nil
	}
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}
	err = s.DB.WithContext(ctx.Request().Context()).Model(model).Update(filedName, value).Error
	if err != nil {
		return err
	}
	if s.cfg.WriteChangelog {
		log := UpdateLog
		if len(typeLog) > 0 {
			log = typeLog[0]
		}
		if err := s.updateLog(ctx, model, log); err != nil {
			return err
		}
	}
	return nil
}

// update where id
func (s *BaseStore[T, U, C]) UpdateWhereID(ctx echo.Context, model *T, idStr string, typeLog ...string) error {
	id, idUUID := domain.GetUUID(idStr)
	if idUUID == uuid.Nil {
		return xerror.EInvalidParameter(nil)
	}

	err := s.DB.WithContext(ctx.Request().Context()).Model(model).Where("id = ?", id).Updates(model).Error
	if err != nil {
		return err
	}
	if s.cfg.WriteChangelog {
		log := UpdateLog
		if len(typeLog) > 0 {
			log = typeLog[0]
		}
		if err := s.updateLog(ctx, model, log); err != nil {
			return err
		}
	}
	return nil
}

// delete base on store
func (s *BaseStore[T, U, C]) Delete(ctx echo.Context, id uuid.UUID) error {
	var model T

	if err := s.DB.WithContext(ctx.Request().Context()).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if err := s.DB.WithContext(ctx.Request().Context()).Delete(&model).Error; err != nil {
		return err
	}
	if s.cfg.WriteChangelog {
		if err := s.deleteLog(ctx, &model); err != nil {
			return err
		}
	}
	return nil
}

// delete if exist base on store
//  1. find old model
//  2. if found is soft delete is exist, then delete
//  3. if not found, then return nil
func (s *BaseStore[T, U, C]) DeleteIfExist(ctx echo.Context, model *T, field any, value any, typeLog ...string) error {
	filedName, err := helper.GetFieldNameByField(model, field)
	if err != nil {
		return err
	}
	var isSoftDelete *gorm.DeletedAt
	var newModel T

	if err := s.DB.WithContext(ctx.Request().Context()).Unscoped().Model(newModel).Select("deleted_at").Where(filedName+" = ?", value).Scan(&isSoftDelete).Error; err != nil {
		return err
	}
	if isSoftDelete != nil && isSoftDelete.Valid {
		if err := s.DB.Unscoped().Model(newModel).Where(filedName+" = ?", value).Delete(newModel).Error; err != nil {
			return err
		}
		if s.cfg.WriteChangelog {
			log := DeleteLog
			if len(typeLog) > 0 {
				log = typeLog[0]
			}
			if err := s.deleteLog(ctx, model, log); err != nil {
				return err
			}
		}
		return nil
	}

	return nil
}

// write log base on store
func (s *BaseStore[T, U, C]) WriteLog(ctx echo.Context, _model any, action string) error {
	gls.Go(func() {
		actionFromCtx := domain.GetActionFromContext(ctx)
		if action == "" {
			action = actionFromCtx
		}
		model, err := convertAnyIntoJSONType(_model)
		domain.ErrLogGlsGo(ctx, err)
		_doer := s.getDoer(ctx)
		doer, err := convertAnyIntoJSONType(_doer)
		domain.ErrLogGlsGo(ctx, err)

		switch _model.(type) {
		case *T:
			log := domain.NewLogs[T]()
			log.Action = action
			log.Model = model
			log.Doer = doer
			if err := s.DB.Create(log).Error; err != nil {
				domain.ErrLogGlsGo(ctx, err)
			}
		case *U:
			log := domain.NewLogs[U]()
			log.Action = action
			log.Model = model
			log.Doer = doer
			if err := s.DB.Create(log).Error; err != nil {
				domain.ErrLogGlsGo(ctx, err)
			}
		case *C:
			log := domain.NewLogs[C]()
			log.Action = action
			log.Model = model
			log.Doer = doer
			if err := s.DB.Create(log).Error; err != nil {
				domain.ErrLogGlsGo(ctx, err)
			}
		default:
			domain.ErrLogGlsGo(ctx, errors.New("model type not found"+reflect.TypeOf(_model).String()))
		}

		var m T
		// table name
		nameOfModel := strcase.SnakeCase(reflect.TypeOf(m).Name())
		if strings.Contains(nameOfModel, "_") {
			// remove last index
			nameOfModel = nameOfModel[:strings.LastIndex(nameOfModel, "_")]
		}
		// ถ้า model นั้นเป็น staff หรือ user จะไม่เขียน logs ไปที่ผู้กระทำซ้ำ

		// เขียน logs ไปที่ผู้กระทำ (จะอยู่ในถัง logs ของผู้กระทำ staff | user)
		if staffCtx := domain.StaffFromContext(ctx); !strings.Contains(nameOfModel, "staff") && staffCtx != nil {

			var logStaff domain.Logs[domain.Staff]
			logStaff.Action = action
			logStaff.FromTable = &nameOfModel
			logStaff.Model = model
			logStaff.Doer = doer
			if err := s.DB.Create(&logStaff).Error; err != nil {
				domain.ErrLogGlsGo(ctx, err)
			}
		} else if userCtx := domain.UserFromContext(ctx); !strings.Contains(nameOfModel, "user") && userCtx != nil {
			var logUser domain.Logs[domain.User]
			logUser.Action = action
			logUser.FromTable = &nameOfModel
			logUser.Model = model
			logUser.Doer = doer
			if err := s.DB.Create(&logUser).Error; err != nil {
				domain.ErrLogGlsGo(ctx, err)
			}
		}

	})
	return nil
}
func (s *BaseStore[T, U, C]) toLogsT(ctx echo.Context, _model T, action string) (*domain.Logs[T], *datatypes.JSON, *datatypes.JSON, error) {
	log := domain.NewLogs[T]()
	actionFromCtx := domain.GetActionFromContext(ctx)
	log.Action = action
	if log.Action == "" {
		log.Action = actionFromCtx
	}
	model, err := convertAnyIntoJSONType(_model)
	if err != nil {
		return nil, nil, nil, err
	}
	log.Model = model
	_doer := s.getDoer(ctx)
	doer, err := convertAnyIntoJSONType(_doer)
	if err != nil {
		return nil, nil, nil, err
	}

	log.Doer = doer
	log.CreatedAt = domain.TimeNow()
	log.UpdatedAt = domain.TimeNow()

	return log, &model, &doer, nil
}
func (s *BaseStore[T, U, C]) toLogsU(ctx echo.Context, _model U, action string) (*domain.Logs[U], *datatypes.JSON, *datatypes.JSON, error) {
	log := domain.NewLogs[U]()
	actionFromCtx := domain.GetActionFromContext(ctx)
	log.Action = action
	if log.Action == "" {
		log.Action = actionFromCtx
	}
	model, err := convertAnyIntoJSONType(_model)
	if err != nil {
		return nil, nil, nil, err
	}
	log.Model = model
	_doer := s.getDoer(ctx)
	doer, err := convertAnyIntoJSONType(_doer)
	if err != nil {
		return nil, nil, nil, err
	}

	log.Doer = doer
	log.CreatedAt = domain.TimeNow()
	log.UpdatedAt = domain.TimeNow()

	return log, &model, &doer, nil
}
func (s *BaseStore[T, U, C]) toLogsC(ctx echo.Context, _model C, action string) (*domain.Logs[C], *datatypes.JSON, *datatypes.JSON, error) {
	log := domain.NewLogs[C]()
	actionFromCtx := domain.GetActionFromContext(ctx)
	log.Action = action
	if log.Action == "" {
		log.Action = actionFromCtx
	}
	model, err := convertAnyIntoJSONType(_model)
	if err != nil {
		return nil, nil, nil, err
	}
	log.Model = model
	_doer := s.getDoer(ctx)
	doer, err := convertAnyIntoJSONType(_doer)
	if err != nil {
		return nil, nil, nil, err
	}

	log.Doer = doer
	log.CreatedAt = domain.TimeNow()
	log.UpdatedAt = domain.TimeNow()

	return log, &model, &doer, nil
}
func (s *BaseStore[T, U, C]) WriteLogs(ctx echo.Context, modelsT *[]T, modelsU *[]U, modelsC *[]C, action string) error {
	gls.Go(func() {
		if modelsT == nil && modelsU == nil && modelsC == nil {
			return
		}
		var logsT []domain.Logs[T]
		var logsU []domain.Logs[U]
		var logsC []domain.Logs[C]
		var logStaffs []domain.Logs[domain.Staff]
		var logUsers []domain.Logs[domain.User]
		var doers []datatypes.JSON
		var models []datatypes.JSON

		if modelsT != nil {
			for _, _model := range *modelsT {

				t, model, doer, err := s.toLogsT(ctx, _model, action)
				if err != nil {
					domain.ErrLogGlsGo(ctx, err)
				}
				logsT = append(logsT, *t)
				doers = append(doers, *doer)
				models = append(models, *model)

			}
		}

		if modelsU != nil {
			for _, _model := range *modelsU {

				t, model, doer, err := s.toLogsU(ctx, _model, action)
				if err != nil {
					domain.ErrLogGlsGo(ctx, err)
				}
				logsU = append(logsU, *t)
				doers = append(doers, *doer)
				models = append(models, *model)

			}
		}

		if modelsC != nil {
			for _, _model := range *modelsC {

				t, model, doer, err := s.toLogsC(ctx, _model, action)
				if err != nil {
					domain.ErrLogGlsGo(ctx, err)
				}
				logsC = append(logsC, *t)
				doers = append(doers, *doer)
				models = append(models, *model)

			}
		}

		for i, doer := range doers {

			var m T
			// table name
			nameOfModel := strcase.SnakeCase(reflect.TypeOf(m).Name())
			if strings.Contains(nameOfModel, "_") {
				// remove last index
				nameOfModel = nameOfModel[:strings.LastIndex(nameOfModel, "_")]
			}
			// ถ้า model นั้นเป็น staff หรือ user จะไม่เขียน logs ไปที่ผู้กระทำซ้ำ

			// เขียน logs ไปที่ผู้กระทำ (จะอยู่ในถัง logs ของผู้กระทำ staff | user)
			if staffCtx := domain.StaffFromContext(ctx); !strings.Contains(nameOfModel, "staff") && staffCtx != nil {
				var logStaff domain.Logs[domain.Staff]
				logStaff.Action = action
				logStaff.FromTable = &nameOfModel
				logStaff.Model = models[i]
				logStaff.Doer = doer
				logStaffs = append(logStaffs, logStaff)
			} else if userCtx := domain.UserFromContext(ctx); !strings.Contains(nameOfModel, "user") && userCtx != nil {
				var logUser domain.Logs[domain.User]
				logUser.Action = action
				logUser.FromTable = &nameOfModel
				logUser.Model = models[i]
				logUser.Doer = doer
				logUsers = append(logUsers, logUser)
			}
		}

		if len(logsT) > 0 {
			if err := s.DB.Create(logsT).Error; err != nil {
				domain.ErrLogGlsGo(ctx, err)
			}
		}
		if len(logStaffs) > 0 {
			if err := s.DB.Create(logStaffs).Error; err != nil {
				domain.ErrLogGlsGo(ctx, err)
			}
		}
		if len(logUsers) > 0 {
			if err := s.DB.Create(logUsers).Error; err != nil {
				domain.ErrLogGlsGo(ctx, err)
			}
		}
	})
	return nil
}

func (s *BaseStore[T, U, C]) getDoer(ctx echo.Context) domain.Doer {
	isUserKey, ok := ctx.Get(string(domain.IsUserKey)).(bool)
	var doer domain.Doer
	if ok && isUserKey {
		user := domain.UserFromContext(ctx)
		if user != nil {
			doer.ID = user.ID
			doer.Name = domain.UserGetName(user)
			doer.Email = string(user.Email)
			doer.Type = domain.DoerTypeUser
		} else {
			doer.Type = domain.DoerTypeSystem
		}
		return doer

	}
	staff := domain.StaffFromContext(ctx)
	if staff != nil {
		doer.ID = staff.ID
		doer.Name = domain.StaffGetName(staff)
		doer.Email = string(staff.Email)
		doer.Type = domain.DoerTypeStaff
		if staff.RoleID != nil {
			var r domain.Role
			doer.RoleID = staff.RoleID
			if err := s.DB.WithContext(ctx.Request().Context()).Model(staff).Association("Roles").Find(&r); err == nil {
				doer.Role = lo.ToPtr(r)
			}
		}
		return doer
	}
	doer.Type = domain.DoerTypeSystem

	return doer
}

// update log base on store
func (s *BaseStore[T, U, C]) updateLog(ctx echo.Context, model *T, typeLog ...string) error {
	if s.cfg.WriteChangelog {
		log := UpdateLog
		if len(typeLog) > 0 {
			log = typeLog[0]
		}
		if err := s.WriteLog(ctx, model, log); err != nil {
			return err
		}
	}
	return nil
}
func (s *BaseStore[T, U, C]) updateULog(ctx echo.Context, model *U, typeLog ...string) error {
	if s.cfg.WriteChangelog {
		log := UpdateLog
		if len(typeLog) > 0 {
			log = typeLog[0]
		}
		if err := s.WriteLog(ctx, model, log); err != nil {
			return err
		}
	}
	return nil
}

// delete log base on store
func (s *BaseStore[T, U, C]) deleteLog(ctx echo.Context, model *T, typeLog ...string) error {
	if s.cfg.WriteChangelog {
		log := DeleteLog
		if len(typeLog) > 0 {
			log = typeLog[0]
		}
		if err := s.WriteLog(ctx, model, log); err != nil {
			return err
		}
	}
	return nil
}

// convert any into json type base on store
func convertAnyIntoJSONType(value any) (datatypes.JSON, error) {
	var result datatypes.JSON
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		value = reflect.ValueOf(value).Elem().Interface()
	}
	marshal, err := json.Marshal(value)
	if err != nil {
		return result, err
	}
	return datatypes.JSON(marshal), nil
}

// FindByID base on store
func (s *BaseStore[T, U, C]) FindByID(ctx echo.Context) (*T, error) {
	var result T
	id, idUid := domain.GetUUIDFromParam(ctx, "id")
	if idUid == uuid.Nil {
		return nil, xerror.EInvalidParameter(nil)
	}

	if err := s.DB.WithContext(ctx.Request().Context()).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// query log jsonb base on store
func (s *BaseStore[T, U, C]) LogsQueryModel(ctx echo.Context, model *domain.Logs[T], fields string, value string) (*domain.Pagination[*domain.Logs[T]], error) {
	// "id": "ee3590d2-512d-4111-8142-63651d70c34a"
	// field = "model->>'" + "id" + "'"
	// val := "ee3590d2-512d-4111-8142-63651d70c34a"
	if !strings.Contains(fields, ".") {
		logger.L().Errorw("dev error : fields should be model.key like doer.id")
	}
	var pagination domain.Pagination[*domain.Logs[T]]
	fieldsArr := strings.Split(fields, ".")
	if len(fieldsArr) != 2 {
		logger.L().Errorw("dev error : fields should be model.key like doer.id")
	}
	field, key := fieldsArr[0], fieldsArr[1]
	wherCon := fmt.Sprintf("%s->>'%s' = ?", field, key)
	DB := s.DB.WithContext(ctx.Request().Context()).
		Model(model).
		Where(wherCon, value)
	return pagination.Paginate(ctx, DB)
}

// DeleteIds
func (s *BaseStore[T, U, C]) DeleteIds(ctx echo.Context, ids domain.Ids) error {
	var models []T
	var count int64
	if err := s.DB.WithContext(ctx.Request().Context()).Where("id in ?", ids.IDs).Find(&models).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	if err := s.DB.WithContext(ctx.Request().Context()).Where("id in ?", ids.IDs).Delete(&models).Error; err != nil {
		return err
	}
	// write log
	if err := s.WriteLogs(ctx, &models, nil, nil, DeleteLog); err != nil {
		return err
	}
	return nil
}

func (s *BaseStore[T, U, C]) getCache(ctx echo.Context, key string) (*map[string]int64, error) {
	if s.cache != nil {
		var g map[string]int64
		if val, err := s.cache.GetCache(ctx.Request().Context(), key); err == nil {
			if err := json.Unmarshal([]byte(val), &g); err == nil {
				return &g, nil
			}
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *BaseStore[T, U, C]) setCache(ctx echo.Context, key string, value map[string]int64) error {
	if s.cache != nil {
		val, err := json.Marshal(value)
		if err != nil {
			return err
		}
		if err := s.cache.SetCache(ctx.Request().Context(), key, string(val), s.cfg.CacheExpire); err != nil {
			return err
		}
	}
	return nil
}

// Count group in array jsonb
func (s *BaseStore[T, U, C]) CountJsonGroup(ctx echo.Context, fieldName string) (*map[string]int64, error) {
	if cache, err := s.getCache(ctx, fmt.Sprintf(groupCache, fieldName)); err == nil {
		return cache, nil
	}

	var count []GroupTypeCount
	var model T
	selectStatement := fmt.Sprintf("jsonb_array_elements_text(%s) as field, count(*)", fieldName)
	if err := s.DB.WithContext(ctx.Request().Context()).Model(&model).Select(selectStatement).Group("field").Scan(&count).Error; err != nil {
		return nil, err
	}

	g := map[string]int64{}
	for _, c := range count {
		g[c.Field] = c.Count
	}

	if err := s.setCache(ctx, fmt.Sprintf(groupCache, fieldName), g); err != nil {
		return nil, err
	}
	return &g, nil
}

// user
// find with user id base on store
func (s *BaseStore[T, U, C]) FindWithUserID(ctx echo.Context, pagination domain.Pagination[T], ignoreRelations ...string) (*domain.Pagination[T], error) {
	user := domain.UserFromContext(ctx)
	if user == nil {
		return nil, xerror.EForbidden()
	}
	iDB := s.DB.WithContext(ctx.Request().Context())
	if len(ignoreRelations) > 0 {
		iDB = iDB.Omit(ignoreRelations...)
	}
	iDB = iDB.Scopes(domain.WithUserID(user.ID))
	result, err := pagination.Paginate(ctx, iDB)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// get by user id base on store
func (s *BaseStore[T, U, C]) GetWithUserID(ctx echo.Context, idStr string) (*T, error) {
	user := domain.UserFromContext(ctx)
	if user == nil {
		return nil, xerror.EForbidden()
	}
	id, idUUID := domain.GetUUID(idStr)
	if idUUID == uuid.Nil {
		return nil, xerror.EInvalidParameter(nil)
	}
	var result T
	if err := s.DB.WithContext(ctx.Request().Context()).Scopes(domain.WithUserID(user.ID)).Where("id = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

// delete with user id base on store
func (s *BaseStore[T, U, C]) DeleteWithUserID(ctx echo.Context, id uuid.UUID) error {
	user := domain.UserFromContext(ctx)
	if user == nil {
		return xerror.EForbidden()
	}
	var model T
	if err := s.DB.WithContext(ctx.Request().Context()).Scopes(domain.WithUserID(user.ID)).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if err := s.DB.WithContext(ctx.Request().Context()).Scopes(domain.WithUserID(user.ID)).Delete(&model).Error; err != nil {
		return err
	}
	if s.cfg.WriteChangelog {
		if err := s.deleteLog(ctx, &model); err != nil {
			return err
		}
	}
	return nil
}

// delete if exist with user id base on store
func (s *BaseStore[T, U, C]) DeleteIfExistWithUserID(ctx echo.Context, model *T, field any, value any, typeLog ...string) error {
	user := domain.UserFromContext(ctx)
	if user == nil {
		return xerror.EForbidden()
	}
	filedName, err := helper.GetFieldNameByField(model, field)
	if err != nil {
		return err
	}
	var isSoftDelete *gorm.DeletedAt
	var newModel T
	if err := s.DB.WithContext(ctx.Request().Context()).Unscoped().Model(newModel).Select("deleted_at").Where(filedName+" = ?", value).Scan(&isSoftDelete).Error; err != nil {
		return err
	}
	if isSoftDelete != nil && isSoftDelete.Valid {
		if err := s.DB.Unscoped().Model(newModel).Where(filedName+" = ?", value).Delete(newModel).Error; err != nil {
			return err
		}
		if s.cfg.WriteChangelog {
			log := DeleteLog
			if len(typeLog) > 0 {
				log = typeLog[0]
			}
			if err := s.deleteLog(ctx, model, log); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

// update with user id base on store
func (s *BaseStore[T, U, C]) UpdateWithUserID(ctx echo.Context, model *U, typeLog ...string) error {
	user := domain.UserFromContext(ctx)
	if user == nil {
		return xerror.EForbidden()
	}
	err := s.DB.WithContext(ctx.Request().Context()).Scopes(domain.WithUserID(user.ID)).Updates(model).Error
	if err != nil {
		return err
	}
	if s.cfg.WriteChangelog {
		log := UpdateLog
		if len(typeLog) > 0 {
			log = typeLog[0]
		}
		if err := s.updateULog(ctx, model, log); err != nil {
			return err
		}
	}
	return nil
}
