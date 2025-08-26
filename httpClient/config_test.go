package httpClient

import (
	"net/http"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	const authToken = "test-token"
	config := DefaultConfig(authToken)

	// Verify the config was created with the correct values
	if config.BaseURL != OpenaiRealtimeAPIURLv1 {
		t.Errorf("Expected BaseURL to be %q, got %q", OpenaiRealtimeAPIURLv1, config.BaseURL)
	}

	if config.APIType != APITypeOpenAI {
		t.Errorf("Expected APIType to be %v, got %v", APITypeOpenAI, config.APIType)
	}

	if config.APIBaseURL != OpenaiAPIURLv1 {
		t.Errorf("Expected APIBaseURL to be %q, got %q", OpenaiAPIURLv1, config.APIBaseURL)
	}
}

func TestDefaultAzureConfig(t *testing.T) {
	const apiKey = "test-api-key"
	const baseURL = "https://test.openai.azure.com/openai"
	config := DefaultAzureConfig(apiKey, baseURL)

	// Verify the config was created with the correct values
	if config.BaseURL != baseURL {
		t.Errorf("Expected BaseURL to be %q, got %q", baseURL, config.BaseURL)
	}

	if config.APIType != APITypeAzure {
		t.Errorf("Expected APIType to be %v, got %v", APITypeAzure, config.APIType)
	}

	if !strings.HasPrefix(config.APIBaseURL, baseURL) {
		t.Errorf("Expected APIBaseURL to start with %q, got %q", baseURL, config.APIBaseURL)
	}
}

func TestGetHeaders(t *testing.T) {
	const authToken = "test-token"
	config := DefaultConfig(authToken)
	headers := GetHeaders(config)

	// Verify the headers were created correctly
	expectedAuthHeader := "Bearer test-token"
	if headers.Get("Authorization") != expectedAuthHeader {
		t.Errorf("Expected Authorization header to be %q, got %q", expectedAuthHeader, headers.Get("Authorization"))
	}

	if headers.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be 'application/json', got %q", headers.Get("Content-Type"))
	}

	if beta := headers.Get("OpenAI-Beta"); beta != "" {
		t.Errorf("Expected no OpenAI-Beta header, got %q", beta)
	}
}

func TestGetHeadersAzure(t *testing.T) {
	const apiKey = "test-api-key"
	const baseURL = "https://test.openai.azure.com/openai"
	config := DefaultAzureConfig(apiKey, baseURL)
	headers := GetHeaders(config)

	// Verify the headers were created correctly for Azure
	expectedApiKeyHeader := "test-api-key"
	if headers.Get("api-key") != expectedApiKeyHeader {
		t.Errorf("Expected api-key header to be %q, got %q", expectedApiKeyHeader, headers.Get("api-key"))
	}

	if headers.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be 'application/json', got %q", headers.Get("Content-Type"))
	}
}

func TestConfigString(t *testing.T) {
	config := ClientConfig{
		BaseURL:    "https://api.example.com/v1",
		APIType:    APITypeOpenAI,
		APIVersion: "2023-05-15",
		HTTPClient: http.DefaultClient,
		APIBaseURL: "https://api.example.com/v1",
	}

	// Call String() method
	str := config.String()

	// Verify the string representation contains the expected information
	expectedSubstrings := []string{
		"BaseURL: https://api.example.com/v1",
		"APIType: openai",
		"APIVersion: 2023-05-15",
		"APIBaseURL: https://api.example.com/v1",
	}

	for _, substr := range expectedSubstrings {
		if !strings.Contains(str, substr) {
			t.Errorf("Expected string representation to contain %q, but it was not found in %q", substr, str)
		}
	}
}
