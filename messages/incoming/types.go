// Package incoming provides message types and parsers for messages received from the OpenAI Realtime API.
// These messages correspond to the various response types and events defined in the OpenAI Realtime API reference.
// The package handles deserialization and structured access to these messages to simplify API interaction.
package incoming

import "github.com/Mliviu79/openai-realtime-go/apierrs"

// RcvdMsgType represents the type of message received from the server
type RcvdMsgType string

//-----------------------------------------------------------------------------
// Message Type Constants
//-----------------------------------------------------------------------------

// Error message type
const (
	RcvdMsgTypeError RcvdMsgType = "error"
)

// Session-related message types
const (
	RcvdMsgTypeSessionCreated RcvdMsgType = "session.created"
	RcvdMsgTypeSessionUpdated RcvdMsgType = "session.updated"
)

// Conversation-related message types
const (
	RcvdMsgTypeConversationCreated                              RcvdMsgType = "conversation.created"
	RcvdMsgTypeConversationItemCreated                          RcvdMsgType = "conversation.item.created"
	RcvdMsgTypeConversationItemInputAudioTranscriptionCompleted RcvdMsgType = "conversation.item.input_audio_transcription.completed"
	RcvdMsgTypeConversationItemInputAudioTranscriptionDelta     RcvdMsgType = "conversation.item.input_audio_transcription.delta"
	RcvdMsgTypeConversationItemInputAudioTranscriptionFailed    RcvdMsgType = "conversation.item.input_audio_transcription.failed"
	RcvdMsgTypeConversationItemTruncated                        RcvdMsgType = "conversation.item.truncated"
	RcvdMsgTypeConversationItemDeleted                          RcvdMsgType = "conversation.item.deleted"
)

// Audio buffer-related message types
const (
	RcvdMsgTypeAudioBufferCommitted     RcvdMsgType = "input_audio_buffer.committed"
	RcvdMsgTypeAudioBufferCleared       RcvdMsgType = "input_audio_buffer.cleared"
	RcvdMsgTypeAudioBufferSpeechStarted RcvdMsgType = "input_audio_buffer.speech_started"
	RcvdMsgTypeAudioBufferSpeechStopped RcvdMsgType = "input_audio_buffer.speech_stopped"
)

// Response-related message types
const (
	RcvdMsgTypeResponseCreated                    RcvdMsgType = "response.created"
	RcvdMsgTypeResponseDone                       RcvdMsgType = "response.done"
	RcvdMsgTypeResponseContentPartAdded           RcvdMsgType = "response.content_part.added"
	RcvdMsgTypeResponseContentPartDone            RcvdMsgType = "response.content_part.done"
	RcvdMsgTypeResponseTextDelta                  RcvdMsgType = "response.text.delta"
	RcvdMsgTypeResponseTextDone                   RcvdMsgType = "response.text.done"
	RcvdMsgTypeResponseOutputItemAdded            RcvdMsgType = "response.output_item.added"
	RcvdMsgTypeResponseOutputItemDone             RcvdMsgType = "response.output_item.done"
	RcvdMsgTypeResponseAudioTranscriptDelta       RcvdMsgType = "response.audio_transcript.delta"
	RcvdMsgTypeResponseAudioTranscriptDone        RcvdMsgType = "response.audio_transcript.done"
	RcvdMsgTypeResponseAudioDelta                 RcvdMsgType = "response.audio.delta"
	RcvdMsgTypeResponseAudioDone                  RcvdMsgType = "response.audio.done"
	RcvdMsgTypeResponseFunctionCallArgumentsDelta RcvdMsgType = "response.function_call_arguments.delta"
	RcvdMsgTypeResponseFunctionCallArgumentsDone  RcvdMsgType = "response.function_call_arguments.done"
)

// Rate limit-related message types
const (
	RcvdMsgTypeRateLimitsUpdated RcvdMsgType = "rate_limits.updated"
)

// RcvdMsg is the interface implemented by all received message types
type RcvdMsg interface {
	// RcvdMsgType returns the type of the message
	RcvdMsgType() RcvdMsgType
}

// RcvdMsgBase contains fields common to all received messages
type RcvdMsgBase struct {
	// ID is an optional server-generated identifier for this message
	ID string `json:"message_id,omitempty"`
	// EventID is the unique ID of the server event (used in error messages)
	EventID string `json:"event_id,omitempty"`
	// Type indicates the specific type of this message
	Type RcvdMsgType `json:"type"`
}

// RcvdMsgType returns the type of the message
func (m RcvdMsgBase) RcvdMsgType() RcvdMsgType {
	return m.Type
}

// String returns the string representation of the message type
func (t RcvdMsgType) String() string {
	return string(t)
}

// ErrorInfo contains detailed information about an error
type ErrorInfo struct {
	// Type indicates the category of error (e.g., "invalid_request_error")
	Type apierrs.ErrorType `json:"type"`
	// Code is a machine-readable error code, can be null
	Code apierrs.ErrorCode `json:"code,omitempty"`
	// Message contains a human-readable error description
	Message string `json:"message"`
	// Param identifies which parameter caused the error, if applicable
	Param *string `json:"param"`
	// EventID references the client event that triggered this error, if applicable
	EventID string `json:"event_id,omitempty"`
}

// ErrorMessage represents an error response from the server
type ErrorMessage struct {
	RcvdMsgBase
	// Error contains detailed information about what went wrong
	Error ErrorInfo `json:"error"`
}
