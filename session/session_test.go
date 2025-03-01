package session

import (
	"errors"
	"testing"
)

// MockSessionClient is a mock implementation of the SessionClient interface
type MockSessionClient struct {
	CreateSessionFunc func(req *CreateRequest) (*CreateResponse, error)
	UpdateSessionFunc func(sessionID string, req *UpdateRequest) (*UpdateResponse, error)
	GetSessionFunc    func(sessionID string) (*Session, error)
	DeleteSessionFunc func(sessionID string) error
}

func (m *MockSessionClient) CreateSession(req *CreateRequest) (*CreateResponse, error) {
	if m.CreateSessionFunc != nil {
		return m.CreateSessionFunc(req)
	}
	return nil, nil
}

func (m *MockSessionClient) UpdateSession(sessionID string, req *UpdateRequest) (*UpdateResponse, error) {
	if m.UpdateSessionFunc != nil {
		return m.UpdateSessionFunc(sessionID, req)
	}
	return nil, nil
}

func (m *MockSessionClient) GetSession(sessionID string) (*Session, error) {
	if m.GetSessionFunc != nil {
		return m.GetSessionFunc(sessionID)
	}
	return nil, nil
}

func (m *MockSessionClient) DeleteSession(sessionID string) error {
	if m.DeleteSessionFunc != nil {
		return m.DeleteSessionFunc(sessionID)
	}
	return nil
}

func TestNewManager(t *testing.T) {
	client := &MockSessionClient{}
	manager := NewManager(client)

	if manager == nil {
		t.Fatal("Expected manager to not be nil")
	}
}

func TestCreate(t *testing.T) {
	// Create a mock client that returns a predefined response
	expectedResponse := &CreateResponse{
		Session: Session{
			ID: "test-session-id",
		},
	}

	mockClient := &MockSessionClient{
		CreateSessionFunc: func(req *CreateRequest) (*CreateResponse, error) {
			// Verify the request was passed correctly
			model := req.Model
			if model == nil {
				t.Error("Expected model to not be nil")
			} else if *model != GPT4oRealtimePreview {
				t.Errorf("Expected model to be %v, got %v", GPT4oRealtimePreview, *model)
			}

			return expectedResponse, nil
		},
	}

	manager := NewManager(mockClient)

	// Create a request
	modelValue := GPT4oRealtimePreview
	req := &CreateRequest{
		SessionRequest: SessionRequest{
			Model: &modelValue,
		},
	}

	// Call Create
	resp, err := manager.Create(req)

	// Verify results
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.ID != expectedResponse.ID {
		t.Errorf("Expected session ID to be %q, got %q", expectedResponse.ID, resp.ID)
	}
}

func TestUpdate(t *testing.T) {
	// Create a mock client that returns a predefined response
	expectedResponse := &UpdateResponse{
		Session: Session{
			ID: "test-session-id",
		},
	}

	mockClient := &MockSessionClient{
		UpdateSessionFunc: func(sessionID string, req *UpdateRequest) (*UpdateResponse, error) {
			// Verify the sessionID and request were passed correctly
			if sessionID != "test-session-id" {
				t.Errorf("Expected session ID to be 'test-session-id', got %q", sessionID)
			}

			return expectedResponse, nil
		},
	}

	manager := NewManager(mockClient)

	// Create a request
	req := &UpdateRequest{}

	// Call Update
	resp, err := manager.Update("test-session-id", req)

	// Verify results
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.ID != expectedResponse.ID {
		t.Errorf("Expected session ID to be %q, got %q", expectedResponse.ID, resp.ID)
	}
}

func TestGet(t *testing.T) {
	// Create a mock client that returns a predefined response
	expectedResponse := &Session{
		ID: "test-session-id",
	}

	mockClient := &MockSessionClient{
		GetSessionFunc: func(sessionID string) (*Session, error) {
			// Verify the sessionID was passed correctly
			if sessionID != "test-session-id" {
				t.Errorf("Expected session ID to be 'test-session-id', got %q", sessionID)
			}

			return expectedResponse, nil
		},
	}

	manager := NewManager(mockClient)

	// Call Get
	resp, err := manager.Get("test-session-id")

	// Verify results
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.ID != expectedResponse.ID {
		t.Errorf("Expected session ID to be %q, got %q", expectedResponse.ID, resp.ID)
	}
}

func TestDelete(t *testing.T) {
	// Create a mock client that returns a predefined response
	mockClient := &MockSessionClient{
		DeleteSessionFunc: func(sessionID string) error {
			// Verify the sessionID was passed correctly
			if sessionID != "test-session-id" {
				t.Errorf("Expected session ID to be 'test-session-id', got %q", sessionID)
			}

			return nil
		},
	}

	manager := NewManager(mockClient)

	// Call Delete
	err := manager.Delete("test-session-id")

	// Verify results
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestErrorHandling(t *testing.T) {
	// Create a mock client that returns an error
	expectedError := errors.New("test error")
	mockClient := &MockSessionClient{
		CreateSessionFunc: func(req *CreateRequest) (*CreateResponse, error) {
			return nil, expectedError
		},
		UpdateSessionFunc: func(sessionID string, req *UpdateRequest) (*UpdateResponse, error) {
			return nil, expectedError
		},
		GetSessionFunc: func(sessionID string) (*Session, error) {
			return nil, expectedError
		},
		DeleteSessionFunc: func(sessionID string) error {
			return expectedError
		},
	}

	manager := NewManager(mockClient)

	// Test Create error handling
	_, err := manager.Create(&CreateRequest{})
	if err != expectedError {
		t.Errorf("Expected error to be %v, got %v", expectedError, err)
	}

	// Test Update error handling
	_, err = manager.Update("test-session-id", &UpdateRequest{})
	if err != expectedError {
		t.Errorf("Expected error to be %v, got %v", expectedError, err)
	}

	// Test Get error handling
	_, err = manager.Get("test-session-id")
	if err != expectedError {
		t.Errorf("Expected error to be %v, got %v", expectedError, err)
	}

	// Test Delete error handling
	err = manager.Delete("test-session-id")
	if err != expectedError {
		t.Errorf("Expected error to be %v, got %v", expectedError, err)
	}
}
