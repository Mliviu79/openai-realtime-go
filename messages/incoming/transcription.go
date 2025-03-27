package incoming

import (
	"github.com/Mliviu79/openai-realtime-go/messages/types"
)

//-----------------------------------------------------------------------------
// Transcription Message Types and Constants
//-----------------------------------------------------------------------------

// Transcription-related message types
const (
	RcvdMsgTypeTranscriptionSessionCreated RcvdMsgType = "transcription_session.created"
	RcvdMsgTypeTranscriptionSessionUpdated RcvdMsgType = "transcription_session.updated"
	RcvdMsgTypeInputAudioTranscription     RcvdMsgType = "input_audio.transcription"
	RcvdMsgTypeTranscriptionDone           RcvdMsgType = "transcription.done"
)

// LogProbItem represents a single token and its associated log probability
type LogProbItem struct {
	// Token is the text representation of the token
	Token string `json:"token"`

	// LogProb is the log probability of the token
	LogProb float64 `json:"logprob"`
}

// InputAudioTranscriptionMessage represents a transcription of audio input
type InputAudioTranscriptionMessage struct {
	RcvdMsgBase

	// Text is the transcribed text
	Text string `json:"text"`

	// Logprobs contains token log probabilities if requested via the include parameter
	Logprobs []LogProbItem `json:"logprobs,omitempty"`
}

// TranscriptionDoneMessage signals the completion of a transcription
type TranscriptionDoneMessage struct {
	RcvdMsgBase
}

// TranscriptionSessionCreatedMessage signals the creation of a transcription session
type TranscriptionSessionCreatedMessage struct {
	RcvdMsgBase

	// Session contains the updated session information
	Session types.TranscriptionSession `json:"session"`
}

// TranscriptionSessionUpdatedMessage signals an update to a transcription session
type TranscriptionSessionUpdatedMessage struct {
	RcvdMsgBase

	// Session contains the updated session information
	Session types.TranscriptionSession `json:"session"`
}
