package outgoing

import (
	"github.com/Mliviu79/openai-realtime-go/session"
)

// SessionUpdateMessage is used to update session configuration
type SessionUpdateMessage struct {
	OutMsgBase
	// Session contains the configuration parameters to update
	Session session.SessionRequest `json:"session"`
}

// NewSessionUpdateMessage creates a new session update message
func NewSessionUpdateMessage(sessionReq session.SessionRequest) SessionUpdateMessage {
	return SessionUpdateMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeSessionUpdate,
		},
		Session: sessionReq,
	}
}
