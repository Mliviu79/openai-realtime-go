// Package ws provides WebSocket functionality for the OpenAI Realtime API.
package ws

import (
	"context"
	"net/http"
)

// WebSocketDialer is the interface for WebSocket dialers.
type WebSocketDialer interface {
	// Dial establishes a new WebSocket connection to the given URL.
	// The ctx can be used to cancel or timeout the request.
	// The header can be used to pass additional HTTP headers with the request.
	Dial(ctx context.Context, url string, header http.Header) (WebSocketConn, error)
}

// DefaultReadLimit is the default maximum size of a message in bytes.
const DefaultReadLimit int64 = 32 * 1024 // 32KB read limit by default

// DialerOptions contains configuration options for WebSocket dialers
type DialerOptions struct {
	// ReadLimit is the maximum size of a message in bytes
	// If set to 0 or negative, the underlying implementation will use its default
	// For Gorilla WebSocket, this means -1 (no limit)
	ReadLimit int64
}

// DefaultDialer returns a default WebSocket dialer
// By default, this uses the multiplexed dialer for better resource utilization
func DefaultDialer() WebSocketDialer {
	return DirectDialer(DialerOptions{})
}

// DirectDialer returns a direct (non-multiplexed) WebSocket dialer
// This is useful for cases where multiplexing is not desired
func DirectDialer(options DialerOptions) WebSocketDialer {
	// Pass the ReadLimit directly to the Gorilla implementation
	// The Gorilla implementation will handle the default value if ReadLimit <= 0
	return NewGorillaWebSocketDialer(GorillaWebSocketOptions{
		ReadLimit: options.ReadLimit,
	})
}
