package domain

import (
	"reflect"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToTableName(m any, isTableNames ...bool) string {
	isTableName := true
	if len(isTableNames) > 0 {
		isTableName = isTableNames[0]
	}
	switch reflect.TypeOf(m).Kind() {
	case reflect.String:
		s, ok := m.(string)
		if !ok {
			return ""
		}
		sneak := ToSnakeCase(s)
		if !strings.HasSuffix(sneak, "s") && isTableName {
			return sneak + "s"
		}
		if !isTableName && strings.HasSuffix(sneak, "s") {
			return sneak[:len(sneak)-1]
		}
		return sneak
	case reflect.Struct:
		t := reflect.TypeOf(m)
		m := reflect.New(t).Elem().Interface()
		modelName := reflect.TypeOf(m).Name()
		sneak := ToSnakeCase(modelName)
		if !strings.HasSuffix(sneak, "s") && isTableName {
			return sneak + "s"
		}
		if !isTableName && strings.HasSuffix(sneak, "s") {
			return sneak[:len(sneak)-1]
		}
		return sneak

	}
	return ""
}

func IsArray(str string) bool {
	return strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]")
}

func IsObject(str string) bool {
	return strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}")
}
