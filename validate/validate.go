package validate

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"go_base/domain"
	"go_base/logger"
	"go_base/xerror"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entrans "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/datatypes"
)

// use a single instance , it caches struct info
var (
	trans ut.Translator
	v     *validator.Validate
)

func init() {
	en := en.New()
	uni := ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ = uni.GetTranslator("en")

	v = validator.New()

	v.RegisterValidation("time", isTime)
	entrans.RegisterDefaultTranslations(v, trans)

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}

		return name
	})

	v.RegisterCustomTypeFunc(validateObjectID, uuid.UUID{})

	v.RegisterValidation("valid_permissions", func(fl validator.FieldLevel) bool {
		return ValidPermissions(fl.Field()) == nil
	})
	v.RegisterCustomTypeFunc(ValidPermissions, domain.PermissionTree{})

	v.RegisterValidation("staff_status", ValidateStaffStatus)
	v.RegisterValidation("phone", ValidatePhone)

	v.RegisterValidation("valid_jsonb", ValidateJsonb)

	v.RegisterValidation("one_of_array", func(fl validator.FieldLevel) bool {
		return ValidateOneOfArray(fl) == nil
	})
	v.RegisterValidation("enum", func(fl validator.FieldLevel) bool {
		return ValidateEnumArray(fl) == nil
	})
}

func New() *validator.Validate {
	return v
}

// Struct validates struct fields from rules defined in struct tag.
func Struct(s any) error {
	err := v.Struct(s)

	if err == nil {
		return nil
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	ferrs := make(domain.FieldErrors, len(errs))
	for i, err := range errs {
		ferrs[i] = &domain.FieldError{
			Field: err.Field(),
			Err:   err.Translate(trans),
		}
	}

	xerr := xerror.EInvalidInput(ferrs).SetMessage("invalid input")
	for _, data := range ferrs {
		_ = xerr.SetExtraInfo(data.Field, data.Err)
	}

	return xerr
}

func isTime(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	// UTC + 0 Only
	if _, err := time.Parse("2006-01-02T15:04:05Z", fl.Field().String()); err == nil {
		return true
	}
	return false
}

func validateObjectID(field reflect.Value) any {
	if val, ok := field.Interface().(uuid.UUID); ok {
		if val == uuid.Nil {
			return nil
		}
		if val.String() == "00000000-0000-0000-0000-000000000000" {
			return nil
		}
		return val.String()
	}
	return nil
}

func ValidPermissions(field reflect.Value) interface{} {
	switch field.Interface().(type) {
	case datatypes.JSON:
		d, ok := field.Interface().(datatypes.JSON)
		if !ok {
			return ok
		}
		var p domain.PermissionTree
		b := []byte(d)
		if err := json.Unmarshal(b, &p); err != nil {
			return fmt.Errorf("invalid permission")
		}
		if p["*"]["*"]["*"] == "*" {
			return fmt.Errorf("invalid permission")
		}
		// last str is true or false
		for _, v := range p {
			for _, v1 := range v {
				for _, v2 := range v1 {
					if v2 != "true" && v2 != "false" {
						return fmt.Errorf("invalid permission")
					}

				}
			}
		}
		return nil
	}
	return nil
}

func ValidateStaffStatus(fl validator.FieldLevel) bool {
	switch fl.Field().String() {
	case string(domain.StaffActive), string(domain.StaffInactive), string(domain.StaffPending):
		return true
	}
	return false
}

func ValidatePhone(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	phone := fl.Field().String()
	// regex phone 0{8,9,6,2}[-. ]?\d{3}[-. ]?\d{4} or without -
	phoneRegexD := regexp.MustCompile(`^0[8,9,6,2][-.\s]?\d{3}[-.\s]?\d{4}$`)
	phoneRegex := regexp.MustCompile(`^0[8,9,6,2]\d{8}$`)
	if phoneRegexD.MatchString(phone) || phoneRegex.MatchString(phone) {
		return true
	}

	return false
}

func ValidateJsonb(fl validator.FieldLevel) bool {
	switch fl.Field().Interface().(type) {
	case datatypes.JSON:
		// check case [] , {}
		d, ok := fl.Field().Interface().(datatypes.JSON)
		if !ok {
			return ok
		}
		var p interface{}
		b := []byte(d)
		if err := json.Unmarshal(b, &p); err != nil {
			return false
		}
		return true
	}
	return false
}

func ValidateOneOfArray(fl validator.FieldLevel) any {
	tagsParam := strings.Split(fl.Param(), " ")
	if len(tagsParam) == 0 || len(fl.Param()) == 0 {
		logger.L().Error("one_of_array is empty in struct")
		return false
	}
	switch fl.Field().Interface().(type) {
	case datatypes.JSON:
		d, ok := fl.Field().Interface().(datatypes.JSON)
		if !ok {
			return ok
		}
		var p interface{}
		b := []byte(d)
		if err := json.Unmarshal(b, &p); err != nil {
			return false
		}
		for _, v := range p.([]interface{}) {
			for _, tag := range tagsParam {
				if v == tag {
					return nil
				}
			}
		}
	}
	return fmt.Errorf("should be one of %s", fl.Param())
}

func ValidateEnumArray(fl validator.FieldLevel) any {
	tagsParam := strings.Split(fl.Param(), " ")
	if len(tagsParam) == 0 || len(fl.Param()) == 0 {
		logger.L().Error("enum is empty in struct")
		return false
	}
	switch fl.Field().Interface().(type) {
	case datatypes.JSON:
		d, ok := fl.Field().Interface().(datatypes.JSON)
		if !ok {
			return ok
		}
		var p interface{}
		b := []byte(d)
		if err := json.Unmarshal(b, &p); err != nil {
			return false
		}
		// find is not match in tagsParam
		for _, v := range p.([]interface{}) {
			if _, ok := lo.Find(tagsParam, func(t string) bool {
				return fmt.Sprintf("%v", v) == t
			}); !ok {
				return fmt.Errorf("should be one of %s", fl.Param())
			}
		}

	}
	return nil
}
