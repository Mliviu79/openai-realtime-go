package outgoing

import (
	"github.com/Mliviu79/go-openai-realtime/messages/types"
)

// ResponseCreateMessage is used to create a new response
type ResponseCreateMessage struct {
	OutMsgBase
	// Response contains the configuration for the response
	Response types.ResponseConfig `json:"response"`
}

// NewResponseCreateMessage creates a new response create message
func NewResponseCreateMessage(config types.ResponseConfig) ResponseCreateMessage {
	return ResponseCreateMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeResponseCreate,
		},
		Response: config,
	}
}

// ResponseCancelMessage is used to cancel an in-progress response
type ResponseCancelMessage struct {
	OutMsgBase
	// ResponseID identifies the response to cancel
	ResponseID string `json:"response_id"`
}

// NewResponseCancelMessage creates a new response cancel message
func NewResponseCancelMessage(responseID string) ResponseCancelMessage {
	return ResponseCancelMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeResponseCancel,
		},
		ResponseID: responseID,
	}
}
