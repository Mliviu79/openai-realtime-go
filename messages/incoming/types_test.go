package incoming

import (
	"testing"
)

func TestRcvdMsgBaseType(t *testing.T) {
	// Test that RcvdMsgBase.RcvdMsgType correctly returns the message type
	testCases := []struct {
		msgType RcvdMsgType
		want    RcvdMsgType
	}{
		{RcvdMsgTypeError, RcvdMsgTypeError},
		{RcvdMsgTypeSessionCreated, RcvdMsgTypeSessionCreated},
		{RcvdMsgTypeConversationCreated, RcvdMsgTypeConversationCreated},
		{RcvdMsgTypeAudioBufferCommitted, RcvdMsgTypeAudioBufferCommitted},
		{RcvdMsgTypeResponseCreated, RcvdMsgTypeResponseCreated},
		{RcvdMsgTypeRateLimitsUpdated, RcvdMsgTypeRateLimitsUpdated},
	}

	for _, tc := range testCases {
		t.Run(string(tc.msgType), func(t *testing.T) {
			base := RcvdMsgBase{Type: tc.msgType}
			if got := base.RcvdMsgType(); got != tc.want {
				t.Errorf("RcvdMsgBase.RcvdMsgType() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestErrorMessageStruct(t *testing.T) {
	// Test that ErrorMessage implements the RcvdMsg interface
	errorMsg := ErrorMessage{
		RcvdMsgBase: RcvdMsgBase{
			Type: RcvdMsgTypeError,
			ID:   "msg_123",
		},
		Error: ErrorInfo{
			Type:    "invalid_request_error",
			Code:    "invalid_parameter",
			Message: "The parameter is invalid",
			Param:   nil,
			EventID: "evt_456",
		},
	}

	// Test the implementation of RcvdMsg interface
	var msg RcvdMsg = &errorMsg
	if msg.RcvdMsgType() != RcvdMsgTypeError {
		t.Errorf("ErrorMessage.RcvdMsgType() = %v, want %v", msg.RcvdMsgType(), RcvdMsgTypeError)
	}

	// Verify direct access to the error info
	if errorMsg.Error.Type != "invalid_request_error" {
		t.Errorf("ErrorMessage.Error.Type = %v, want %v", errorMsg.Error.Type, "invalid_request_error")
	}

	if errorMsg.Error.Code != "invalid_parameter" {
		t.Errorf("ErrorMessage.Error.Code = %v, want %v", errorMsg.Error.Code, "invalid_parameter")
	}

	if errorMsg.Error.Message != "The parameter is invalid" {
		t.Errorf("ErrorMessage.Error.Message = %v, want %v", errorMsg.Error.Message, "The parameter is invalid")
	}

	if errorMsg.Error.Param != nil {
		t.Errorf("ErrorMessage.Error.Param = %v, want nil", errorMsg.Error.Param)
	}

	if errorMsg.Error.EventID != "evt_456" {
		t.Errorf("ErrorMessage.Error.EventID = %v, want %v", errorMsg.Error.EventID, "evt_456")
	}
}

func TestMessageTypeConsistency(t *testing.T) {
	// Test that all message types in the registry match their constant values
	// This helps ensure the constants are being used correctly
	for msgType, factory := range MessageTypeRegistry {
		msg := factory()
		if msg.RcvdMsgType() != msgType {
			t.Errorf("Message type mismatch for %q: factory produced message with type %q",
				msgType, msg.RcvdMsgType())
		}
	}
}
