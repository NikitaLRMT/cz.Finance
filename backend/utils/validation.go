package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct проверяет структуру на основе тегов валидации
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		var sb strings.Builder
		for _, err := range err.(validator.ValidationErrors) {
			sb.WriteString(formatValidationError(err))
			sb.WriteString("; ")
		}
		return errors.New(strings.TrimSuffix(sb.String(), "; "))
	}
	return nil
}

// formatValidationError форматирует ошибку валидации в понятное сообщение
func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("Поле '%s' обязательно", err.Field())
	case "email":
		return fmt.Sprintf("Поле '%s' должно содержать корректный email", err.Field())
	case "min":
		if err.Type().Kind().String() == "string" {
			return fmt.Sprintf("Поле '%s' должно содержать минимум %s символов", err.Field(), err.Param())
		}
		return fmt.Sprintf("Поле '%s' должно быть не менее %s", err.Field(), err.Param())
	case "max":
		if err.Type().Kind().String() == "string" {
			return fmt.Sprintf("Поле '%s' должно содержать максимум %s символов", err.Field(), err.Param())
		}
		return fmt.Sprintf("Поле '%s' должно быть не более %s", err.Field(), err.Param())
	case "gt":
		return fmt.Sprintf("Поле '%s' должно быть больше %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("Поле '%s' не соответствует правилу '%s'", err.Field(), err.Tag())
	}
}
