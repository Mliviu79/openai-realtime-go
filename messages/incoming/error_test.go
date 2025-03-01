package incoming

import (
	"encoding/json"
	"testing"
)

func TestErrorMessage(t *testing.T) {
	// Example error message from the API
	jsonData := []byte(`{
		"event_id": "event_890",
		"type": "error",
		"error": {
			"type": "invalid_request_error",
			"code": "invalid_event",
			"message": "The 'type' field is missing.",
			"param": null,
			"event_id": "event_567"
		}
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal error message: %v", err)
	}

	// Verify it's an error message
	if msg.RcvdMsgType() != RcvdMsgTypeError {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeError, msg.RcvdMsgType())
	}

	// Cast to ErrorMessage
	errMsg, ok := msg.(*ErrorMessage)
	if !ok {
		t.Fatalf("Failed to cast message to ErrorMessage")
	}

	// Verify the fields
	if errMsg.EventID != "event_890" {
		t.Errorf("Expected EventID to be %q, got %q", "event_890", errMsg.EventID)
	}

	if errMsg.Error.Type != "invalid_request_error" {
		t.Errorf("Expected Error.Type to be %q, got %q", "invalid_request_error", errMsg.Error.Type)
	}

	if errMsg.Error.Code != "invalid_event" {
		t.Errorf("Expected Error.Code to be %q, got %q", "invalid_event", errMsg.Error.Code)
	}

	if errMsg.Error.Message != "The 'type' field is missing." {
		t.Errorf("Expected Error.Message to be %q, got %q", "The 'type' field is missing.", errMsg.Error.Message)
	}

	if errMsg.Error.EventID != "event_567" {
		t.Errorf("Expected Error.EventID to be %q, got %q", "event_567", errMsg.Error.EventID)
	}

	// Param should be nil
	if errMsg.Error.Param != nil {
		t.Errorf("Expected Error.Param to be nil, got %v", *errMsg.Error.Param)
	}

	// Test marshaling back to JSON
	_, err = json.Marshal(errMsg)
	if err != nil {
		t.Fatalf("Failed to marshal error message: %v", err)
	}
}
