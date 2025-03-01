package session

//-----------------------------------------------------------------------------
// Session Types
//-----------------------------------------------------------------------------

// ClientSecret contains authentication information for client-side connections
type ClientSecret struct {
	// Value is an ephemeral key usable in client environments to authenticate connections to the Realtime API.
	// Use this in client-side environments rather than a standard API token, which should only be used server-side.
	Value string `json:"value"`

	// ExpiresAt is the timestamp for when the token expires.
	// Currently, all tokens expire after one minute.
	ExpiresAt int64 `json:"expires_at"`
}

// ClientSecretInfo is a wrapper for ClientSecret
// Deprecated: This type is no longer used as the API directly returns the ClientSecret structure
// without this additional nesting level.
type ClientSecretInfo struct {
	ClientSecret ClientSecret `json:"client_secret"`
}

// Session represents a complete session with the OpenAI Realtime API
type Session struct {
	// Server-assigned fields
	ID           string        `json:"id,omitempty"`
	Object       string        `json:"object,omitempty"` // Always "realtime.session" when present
	ClientSecret *ClientSecret `json:"client_secret,omitempty"`

	// Common configuration fields
	Modalities              *[]Modality              `json:"modalities,omitempty"`
	Model                   *Model                   `json:"model,omitempty"`
	Instructions            *string                  `json:"instructions,omitempty"`
	Voice                   *Voice                   `json:"voice,omitempty"`
	InputAudioFormat        *AudioFormat             `json:"input_audio_format,omitempty"`
	OutputAudioFormat       *AudioFormat             `json:"output_audio_format,omitempty"`
	InputAudioTranscription *InputAudioTranscription `json:"input_audio_transcription,omitempty"`
	TurnDetection           *TurnDetection           `json:"turn_detection,omitempty"`
	Tools                   *[]Tool                  `json:"tools,omitempty"`
	ToolChoice              *ToolChoiceObj           `json:"tool_choice,omitempty"`
	Temperature             *float64                 `json:"temperature,omitempty"`
	MaxResponseOutputTokens *IntOrInf                `json:"max_response_output_tokens,omitempty"`
}

// NewSession creates a new Session with default values
func NewSession() *Session {
	return &Session{
		Object: "realtime.session",
	}
}

// SessionRequest represents both create and update requests
// All fields are pointers to make them optional
type SessionRequest struct {
	// Modalities specifies the types of input/output the model can handle
	Modalities *[]Modality `json:"modalities,omitempty"`

	// Model specifies which model to use for the session
	Model *Model `json:"model,omitempty"`

	// Instructions provide system instructions to guide the model
	Instructions *string `json:"instructions,omitempty"`

	// Voice specifies which voice to use for audio responses
	Voice *Voice `json:"voice,omitempty"`

	// InputAudioFormat specifies the format for audio input
	InputAudioFormat *AudioFormat `json:"input_audio_format,omitempty"`

	// OutputAudioFormat specifies the format for audio output
	OutputAudioFormat *AudioFormat `json:"output_audio_format,omitempty"`

	// InputAudioTranscription configures audio transcription settings
	InputAudioTranscription *InputAudioTranscription `json:"input_audio_transcription,omitempty"`

	// TurnDetection configures how turns are detected in conversation
	TurnDetection *TurnDetection `json:"turn_detection,omitempty"`

	// Tools specifies the available functions the model can call
	Tools *[]Tool `json:"tools,omitempty"`

	// ToolChoice controls how the model selects tools
	ToolChoice *ToolChoiceObj `json:"tool_choice,omitempty"`

	// Temperature controls the randomness of the model's output
	Temperature *float64 `json:"temperature,omitempty"`

	// MaxResponseOutputTokens limits the length of responses
	MaxResponseOutputTokens *IntOrInf `json:"max_response_output_tokens,omitempty"`
}
