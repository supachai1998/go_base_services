// postgresql filter

package domain

import (
	"go_base/logger"
	"go_base/xerror"
	"reflect"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Filter[T any] struct {
	PaginationSwagger

	FieldFilter T
}

func FilterOptions(filter any) (tx *gorm.DB, err error) {
	v := reflect.ValueOf(filter)
	t := reflect.TypeOf(filter)
	if v.Kind() == reflect.Pointer {
		st := v.Elem().Interface()
		v = reflect.ValueOf(st)
		t = reflect.TypeOf(st)
	}
	if v.Kind() != reflect.Struct {
		return nil, xerror.EInternalError().SetMessage("filter must be struct")
	}
	for i := 0; i < v.NumField(); i++ {
		fieldVal := v.Field(i)

		if !fieldVal.IsValid() {
			continue
		}

		fieldType := t.Field(i)
		if fieldType.Anonymous {
			continue
		}

		if fieldType.Type.Kind() == reflect.Struct {
			tx, _ = FilterOptions(fieldVal.Interface())
			continue
		}

		if fieldType.Type.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				continue
			}
			fieldVal = fieldVal.Elem()
		}

		if !fieldVal.IsValid() {
			continue
		}

		filter := t.Field(i).Tag.Get("filter")
		if filter == "" || filter == "-" {
			continue
		}
		// if we found . in filter tag, it means we want to filter by relation
		if strings.Contains(filter, ".") {
			splitRelations := strings.Split(filter, ".")
			if len(splitRelations)%2 != 0 {
				logger.L().Error("invalid filter", zap.String("filter", filter))
				return nil, xerror.EInternalError().SetMessage("invalid filter")
			}
			for i := 0; i < len(splitRelations); i += 2 {
				relation := splitRelations[i]
				operator := splitRelations[i+1]
				if !operators[operator] {
					logger.L().Error("invalid operator", zap.String("operator", operator))
					return nil, xerror.ErrInvalidOperator(operators)
				}
				operator, err := operatorCase(operator)
				if err != nil {
					logger.L().Error("invalid operator", zap.String("operator", operator))
					return nil, err
				}
				tx = tx.Joins(relation).Where(relation+".id "+operator+" ?", fieldVal.Interface())
			}
		}
		allowzero := false
		if idx := strings.LastIndex(filter, ","); idx > 0 {
			allowzero = filter[idx+1:] == "allowzero"
			filter = filter[:idx]
		}

		if !allowzero && fieldVal.IsZero() {
			continue
		}
		field, operator := reflect.TypeOf(v.Elem().Interface()).Field(i).Name, filter

		where, err := FilterFieldOperator(field, operator)
		if err != nil {
			logger.L().Error("FilterFieldOperatorValue", zap.Error(err))
			return nil, err
		}

		tx = tx.Where(where, fieldVal.Interface())

	}
	return tx, nil

}

func FilterFieldOperator(filed, operator string) (string, error) {
	if !operators[operator] {
		logger.L().Warn("invalid operator", zap.String("operator", operator))
		return "", xerror.ErrInvalidOperator(operators)
	}
	operator, err := operatorCase(operator)
	if err != nil {
		logger.L().Warn("invalid operator", zap.String("operator", operator))
		return "", err
	}
	return filed + " " + operator + " ?", nil
}
