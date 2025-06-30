package exceptions

import (
	"context" // Import context
	"fmt"
	"mini-evv-logger-backend/responses"
	"net/http"

	"github.com/gofiber/fiber/v2" // Import Fiber for c *fiber.Ctx
)

// CustomError represents a custom application error
type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface for CustomError
func (e *CustomError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("Error %d: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

// WithDetails allows adding more detailed information to an existing CustomError
func (e *CustomError) WithDetails(details string) *CustomError {
	// Create a new CustomError instance to avoid modifying the original global error variable
	newErr := *e
	newErr.Details = details
	return &newErr
}

// NewCustomError creates a new CustomError instance
func NewCustomError(code int, message string, details ...string) *CustomError {
	ce := &CustomError{
		Code:    code,
		Message: message,
	}
	if len(details) > 0 {
		ce.Details = details[0]
	}
	return ce
}

// IsContextError checks if the given error is a context.Canceled or context.DeadlineExceeded error.
// If it is, it returns true and a corresponding CustomError. Otherwise, it returns false and nil.
func IsContextError(err error) (bool, *CustomError) {
	switch err {
	case context.Canceled:
		return true, NewCustomError(http.StatusServiceUnavailable, "Request cancelled", "The client cancelled the request.")
	case context.DeadlineExceeded:
		return true, NewCustomError(http.StatusGatewayTimeout, "Request timed out", "The request exceeded its allotted time.")
	default:
		return false, nil
	}
}

// HandleError centralizes error handling logic for Fiber controllers.
// It checks for CustomError, context errors, and falls back to a generic internal server error.
func HandleError(c *fiber.Ctx, err error) error {
	if customErr, ok := err.(*CustomError); ok {
		return responses.Error(c, customErr.Code, customErr.Message, customErr)
	}

	if isContextErr, ctxErr := IsContextError(err); isContextErr {
		return responses.Error(c, ctxErr.Code, ctxErr.Message, ctxErr)
	}

	// Default to internal server error for unhandled errors
	return responses.Error(c, http.StatusInternalServerError, "An unexpected error occurred", err.Error())
}

// Common errors
var (
	ErrNotFound            = NewCustomError(http.StatusNotFound, "Resource not found")
	ErrBadRequest          = NewCustomError(http.StatusBadRequest, "Bad request")
	ErrInternalError       = NewCustomError(http.StatusInternalServerError, "Internal server error")
	ErrUnauthorized        = NewCustomError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden           = NewCustomError(http.StatusForbidden, "Forbidden")
	ErrConflict            = NewCustomError(http.StatusConflict, "Conflict")
	ErrUnprocessableEntity = NewCustomError(http.StatusUnprocessableEntity, "Unprocessable Entity")
)
