package types

import (
	"github.com/Mliviu79/openai-realtime-go/apierrs"
	"github.com/Mliviu79/openai-realtime-go/session"
)

//-----------------------------------------------------------------------------
// Response Status Types
//-----------------------------------------------------------------------------

// ResponseStatus represents the status of a response
type ResponseStatus string

const (
	// ResponseStatusInProgress indicates the response is still being generated
	ResponseStatusInProgress ResponseStatus = "in_progress"

	// ResponseStatusCompleted indicates the response has been fully generated
	ResponseStatusCompleted ResponseStatus = "completed"

	// ResponseStatusFailed indicates the response generation failed
	ResponseStatusFailed ResponseStatus = "failed"

	// ResponseStatusCancelled indicates the response generation was cancelled
	ResponseStatusCancelled ResponseStatus = "cancelled"

	// ResponseStatusIncomplete indicates the response generation was truncated
	ResponseStatusIncomplete ResponseStatus = "incomplete"
)

//-----------------------------------------------------------------------------
// Response Error Types
//-----------------------------------------------------------------------------

// ResponseErrorType represents the type of error in a response status
type ResponseErrorType string

const (
	// ResponseErrorTypeCompleted indicates the response completed successfully
	ResponseErrorTypeCompleted ResponseErrorType = "completed"

	// ResponseErrorTypeCancelled indicates the response was cancelled
	ResponseErrorTypeCancelled ResponseErrorType = "cancelled"

	// ResponseErrorTypeFailed indicates the response failed
	ResponseErrorTypeFailed ResponseErrorType = "failed"

	// ResponseErrorTypeIncomplete indicates the response was incomplete
	ResponseErrorTypeIncomplete ResponseErrorType = "incomplete"
)

// ResponseMessageItem extends MessageItem with additional response-specific fields
type ResponseMessageItem struct {
	MessageItem

	// Object is always "realtime.item" when present
	Object string `json:"object,omitempty"`
}

// ResponseStatusDetails provides additional information about the response status
type ResponseStatusDetails struct {
	// Type is the type of error that caused the response to fail
	// Values: "completed", "cancelled", "incomplete", "failed"
	Type ResponseErrorType `json:"type,omitempty"`

	// Reason is the reason the Response did not complete
	// For "cancelled" status: "turn_detected" or "client_cancelled"
	// For "incomplete" status: "max_output_tokens" or "content_filter"
	Reason string `json:"reason,omitempty"`

	// Error contains details about the error when status is "failed"
	Error *ResponseError `json:"error,omitempty"`
}

// ResponseError describes an error that caused a response to fail
type ResponseError struct {
	// Type is the type of error
	Type apierrs.ErrorType `json:"type,omitempty"`

	// Code is the error code, if any
	Code apierrs.ErrorCode `json:"code,omitempty"`
}

// OutputItem represents an item generated in a response
type OutputItem struct {
	// ID is the unique identifier for this item
	ID string `json:"id,omitempty"`

	// Type is the type of the item
	Type MessageItemType `json:"type,omitempty"`

	// Object is always "realtime.item"
	Object string `json:"object,omitempty"`

	// Status is the status of the item
	Status ItemStatus `json:"status,omitempty"`

	// Role is the role of the message sender, only applicable for message items
	Role MessageRole `json:"role,omitempty"`

	// Content is the content of the message, applicable for message items
	Content []MessageContentPart `json:"content,omitempty"`

	// CallID is the ID of the function call
	CallID string `json:"call_id,omitempty"`

	// Name is the name of the function being called
	Name string `json:"name,omitempty"`

	// Arguments are the arguments of the function call
	Arguments string `json:"arguments,omitempty"`

	// Output is the output of the function call
	Output string `json:"output,omitempty"`
}

// NewOutputItem creates a new OutputItem with default values
func NewOutputItem() *OutputItem {
	return &OutputItem{
		Object: "realtime.item",
	}
}

// Response represents a complete response from the API
type Response struct {
	// ID is the unique identifier for this response
	ID string `json:"id"`

	// Object is always "realtime.response"
	Object string `json:"object,omitempty"`

	// Status indicates the current status of the response
	Status ResponseStatus `json:"status"`

	// StatusDetails provides additional information about the status
	StatusDetails *ResponseStatusDetails `json:"status_details,omitempty"`

	// Output contains the items generated in this response
	Output []OutputItem `json:"output"`

	// Metadata contains additional information about the response
	Metadata map[string]string `json:"metadata,omitempty"`

	// Usage contains token usage statistics
	Usage *Usage `json:"usage,omitempty"`

	// ConversationID identifies which conversation the response is added to
	ConversationID string `json:"conversation_id,omitempty"`

	// Voice specifies which voice was used for audio responses
	Voice session.Voice `json:"voice,omitempty"`

	// Modalities indicates the types of output the model used to respond
	Modalities []session.Modality `json:"modalities,omitempty"`

	// OutputAudioFormat specifies the format of output audio
	OutputAudioFormat session.AudioFormat `json:"output_audio_format,omitempty"`

	// Temperature is the sampling temperature for the model
	Temperature float64 `json:"temperature,omitempty"`

	// MaxOutputTokens is the maximum number of output tokens for a single response
	MaxOutputTokens session.IntOrInf `json:"max_output_tokens,omitempty"`
}

// NewResponse creates a new Response with default values
func NewResponse() *Response {
	return &Response{
		Object: "realtime.response",
		Output: []OutputItem{},
	}
}

// ResponseConfig contains configuration options for generating a response
type ResponseConfig struct {
	// Modalities specifies the types of output the model can generate
	// Example: [session.ModalityText, session.ModalityAudio]
	Modalities []session.Modality `json:"modalities"`

	// Instructions provides system instructions to guide the model
	Instructions *string `json:"instructions,omitempty"`

	// Voice specifies which voice to use for audio responses
	// Options: VoiceAlloy, VoiceAsh, VoiceBallad, VoiceCoral, VoiceEcho, VoiceSage, VoiceShimmer, VoiceVerse
	Voice *session.Voice `json:"voice,omitempty"`

	// OutputAudioFormat specifies the format for audio responses
	// Options: AudioFormatPCM16, AudioFormatG711ULaw, AudioFormatG711ALaw
	OutputAudioFormat *session.AudioFormat `json:"output_audio_format,omitempty"`

	// Tools specifies the available functions the model can call
	Tools []session.Tool `json:"tools,omitempty"`

	// ToolChoice controls how the model selects tools
	// Options: ToolChoiceAuto, ToolChoiceNone, ToolChoiceRequired
	ToolChoice *session.ToolChoiceObj `json:"tool_choice,omitempty"`

	// Temperature controls the randomness of the model's output
	// Range: [0.6, 1.2], default 0.8
	Temperature *float64 `json:"temperature,omitempty"`

	// MaxResponseOutputTokens limits the length of the response
	// Range: 1-4096 or "inf", default "inf"
	MaxResponseOutputTokens *session.IntOrInf `json:"max_response_output_tokens,omitempty"`

	// Conversation controls whether to use conversation history
	// Options: "auto" or "none", default "auto"
	Conversation *string `json:"conversation,omitempty"`

	// Metadata contains optional key-value pairs for tracking
	// Maximum 16 pairs
	Metadata map[string]string `json:"metadata,omitempty"`

	// Input provides additional items for model context
	Input []ConversationItem `json:"input,omitempty"`
}
