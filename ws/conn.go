// Package ws provides WebSocket connection handling for the OpenAI Realtime API.
// It implements low-level connection management, message reading/writing,
// and event handling for the WebSocket protocol.
//
// This package serves as the transport layer for the OpenAI Realtime API, handling:
//   - WebSocket connection establishment and management
//   - Message serialization and deserialization
//   - Binary and text message support
//   - Connection lifecycle (opening, maintaining, closing)
//   - Error handling and recovery
//
// The ws package is designed to be a low-level component used by higher-level
// packages like messaging. Most users should not need to interact with this
// package directly unless implementing custom connection management.
//
// Example usage (advanced use case):
//
//	// Create a connection from a websocket.Conn
//	wsConn := ws.NewConn(rawConn)
//
//	// Set a logger
//	wsConn.SetLogger(logger.Default)
//
//	// Send a message
//	err := wsConn.SendRaw(ctx, ws.MessageText, []byte(`{"type":"ping"}`))
//
//	// Read a message
//	msgType, data, err := wsConn.ReadRaw(ctx)
//
// The package abstracts away the differences between different WebSocket implementations,
// providing a consistent interface for the OpenAI Realtime API client.
package ws

import (
	"context"
	"sync"

	"github.com/Mliviu79/openai-realtime-go/logger"
)

// Conn is a generic WebSocket connection wrapper.
// It provides thread-safe methods for sending and receiving messages over a WebSocket connection.
// Conn implements connection management, including thread safety, logging, and error handling.
type Conn struct {
	mu     sync.RWMutex
	logger logger.Logger
	conn   WebSocketConn
}

// NewConn creates a new Conn instance
// It wraps the provided WebSocketConn with thread-safe methods and optional logging.
func NewConn(conn WebSocketConn) *Conn {
	return &Conn{
		conn: conn,
	}
}

// SetLogger sets the logger for the connection
// The logger is used to log WebSocket operations for debugging purposes.
// If nil, no logging is performed.
func (c *Conn) SetLogger(logger logger.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
}

// Close closes the connection.
// This method is thread-safe and can be called from any goroutine.
// After closing, no more messages can be sent or received.
func (c *Conn) Close() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn.Close()
}

// SendRaw sends a raw message to the server.
// This is a low-level method that takes a message type (text or binary) and raw byte data.
// Most users should use higher-level methods that handle serialization.
// This method is thread-safe and can be called from any goroutine.
func (c *Conn) SendRaw(ctx context.Context, messageType MessageType, data []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.logger != nil {
		c.logger.Debugf("sending raw message: type=%s data=%s", messageType.String(), string(data))
	}

	return c.conn.WriteMessage(ctx, messageType, data)
}

// ReadRaw reads a raw message from the server.
// This is a low-level method that returns the message type and raw byte data.
// Most users should use higher-level methods that handle deserialization.
// This method is thread-safe and can be called from any goroutine.
// It will block until a message is received, the context is canceled, or an error occurs.
func (c *Conn) ReadRaw(ctx context.Context) (MessageType, []byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	messageType, data, err := c.conn.ReadMessage(ctx)
	if err != nil {
		return 0, nil, err
	}

	if c.logger != nil {
		c.logger.Debugf("received raw message: type=%s data=%s", messageType.String(), string(data))
	}

	return messageType, data, nil
}

// Ping sends a ping message to the WebSocket connection.
// This can be used to keep the connection alive or check if it's still operational.
// This method is thread-safe and can be called from any goroutine.
func (c *Conn) Ping(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn.Ping(ctx)
}
