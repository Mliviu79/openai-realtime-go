package incoming

import (
	"testing"
)

func TestUnmarshalRcvdMsg(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
		msgType RcvdMsgType
	}{
		{
			name: "error message",
			json: `{
				"type": "error",
				"event_id": "evt_123",
				"error": {
					"type": "invalid_request_error",
					"code": "parameter_invalid",
					"message": "Invalid parameter",
					"param": "model"
				}
			}`,
			wantErr: false,
			msgType: RcvdMsgTypeError,
		},
		{
			name: "session created message",
			json: `{
				"type": "session.created",
				"message_id": "msg_123",
				"session": {
					"id": "session_123",
					"title": "New Session"
				}
			}`,
			wantErr: false,
			msgType: RcvdMsgTypeSessionCreated,
		},
		{
			name: "conversation created message",
			json: `{
				"type": "conversation.created",
				"message_id": "msg_456",
				"conversation": {
					"id": "conv_123"
				}
			}`,
			wantErr: false,
			msgType: RcvdMsgTypeConversationCreated,
		},
		{
			name:    "invalid json",
			json:    `{invalid json`,
			wantErr: true,
		},
		{
			name:    "missing type field",
			json:    `{"message_id": "msg_123"}`,
			wantErr: true,
		},
		{
			name:    "unknown message type",
			json:    `{"type": "unknown.type", "message_id": "msg_123"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := UnmarshalRcvdMsg([]byte(tt.json))

			// Check error conditions
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalRcvdMsg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we don't expect an error, verify the message type
			if !tt.wantErr {
				if msg == nil {
					t.Errorf("UnmarshalRcvdMsg() returned nil message but no error")
					return
				}

				if msg.RcvdMsgType() != tt.msgType {
					t.Errorf("UnmarshalRcvdMsg() message type = %v, want %v", msg.RcvdMsgType(), tt.msgType)
				}

				// Type-specific checks could be added here, e.g.:
				if tt.msgType == RcvdMsgTypeError {
					errMsg, ok := msg.(*ErrorMessage)
					if !ok {
						t.Errorf("Expected message to be *ErrorMessage but was %T", msg)
						return
					}
					if errMsg.Error.Message == "" {
						t.Errorf("Expected error message to have a non-empty Error.Message")
					}
				}
			}
		})
	}
}
