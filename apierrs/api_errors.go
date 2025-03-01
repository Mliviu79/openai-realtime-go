// Package apierrs provides error types and handling for the OpenAI Realtime API.
// It defines structured error types that match the OpenAI API error format and provides
// utility functions for error classification, creation, and handling.
//
// The package offers:
//   - Type-safe error types for different OpenAI API error categories
//   - Error creation and handling utilities
//   - Classification methods to determine error types (e.g., rate limits, server errors)
//   - Methods to determine if errors are transient and can be retried
//
// Error Structure:
// APIError implements the standard Go error interface and includes:
//   - Type: The category of error (e.g., "invalid_request_error", "rate_limit_error")
//   - Code: A specific error code for more detailed error information
//   - Message: A human-readable error message
//   - Param: The parameter that caused the error (if applicable)
//   - EventID: The ID of the client event that triggered the error (if applicable)
//
// Example Usage:
//
//	// Basic error checking
//	resp, err := client.CreateSession(ctx, req)
//	if err != nil {
//		var apiErr *apierrs.APIError
//		if errors.As(err, &apiErr) {
//			if apiErr.IsRateLimit() {
//				fmt.Println("Rate limit exceeded, retry after a delay")
//			} else if apiErr.IsInvalidRequest() {
//				fmt.Printf("Invalid request: %s\n", apiErr.Message)
//			}
//		}
//		return err
//	}
//
//	// Creating custom API errors
//	err := apierrs.NewAPIError(
//		apierrs.ErrorTypeInvalidRequest,
//		string(apierrs.ErrorCodeInvalidInput),
//		"The model parameter is required",
//	).WithParam("model")
//
// The error types and codes in this package match those returned by the OpenAI Realtime API.
package apierrs

import "fmt"

// ErrorType represents the type of error returned by the API
type ErrorType string

const (
	// ErrorTypeInvalidRequest indicates an error with the request format or parameters
	ErrorTypeInvalidRequest ErrorType = "invalid_request_error"
	// ErrorTypeRateLimit indicates the client has sent too many requests
	ErrorTypeRateLimit ErrorType = "rate_limit_error"
	// ErrorTypeServer indicates an internal server error
	ErrorTypeServer ErrorType = "server_error"
	// ErrorTypeAuthentication indicates issues with API keys or authentication
	ErrorTypeAuthentication ErrorType = "authentication_error"
	// ErrorTypePermission indicates issues with permissions for the requested operation
	ErrorTypePermission ErrorType = "permission_error"
)

// ErrorCode represents specific error codes returned by the API
type ErrorCode string

const (
	// Request-related errors
	ErrorCodeInvalidInput ErrorCode = "invalid_input"
	ErrorCodeMissingField ErrorCode = "missing_field"
	ErrorCodeInvalidField ErrorCode = "invalid_field"
	ErrorCodeInvalidEvent ErrorCode = "invalid_event"

	// Rate limit errors
	ErrorCodeRateLimitExceeded ErrorCode = "rate_limit_exceeded"
	ErrorCodeTooManyRequests   ErrorCode = "too_many_requests"

	// Server errors
	ErrorCodeInternalError ErrorCode = "internal_error"
	ErrorCodeServiceDown   ErrorCode = "service_unavailable"

	// Authentication errors
	ErrorCodeInvalidAPIKey     ErrorCode = "invalid_api_key"
	ErrorCodeMissingAPIKey     ErrorCode = "missing_api_key"
	ErrorCodeInsufficientQuota ErrorCode = "insufficient_quota"
)

// ErrorDetails represents the nested error details in an error response
type ErrorDetails struct {
	// The type of error (e.g., "invalid_request_error", "server_error")
	Type ErrorType `json:"type"`
	// Error code, providing more specific information about what went wrong
	Code ErrorCode `json:"code,omitempty"`
	// A human-readable error message
	Message string `json:"message"`
	// Parameter related to the error, if any
	Param *string `json:"param"`
	// The event_id of the client event that caused the error, if applicable
	EventID string `json:"event_id,omitempty"`
}

// ErrorResponse represents the complete error response from the API
type ErrorResponse struct {
	// The unique ID of the server event
	EventID string `json:"event_id"`
	// The type of event, always "error" for error responses
	Type string `json:"type"`
	// The error details
	Error ErrorDetails `json:"error"`
}

// APIError represents a structured error from the API
// This is the primary error type that should be used for all API-related errors
type APIError struct {
	// The complete error response
	Response ErrorResponse
}

// Error implements the error interface
func (e *APIError) Error() string {
	details := &e.Response.Error

	if details.Param != nil && details.EventID != "" {
		return fmt.Sprintf("%s: %s (code: %s, param: %s, event_id: %s)",
			details.Type, details.Message, details.Code, *details.Param, details.EventID)
	}

	if details.Param != nil && *details.Param != "" {
		return fmt.Sprintf("%s: %s (code: %s, param: %s)",
			details.Type, details.Message, details.Code, *details.Param)
	}

	if details.EventID != "" {
		return fmt.Sprintf("%s: %s (code: %s, event_id: %s)",
			details.Type, details.Message, details.Code, details.EventID)
	}

	return fmt.Sprintf("%s: %s (code: %s)", details.Type, details.Message, details.Code)
}

// NewAPIError creates a new APIError with the given parameters
func NewAPIError(errType ErrorType, code string, message string) *APIError {
	return &APIError{
		Response: ErrorResponse{
			EventID: "", // Will be set by the server
			Type:    "error",
			Error: ErrorDetails{
				Type:    errType,
				Code:    ErrorCode(code),
				Message: message,
				Param:   nil,
				EventID: "",
			},
		},
	}
}

// WithParam adds a parameter to the error
func (e *APIError) WithParam(param string) *APIError {
	e.Response.Error.Param = &param
	return e
}

// WithEventID adds an event ID to the error
func (e *APIError) WithEventID(eventID string) *APIError {
	e.Response.Error.EventID = eventID
	return e
}

// WithResponseEventID sets the event ID for the response
func (e *APIError) WithResponseEventID(eventID string) *APIError {
	e.Response.EventID = eventID
	return e
}

// Error Classification Methods

// IsInvalidRequest returns true if the error is an invalid request error
func (e *APIError) IsInvalidRequest() bool {
	return e.Response.Error.Type == ErrorTypeInvalidRequest
}

// IsRateLimit returns true if the error is a rate limit error
func (e *APIError) IsRateLimit() bool {
	return e.Response.Error.Type == ErrorTypeRateLimit
}

// IsServerError returns true if the error is a server error
func (e *APIError) IsServerError() bool {
	return e.Response.Error.Type == ErrorTypeServer
}

// IsAuthenticationError returns true if the error is an authentication error
func (e *APIError) IsAuthenticationError() bool {
	return e.Response.Error.Type == ErrorTypeAuthentication
}

// IsPermissionError returns true if the error is a permission error
func (e *APIError) IsPermissionError() bool {
	return e.Response.Error.Type == ErrorTypePermission
}

// IsTransient returns true if the error is likely transient and can be retried
// Retryable errors include rate limits and server errors
func (e *APIError) IsTransient() bool {
	return e.IsRateLimit() || e.IsServerError()
}
