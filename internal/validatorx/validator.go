package validatorx

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
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

	return &Validator{validate: v}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}
