package apierrs

import (
	"errors"
	"fmt"
)

// Map of error types to HTTP status codes
var errorTypeToStatusCode = map[ErrorType]int{
	ErrorTypeInvalidRequest: 400, // Bad Request
	ErrorTypeRateLimit:      429, // Too Many Requests
	ErrorTypeServer:         500, // Internal Server Error
	ErrorTypeAuthentication: 401, // Unauthorized
	ErrorTypePermission:     403, // Forbidden
}

// Common error creation helpers

// NewInvalidRequest creates an invalid request error with the given message
func NewInvalidRequest(message string) *APIError {
	return &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type:    ErrorTypeInvalidRequest,
				Code:    ErrorCodeInvalidInput,
				Message: message,
			},
		},
	}
}

// NewInvalidField creates an invalid request error for a specific field
func NewInvalidField(field, message string) *APIError {
	fieldParam := field
	return &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type:    ErrorTypeInvalidRequest,
				Code:    ErrorCodeInvalidField,
				Message: message,
				Param:   &fieldParam,
			},
		},
	}
}

// NewMissingField creates an error for a missing required field
func NewMissingField(field string) *APIError {
	fieldParam := field
	return &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type:    ErrorTypeInvalidRequest,
				Code:    ErrorCodeMissingField,
				Message: fmt.Sprintf("Missing required field: %s", field),
				Param:   &fieldParam,
			},
		},
	}
}

// NewServerError creates a server error with the given message
func NewServerError(message string) *APIError {
	return &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type:    ErrorTypeServer,
				Code:    ErrorCodeInternalError,
				Message: message,
			},
		},
	}
}

// NewAuthenticationError creates an authentication error with the given message
func NewAuthenticationError(message string, code ErrorCode) *APIError {
	return &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type:    ErrorTypeAuthentication,
				Code:    code,
				Message: message,
			},
		},
	}
}

// NewPermissionError creates a permission error with the given message
func NewPermissionError(message string) *APIError {
	return &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type:    ErrorTypePermission,
				Message: message,
			},
		},
	}
}

// NewRateLimitError creates a rate limit error with the given message
func NewRateLimitError(message string) *APIError {
	return &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type:    ErrorTypeRateLimit,
				Code:    ErrorCodeRateLimitExceeded,
				Message: message,
			},
		},
	}
}

// Error inspection helpers

// IsAPIError checks if an error is or wraps an APIError
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// GetAPIError extracts an APIError from an error if possible
// Returns nil if the error is not an APIError
func GetAPIError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return nil
}

// HTTPStatusForError returns an appropriate HTTP status code for the given error
func HTTPStatusForError(err error) int {
	if err == nil {
		return 200 // OK
	}

	apiErr := GetAPIError(err)
	if apiErr == nil {
		// Default for non-API errors
		return 500 // Internal Server Error
	}

	// Look up the status code from the map
	statusCode, exists := errorTypeToStatusCode[apiErr.Response.Error.Type]
	if !exists {
		return 500 // Default to Internal Server Error if type is unknown
	}

	// Special case for service unavailable
	if apiErr.Response.Error.Code == ErrorCodeServiceDown {
		return 503 // Service Unavailable
	}

	return statusCode
}
