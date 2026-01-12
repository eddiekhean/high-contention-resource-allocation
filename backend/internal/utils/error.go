package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrOutOfBounds     = errors.New("point out of bounds")
	ErrMazeUnreachable = errors.New("maze unreachable")
	ErrInvalidStrategy = errors.New("invalid solve strategy")
)

type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s %s", e.Field, e.Msg)
}

// ErrorDetail represents a single validation error
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// APIError represents the standard error response structure
type APIError struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

// NewAPIError creates a new APIError with the given code and message
func NewAPIError(code int, message string) APIError {
	return APIError{
		Code:    code,
		Message: message,
	}
}

// FormatValidationError converts validator.ValidationErrors into a standardized APIError
func FormatValidationError(err error) APIError {
	var errors []ErrorDetail

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors = append(errors, ErrorDetail{
				Field:   e.Field(),
				Message: msgForTag(e),
			})
		}
	} else {
		// Fallback for non-validation errors (shouldn't happen with ShouldBindJSON if triggered by validation)
		errors = append(errors, ErrorDetail{
			Field:   "global",
			Message: err.Error(),
		})
	}

	return APIError{
		Code:    400,
		Message: "Validation failed",
		Errors:  errors,
	}
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", strings.ReplaceAll(fe.Param(), " ", ", "))
	case "gt":
		return fmt.Sprintf("Must be greater than %s", fe.Param())
	default:
		return fe.Error() // Default formatted error
	}
}
