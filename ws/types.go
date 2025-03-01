package ws

import (
	"context"
)

// MessageType represents the type of WebSocket message
type MessageType int

const (
	// MessageText denotes a text data message.
	MessageText MessageType = iota + 1
	// MessageBinary denotes a binary data message.
	MessageBinary
	// MessageClose denotes a close control message.
	MessageClose
	// MessagePing denotes a ping control message.
	MessagePing
	// MessagePong denotes a pong control message.
	MessagePong
)

// messageTypeNames maps MessageType values to their string representations.
var messageTypeNames = map[MessageType]string{
	MessageText:   "text",
	MessageBinary: "binary",
	MessageClose:  "close",
	MessagePing:   "ping",
	MessagePong:   "pong",
}

// String returns a string representation of the MessageType.
func (m MessageType) String() string {
	if name, ok := messageTypeNames[m]; ok {
		return name
	}
	return "unknown"
}

// WebSocketConn represents a WebSocket connection interface
type WebSocketConn interface {
	WriteMessage(ctx context.Context, messageType MessageType, data []byte) error
	ReadMessage(ctx context.Context) (messageType MessageType, data []byte, err error)
	Close() error
	Ping(ctx context.Context) error
}

// Logger interface for connection logging
type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}
