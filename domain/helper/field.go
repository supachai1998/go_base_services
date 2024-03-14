package domain

import (
	"encoding/json"
	"go_base/xerror"
	"math/big"
	"reflect"
)

func GetFieldValueByTag(field interface{}, tag string) (interface{}, error) {
	value := reflect.ValueOf(field)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil, xerror.EInternalError().SetMessage("field is not struct")
	}
	fieldValue := value.FieldByName(tag)
	if !fieldValue.IsValid() {
		return nil, xerror.EInternalError().SetMessage("field not found")
	}
	return fieldValue.Interface(), nil
}

func GetFieldNameByTag(field interface{}, tag string) (string, error) {
	value := reflect.ValueOf(field)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return "", xerror.EInternalError().SetMessage("field is not struct")
	}
	fieldType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := fieldType.Field(i)
		if field.Tag.Get("json") == tag {
			return ToSnakeCase(field.Name), nil
		}
	}
	return "", xerror.EInternalError().SetMessage("field not found")
}

func GetFieldNameByField(field interface{}, fieldVal interface{}) (string, error) {
	value := reflect.ValueOf(field)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	fieldValue := reflect.ValueOf(fieldVal)
	if value.Kind() != reflect.Struct {
		return "", xerror.EInternalError().SetMessage("field is not struct")
	}
	fieldType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := fieldType.Field(i)
		if field.Type == fieldValue.Type() {
			return ToSnakeCase(field.Name), nil
		}
	}
	return "", xerror.EInternalError().SetMessage("field not found")
}

func GetStructNameByStruct(structVal interface{}) (string, error) {

	if structVal == nil {
		return "", xerror.EInternalError().SetMessage("structVal is nil")
	}
	if reflect.ValueOf(structVal).Kind() == reflect.Ptr {
		structVal = reflect.ValueOf(structVal).Elem().Interface()
	}
	if reflect.ValueOf(structVal).Kind() == reflect.Slice {
		return ToSnakeCase(reflect.TypeOf(structVal).Elem().Name()), nil
	}
	t := reflect.TypeOf(structVal)
	if t.Kind() != reflect.Struct {
		return "", xerror.EInternalError().SetMessage("structVal is not struct")
	}
	return ToSnakeCase(t.Name()), nil

}

func GetValueFromStructByFieldName(structVal interface{}, fieldName string) (interface{}, error) {
	if structVal == nil {
		return nil, xerror.EInternalError().SetMessage("structVal is nil")
	}
	if reflect.ValueOf(structVal).Kind() == reflect.Ptr {
		structVal = reflect.ValueOf(structVal).Elem().Interface()
	}
	value := reflect.ValueOf(structVal)
	if value.Kind() != reflect.Struct {
		return nil, xerror.EInternalError().SetMessage("structVal is not struct")
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if ToSnakeCase(field.Name) == fieldName {
			fieldValue := value.Field(i)
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue = fieldValue.Elem()
			}
			if !fieldValue.IsValid() {
				return nil, xerror.EInternalError().SetMessage("field not found")
			}
			return fieldValue.Interface(), nil
		}
	}
	return nil, xerror.EInternalError().SetMessage("field not found")
}

func GetTableNameByStruct(structVal interface{}) (string, error) {
	table, err := GetStructNameByStruct(structVal)
	if err != nil {
		return "", err
	}
	return ToSnakeCase(table) + "s", nil
}

func StructHiddenFieldFromResponse[T any](_struct any) (T, error) {
	var result T
	if reflect.TypeOf(_struct).Kind() == reflect.Ptr {
		_struct = reflect.ValueOf(_struct).Elem().Interface()
	}

	b, err := json.MarshalIndent(_struct, "", "  ")
	if err != nil {
		return result, err
	}
	// assgin value to result
	json.Unmarshal(b, &result)
	return result, nil
}

// SerializeData is used to serialize data from interface to struct
//
// example:
//
//	staffSerialize, _ := helper.SerializeData[domain.Staff](staffCreate)
func SerializeData[T any](original interface{}) (*T, error) {
	var dest T
	b, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &dest)
	if err != nil {
		return nil, err
	}
	return &dest, nil
}

func Copy[T any](original interface{}) T {
	var dest T
	b, _ := json.Marshal(original)
	json.Unmarshal(b, &dest)
	return dest
}

func ToBigFloat(n interface{}) *big.Float {
	bf := new(big.Float)
	switch n.(type) {
	case int:
		bf.SetInt64(int64(n.(int)))
	case int64:
		bf.SetInt64(n.(int64))
	case float64:
		bf.SetFloat64(n.(float64))
	case float32:
		bf.SetFloat64(float64(n.(float32)))
	}
	return bf
}
