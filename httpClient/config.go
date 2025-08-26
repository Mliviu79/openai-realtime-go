package httpClient

import (
	"fmt"
	"net/http"
	"net/url"
)

// ClientConfig holds the configuration for the HTTP client
// It contains settings for authentication, API endpoints, and HTTP client configuration.
type ClientConfig struct {
	authToken  string       // Authentication token for the API
	BaseURL    string       // Base URL for the API
	APIType    APIType      // Type of API (OpenAI or Azure)
	APIVersion string       // API version (used for Azure)
	HTTPClient *http.Client // HTTP client to use for requests
	APIBaseURL string       // Base URL for the REST API
}

// DefaultConfig creates a default configuration with the given auth token
// for the OpenAI API.
//
// Parameters:
//   - authToken: The authentication token for the OpenAI API
//
// Returns:
//   - ClientConfig: A configuration for the OpenAI API
func DefaultConfig(authToken string) ClientConfig {
	return ClientConfig{
		authToken:  authToken,
		BaseURL:    OpenaiRealtimeAPIURLv1,
		APIType:    APITypeOpenAI,
		HTTPClient: http.DefaultClient,
		APIBaseURL: OpenaiAPIURLv1,
	}
}

// DefaultAzureConfig creates a default configuration for Azure OpenAI API
//
// Parameters:
//   - apiKey: The API key for Azure OpenAI
//   - baseURL: The base URL for the Azure OpenAI endpoint
//
// Returns:
//   - ClientConfig: A configuration for the Azure OpenAI API
func DefaultAzureConfig(apiKey, baseURL string) ClientConfig {
	return ClientConfig{
		authToken:  apiKey,
		BaseURL:    baseURL,
		APIType:    APITypeAzure,
		APIVersion: azureAPIVersion20241001Preview,
		HTTPClient: http.DefaultClient,
		APIBaseURL: baseURL,
	}
}

// GetHeaders returns the appropriate headers based on API type
//
// Parameters:
//   - config: The client configuration
//
// Returns:
//   - http.Header: The headers to use for the request
func GetHeaders(config ClientConfig) http.Header {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	// Use a map of API types to header setting functions
	headerSetters := map[APIType]func(http.Header, string){
		APITypeAzure: func(h http.Header, token string) {
			h.Set("api-key", token)
		},
		APITypeOpenAI: func(h http.Header, token string) {
			h.Set("Authorization", "Bearer "+token)
		},
	}

	// Apply the appropriate header setter based on API type
	if setter, exists := headerSetters[config.APIType]; exists {
		setter(headers, config.authToken)
	}

	return headers
}

// GetURL constructs the appropriate URL based on API type and model
//
// Parameters:
//   - config: The client configuration
//   - model: The model to use
//
// Returns:
//   - string: The URL to use for the request
func GetURL(config ClientConfig, model string) string {
	query := url.Values{}

	// Use a map of API types to query parameter setting functions
	querySetters := map[APIType]func(url.Values, string, ClientConfig){
		APITypeAzure: func(q url.Values, modelName string, cfg ClientConfig) {
			q.Set("api-version", cfg.APIVersion)
			q.Set("deployment", modelName)
		},
		APITypeOpenAI: func(q url.Values, modelName string, cfg ClientConfig) {
			q.Set("model", modelName)
		},
	}

	// Apply the appropriate query setter based on API type
	if setter, exists := querySetters[config.APIType]; exists {
		setter(query, model, config)
	}

	return config.BaseURL + "?" + query.Encode()
}

// String returns a string representation of the ClientConfig
// It masks sensitive information like auth tokens for security
func (c ClientConfig) String() string {
	// Include HTTP client timeout if available
	timeout := "default"
	if c.HTTPClient != nil && c.HTTPClient.Timeout > 0 {
		timeout = c.HTTPClient.Timeout.String()
	}

	// Get a formatted string for the API type
	apiTypeStr := string(c.APIType)
	if c.APIType == APITypeOpenAI {
		apiTypeStr = "openai"
	} else if c.APIType == APITypeAzure {
		apiTypeStr = "azure"
	}

	// Build a more comprehensive string representation
	return fmt.Sprintf(
		"ClientConfig{\n"+
			"  AuthToken: %s\n"+
			"  BaseURL: %s\n"+
			"  APIType: %s\n"+
			"  APIVersion: %s\n"+
			"  APIBaseURL: %s\n"+
			"  HTTPClient: {Timeout: %s}\n"+
			"}",
		valueOrEmpty(c.authToken),
		valueOrEmpty(c.BaseURL),
		valueOrEmpty(apiTypeStr),
		valueOrEmpty(c.APIVersion),
		valueOrEmpty(c.APIBaseURL),
		valueOrEmpty(timeout),
	)
}

// Helper function to handle empty strings
func valueOrEmpty(s string) string {
	if s == "" {
		return "<empty>"
	}
	return s
}
