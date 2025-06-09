package ws

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Error definitions
var (
	ErrUnsupportedMessageType = errors.New("unsupported message type")
)

// Mapping between gorilla message types and our internal message types
var (
	// gorillaToMessageType maps gorilla's message types to our MessageType
	gorillaToMessageType = map[int]MessageType{
		websocket.TextMessage:   MessageText,
		websocket.BinaryMessage: MessageBinary,
	}

	// messageTypeToGorilla maps our MessageType to gorilla's message types
	messageTypeToGorilla = map[MessageType]int{
		MessageText:   websocket.TextMessage,
		MessageBinary: websocket.BinaryMessage,
	}
)

// GorillaWebSocketOptions is the options for GorillaWebSocketConn.
type GorillaWebSocketOptions struct {
	// ReadLimit is the maximum size of a message in bytes. -1 means no limit. Default is -1.
	ReadLimit int64
	// Dialer is the websocket dialer to use. If nil, websocket.DefaultDialer will be used.
	Dialer *websocket.Dialer
}

// GorillaWebSocketDialer is a WebSocket dialer implementation based on gorilla/websocket.
type GorillaWebSocketDialer struct {
	options GorillaWebSocketOptions
}

// NewGorillaWebSocketDialer creates a new GorillaWebSocketDialer.
func NewGorillaWebSocketDialer(options GorillaWebSocketOptions) *GorillaWebSocketDialer {
	// set default read limit
	if options.ReadLimit <= 0 {
		options.ReadLimit = -1
	}
	return &GorillaWebSocketDialer{
		options: options,
	}
}

// Dial establishes a new WebSocket connection to the given URL.
func (d *GorillaWebSocketDialer) Dial(ctx context.Context, url string, header http.Header) (WebSocketConn, error) {
	dialer := d.options.Dialer
	if dialer == nil {
		dialer = websocket.DefaultDialer
	}

	conn, resp, err := dialer.DialContext(ctx, url, header)
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return nil, err
	}

	if d.options.ReadLimit > 0 {
		conn.SetReadLimit(d.options.ReadLimit)
	}

	return &GorillaWebSocketConn{conn: conn, resp: resp, options: d.options}, nil
}

// GorillaWebSocketConn is a WebSocket connection implementation based on gorilla/websocket.
type GorillaWebSocketConn struct {
	conn    *websocket.Conn
	resp    *http.Response
	options GorillaWebSocketOptions
}

// ReadMessage reads a message from the WebSocket connection.
func (c *GorillaWebSocketConn) ReadMessage(ctx context.Context) (MessageType, []byte, error) {
	// Set up context cancellation
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			c.conn.Close()
		case <-done:
		}
	}()
	defer close(done)

	messageType, data, err := c.conn.ReadMessage()
	if err != nil {
		return 0, nil, err
	}

	// Map gorilla message type to our message type
	if ourType, ok := gorillaToMessageType[messageType]; ok {
		return ourType, data, nil
	}

	return 0, nil, ErrUnsupportedMessageType
}

// WriteMessage writes a message to the WebSocket connection.
func (c *GorillaWebSocketConn) WriteMessage(ctx context.Context, messageType MessageType, data []byte) error {
	// Set up context cancellation
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			c.conn.Close()
		case <-done:
		}
	}()
	defer close(done)

	// Map our message type to gorilla message type
	gorillaType, ok := messageTypeToGorilla[messageType]
	if !ok {
		return ErrUnsupportedMessageType
	}

	return c.conn.WriteMessage(gorillaType, data)
}

// Close closes the WebSocket connection.
func (c *GorillaWebSocketConn) Close() error {
	return c.conn.Close()
}

// Ping sends a ping message to the WebSocket connection.
func (c *GorillaWebSocketConn) Ping(ctx context.Context) error {
	// Set up context cancellation
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			c.conn.Close()
		case <-done:
		}
	}()
	defer close(done)

	deadline := time.Now().Add(59 * time.Second)
	return c.conn.WriteControl(websocket.PingMessage, []byte{}, deadline)
}
