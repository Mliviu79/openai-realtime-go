package httpClient

import (
	"net/http"
	"time"
)

// Default configuration constants
const (
	DefaultRequestTimeout = 500 * time.Millisecond // Fast timeout for API calls
)

// RetryConfig defines the retry behavior for HTTP requests
type RetryConfig struct {
	MaxRetries           int           // Maximum number of retry attempts
	RetryDelay           time.Duration // Base delay between retries
	MaxDelay             time.Duration // Maximum delay between retries
	RetryableStatusCodes []int         // HTTP status codes that should trigger a retry
}

// DefaultRetryConfig returns a sensible default retry configuration
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:           3,
		RetryDelay:           1 * time.Second,
		MaxDelay:             30 * time.Second,
		RetryableStatusCodes: []int{408, 429, 500, 502, 503, 504},
	}
}

// option holds the configuration for an HTTP request
type option struct {
	client      *http.Client
	headers     http.Header
	method      string
	timeout     time.Duration
	retryConfig RetryConfig
}

// defaultOption returns an option with sensible defaults
// This function is used in the Do function in client.go to initialize options
func defaultOption() option {
	return option{
		client:      &http.Client{Timeout: DefaultRequestTimeout},
		headers:     http.Header{},
		method:      http.MethodPost,
		timeout:     DefaultRequestTimeout,
		retryConfig: DefaultRetryConfig(),
	}
}

// HTTPOption is a function that configures an HTTP request option
type HTTPOption func(*option)

// WithHeaders sets the headers for the HTTP request
// Parameters:
//   - headers: The HTTP headers to use for the request
func WithHeaders(headers http.Header) HTTPOption {
	return func(o *option) {
		// Use WithHeader for each header to ensure consistent behavior
		for key, values := range headers {
			WithHeader(key, values[0])(o) // Apply the returned function to 'o'
		}
	}
}

// WithHeader sets a single header for the HTTP request
// Parameters:
//   - key: The header key
//   - value: The header value
func WithHeader(key, value string) HTTPOption {
	return func(o *option) {
		if o.headers == nil {
			o.headers = http.Header{}
		}
		o.headers.Set(key, value)
	}
}

// WithClient sets the HTTP client for the request
// Parameters:
//   - client: The HTTP client to use for the request
func WithClient(client *http.Client) HTTPOption {
	return func(o *option) {
		if client != nil {
			o.client = client
		}
	}
}

// WithMethod sets the HTTP method for the request
// Parameters:
//   - method: The HTTP method to use (e.g., GET, POST, PUT)
func WithMethod(method string) HTTPOption {
	return func(o *option) {
		if method != "" {
			o.method = method
		}
	}
}

// WithTimeout sets a timeout for the HTTP request
// Parameters:
//   - timeout: The timeout duration for the request
func WithTimeout(timeout time.Duration) HTTPOption {
	return func(o *option) {
		if timeout > 0 {
			o.timeout = timeout

			// If we're using the default client, create a new one with the timeout
			if o.client == http.DefaultClient {
				o.client = &http.Client{Timeout: timeout}
			} else if o.client != nil {
				// Otherwise, update the existing client's timeout
				client := *o.client
				client.Timeout = timeout
				o.client = &client
			}
		}
	}
}

// WithRetryConfig sets the retry configuration for the HTTP request
// Parameters:
//   - config: The retry configuration to use
func WithRetryConfig(config RetryConfig) HTTPOption {
	return func(o *option) {
		o.retryConfig = config
	}
}

// WithBasicAuth sets basic authentication headers for the HTTP request
// Parameters:
//   - username: The username for basic auth
//   - password: The password for basic auth
func WithBasicAuth(username, password string) HTTPOption {
	return func(o *option) {
		if o.headers == nil {
			o.headers = http.Header{}
		}
		req := &http.Request{Header: o.headers}
		req.SetBasicAuth(username, password)
		o.headers = req.Header
	}
}

// WithBearerAuth sets bearer token authentication for the HTTP request
// Parameters:
//   - token: The bearer token for authentication
func WithBearerAuth(token string) HTTPOption {
	return func(o *option) {
		if token != "" {
			if o.headers == nil {
				o.headers = http.Header{}
			}
			o.headers.Set("Authorization", "Bearer "+token)
		}
	}
}
