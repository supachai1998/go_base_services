package domain

import (
	"fmt"
	helper "go_base/domain/helper"
	"go_base/logger"
	"go_base/xerror"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stoewer/go-strcase"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaginationSwagger struct {
	// query params
	Page     *int `query:"page" swagger:"default=1" json:"page" validate:"omitempty,gt=0"`
	PageSize *int `query:"limit" json:"limit" validate:"omitempty,gt=0"`
	// ex : created_at,desc|updated_at,asc
	Sort *string `query:"sort" swagger:"desc(desc)" json:"-" validate:"omitempty,excludes=."`

	SortArray []string `query:"sort[]" json:"-" validate:"omitempty,dive,excludes=."`

	// ex : name:like:john|email:like:john|age:gt:18
	Search *string `query:"search" swagger:"desc(name:like:john|email:like:john|age:gt:18)" json:"-" validate:"omitempty,excludesrune=;"`

	SearchArray []string `query:"search[]" json:"-" validate:"omitempty,dive,excludesrune=;"`

	Find         *string  `query:"find" json:"-"`
	Finds        []string `query:"find[]" json:"-"`
	OperatorFind *string  `query:"operator_find" json:"-" validate:"omitempty,oneof=or and"`

	NoLimit bool `query:"-" json:"-"`
}
type Pagination[T any] struct {
	PaginationSwagger

	// response
	TotalCount int `json:"total_count"`
	TotalPage  int `json:"total_page"`
	Items      []T `json:"items"`

	// meta optional
	MetaCount any `json:"meta_count,omitempty"`
}

var (
	limitPerPage = 10
)

func canParseTypeInStatementSQL(_type string, value string) bool {
	// case if found uuid
	if strings.Contains(_type, "uuid") {
		if _, err := uuid.Parse(value); err != nil {
			return false
		}
	}
	// number
	if strings.Contains(_type, "int") {
		if _, err := strconv.Atoi(value); err != nil {
			return false
		}
	}
	// float
	if strings.Contains(_type, "float") {
		if _, err := strconv.ParseFloat(value, 64); err != nil {
			return false
		}
	}

	if strings.Contains(_type, "datatypes.JSON") {
		return (helper.IsArray(value) || helper.IsObject(value))
	}

	return true
}

func (p Pagination[T]) Offset(tx *gorm.DB) (Pagination[T], *gorm.DB) {
	if p.Page == nil {
		p.Page = lo.ToPtr(int(1))
	}
	if p.PageSize == nil {
		p.PageSize = lo.ToPtr(int(limitPerPage))
	}
	if p.Page != nil && *p.Page < 1 {
		p.Page = lo.ToPtr(int(1))
	}
	if *p.PageSize < 0 {
		return p, tx
	}
	return p, tx.Offset((*p.Page - 1) * *p.PageSize)
}

func (p Pagination[T]) Limit(tx *gorm.DB) (Pagination[T], *gorm.DB) {
	if p.PageSize == nil {
		p.PageSize = lo.ToPtr(int(limitPerPage))
	}
	if p.NoLimit {
		return p, tx.Limit(-1)
	}
	return p, tx.Limit(*p.PageSize)
}

func (p Pagination[T]) SortBy(tx *gorm.DB) (Pagination[T], *gorm.DB) {
	if p.Sort == nil || lo.IsEmpty(p.Sort) {
		p.Sort = lo.ToPtr("created_at,desc")
	}
	order := *p.Sort
	if p.SortArray != nil && len(p.SortArray) > 0 {
		order = strings.Join(p.SortArray, "|")
	}
	orders := strings.Split(order, "|")
	if len(orders) == 1 {
		oo := strings.Split(order, ",")
		if len(oo) == 1 { // input is asc, desc
			if !sortDirections[oo[0]] {
				logger.L().Warn("invalid sort direction", zap.String("direction", oo[0]))
				return p, tx
			}
			if len(oo) == 1 {
				tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: oo[0] == "desc"})
			}
		} else if len(oo) == 2 { // input is created_at,asc
			if !sortDirections[oo[1]] {
				logger.L().Warn("invalid sort direction", zap.String("direction", oo[1]))
				return p, tx
			}
			tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: oo[0]}, Desc: oo[1] == "desc"})
		}
		return p, tx
	}

	for _, o := range orders { // input is created_at,asc;updated_at,desc
		if o == "" {
			continue
		}
		oo := strings.Split(o, ",")
		if len(oo) != 2 {
			logger.L().Warn("invalid sort len not 2", zap.String("sort", order))
			continue
		}
		field, order := oo[0], oo[1]
		order = strings.ToLower(order)

		tx = tx.Order(clause.OrderByColumn{Column: clause.Column{Name: field}, Desc: order == "desc"})

	}

	return p, tx
}

var operators map[string]bool = map[string]bool{
	"eq":         true,
	"neq":        true,
	"gt":         true,
	"gte":        true,
	"lt":         true,
	"lte":        true,
	"like":       true,
	"in":         true,
	"nnull":      true,
	"is_deleted": true,
}

var sortDirections map[string]bool = map[string]bool{
	"asc":  true,
	"desc": true,
}

func operatorCase(operator string) (string, error) {
	if !operators[operator] {
		logger.L().Warn("invalid operator", zap.String("operator", operator))
		return "", xerror.ErrInvalidOperator(operators)
	}
	switch operator {
	case "eq":
		return "=", nil
	case "neq":
		return "<>", nil
	case "nnull":
		return "nnull", nil
	case "gt":
		return ">", nil
	case "gte":
		return ">=", nil
	case "lt":
		return "<", nil
	case "lte":
		return "<=", nil
	case "like":
		return "like", nil
	case "in":
		return "in", nil
	case "is_deleted":
		return "is_deleted", nil
	}
	return "", xerror.ErrInvalidOperator(operators)
}

var walkOP map[string]bool = map[string]bool{
	"|": true,
	"&": true,
}

// email,like,admin.com|email,like,dedecms.com&name,like,admin
// return conditions, error
func walkStrFindAndOrOperators(str string) ([]string, error) {
	var conditions []string
	if str == "" {
		return conditions, nil
	}
	for i := 0; i < len(str); i++ {
		if walkOP[string(str[i])] {
			con := fmt.Sprintf("%s,%s", str[:i], str[i:i+1])
			conditions = append(conditions, con) // -> [email,like,admin.com,|]
			str = str[i+1:]                      // -> split [email,like,admin.com] out --> email,like,dedecms.com&name,like,admin
			i = 0
		}
	}
	if str != "" {
		conditions = append(conditions, str)
	}
	return conditions, nil // -> [email,like,admin.com,| email,like,dedecms.com name,like,admin]
}

func (p Pagination[T]) SearchBy(tx *gorm.DB) (Pagination[T], *gorm.DB, error) {
	if p.SearchArray != nil && len(p.SearchArray) > 0 {
		p.Search = lo.ToPtr(strings.Join(p.SearchArray, "&"))
	}
	if p.Search == nil || lo.IsEmpty(p.Search) {
		return p, tx, nil
	}
	_tx := tx.Session(&gorm.Session{QueryFields: true})
	search := *p.Search
	// fmt.Println("search", search)
	conditions, _ := walkStrFindAndOrOperators(search)
	var _opOr bool

	var model T
	modelName := helper.ToTableName(reflect.TypeOf(model).Name())
	for _, c := range conditions {
		if c == "" {
			continue
		}
		cc := strings.Split(c, ",")
		// fmt.Println("cc", cc)
		// if cc is 4 then last index is tx.Or | or tx.Where & in next loop
		if len(cc) >= 3 && len(cc) <= 4 {
			field, operator, value := cc[0], cc[1], cc[2] // -> email,like,admin.com
			where, err := FilterFieldOperator(field, operator)
			if err != nil {
				return p, tx, err
			}
			split := strings.Split(where, " ") // -> [email like ?]
			// fmt.Println("split", split, "value", value)
			f, op, _ := split[0], split[1], split[2] // -> email,like,?

			if strings.Contains(f, ".") {
				tableNameRelation := strings.Split(f, ".")[0]
				_field := strings.Split(f, ".")[1]
				if modelName == helper.ToTableName(tableNameRelation) {
					f = fmt.Sprintf(`%s.%s`, modelName, _field) // -> table.email
				} else {
					f = fmt.Sprintf(`"%s".%s`, strcase.UpperCamelCase(strings.Split(f, ".")[0]), strings.Split(f, ".")[1])
				}
			}

			if op == "like" {
				value = "%" + value + "%"
			}
			if _opOr {
				_tx = _tx.Or(fmt.Sprintf("%s %s ?", f, op), value)
			}
			switch op {
			case "in":
				_tx = _tx.Where(fmt.Sprintf("%s @> ?", f), value)
			case "nnull":
				_tx = _tx.Where(fmt.Sprintf("%s IS NOT NULL", f))
			case "is_deleted":
				_tx = _tx.Unscoped().Where(fmt.Sprintf("%s IS NOT NULL", f))
			default:
				_tx = _tx.Where(fmt.Sprintf("%s %s ?", f, op), value)
			}
			if len(cc) == 4 {
				_, _, _, found := cc[0], cc[1], cc[2], cc[3]
				switch found {
				case "|":
					_opOr = true
				case "&":
					_opOr = false
				default:
					_opOr = false
				}
			} else {
				_opOr = false
			}

		}
	}

	return p, _tx, nil
}

func (p Pagination[T]) SearchFilter(tx *gorm.DB) (Pagination[T], *gorm.DB, error) {
	if (p.Find == nil || lo.IsEmpty(p.Find)) && (p.Finds == nil || len(p.Finds) == 0) {
		return p, tx, nil
	}
	var finds []string
	if p.Find != nil {
		find := *p.Find // value
		finds = append(finds, find)
	}
	if len(p.Finds) > 0 {
		finds = append(finds, p.Finds...)
	}
	_tx := tx.Session(&gorm.Session{QueryFields: true})
	var model T
	t := reflect.ValueOf(&model).Elem()
	modelName := helper.ToTableName(reflect.TypeOf(model).Name())
	for i := 0; i < t.NumField(); i++ {

		// field operator value
		for _, find := range finds {
			field := t.Type().Field(i).Name
			if field == "" {
				continue
			}
			_type := t.Type().Field(i).Type.String()
			filter := t.Type().Field(i).Tag.Get("filter")
			if filter == "" {
				continue
			}

			field = strcase.SnakeCase(field)
			if !canParseTypeInStatementSQL(_type, find) {
				continue
			}

			if strings.Contains(filter, ".") {
				ops := strings.Split(filter, ".") // -> [table, email, like]
				if len(ops) != 3 {
					logger.L().Warn("invalid filter", zap.String("filter", filter))
					continue
				}
				tableNameRelation, _field, _filter := ops[0], ops[1], ops[2]
				if modelName == tableNameRelation {
					field = fmt.Sprintf(`%s.%s`, modelName, _field) // -> table.email
					filter = _filter
				} else {
					field = fmt.Sprintf(`"%s".%s`, strcase.UpperCamelCase(helper.ToTableName(tableNameRelation, false)), _field) // -> "Table".email
					filter = _filter                                                                                             // -> like
				}
			}

			switch filter {
			case "in":
				if helper.IsArray(find) {
					_tx = p.txOperatorFindAndOr(_tx, fmt.Sprintf("%s @> ?", field), find)
					continue
				}
			case "like":
				_tx = p.txOperatorFindAndOr(_tx, fmt.Sprintf("%s %s ?", field, filter), "%"+find+"%")
				continue
			}
			_tx = p.txOperatorFindAndOr(_tx, fmt.Sprintf("%s %s ?", field, filter), find)

		}
	}
	return p, _tx, nil
}
func (p Pagination[T]) txOperatorFindAndOr(tx *gorm.DB, query interface{}, args ...interface{}) *gorm.DB {
	if p.OperatorFind == nil || lo.IsEmpty(p.OperatorFind) {
		return tx.Or(query, args...)
	}

	if *p.OperatorFind == "or" {
		return tx.Or(query, args...)
	}
	return tx.Where(query, args...)
}
func (p Pagination[T]) Paginate(ctx echo.Context, db *gorm.DB, isJoinArgs ...bool) (*Pagination[T], error) {
	var items []T
	var count int64
	var err error
	isJoin := true
	if len(isJoinArgs) > 0 {
		isJoin = isJoinArgs[0]
	}

	ctx.Bind(&p)

	if err := ctx.Validate(&p); err != nil {
		return nil, err
	}
	var model T

	t := reflect.ValueOf(&model).Elem()

	db = db.Session(&gorm.Session{QueryFields: true})
	// ถ้าเป็น relation ให้เช็คว่า domain แล้ว join กับ relation นั้น เพื่อค้นหา
	if isJoin && reflect.TypeOf(model).Kind() == reflect.Struct {
		fmt.Println(t)
		for i := 0; i < t.NumField(); i++ {
			fieldName := t.Type().Field(i).Name
			if lo.Contains(TableNames, fieldName) {
				db = db.Joins(fieldName)
			}
		}
	}

	isFind := (p.Find != nil && !lo.IsEmpty(p.Find)) || (p.Finds != nil && len(p.Finds) > 0)
	isSearch := (p.Search != nil && !lo.IsEmpty(p.Search)) || (p.SearchArray != nil && len(p.SearchArray) > 0)

	// ในกรณีที่มีทั้ง find และ search ให้เป็น And เช่น
	if isFind && isSearch {
		_, dbSearchFilter, err := p.SearchFilter(db)
		if err != nil {
			return nil, err
		}

		_, dbSearchBy, err := p.SearchBy(db)
		if err != nil {
			return nil, err
		}
		db = db.Where(dbSearchFilter, dbSearchBy)
	} else {
		p, db, err = p.SearchBy(db)
		if err != nil {
			return nil, err
		}
		p, db, err = p.SearchFilter(db)
		if err != nil {
			return nil, err
		}
	}
	if len(db.Statement.Omits) == 0 {
		db = db.Preload(clause.Associations)
	}
	if err := db.Model(&items).Count(&count).Error; err != nil {
		return nil, err
	}

	p, db = p.Offset(db)
	p, db = p.Limit(db)
	p, db = p.SortBy(db)

	p.TotalCount = int(count)
	p.TotalPage = int(count) / *p.PageSize
	if *p.PageSize < 0 {
		p.TotalPage = 1
		*p.PageSize = int(count)
	}
	if int(count)%*p.PageSize != 0 {
		p.TotalPage++
	}

	if err := db.Find(&items).Error; err != nil {
		return nil, xerror.E(err).SetDebugInfo("pagination", p)
	}
	p.Items = items
	return &p, nil
}

func PaginationFromCtx[T any](ctx echo.Context) Pagination[T] {
	var p PaginationSwagger
	if err := ctx.Bind(&p); err != nil {
		return Pagination[T]{}
	}

	return Pagination[T]{
		PaginationSwagger: p,
	}
}
