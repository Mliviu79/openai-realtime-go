package incoming

import (
	"github.com/Mliviu79/go-openai-realtime/session"
)

// SessionCreatedMessage is sent when a new session is created
type SessionCreatedMessage struct {
	RcvdMsgBase
	// Session contains the details of the newly created session
	Session session.Session `json:"session"`
}

// SessionUpdatedMessage is sent when an existing session is updated
type SessionUpdatedMessage struct {
	RcvdMsgBase
	// Session contains the updated session information
	Session session.Session `json:"session"`
}
