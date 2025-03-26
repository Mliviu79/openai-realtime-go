package incoming

import (
	"encoding/json"
	"testing"
)

func TestSessionCreatedMessage(t *testing.T) {
	// Test JSON unmarshal
	jsonData := []byte(`{
		"type": "session.created",
		"message_id": "msg_123",
		"session": {
			"id": "sess_abc123",
			"object": "session",
			"client_secret": {
				"value": "secret-value",
				"expires_at": 1677721600
			}
		}
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal session.created message: %v", err)
	}

	// Verify it's a session created message
	if msg.RcvdMsgType() != RcvdMsgTypeSessionCreated {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeSessionCreated, msg.RcvdMsgType())
	}

	// Cast to specific message type
	sessionMsg, ok := msg.(*SessionCreatedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to SessionCreatedMessage, got %T", msg)
	}

	// Verify fields
	if sessionMsg.ID != "msg_123" {
		t.Errorf("Expected ID to be %q, got %q", "msg_123", sessionMsg.ID)
	}

	if sessionMsg.Session.ID != "sess_abc123" {
		t.Errorf("Expected Session.ID to be %q, got %q", "sess_abc123", sessionMsg.Session.ID)
	}

	if sessionMsg.Session.Object != "session" {
		t.Errorf("Expected Session.Object to be %q, got %q", "session", sessionMsg.Session.Object)
	}

	if sessionMsg.Session.ClientSecretInfo.ClientSecret.Value != "secret-value" {
		t.Errorf("Expected Session.ClientSecret.Value to be %q, got %q", "secret-value", sessionMsg.Session.ClientSecretInfo.ClientSecret.Value)
	}

	// Test marshaling back to JSON
	_, err = json.Marshal(sessionMsg)
	if err != nil {
		t.Fatalf("Failed to marshal session created message: %v", err)
	}
}

func TestSessionUpdatedMessage(t *testing.T) {
	// Test JSON unmarshal
	jsonData := []byte(`{
		"type": "session.updated",
		"message_id": "msg_456",
		"session": {
			"id": "sess_abc123",
			"object": "session",
			"instructions": "Updated instructions"
		}
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal session.updated message: %v", err)
	}

	// Verify it's a session updated message
	if msg.RcvdMsgType() != RcvdMsgTypeSessionUpdated {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeSessionUpdated, msg.RcvdMsgType())
	}

	// Cast to specific message type
	sessionMsg, ok := msg.(*SessionUpdatedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to SessionUpdatedMessage, got %T", msg)
	}

	// Verify fields
	if sessionMsg.ID != "msg_456" {
		t.Errorf("Expected ID to be %q, got %q", "msg_456", sessionMsg.ID)
	}

	if sessionMsg.Session.ID != "sess_abc123" {
		t.Errorf("Expected Session.ID to be %q, got %q", "sess_abc123", sessionMsg.Session.ID)
	}

	if sessionMsg.Session.Object != "session" {
		t.Errorf("Expected Session.Object to be %q, got %q", "session", sessionMsg.Session.Object)
	}

	if sessionMsg.Session.Instructions == nil {
		t.Errorf("Expected Session.Instructions to not be nil")
	} else if *sessionMsg.Session.Instructions != "Updated instructions" {
		t.Errorf("Expected Session.Instructions to be %q, got %q", "Updated instructions", *sessionMsg.Session.Instructions)
	}

	// Test marshaling back to JSON
	_, err = json.Marshal(sessionMsg)
	if err != nil {
		t.Fatalf("Failed to marshal session updated message: %v", err)
	}
}
