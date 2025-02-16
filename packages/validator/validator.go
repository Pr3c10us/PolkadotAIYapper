package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	StatusCode   int            `json:"statusCode"`
	Message      string         `json:"message"`
	ErrorMessage []ErrorMessage `json:"error"`
}

func (err *ValidationError) Error() string {
	return "validation failed"
}

func getErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fieldError.Field() + " is required"
	case "lte":
		return fieldError.Field() + " should be less than " + fieldError.Param()
	case "gte":
		return fieldError.Field() + " should be greater than " + fieldError.Param()
	case "email":
		return "Invalid email format"
	case "min":
		return fieldError.Field() + " should be at least " + fieldError.Param() + " characters long"
	case "e164":
		return fieldError.Field() + " should be in valid E.164 format"
	default:
		return "Unknown error on field " + fieldError.Field()
	}
}

func ValidateRequest(err error) error {
	if err == nil {
		return nil
	}
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorMessages := make([]ErrorMessage, len(validationErrors))
		for i, fieldError := range validationErrors {
			errorMessages[i] = ErrorMessage{fieldError.Field(), getErrorMessage(fieldError)}
		}
		return &ValidationError{
			StatusCode:   http.StatusNotAcceptable,
			Message:      "validation failed",
			ErrorMessage: errorMessages,
		}
	}
	return err
}
