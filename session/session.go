// Package session provides session management for the OpenAI Realtime API.
// It defines types for creating, updating, retrieving, and deleting sessions,
// as well as various configuration options for those sessions.
//
// The session package is responsible for:
//   - Session lifecycle management (creation, updating, retrieval, deletion)
//   - Session configuration (model, modalities, voice, tools, etc.)
//   - Type definitions for all session-related parameters
//
// A session in the OpenAI Realtime API represents a persistent connection with specific
// parameters like the model, voice settings, and tools. Sessions can be created, updated,
// and managed independently of the WebSocket connection.
//
// Example usage:
//
//	// Create a new session
//	model := session.GPT4oRealtimePreview
//	createReq := &session.CreateRequest{
//		SessionRequest: session.SessionRequest{
//			Model: &model,
//		},
//	}
//
//	sessionResp, err := client.CreateSession(ctx, createReq)
//	if err != nil {
//		log.Fatalf("Failed to create session: %v", err)
//	}
//
//	// Use the session ID to connect
//	conn, err := client.Connect(ctx,
//		openai.WithSessionID(sessionResp.ID),
//	)
//
// The session package works closely with the openai package, which provides
// client implementations for the session management operations.
package session

// CreateRequest represents a request to create a new session.
// It contains all the configuration options for the new session.
type CreateRequest struct {
	SessionRequest
}

// CreateResponse represents the response from creating a new session.
// It contains the details of the newly created session, including its ID.
type CreateResponse struct {
	Session
}

// UpdateRequest represents a request to update an existing session.
// It contains the configuration options to update in the session.
type UpdateRequest struct {
	SessionRequest
}

// UpdateResponse represents the response from updating a session.
// It contains the updated details of the session.
type UpdateResponse struct {
	Session
}

// Manager provides methods for managing sessions with the OpenAI Realtime API.
// It offers a high-level interface for session operations, abstracting away the
// underlying API calls.
type Manager interface {
	// Create creates a new session with the given configuration.
	// Returns the created session details or an error if creation fails.
	Create(req *CreateRequest) (*CreateResponse, error)

	// Update updates an existing session with the given configuration.
	// The sessionID identifies which session to update.
	// Returns the updated session details or an error if the update fails.
	Update(sessionID string, req *UpdateRequest) (*UpdateResponse, error)

	// Get retrieves information about an existing session.
	// The sessionID identifies which session to retrieve.
	// Returns the session details or an error if retrieval fails.
	Get(sessionID string) (*Session, error)

	// Delete deletes an existing session.
	// The sessionID identifies which session to delete.
	// Returns an error if deletion fails.
	Delete(sessionID string) error
}

// DefaultManager implements the Manager interface using the OpenAI API.
// It delegates session operations to a SessionClient which handles the
// actual API requests.
type DefaultManager struct {
	client SessionClient
}

// SessionClient defines the interface for session-related API operations.
// This interface allows different client implementations (e.g., for testing
// or for different API endpoints).
type SessionClient interface {
	// CreateSession creates a new session with the given configuration.
	CreateSession(req *CreateRequest) (*CreateResponse, error)

	// UpdateSession updates an existing session with the given configuration.
	UpdateSession(sessionID string, req *UpdateRequest) (*UpdateResponse, error)

	// GetSession retrieves information about an existing session.
	GetSession(sessionID string) (*Session, error)

	// DeleteSession deletes an existing session.
	DeleteSession(sessionID string) error
}

// NewManager creates a new session manager with the given client.
// The client handles the actual API requests for session operations.
//
// Parameters:
//   - client: An implementation of SessionClient that handles API requests
//
// Returns:
//   - A DefaultManager instance that implements the Manager interface
func NewManager(client SessionClient) *DefaultManager {
	return &DefaultManager{
		client: client,
	}
}

// Create creates a new session with the given configuration.
// It delegates the API call to the underlying client.
//
// Parameters:
//   - req: The configuration for the new session
//
// Returns:
//   - The created session details
//   - An error if creation fails
func (m *DefaultManager) Create(req *CreateRequest) (*CreateResponse, error) {
	return m.client.CreateSession(req)
}

// Update updates an existing session with the given configuration.
// It delegates the API call to the underlying client.
//
// Parameters:
//   - sessionID: The ID of the session to update
//   - req: The new configuration options for the session
//
// Returns:
//   - The updated session details
//   - An error if the update fails
func (m *DefaultManager) Update(sessionID string, req *UpdateRequest) (*UpdateResponse, error) {
	return m.client.UpdateSession(sessionID, req)
}

// Get retrieves information about an existing session.
// It delegates the API call to the underlying client.
//
// Parameters:
//   - sessionID: The ID of the session to retrieve
//
// Returns:
//   - The session details
//   - An error if retrieval fails
func (m *DefaultManager) Get(sessionID string) (*Session, error) {
	return m.client.GetSession(sessionID)
}

// Delete deletes an existing session.
// It delegates the API call to the underlying client.
//
// Parameters:
//   - sessionID: The ID of the session to delete
//
// Returns:
//   - An error if deletion fails
func (m *DefaultManager) Delete(sessionID string) error {
	return m.client.DeleteSession(sessionID)
}
