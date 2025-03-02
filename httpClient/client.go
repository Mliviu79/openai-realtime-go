// Package httputil provides HTTP client utilities for making API requests.
// It handles common HTTP operations like request creation, response parsing, and error handling.
package httpClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"

	"github.com/Mliviu79/openai-realtime-go/apierrs"
	"github.com/rs/zerolog/log"
)

// prepareRequest creates and configures an HTTP request
func prepareRequest[Q any](ctx context.Context, method, url string, req *Q, headers http.Header) (*http.Request, error) {
	var requestBody io.Reader
	if req != nil {
		data, err := json.Marshal(req)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		requestBody = bytes.NewReader(data)
	}

	request, err := http.NewRequestWithContext(ctx, method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Content-Type header if not already set and we have a request body
	if req != nil && request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	// Merge headers
	maps.Copy(request.Header, headers)

	return request, nil
}

// processResponse handles the HTTP response and unmarshals the body
func processResponse[R any](response *http.Response) (*R, error) {
	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle non-200 status codes
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		// Try to unmarshal as an APIError
		var apiErr apierrs.APIError
		if err := json.Unmarshal(body, &apiErr.Response); err == nil {
			// Check if this looks like a valid API error
			if apiErr.Response.Type == "error" &&
				apiErr.Response.Error.Message != "" {
				return nil, &apiErr
			}
		}

		return nil, fmt.Errorf("request failed with status %d: %s", response.StatusCode, string(body))
	}

	// Skip unmarshaling for empty responses
	if len(body) == 0 {
		var resp R
		return &resp, nil
	}

	var resp R
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w, body: %s", err, string(body))
	}

	return &resp, nil
}

// Do performs an HTTP request with the given options.
// It handles request creation, execution, response parsing, and error handling.
//
// Generic type parameters:
//   - Q: The request type to be marshaled to JSON
//   - R: The response type to be unmarshaled from JSON
//
// Parameters:
//   - ctx: The context for the request
//   - url: The URL to send the request to
//   - req: The request body to be marshaled to JSON
//   - opts: Optional configuration for the request
//
// Returns:
//   - *R: A pointer to the unmarshaled response
//   - error: An error if the request failed
func Do[Q any, R any](ctx context.Context, url string, req *Q, opts ...HTTPOption) (*R, error) {
	// Use defaultOption() instead of direct initialization
	opt := defaultOption()

	// Apply any custom options
	for _, o := range opts {
		o(&opt)
	}

	// Prepare the request
	request, err := prepareRequest(ctx, opt.method, url, req, opt.headers)
	if err != nil {
		return nil, err
	}

	// Log request details
	log.Debug().
		Str("url", url).
		Str("method", opt.method).
		Interface("request", req).
		Msg("Sending request")

	// Execute the request
	response, err := opt.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("http failed: %w", err)
	}
	defer response.Body.Close()

	// Process the response
	resp, err := processResponse[R](response)
	if err != nil {
		return nil, err
	}

	// Log response details
	log.Debug().
		Int("status", response.StatusCode).
		Interface("response", resp).
		Msg("Received response")

	return resp, nil
}
