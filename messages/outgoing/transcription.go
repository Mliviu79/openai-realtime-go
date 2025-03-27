package outgoing

import (
	"github.com/Mliviu79/openai-realtime-go/session"
)

// TranscriptionSessionUpdateMessage is used to update transcription session configuration
type TranscriptionSessionUpdateMessage struct {
	OutMsgBase
	// Session contains the configuration parameters to update
	Session session.TranscriptionSessionRequest `json:"session"`
}

// NewTranscriptionSessionUpdateMessage creates a new transcription session update message
func NewTranscriptionSessionUpdateMessage(sessionReq session.TranscriptionSessionRequest) TranscriptionSessionUpdateMessage {
	return TranscriptionSessionUpdateMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeTranscriptionSessionUpdate,
		},
		Session: sessionReq,
	}
}

// NewTranscriptionSessionUpdateMessageWithID creates a new transcription session update message with a custom ID
func NewTranscriptionSessionUpdateMessageWithID(id string, sessionReq session.TranscriptionSessionRequest) TranscriptionSessionUpdateMessage {
	msg := NewTranscriptionSessionUpdateMessage(sessionReq)
	msg.ID = id
	return msg
}
