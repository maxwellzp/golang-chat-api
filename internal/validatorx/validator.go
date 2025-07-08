package validatorx

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"unicode"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		jsonTag := field.Tag.Get("json")
		name := strings.Split(jsonTag, ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validations
	v.RegisterValidation("containsuppercase", containsUpperCase)
	v.RegisterValidation("containslowercase", containsLowerCase)
	v.RegisterValidation("containsnumber", containsNumber)
	v.RegisterValidation("containsspecial", containsSpecial)

	return &Validator{validate: v}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

func containsUpperCase(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

func containsLowerCase(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

func containsNumber(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsNumber(char) {
			return true
		}
	}
	return false
}

func containsSpecial(fl validator.FieldLevel) bool {
	specialChars := "!@#$%^&*"
	fieldValue := fl.Field().String()

	for _, char := range specialChars {
		if strings.ContainsRune(fieldValue, char) {
			return true
		}
	}
	return false
}
