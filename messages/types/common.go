// Package types provides common data structures and types for OpenAI Realtime API messages.
// These types are shared between incoming and outgoing messages and represent the core
// data models used in the API communication.
//
// The package defines:
//   - Message roles (system, user, assistant)
//   - Content types (text, audio, etc.)
//   - Message structure
//   - Rate limit information
//   - Transcription models and settings
//
// These types form the foundation for creating and handling messages in the OpenAI Realtime API.
// They are used by other packages such as factory (for message creation), incoming (for parsing
// server responses), and outgoing (for creating client requests).
//
// Message Roles:
//   - MessageRoleSystem: For system instructions that guide the assistant's behavior
//   - MessageRoleUser: For messages from the user
//   - MessageRoleAssistant: For messages from the assistant
//
// Content Types:
//   - MessageContentTypeText: For plain text content from the assistant
//   - MessageContentTypeInputText: For text input from the user
//   - MessageContentTypeInputAudio: For audio input from the user
//   - MessageContentTypeItemReference: For references to other items
//   - MessageContentTypeAudio: For audio content
//   - MessageContentTypeTranscript: For transcripts of audio content
//
// Example Usage:
//
//	// Create a message with text content from the user
//	message := types.Message{
//		Role: types.MessageRoleUser,
//		Content: []types.MessageContentPart{
//			{
//				Type: types.MessageContentTypeInputText,
//				Text: "What's the weather like today?",
//			},
//		},
//	}
//
// Most users will not need to create these types directly but will use the factory
// package functions to construct messages in the correct format.
package types

import (
	"github.com/Mliviu79/openai-realtime-go/session"
)

//-----------------------------------------------------------------------------
// Message Role Types
//-----------------------------------------------------------------------------

// MessageRole represents the role of a message participant
type MessageRole string

const (
	// MessageRoleSystem represents a system message
	MessageRoleSystem MessageRole = "system"

	// MessageRoleUser represents a message from the user
	MessageRoleUser MessageRole = "user"

	// MessageRoleAssistant represents a message from the assistant
	MessageRoleAssistant MessageRole = "assistant"
)

//-----------------------------------------------------------------------------
// Message Content Types
//-----------------------------------------------------------------------------

// MessageContentType represents the type of content in a message
type MessageContentType string

const (
	// MessageContentTypeText represents plain text content from the assistant
	MessageContentTypeText MessageContentType = "text"

	// MessageContentTypeInputText represents text input from the user
	MessageContentTypeInputText MessageContentType = "input_text"

	// MessageContentTypeInputAudio represents audio input from the user
	MessageContentTypeInputAudio MessageContentType = "input_audio"

	// MessageContentTypeItemReference represents a reference to another item
	MessageContentTypeItemReference MessageContentType = "item_reference"

	// MessageContentTypeAudio represents audio content
	// Note: This is not in the official API specs, but is used internally
	MessageContentTypeAudio MessageContentType = "audio"

	// MessageContentTypeTranscript represents a transcript of audio content
	// Note: This is not in the official API specs, but is used internally
	MessageContentTypeTranscript MessageContentType = "transcript"
)

// MessageOption is a function that configures a Message
type MessageOption func(*Message)

// Message represents a complete message to be sent
type Message struct {
	// Role specifies who created this message
	Role MessageRole `json:"role"`

	// Content contains the actual content of the message
	Content []MessageContentPart `json:"content"`

	// Name specifies the function being called (for function messages)
	Name string `json:"name,omitempty"`

	// Tools specifies the available functions the model can call
	Tools []session.Tool `json:"tools,omitempty"`

	// EndTurn indicates whether this message ends the conversation turn
	EndTurn bool `json:"end_turn,omitempty"`

	// Metadata contains optional key-value pairs for tracking
	Metadata any `json:"metadata,omitempty"`
}

// TranscriptionModel specifies which model to use for transcription
type TranscriptionModel string

const (
	// TranscriptionModelWhisper1 is the Whisper v1 model
	TranscriptionModelWhisper1 TranscriptionModel = "whisper-1"
)

// InputAudioTranscription contains options for audio transcription
type InputAudioTranscription struct {
	// Model specifies which model to use for transcription
	// Currently only "whisper-1" is supported
	Model TranscriptionModel `json:"model,omitempty"`

	// Language specifies the language of the audio in ISO-639-1 format
	Language string `json:"language,omitempty"`

	// Prompt provides optional text to guide the model's style
	Prompt string `json:"prompt,omitempty"`
}

// RateLimit contains information about API rate limits
type RateLimit struct {
	// Name specifies the type of rate limit ("requests" or "tokens")
	Name string `json:"name"`

	// Limit specifies the maximum allowed value for the rate limit
	Limit int `json:"limit"`

	// Remaining specifies the remaining value before the limit is reached
	Remaining int `json:"remaining"`

	// ResetSeconds specifies the seconds until the rate limit resets
	ResetSeconds float64 `json:"reset_seconds"`
}
