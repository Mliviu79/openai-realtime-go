package incoming

import (
	"encoding/json"
	"fmt"
)

// UnmarshalRcvdMsg unmarshals a JSON message into the appropriate message type
func UnmarshalRcvdMsg(data []byte) (RcvdMsg, error) {
	// First, unmarshal just enough to get the message type
	var base struct {
		Type    RcvdMsgType `json:"type"`
		EventID string      `json:"event_id,omitempty"`
	}

	if err := json.Unmarshal(data, &base); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message base: %w", err)
	}

	// Special handling for error messages which have a type of "error"
	if base.Type == "error" {
		errMsg := &ErrorMessage{}
		if err := json.Unmarshal(data, errMsg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal error message: %w", err)
		}
		return errMsg, nil
	}

	// Use the registry to create the appropriate message type
	msgType := RcvdMsgType(base.Type)
	msg, exists := CreateMessage(msgType)
	if !exists {
		// For unknown message types, try to unmarshal as an error message as a fallback
		// This is for backward compatibility
		errMsg := &ErrorMessage{}
		if err := json.Unmarshal(data, errMsg); err == nil && errMsg.Error.Message != "" {
			return errMsg, nil
		}
		return nil, fmt.Errorf("unknown message type: %s", base.Type)
	}

	// Unmarshal the full message
	if err := json.Unmarshal(data, msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message of type %s: %w", base.Type, err)
	}

	return msg, nil
}
