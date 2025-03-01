package ws

import (
	"context"
	"testing"
)

func TestNewConn(t *testing.T) {
	// Create a mock websocket connection
	mockConn := &MockWebSocketConn{}

	// Create a new Conn
	conn := NewConn(mockConn)

	if conn == nil {
		t.Fatal("Expected conn to not be nil")
	}
}

func TestConnSendRaw(t *testing.T) {
	// Create a mock websocket connection that records the sent messages
	var capturedMessageType MessageType
	var capturedData []byte
	mockConn := &MockWebSocketConn{
		WriteMessageFunc: func(ctx context.Context, messageType MessageType, data []byte) error {
			capturedMessageType = messageType
			capturedData = data
			return nil
		},
	}

	// Create a new Conn with the mock
	conn := NewConn(mockConn)

	// Send a text message
	ctx := context.Background()
	err := conn.SendRaw(ctx, MessageText, []byte("test message"))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if capturedMessageType != MessageText {
		t.Errorf("Expected message type to be MessageText, got %v", capturedMessageType)
	}

	if string(capturedData) != "test message" {
		t.Errorf("Expected message to be 'test message', got %q", string(capturedData))
	}
}

func TestConnReadRaw(t *testing.T) {
	// Create a mock websocket connection that returns a predefined message
	mockConn := &MockWebSocketConn{
		ReadMessageFunc: func(ctx context.Context) (MessageType, []byte, error) {
			return MessageBinary, []byte("binary data"), nil
		},
	}

	// Create a new Conn with the mock
	conn := NewConn(mockConn)

	// Read a message
	ctx := context.Background()
	messageType, data, err := conn.ReadRaw(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if messageType != MessageBinary {
		t.Errorf("Expected message type to be MessageBinary, got %v", messageType)
	}

	if string(data) != "binary data" {
		t.Errorf("Expected data to be 'binary data', got %q", string(data))
	}
}

func TestConnClose(t *testing.T) {
	// Create a mock websocket connection that records the close
	closeWasCalled := false
	mockConn := &MockWebSocketConn{
		CloseFunc: func() error {
			closeWasCalled = true
			return nil
		},
	}

	// Create a new Conn with the mock
	conn := NewConn(mockConn)

	// Close the connection
	err := conn.Close()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !closeWasCalled {
		t.Error("Expected Close to be called on the underlying connection")
	}
}

func TestConnPing(t *testing.T) {
	// Create a mock websocket connection that records the ping
	pingWasCalled := false
	mockConn := &MockWebSocketConn{
		PingFunc: func(ctx context.Context) error {
			pingWasCalled = true
			return nil
		},
	}

	// Create a new Conn with the mock
	conn := NewConn(mockConn)

	// Ping the connection
	ctx := context.Background()
	err := conn.Ping(ctx)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !pingWasCalled {
		t.Error("Expected Ping to be called on the underlying connection")
	}
}

// MockWebSocketConn implements the WebSocketConn interface for testing
type MockWebSocketConn struct {
	WriteMessageFunc func(ctx context.Context, messageType MessageType, data []byte) error
	ReadMessageFunc  func(ctx context.Context) (MessageType, []byte, error)
	CloseFunc        func() error
	PingFunc         func(ctx context.Context) error
}

func (m *MockWebSocketConn) WriteMessage(ctx context.Context, messageType MessageType, data []byte) error {
	if m.WriteMessageFunc != nil {
		return m.WriteMessageFunc(ctx, messageType, data)
	}
	return nil
}

func (m *MockWebSocketConn) ReadMessage(ctx context.Context) (MessageType, []byte, error) {
	if m.ReadMessageFunc != nil {
		return m.ReadMessageFunc(ctx)
	}
	return MessageText, nil, nil
}

func (m *MockWebSocketConn) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *MockWebSocketConn) Ping(ctx context.Context) error {
	if m.PingFunc != nil {
		return m.PingFunc(ctx)
	}
	return nil
}
