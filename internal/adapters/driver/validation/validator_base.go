package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, ", ")
}

func Validate(message interface{}) error {
	validate := validator.New()
	err := validate.Struct(message)
	if err != nil {
		var validationErrors ValidationErrors

		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, err := range errs {
				validationErrors = append(validationErrors, ValidationError{
					Field:   err.Field(),
					Message: getErrorMessage(err),
				})
			}
			return validationErrors
		}
		return err
	}
	return nil
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Please enter a valid email address"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters long", err.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", err.Param())
	case "uuid":
		return "Must be a valid UUID"
	case "url":
		return "Must be a valid URL"
	case "datetime":
		return "Must be a valid date/time"
	case "numeric":
		return "Must be a numeric value"
	case "alphanum":
		return "Must contain only letters and numbers"
	case "json":
		return "Must be a valid JSON"
	default:
		return fmt.Sprintf("Failed %s validation", err.Tag())
	}
}
