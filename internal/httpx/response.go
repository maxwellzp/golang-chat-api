package httpx

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"unicode"
)

type ValidationErrorMap map[string]string

func (v ValidationErrorMap) Error() string {
	return "validation failed"
}

func WriteValidationError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case validator.ValidationErrors:
		errorMap := make(map[string]string)
		for _, ve := range e {
			field := toSnakeCase(ve.Field())
			message := validationMessage(ve.Tag(), field, ve.Param())
			errorMap[field] = message
		}
		WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "Validation failed",
			"errors":  errorMap,
		})
	case ValidationErrorMap:
		WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "Validation failed",
			"errors":  e,
		})
	default:
		WriteError(w, http.StatusBadRequest, "Invalid input")
	}
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_', unicode.ToLower(r))
		} else {
			result = append(result, unicode.ToLower(r))
		}
	}
	return string(result)
}

func validationMessage(tag, field, param string) string {
	var msg string
	switch tag {
	case "required":
		msg = fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		msg = fmt.Sprintf("%s must be at least %s characters", field, param)
	case "max":
		msg = fmt.Sprintf("%s must be at most %s characters", field, param)
	case "containsuppercase":
		msg = fmt.Sprintf("%s must contain at least one uppercase letter", field)
	case "containslowercase":
		msg = fmt.Sprintf("%s must contain at least one lowercase letter", field)
	case "containsnumber":
		msg = fmt.Sprintf("%s must contain at least one number", field)
	case "containsspecial":
		msg = fmt.Sprintf("%s must contain at least one special character", field) + ` (!@#$%^&*)`
	default:
		msg = fmt.Sprintf("%s is invalid", field)
	}
	return msg
}
