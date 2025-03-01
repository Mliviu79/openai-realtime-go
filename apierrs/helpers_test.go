package apierrs

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorCreationHelpers(t *testing.T) {
	tests := []struct {
		name         string
		createFunc   func() *APIError
		expectedType ErrorType
		expectedCode ErrorCode
	}{
		{
			name: "NewInvalidRequest",
			createFunc: func() *APIError {
				return NewInvalidRequest("Invalid request message")
			},
			expectedType: ErrorTypeInvalidRequest,
			expectedCode: ErrorCodeInvalidInput,
		},
		{
			name: "NewInvalidField",
			createFunc: func() *APIError {
				return NewInvalidField("email", "Invalid email format")
			},
			expectedType: ErrorTypeInvalidRequest,
			expectedCode: ErrorCodeInvalidField,
		},
		{
			name: "NewMissingField",
			createFunc: func() *APIError {
				return NewMissingField("password")
			},
			expectedType: ErrorTypeInvalidRequest,
			expectedCode: ErrorCodeMissingField,
		},
		{
			name: "NewServerError",
			createFunc: func() *APIError {
				return NewServerError("Internal server error")
			},
			expectedType: ErrorTypeServer,
			expectedCode: ErrorCodeInternalError,
		},
		{
			name: "NewAuthenticationError",
			createFunc: func() *APIError {
				return NewAuthenticationError("Invalid API key", ErrorCodeInvalidAPIKey)
			},
			expectedType: ErrorTypeAuthentication,
			expectedCode: ErrorCodeInvalidAPIKey,
		},
		{
			name: "NewPermissionError",
			createFunc: func() *APIError {
				return NewPermissionError("Insufficient permissions")
			},
			expectedType: ErrorTypePermission,
			expectedCode: "",
		},
		{
			name: "NewRateLimitError",
			createFunc: func() *APIError {
				return NewRateLimitError("Too many requests")
			},
			expectedType: ErrorTypeRateLimit,
			expectedCode: ErrorCodeRateLimitExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.createFunc()

			if err.Response.Error.Type != tt.expectedType {
				t.Errorf("Expected Type to be %v, got %v", tt.expectedType, err.Response.Error.Type)
			}

			if err.Response.Error.Code != tt.expectedCode {
				t.Errorf("Expected Code to be %v, got %v", tt.expectedCode, err.Response.Error.Code)
			}

			// Ensure the error message is not empty
			if err.Response.Error.Message == "" {
				t.Error("Expected non-empty error message")
			}

			// Check for specific fields in certain error types
			if tt.name == "NewInvalidField" {
				if err.Response.Error.Param == nil || *err.Response.Error.Param != "email" {
					t.Errorf("Expected Param to be 'email', got %v", err.Response.Error.Param)
				}
			}

			if tt.name == "NewMissingField" {
				if err.Response.Error.Param == nil || *err.Response.Error.Param != "password" {
					t.Errorf("Expected Param to be 'password', got %v", err.Response.Error.Param)
				}
			}

			// Check that the response type is set correctly
			if err.Response.Type != "error" {
				t.Errorf("Expected Response.Type to be 'error', got %q", err.Response.Type)
			}
		})
	}
}

func TestIsAPIError(t *testing.T) {
	// Test with APIError
	apiErr := NewInvalidRequest("Test error")
	if !IsAPIError(apiErr) {
		t.Error("IsAPIError should return true for APIError")
	}

	// Test with wrapped APIError
	wrappedErr := fmt.Errorf("wrapped: %w", apiErr)
	if !IsAPIError(wrappedErr) {
		t.Error("IsAPIError should return true for wrapped APIError")
	}

	// Test with non-APIError
	regularErr := errors.New("regular error")
	if IsAPIError(regularErr) {
		t.Error("IsAPIError should return false for non-APIError")
	}

	// Test with nil
	if IsAPIError(nil) {
		t.Error("IsAPIError should return false for nil")
	}
}

func TestGetAPIError(t *testing.T) {
	// Test with APIError
	apiErr := NewInvalidRequest("Test error")
	extracted := GetAPIError(apiErr)
	if extracted != apiErr {
		t.Error("GetAPIError should return the same APIError")
	}

	// Test with wrapped APIError
	wrappedErr := fmt.Errorf("wrapped: %w", apiErr)
	extracted = GetAPIError(wrappedErr)
	if extracted != apiErr {
		t.Error("GetAPIError should extract the APIError from wrapped error")
	}

	// Test with non-APIError
	regularErr := errors.New("regular error")
	extracted = GetAPIError(regularErr)
	if extracted != nil {
		t.Error("GetAPIError should return nil for non-APIError")
	}

	// Test with nil
	extracted = GetAPIError(nil)
	if extracted != nil {
		t.Error("GetAPIError should return nil for nil")
	}
}

func TestHTTPStatusForError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{
			name:         "Nil error",
			err:          nil,
			expectedCode: 200, // OK
		},
		{
			name:         "InvalidRequest error",
			err:          NewInvalidRequest("Invalid request"),
			expectedCode: 400, // Bad Request
		},
		{
			name:         "RateLimit error",
			err:          NewRateLimitError("Too many requests"),
			expectedCode: 429, // Too Many Requests
		},
		{
			name:         "Server error",
			err:          NewServerError("Internal server error"),
			expectedCode: 500, // Internal Server Error
		},
		{
			name: "Service Unavailable error",
			err: &APIError{
				Response: ErrorResponse{
					Type: "error",
					Error: ErrorDetails{
						Type:    ErrorTypeServer,
						Code:    ErrorCodeServiceDown,
						Message: "Service unavailable",
					},
				},
			},
			expectedCode: 503, // Service Unavailable
		},
		{
			name:         "Authentication error",
			err:          NewAuthenticationError("Invalid API key", ErrorCodeInvalidAPIKey),
			expectedCode: 401, // Unauthorized
		},
		{
			name:         "Permission error",
			err:          NewPermissionError("Insufficient permissions"),
			expectedCode: 403, // Forbidden
		},
		{
			name:         "Non-API error",
			err:          errors.New("generic error"),
			expectedCode: 500, // Default to Internal Server Error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode := HTTPStatusForError(tt.err)
			if statusCode != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, statusCode)
			}
		})
	}
}
