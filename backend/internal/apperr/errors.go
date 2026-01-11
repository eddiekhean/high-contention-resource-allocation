package apperr

import (
	"errors"
	"net/http"
)

type ErrorCode string

const (
	CodeInternal       ErrorCode = "INTERNAL_ERROR"
	CodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	CodeNotFound       ErrorCode = "NOT_FOUND"
	CodeDuplicate      ErrorCode = "DUPLICATE_RESOURCE"
	CodeUnauthorized   ErrorCode = "UNAUTHORIZED"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined errors
var (
	ErrInternal       = New(CodeInternal, "Internal server error", nil)
	ErrInvalidRequest = New(CodeInvalidRequest, "Invalid request body", nil)
	ErrNotFound       = New(CodeNotFound, "Resource not found", nil)
	ErrDuplicate      = New(CodeDuplicate, "Resource already exists", nil)
)

func ToHTTPStatus(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case CodeInvalidRequest:
			return http.StatusBadRequest
		case CodeNotFound:
			return http.StatusNotFound
		case CodeDuplicate:
			return http.StatusConflict
		case CodeUnauthorized:
			return http.StatusUnauthorized
		default:
			return http.StatusInternalServerError
		}
	}
	return http.StatusInternalServerError
}

func ToResponse(err error) (int, interface{}) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return ToHTTPStatus(appErr), appErr
	}
	return http.StatusInternalServerError, AppError{
		Code:    CodeInternal,
		Message: "Internal server error",
	}
}
