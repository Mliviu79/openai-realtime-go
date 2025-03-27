package session

//-----------------------------------------------------------------------------
// Audio and Transcription Types
//-----------------------------------------------------------------------------

// TranscriptionModel represents the model used for audio transcription
type TranscriptionModel string

const (
	// TranscriptionModelWhisper1 is the Whisper-1 transcription model
	TranscriptionModelWhisper1 TranscriptionModel = "whisper-1"

	// TranscriptionModelGPT4oTranscribe is the GPT-4o transcription model
	TranscriptionModelGPT4oTranscribe TranscriptionModel = "gpt-4o-transcribe"

	// TranscriptionModelGPT4oMiniTranscribe is the GPT-4o mini transcription model
	TranscriptionModelGPT4oMiniTranscribe TranscriptionModel = "gpt-4o-mini-transcribe"
)

// InputAudioTranscription represents configuration for audio transcription
type InputAudioTranscription struct {
	// Model specifies which model to use for transcription
	// Currently supported models are "whisper-1", "gpt-4o-transcribe", and "gpt-4o-mini-transcribe"
	Model TranscriptionModel `json:"model,omitempty"`

	// Language specifies the language of the audio in ISO-639-1 format
	Language string `json:"language,omitempty"`

	// Prompt provides optional text to guide the model's style
	Prompt string `json:"prompt,omitempty"`
}
