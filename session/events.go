package session

//-----------------------------------------------------------------------------
// Session Event Types
//-----------------------------------------------------------------------------

// SessionEventType represents the type of session event
type SessionEventType string

const (
	// SessionEventTypeCreated indicates a session was created
	SessionEventTypeCreated SessionEventType = "session.created"

	// SessionEventTypeUpdated indicates a session was updated
	SessionEventTypeUpdated SessionEventType = "session.updated"
	// SessionEventTypeExpired indicates a session has expired
	SessionEventTypeExpired SessionEventType = "session.expired"
)

// SessionEvent represents an event related to a session
type SessionEvent struct {
	// Type is the type of event
	Type SessionEventType `json:"type"`

	// ID is the session ID
	ID string `json:"id"`

	// Object is always "realtime.session"
	Object string `json:"object"`

	// Data contains the session data
	Data *Session `json:"data,omitempty"`
}
