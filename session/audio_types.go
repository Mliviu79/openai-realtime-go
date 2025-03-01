package session

//-----------------------------------------------------------------------------
// Audio and Transcription Types
//-----------------------------------------------------------------------------

// TranscriptionModel represents the model used for audio transcription
type TranscriptionModel string

const (
	// TranscriptionModelWhisper1 is the Whisper-1 transcription model
	TranscriptionModelWhisper1 TranscriptionModel = "whisper-1"
)

// InputAudioTranscription represents configuration for audio transcription
type InputAudioTranscription struct {
	// Model specifies which model to use for transcription
	// Currently only "whisper-1" is supported
	Model TranscriptionModel `json:"model,omitempty"`

	// Language specifies the language of the audio in ISO-639-1 format
	Language string `json:"language,omitempty"`

	// Prompt provides optional text to guide the model's style
	Prompt string `json:"prompt,omitempty"`
}
