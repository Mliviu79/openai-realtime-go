package messaging

import (
	"context"
	"strings"
	"testing"

	"github.com/Mliviu79/openai-realtime-go/logger"
	"github.com/Mliviu79/openai-realtime-go/messages/incoming"
	"github.com/Mliviu79/openai-realtime-go/ws"
)

// MockConn implements the ws.WebSocketConn interface for testing
type MockConn struct {
	WriteMessageFunc func(ctx context.Context, messageType ws.MessageType, data []byte) error
	ReadMessageFunc  func(ctx context.Context) (ws.MessageType, []byte, error)
	CloseFunc        func() error
	PingFunc         func(ctx context.Context) error
}

func (m *MockConn) WriteMessage(ctx context.Context, messageType ws.MessageType, data []byte) error {
	if m.WriteMessageFunc != nil {
		return m.WriteMessageFunc(ctx, messageType, data)
	}
	return nil
}

func (m *MockConn) ReadMessage(ctx context.Context) (ws.MessageType, []byte, error) {
	if m.ReadMessageFunc != nil {
		return m.ReadMessageFunc(ctx)
	}
	return ws.MessageText, nil, nil
}

func (m *MockConn) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *MockConn) Ping(ctx context.Context) error {
	if m.PingFunc != nil {
		return m.PingFunc(ctx)
	}
	return nil
}

// MockLogger implements the logger.Logger interface for testing
type MockLogger struct {
	DebugfFunc     func(format string, args ...any)
	InfofFunc      func(format string, args ...any)
	WarnfFunc      func(format string, args ...any)
	ErrorfFunc     func(format string, args ...any)
	WithFieldFunc  func(key string, value any) logger.Logger
	WithFieldsFunc func(fields map[string]any) logger.Logger
}

func (m *MockLogger) Debugf(format string, args ...any) {
	if m.DebugfFunc != nil {
		m.DebugfFunc(format, args...)
	}
}

func (m *MockLogger) Infof(format string, args ...any) {
	if m.InfofFunc != nil {
		m.InfofFunc(format, args...)
	}
}

func (m *MockLogger) Warnf(format string, args ...any) {
	if m.WarnfFunc != nil {
		m.WarnfFunc(format, args...)
	}
}

func (m *MockLogger) Errorf(format string, args ...any) {
	if m.ErrorfFunc != nil {
		m.ErrorfFunc(format, args...)
	}
}

func (m *MockLogger) WithField(key string, value any) logger.Logger {
	if m.WithFieldFunc != nil {
		return m.WithFieldFunc(key, value)
	}
	return m
}

func (m *MockLogger) WithFields(fields map[string]any) logger.Logger {
	if m.WithFieldsFunc != nil {
		return m.WithFieldsFunc(fields)
	}
	return m
}

func TestNewClient(t *testing.T) {
	// Create mock connection
	conn := ws.NewConn(&MockConn{})

	// Create a client
	client := NewClient(conn)

	if client == nil {
		t.Fatal("Expected client to not be nil")
	}
}

func TestSetLogger(t *testing.T) {
	// Create a mock connection
	mockConn := &MockConn{}
	conn := ws.NewConn(mockConn)

	// Create a client
	client := NewClient(conn)

	// Create a mock logger
	mockLogger := &MockLogger{}

	// Set the logger
	client.SetLogger(mockLogger)

	// Indirectly verify by ensuring no panic occurred
	if !client.mu.TryLock() {
		t.Error("Expected mutex to be unlocked")
	} else {
		client.mu.Unlock()
	}
}

func TestReadMessage(t *testing.T) {
	// Create a mock connection that returns a predefined message
	mockConn := &MockConn{
		ReadMessageFunc: func(ctx context.Context) (ws.MessageType, []byte, error) {
			// Return a message that will be parsed as an error message
			return ws.MessageText, []byte(`{"type":"error","error":{"type":"server_error","message":"test error"}}`), nil
		},
	}
	conn := ws.NewConn(mockConn)

	// Create a client
	client := NewClient(conn)

	// Read a message
	ctx := context.Background()
	msg, err := client.ReadMessage(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that we got the expected message type
	if msg.RcvdMsgType() != incoming.RcvdMsgTypeError {
		t.Errorf("Expected message type to be error, got %v", msg.RcvdMsgType())
	}
}

func TestSendText(t *testing.T) {
	// Create a mock connection that verifies the sent message
	textMessageSent := false
	mockConn := &MockConn{
		WriteMessageFunc: func(ctx context.Context, messageType ws.MessageType, data []byte) error {
			if messageType != ws.MessageText {
				t.Errorf("Expected message type to be MessageText, got %v", messageType)
			}

			// Verify the message is a conversation item create message
			dataStr := string(data)
			if strings.Contains(dataStr, `"type":"conversation.item.create"`) &&
				strings.Contains(dataStr, `"content":[{"type":"input_text","text":"test message"}]`) &&
				strings.Contains(dataStr, `"role":"user"`) {
				textMessageSent = true
			}

			return nil
		},
	}
	conn := ws.NewConn(mockConn)

	// Create a client
	client := NewClient(conn)

	// Send a text message
	err := client.SendText(context.Background(), "test message")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !textMessageSent {
		t.Error("Expected text message to be sent, but it wasn't")
	}
}

func TestClose(t *testing.T) {
	// Create a mock connection that verifies Close is called
	closeWasCalled := false
	mockConn := &MockConn{
		CloseFunc: func() error {
			closeWasCalled = true
			return nil
		},
	}
	conn := ws.NewConn(mockConn)

	// Create a client
	client := NewClient(conn)

	// Close the client
	err := client.Close()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !closeWasCalled {
		t.Error("Expected Close to be called on the underlying connection, but it wasn't")
	}
}
