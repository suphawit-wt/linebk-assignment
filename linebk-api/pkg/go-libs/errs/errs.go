package errs

import (
	"net/http"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/dto"
	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/enums/api"
)

// AppError represents a standard application domain error.
type AppError struct {
	Status     string
	StatusCode int
	Message    string
}

// ValidationError represents an error containing input validation field errors.
type ValidationError struct {
	AppError
	Errors []dto.ErrResponse
}

// Error returns the underlying message string for AppError.
func (e AppError) Error() string {
	return e.Message
}

// NewBadRequestError constructs a new 400 Bad Request application error.
func NewBadRequestError(message string) error {
	return AppError{
		Status:     api.StatusBadRequest,
		StatusCode: http.StatusBadRequest,
		Message:    message,
	}
}

// NewValidationError constructs a new ValidationError containing field-level validation errors.
func NewValidationError(message string, errors []dto.ErrResponse) ValidationError {
	return ValidationError{
		AppError: AppError{
			Status:     api.StatusInvalidInput,
			StatusCode: http.StatusBadRequest,
			Message:    message,
		},
		Errors: errors,
	}
}

// NewUnauthorizedError constructs a new 401 Unauthorized application error.
func NewUnauthorizedError(message string) error {
	return AppError{
		Status:     api.StatusUnauthorized,
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

// NewExpiredOTP constructs an application error indicating that an OTP has expired.
func NewExpiredOTP(message string) error {
	return AppError{
		Status:     api.StatusExpiredOTP,
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

// NewTokenNotMatched constructs an application error indicating a token mismatch.
func NewTokenNotMatched(message string) error {
	return AppError{
		Status:     api.StatusTokenNotMatched,
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

// NewForbiddenError constructs a new 403 Forbidden application error.
func NewForbiddenError(message string) error {
	return AppError{
		Status:     api.StatusForbidden,
		StatusCode: http.StatusForbidden,
		Message:    message,
	}
}

// NewNotFoundError constructs a new 404 Not Found application error.
func NewNotFoundError(message string) error {
	return AppError{
		Status:     api.StatusNotFound,
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

// NewRequestTimeOutError constructs a new 408 Request Timeout application error.
func NewRequestTimeOutError(message string) error {
	return AppError{
		Status:     api.StatusRequestTimeout,
		StatusCode: http.StatusRequestTimeout,
		Message:    message,
	}
}

// NewConflictError constructs a new 409 Conflict application error.
func NewConflictError(message string) error {
	return AppError{
		Status:     api.StatusConflict,
		StatusCode: http.StatusConflict,
		Message:    message,
	}
}

// NewRequestEntityTooLargeError constructs a new 413 Payload Too Large application error.
func NewRequestEntityTooLargeError(message string) error {
	return AppError{
		Status:     api.StatusFileTooLarge,
		StatusCode: http.StatusRequestEntityTooLarge,
		Message:    message,
	}
}

// NewUnsupportedMediaTypeError constructs a new 415 Unsupported Media Type application error.
func NewUnsupportedMediaTypeError(message string) error {
	return AppError{
		Status:     api.StatusUnsupportedMediaType,
		StatusCode: http.StatusUnsupportedMediaType,
		Message:    message,
	}
}

// NewUnprocessableContentError constructs a new 422 Unprocessable Content application error.
func NewUnprocessableContentError(message string) error {
	return AppError{
		Status:     api.StatusUnprocessableContent,
		StatusCode: http.StatusUnprocessableEntity,
		Message:    message,
	}
}

// NewUnexpectedError constructs a generic 500 Internal Server Error application error.
func NewUnexpectedError() error {
	return AppError{
		Status:     api.StatusUnexpected,
		StatusCode: http.StatusInternalServerError,
		Message:    "Unexpected error occurred.",
	}
}
