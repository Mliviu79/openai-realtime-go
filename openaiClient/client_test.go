package openaiClient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mliviu79/openai-realtime-go/httpClient"
	"github.com/Mliviu79/openai-realtime-go/session"
)

func TestNewClient(t *testing.T) {
	const authToken = "test-token"
	client := NewClient(authToken)

	if client == nil {
		t.Fatal("Expected client to not be nil")
	}

	// Cannot directly access authToken as it's a private field
	// Test indirectly by checking if the client can be created
}

func TestNewClientWithConfig(t *testing.T) {
	config := httpClient.ClientConfig{
		BaseURL:    "wss://custom.example.com/v1",
		APIBaseURL: "https://custom-api.example.com",
		HTTPClient: &http.Client{},
	}

	client := NewClientWithConfig(config)

	if client == nil {
		t.Fatal("Expected client to not be nil")
	}

	if client.config.BaseURL != config.BaseURL {
		t.Errorf("Expected base URL to be %q, got %q", config.BaseURL, client.config.BaseURL)
	}

	if client.config.APIBaseURL != config.APIBaseURL {
		t.Errorf("Expected API base URL to be %q, got %q", config.APIBaseURL, client.config.APIBaseURL)
	}
}

func TestCreateSession(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/realtime/sessions" {
			t.Errorf("Expected request to '/realtime/sessions', got %q", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header to be 'Bearer test-token', got %q", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "test-session-id"}`))
	}))
	defer server.Close()

	// Create a client that points to the test server
	config := httpClient.ClientConfig{
		BaseURL:    "wss://api.openai.com/v1",
		APIBaseURL: server.URL,
		HTTPClient: server.Client(),
	}

	// Set the auth token using DefaultConfig and then override the URLs
	customConfig := httpClient.DefaultConfig("test-token")
	customConfig.BaseURL = config.BaseURL
	customConfig.APIBaseURL = config.APIBaseURL
	customConfig.HTTPClient = config.HTTPClient

	client := NewClientWithConfig(customConfig)

	// Make the request with a model specified
	modelName := session.Model("gpt-4o")
	req := &session.CreateRequest{
		SessionRequest: session.SessionRequest{
			Model: &modelName,
		},
	}

	resp, err := client.CreateSession(context.Background(), req)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.ID != "test-session-id" {
		t.Errorf("Expected session ID to be 'test-session-id', got %q", resp.ID)
	}
}

func TestConnectOptions(t *testing.T) {
	tests := []struct {
		name           string
		option         ConnectOption
		expectedField  string
		expectedValue  interface{}
		initialOptions *connectOptions
	}{
		{
			name:           "WithModel",
			option:         WithModel("gpt-4o"),
			expectedField:  "model",
			expectedValue:  "gpt-4o",
			initialOptions: &connectOptions{},
		},
		{
			name:           "WithSessionID",
			option:         WithSessionID("test-session"),
			expectedField:  "sessionID",
			expectedValue:  "test-session",
			initialOptions: &connectOptions{},
		},
		{
			name:           "WithReadLimit",
			option:         WithReadLimit(1024),
			expectedField:  "readLimit",
			expectedValue:  int64(1024),
			initialOptions: &connectOptions{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.option(tt.initialOptions)

			switch tt.expectedField {
			case "model":
				if tt.initialOptions.model != tt.expectedValue {
					t.Errorf("Expected model to be %v, got %v", tt.expectedValue, tt.initialOptions.model)
				}
			case "sessionID":
				if tt.initialOptions.sessionID != tt.expectedValue {
					t.Errorf("Expected sessionID to be %v, got %v", tt.expectedValue, tt.initialOptions.sessionID)
				}
			case "readLimit":
				if tt.initialOptions.readLimit != tt.expectedValue {
					t.Errorf("Expected readLimit to be %v, got %v", tt.expectedValue, tt.initialOptions.readLimit)
				}
			}
		})
	}
}
