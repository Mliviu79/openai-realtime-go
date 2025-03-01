package apierrs

import (
	"encoding/json"
	"testing"
)

func TestAPIErrorError(t *testing.T) {
	// Test basic error with only type, code, and message
	basic := &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type:    ErrorTypeInvalidRequest,
				Code:    ErrorCodeInvalidInput,
				Message: "Invalid input provided",
			},
		},
	}
	expected := "invalid_request_error: Invalid input provided (code: invalid_input)"
	if basic.Error() != expected {
		t.Errorf("Expected error string %q, got %q", expected, basic.Error())
	}

	// Test with param
	paramValue := "username"
	withParam := &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type:    ErrorTypeInvalidRequest,
				Code:    ErrorCodeInvalidField,
				Message: "Invalid field value",
				Param:   &paramValue,
			},
		},
	}
	expected = "invalid_request_error: Invalid field value (code: invalid_field, param: username)"
	if withParam.Error() != expected {
		t.Errorf("Expected error string %q, got %q", expected, withParam.Error())
	}

	// Test with event ID
	withEvent := &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type:    ErrorTypeServer,
				Code:    ErrorCodeInternalError,
				Message: "Internal server error",
				EventID: "evt_123456",
			},
		},
	}
	expected = "server_error: Internal server error (code: internal_error, event_id: evt_123456)"
	if withEvent.Error() != expected {
		t.Errorf("Expected error string %q, got %q", expected, withEvent.Error())
	}

	// Test with both param and event ID
	bothParam := "requests_per_minute"
	withBoth := &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type:    ErrorTypeRateLimit,
				Code:    ErrorCodeRateLimitExceeded,
				Message: "Rate limit exceeded",
				Param:   &bothParam,
				EventID: "evt_123456",
			},
		},
	}
	expected = "rate_limit_error: Rate limit exceeded (code: rate_limit_exceeded, param: requests_per_minute, event_id: evt_123456)"
	if withBoth.Error() != expected {
		t.Errorf("Expected error string %q, got %q", expected, withBoth.Error())
	}
}

func TestNewAPIError(t *testing.T) {
	err := NewAPIError(ErrorTypeInvalidRequest, string(ErrorCodeInvalidInput), "Test message")

	if err.Response.Error.Type != ErrorTypeInvalidRequest {
		t.Errorf("Expected Type to be %v, got %v", ErrorTypeInvalidRequest, err.Response.Error.Type)
	}

	if err.Response.Error.Code != ErrorCodeInvalidInput {
		t.Errorf("Expected Code to be %v, got %v", ErrorCodeInvalidInput, err.Response.Error.Code)
	}

	if err.Response.Error.Message != "Test message" {
		t.Errorf("Expected Message to be %q, got %q", "Test message", err.Response.Error.Message)
	}

	if err.Response.Type != "error" {
		t.Errorf("Expected Response.Type to be 'error', got %q", err.Response.Type)
	}
}

func TestWithParam(t *testing.T) {
	err := NewAPIError(ErrorTypeInvalidRequest, string(ErrorCodeInvalidField), "Invalid field")
	err = err.WithParam("username")

	if err.Response.Error.Param == nil || *err.Response.Error.Param != "username" {
		t.Errorf("Expected Param to be 'username', got %v", err.Response.Error.Param)
	}
}

func TestWithEventID(t *testing.T) {
	err := NewAPIError(ErrorTypeServer, string(ErrorCodeInternalError), "Server error")
	err = err.WithEventID("evt_123456")

	if err.Response.Error.EventID != "evt_123456" {
		t.Errorf("Expected EventID to be %q, got %q", "evt_123456", err.Response.Error.EventID)
	}
}

func TestWithResponseEventID(t *testing.T) {
	err := NewAPIError(ErrorTypeServer, string(ErrorCodeInternalError), "Server error")
	err = err.WithResponseEventID("evt_789012")

	if err.Response.EventID != "evt_789012" {
		t.Errorf("Expected Response EventID to be %q, got %q", "evt_789012", err.Response.EventID)
	}
}

func TestErrorClassification(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		isInvReq bool
		isRate   bool
		isServer bool
		isAuth   bool
		isPerm   bool
	}{
		{
			name: "InvalidRequest",
			err: &APIError{
				Response: ErrorResponse{
					Type: "error",
					Error: ErrorDetails{
						Type: ErrorTypeInvalidRequest,
					},
				},
			},
			isInvReq: true,
		},
		{
			name: "RateLimit",
			err: &APIError{
				Response: ErrorResponse{
					Type: "error",
					Error: ErrorDetails{
						Type: ErrorTypeRateLimit,
					},
				},
			},
			isRate: true,
		},
		{
			name: "ServerError",
			err: &APIError{
				Response: ErrorResponse{
					Type: "error",
					Error: ErrorDetails{
						Type: ErrorTypeServer,
					},
				},
			},
			isServer: true,
		},
		{
			name: "AuthenticationError",
			err: &APIError{
				Response: ErrorResponse{
					Type: "error",
					Error: ErrorDetails{
						Type: ErrorTypeAuthentication,
					},
				},
			},
			isAuth: true,
		},
		{
			name: "PermissionError",
			err: &APIError{
				Response: ErrorResponse{
					Type: "error",
					Error: ErrorDetails{
						Type: ErrorTypePermission,
					},
				},
			},
			isPerm: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.IsInvalidRequest() != tt.isInvReq {
				t.Errorf("IsInvalidRequest() = %v, want %v", tt.err.IsInvalidRequest(), tt.isInvReq)
			}
			if tt.err.IsRateLimit() != tt.isRate {
				t.Errorf("IsRateLimit() = %v, want %v", tt.err.IsRateLimit(), tt.isRate)
			}
			if tt.err.IsServerError() != tt.isServer {
				t.Errorf("IsServerError() = %v, want %v", tt.err.IsServerError(), tt.isServer)
			}
			if tt.err.IsAuthenticationError() != tt.isAuth {
				t.Errorf("IsAuthenticationError() = %v, want %v", tt.err.IsAuthenticationError(), tt.isAuth)
			}
			if tt.err.IsPermissionError() != tt.isPerm {
				t.Errorf("IsPermissionError() = %v, want %v", tt.err.IsPermissionError(), tt.isPerm)
			}
		})
	}
}

func TestIsTransient(t *testing.T) {
	// Transient errors: rate limits and server errors
	rateLimit := &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type: ErrorTypeRateLimit,
			},
		},
	}
	if !rateLimit.IsTransient() {
		t.Error("Expected rate limit error to be transient")
	}

	server := &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type: ErrorTypeServer,
			},
		},
	}
	if !server.IsTransient() {
		t.Error("Expected server error to be transient")
	}

	// Non-transient errors
	invalidRequest := &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type: ErrorTypeInvalidRequest,
			},
		},
	}
	if invalidRequest.IsTransient() {
		t.Error("Expected invalid request error to be non-transient")
	}

	auth := &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type: ErrorTypeAuthentication,
			},
		},
	}
	if auth.IsTransient() {
		t.Error("Expected authentication error to be non-transient")
	}

	permission := &APIError{
		Response: ErrorResponse{
			Type: "error",
			Error: ErrorDetails{
				Type: ErrorTypePermission,
			},
		},
	}
	if permission.IsTransient() {
		t.Error("Expected permission error to be non-transient")
	}
}

func TestAPIErrorJSON(t *testing.T) {
	// Test JSON marshaling and unmarshaling
	paramValue := "username"
	apiErr := &APIError{
		Response: ErrorResponse{
			EventID: "event_890",
			Type:    "error",
			Error: ErrorDetails{
				Type:    ErrorTypeInvalidRequest,
				Code:    ErrorCodeInvalidEvent,
				Message: "The 'type' field is missing.",
				Param:   &paramValue,
				EventID: "event_567",
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(apiErr.Response)
	if err != nil {
		t.Fatalf("Failed to marshal APIError to JSON: %v", err)
	}

	// Unmarshal from JSON
	var unmarshaledResp ErrorResponse
	if err := json.Unmarshal(jsonData, &unmarshaledResp); err != nil {
		t.Fatalf("Failed to unmarshal JSON to ErrorResponse: %v", err)
	}

	// Compare original and unmarshaled values
	if unmarshaledResp.Type != apiErr.Response.Type {
		t.Errorf("Expected Type to be %v, got %v", apiErr.Response.Type, unmarshaledResp.Type)
	}

	if unmarshaledResp.EventID != apiErr.Response.EventID {
		t.Errorf("Expected EventID to be %v, got %v", apiErr.Response.EventID, unmarshaledResp.EventID)
	}

	if unmarshaledResp.Error.Type != apiErr.Response.Error.Type {
		t.Errorf("Expected Error.Type to be %v, got %v", apiErr.Response.Error.Type, unmarshaledResp.Error.Type)
	}

	if unmarshaledResp.Error.Code != apiErr.Response.Error.Code {
		t.Errorf("Expected Error.Code to be %v, got %v", apiErr.Response.Error.Code, unmarshaledResp.Error.Code)
	}

	if unmarshaledResp.Error.Message != apiErr.Response.Error.Message {
		t.Errorf("Expected Error.Message to be %q, got %q", apiErr.Response.Error.Message, unmarshaledResp.Error.Message)
	}

	if *unmarshaledResp.Error.Param != *apiErr.Response.Error.Param {
		t.Errorf("Expected Error.Param to be %q, got %q", *apiErr.Response.Error.Param, *unmarshaledResp.Error.Param)
	}

	if unmarshaledResp.Error.EventID != apiErr.Response.Error.EventID {
		t.Errorf("Expected Error.EventID to be %q, got %q", apiErr.Response.Error.EventID, unmarshaledResp.Error.EventID)
	}

	// Check if the JSON matches the expected format from the OpenAI API
	expectedJSON := `{"event_id":"event_890","type":"error","error":{"type":"invalid_request_error","code":"invalid_event","message":"The 'type' field is missing.","param":"username","event_id":"event_567"}}`
	actualJSON := string(jsonData)
	if actualJSON != expectedJSON {
		t.Errorf("Expected JSON %q, got %q", expectedJSON, actualJSON)
	}
}
