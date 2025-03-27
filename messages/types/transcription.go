package types

import (
	"github.com/Mliviu79/openai-realtime-go/session"
)

//-----------------------------------------------------------------------------
// Transcription Session Types
//-----------------------------------------------------------------------------
// These types are now imported from the session package to avoid duplication.
// Type aliases are provided for backward compatibility.
//-----------------------------------------------------------------------------

// TranscriptionModel represents the model used for audio transcription
// It's now an alias to session.TranscriptionModel
type TranscriptionModel = session.TranscriptionModel

const (
	// TranscriptionModelWhisper1 is the Whisper-1 transcription model
	TranscriptionModelWhisper1 TranscriptionModel = session.TranscriptionModelWhisper1

	// TranscriptionModelGPT4oTranscribe is the GPT-4o transcription model
	TranscriptionModelGPT4oTranscribe TranscriptionModel = session.TranscriptionModelGPT4oTranscribe

	// TranscriptionModelGPT4oMiniTranscribe is the GPT-4o mini transcription model
	TranscriptionModelGPT4oMiniTranscribe TranscriptionModel = session.TranscriptionModelGPT4oMiniTranscribe
)

// InputAudioTranscription represents configuration for audio transcription
// It's now an alias to session.InputAudioTranscription
type InputAudioTranscription = session.InputAudioTranscription

// TurnDetection represents turn detection configuration
// It's now an alias to session.TurnDetection
type TurnDetection = session.TurnDetection

// InputAudioNoiseReduction represents configuration for noise reduction on input audio
// It's now an alias to session.InputAudioNoiseReduction
type InputAudioNoiseReduction = session.InputAudioNoiseReduction

// ClientSecret contains authentication information for client-side connections
// It's now an alias to session.ClientSecret
type ClientSecret = session.ClientSecret

// TranscriptionSession represents a transcription session in the Realtime API
// This type is still maintained in this package as it's primarily a message type
type TranscriptionSession struct {
	// ID is the server-assigned identifier for this session
	ID string `json:"id,omitempty"`

	// Object indicates the type of this object (always "realtime.transcription_session")
	Object string `json:"object,omitempty"`

	// ExpiresAt is the timestamp when this session expires
	ExpiresAt int64 `json:"expires_at,omitempty"`

	// Modalities specifies the types of input/output the model can handle
	Modalities []string `json:"modalities,omitempty"`

	// InputAudioFormat specifies the format for audio input
	InputAudioFormat string `json:"input_audio_format,omitempty"`

	// InputAudioTranscription configures audio transcription settings
	InputAudioTranscription *InputAudioTranscription `json:"input_audio_transcription,omitempty"`

	// TurnDetection configures how turns are detected in conversation
	TurnDetection *TurnDetection `json:"turn_detection,omitempty"`

	// InputAudioNoiseReduction configures noise reduction on input audio
	InputAudioNoiseReduction *InputAudioNoiseReduction `json:"input_audio_noise_reduction,omitempty"`

	// ClientSecret contains authentication information for client-side connections
	// Only present when the session is created via REST API
	ClientSecret *ClientSecret `json:"client_secret,omitempty"`
}
