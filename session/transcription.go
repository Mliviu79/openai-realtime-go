package session

//-----------------------------------------------------------------------------
// Transcription Session Types
//-----------------------------------------------------------------------------

// TranscriptionSessionInclude represents the items that can be included in the transcription results
type TranscriptionSessionInclude string

const (
	// TranscriptionSessionIncludeLogprobs includes log probabilities in transcription results
	TranscriptionSessionIncludeLogprobs TranscriptionSessionInclude = "item.input_audio_transcription.logprobs"
)

// TranscriptionSession represents a complete transcription session with the OpenAI Realtime API
type TranscriptionSession struct {
	// Server-assigned fields
	ID        string `json:"id,omitempty"`
	Object    string `json:"object,omitempty"` // Always "realtime.transcription_session" when present
	ExpiresAt int64  `json:"expires_at,omitempty"`

	// Client secret information for authentication
	ClientSecret *ClientSecret `json:"client_secret,omitempty"`

	// Embedded session request fields
	TranscriptionSessionRequest
}

// NewTranscriptionSession creates a new TranscriptionSession with default values
func NewTranscriptionSession() *TranscriptionSession {
	return &TranscriptionSession{
		Object: "realtime.transcription_session",
	}
}

// TranscriptionSessionRequest represents both create and update requests for transcription sessions
type TranscriptionSessionRequest struct {
	// Modalities specifies the types of input/output the model can handle
	Modalities *[]Modality `json:"modalities,omitempty"`

	// InputAudioFormat specifies the format for audio input
	InputAudioFormat *AudioFormat `json:"input_audio_format,omitempty"`

	// InputAudioTranscription configures audio transcription settings
	InputAudioTranscription *InputAudioTranscription `json:"input_audio_transcription,omitempty"`

	// TurnDetection configures how turns are detected in conversation
	TurnDetection *TurnDetection `json:"turn_detection,omitempty"`

	// InputAudioNoiseReduction configures noise reduction on input audio
	InputAudioNoiseReduction *InputAudioNoiseReduction `json:"input_audio_noise_reduction,omitempty"`

	// Include specifies additional items to include in transcription results
	Include *[]TranscriptionSessionInclude `json:"include,omitempty"`
}

// CreateTranscriptionSessionRequest represents a request to create a new transcription session
type CreateTranscriptionSessionRequest struct {
	TranscriptionSessionRequest
}

// CreateTranscriptionSessionResponse represents the response from creating a transcription session
type CreateTranscriptionSessionResponse struct {
	TranscriptionSession
}
