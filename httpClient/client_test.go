package httpClient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Mliviu79/go-openai-realtime/apierrs"
)

type testRequest struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

type testResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func TestDo(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test" {
			t.Errorf("Expected request to '/test', got %q", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be 'application/json', got %q", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("Custom-Header") != "custom-value" {
			t.Errorf("Expected Custom-Header to be 'custom-value', got %q", r.Header.Get("Custom-Header"))
		}

		// Read the request body
		var req testRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify request fields
		if req.Field1 != "test" {
			t.Errorf("Expected Field1 to be 'test', got %q", req.Field1)
		}

		if req.Field2 != 123 {
			t.Errorf("Expected Field2 to be 123, got %d", req.Field2)
		}

		// Write response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "message": "request processed"}`))
	}))
	defer server.Close()

	// Create request
	req := &testRequest{
		Field1: "test",
		Field2: 123,
	}

	// Create custom header
	headers := http.Header{}
	headers.Set("Custom-Header", "custom-value")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Make the request
	resp, err := Do[testRequest, testResponse](
		ctx,
		server.URL+"/test",
		req,
		WithHeaders(headers),
		WithClient(server.Client()),
	)

	// Verify results
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	if !resp.Success {
		t.Errorf("Expected Success to be true, got false")
	}

	if resp.Message != "request processed" {
		t.Errorf("Expected Message to be 'request processed', got %q", resp.Message)
	}
}

func TestDoWithError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Format the error exactly like the OpenAI API would return it
		w.Write([]byte(`{"type":"error","event_id":"evt_123","error":{"type":"invalid_request_error","message":"invalid request","code":"invalid_input"}}`))
	}))
	defer server.Close()

	// Create request
	req := &testRequest{
		Field1: "test",
		Field2: 123,
	}

	// Make the request
	_, err := Do[testRequest, testResponse](
		context.Background(),
		server.URL+"/test",
		req,
	)

	// Verify error
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Print the error for debugging
	t.Logf("Error: %v (type: %T)", err, err)

	// Check if it's an API error
	if !apierrs.IsAPIError(err) {
		t.Fatalf("Expected API error, got %T: %v", err, err)
	}

	apiErr := apierrs.GetAPIError(err)
	if apiErr == nil {
		t.Fatal("Failed to get APIError from error")
	}

	// Print the API error structure for debugging
	t.Logf("APIError: %+v", apiErr)
	t.Logf("Response: %+v", apiErr.Response)
	t.Logf("Error Details: %+v", apiErr.Response.Error)

	if apiErr.Response.Error.Message != "invalid request" {
		t.Errorf("Expected error message to be 'invalid request', got %q", apiErr.Response.Error.Message)
	}

	if apiErr.Response.Error.Type != "invalid_request_error" {
		t.Errorf("Expected error type to be 'invalid_request_error', got %q", apiErr.Response.Error.Type)
	}

	if apiErr.Response.EventID != "evt_123" {
		t.Errorf("Expected event ID to be 'evt_123', got %q", apiErr.Response.EventID)
	}

	if apiErr.Response.Type != "error" {
		t.Errorf("Expected response type to be 'error', got %q", apiErr.Response.Type)
	}

	if apiErr.Response.Error.Code != "invalid_input" {
		t.Errorf("Expected error code to be 'invalid_input', got %q", apiErr.Response.Error.Code)
	}
}

func TestDoWithContextCancellation(t *testing.T) {
	// Create a test server with a delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sleep to simulate a slow response
		time.Sleep(500 * time.Millisecond)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true, "message": "request processed"}`))
	}))
	defer server.Close()

	// Create request
	req := &testRequest{
		Field1: "test",
		Field2: 123,
	}

	// Create a context with a short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Make the request
	_, err := Do[testRequest, testResponse](
		ctx,
		server.URL+"/test",
		req,
		WithClient(server.Client()),
	)

	// Verify error
	if err == nil {
		t.Fatal("Expected error due to context cancellation, got nil")
	}

	// Check if the error contains context.DeadlineExceeded
	if err.Error() == "" || !containsDeadlineExceeded(err.Error()) {
		t.Errorf("Expected error to include context deadline exceeded, got: %v", err)
	}
}

// Helper function to check if error string contains deadline exceeded message
func containsDeadlineExceeded(errStr string) bool {
	return strings.Contains(errStr, "context deadline exceeded")
}
