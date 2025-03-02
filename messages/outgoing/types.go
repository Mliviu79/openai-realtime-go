// Package outgoing provides message types and constructors for messages sent to the OpenAI Realtime API.
// These messages correspond to the various message types defined in the OpenAI Realtime API reference.
// The package handles construction and serialization of these messages to ensure compatibility with the API.
package outgoing

// OutMsgType represents the type of message being sent to the server
type OutMsgType string

// Session-related message types
const (
	OutMsgTypeSessionUpdate OutMsgType = "session.update"
)

// Audio buffer-related message types
const (
	OutMsgTypeAudioBufferAppend OutMsgType = "input_audio_buffer.append"
	OutMsgTypeAudioBufferCommit OutMsgType = "input_audio_buffer.commit"
	OutMsgTypeAudioBufferClear  OutMsgType = "input_audio_buffer.clear"
)

// Conversation-related message types
const (
	OutMsgTypeConversationCreate   OutMsgType = "conversation.item.create"
	OutMsgTypeConversationTruncate OutMsgType = "conversation.item.truncate"
	OutMsgTypeConversationDelete   OutMsgType = "conversation.item.delete"
)

// Response-related message types
const (
	OutMsgTypeResponseCreate OutMsgType = "response.create"
	OutMsgTypeResponseCancel OutMsgType = "response.cancel"
)

// OutMsg is the interface implemented by all outgoing message types
type OutMsg interface {
	// OutMsgType returns the type of the message as a string
	OutMsgType() string

	// OutMsgID returns the client-generated ID for this message
	OutMsgID() string
}

// OutMsgBase contains fields common to all outgoing messages
type OutMsgBase struct {
	// ID is an optional client-generated identifier for this message
	// Can be used to correlate responses with requests
	ID string `json:"event_id,omitempty"`

	// Type indicates the specific type of this message
	Type OutMsgType `json:"type"`
}

// OutMsgType returns the type of the message as a string
func (m OutMsgBase) OutMsgType() string {
	return string(m.Type)
}

// String returns the string representation of the message type
func (t OutMsgType) String() string {
	return string(t)
}

// OutMsgID returns the client-generated ID for this message
func (m OutMsgBase) OutMsgID() string {
	return m.ID
}
