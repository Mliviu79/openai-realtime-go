package incoming

import (
	"encoding/json"
	"testing"
)

func TestAudioBufferCommittedMessage(t *testing.T) {
	// Example input_audio_buffer.committed message from the API
	jsonData := []byte(`{
		"event_id": "event_1121",
		"type": "input_audio_buffer.committed",
		"previous_item_id": "msg_001",
		"item_id": "msg_002"
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal input_audio_buffer.committed message: %v", err)
	}

	// Verify it's an input_audio_buffer.committed message
	if msg.RcvdMsgType() != RcvdMsgTypeAudioBufferCommitted {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeAudioBufferCommitted, msg.RcvdMsgType())
	}

	// Cast to AudioBufferCommittedMessage
	committedMsg, ok := msg.(*AudioBufferCommittedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to AudioBufferCommittedMessage")
	}

	// Verify the fields
	if committedMsg.EventID != "event_1121" {
		t.Errorf("Expected EventID to be %q, got %q", "event_1121", committedMsg.EventID)
	}

	if committedMsg.PreviousItemID != "msg_001" {
		t.Errorf("Expected PreviousItemID to be %q, got %q", "msg_001", committedMsg.PreviousItemID)
	}

	if committedMsg.ItemID != "msg_002" {
		t.Errorf("Expected ItemID to be %q, got %q", "msg_002", committedMsg.ItemID)
	}

	// Test marshaling back to JSON
	marshaled, err := json.Marshal(committedMsg)
	if err != nil {
		t.Fatalf("Failed to marshal input_audio_buffer.committed message: %v", err)
	}

	// Unmarshal the marshaled data to verify it's valid
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify fields in the marshaled JSON
	previousItemID, ok := unmarshaled["previous_item_id"].(string)
	if !ok || previousItemID != "msg_001" {
		t.Errorf("Expected previous_item_id to be %q, got %v", "msg_001", unmarshaled["previous_item_id"])
	}

	itemID, ok := unmarshaled["item_id"].(string)
	if !ok || itemID != "msg_002" {
		t.Errorf("Expected item_id to be %q, got %v", "msg_002", unmarshaled["item_id"])
	}
}

func TestAudioBufferClearedMessage(t *testing.T) {
	// Example input_audio_buffer.cleared message from the API
	jsonData := []byte(`{
		"event_id": "event_1122",
		"type": "input_audio_buffer.cleared"
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal input_audio_buffer.cleared message: %v", err)
	}

	// Verify it's an input_audio_buffer.cleared message
	if msg.RcvdMsgType() != RcvdMsgTypeAudioBufferCleared {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeAudioBufferCleared, msg.RcvdMsgType())
	}

	// Cast to AudioBufferClearedMessage
	clearedMsg, ok := msg.(*AudioBufferClearedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to AudioBufferClearedMessage")
	}

	// Verify the fields
	if clearedMsg.EventID != "event_1122" {
		t.Errorf("Expected EventID to be %q, got %q", "event_1122", clearedMsg.EventID)
	}

	// Test marshaling back to JSON
	marshaled, err := json.Marshal(clearedMsg)
	if err != nil {
		t.Fatalf("Failed to marshal input_audio_buffer.cleared message: %v", err)
	}

	// Unmarshal the marshaled data to verify it's valid
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify the type field is included
	typeVal, ok := unmarshaled["type"].(string)
	if !ok || typeVal != "input_audio_buffer.cleared" {
		t.Errorf("Expected type to be %q, got %v", "input_audio_buffer.cleared", unmarshaled["type"])
	}
}

func TestAudioBufferSpeechStartedMessage(t *testing.T) {
	// Example input_audio_buffer.speech_started message from the API
	jsonData := []byte(`{
		"event_id": "event_1516",
		"type": "input_audio_buffer.speech_started",
		"audio_start_ms": 1000,
		"item_id": "msg_003"
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal input_audio_buffer.speech_started message: %v", err)
	}

	// Verify it's an input_audio_buffer.speech_started message
	if msg.RcvdMsgType() != RcvdMsgTypeAudioBufferSpeechStarted {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeAudioBufferSpeechStarted, msg.RcvdMsgType())
	}

	// Cast to AudioBufferSpeechStartedMessage
	speechStartedMsg, ok := msg.(*AudioBufferSpeechStartedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to AudioBufferSpeechStartedMessage")
	}

	// Verify the fields
	if speechStartedMsg.EventID != "event_1516" {
		t.Errorf("Expected EventID to be %q, got %q", "event_1516", speechStartedMsg.EventID)
	}

	if speechStartedMsg.AudioStartMs != 1000 {
		t.Errorf("Expected AudioStartMs to be %d, got %d", 1000, speechStartedMsg.AudioStartMs)
	}

	if speechStartedMsg.ItemID != "msg_003" {
		t.Errorf("Expected ItemID to be %q, got %q", "msg_003", speechStartedMsg.ItemID)
	}

	// Test marshaling back to JSON
	marshaled, err := json.Marshal(speechStartedMsg)
	if err != nil {
		t.Fatalf("Failed to marshal input_audio_buffer.speech_started message: %v", err)
	}

	// Unmarshal the marshaled data to verify it's valid
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify fields in the marshaled JSON
	audioStartMs, ok := unmarshaled["audio_start_ms"].(float64)
	if !ok || int64(audioStartMs) != 1000 {
		t.Errorf("Expected audio_start_ms to be %d, got %v", 1000, unmarshaled["audio_start_ms"])
	}

	itemID, ok := unmarshaled["item_id"].(string)
	if !ok || itemID != "msg_003" {
		t.Errorf("Expected item_id to be %q, got %v", "msg_003", unmarshaled["item_id"])
	}
}

func TestAudioBufferSpeechStoppedMessage(t *testing.T) {
	// Example input_audio_buffer.speech_stopped message from the API
	jsonData := []byte(`{
		"event_id": "event_1718",
		"type": "input_audio_buffer.speech_stopped",
		"audio_end_ms": 2000,
		"item_id": "msg_003"
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal input_audio_buffer.speech_stopped message: %v", err)
	}

	// Verify it's an input_audio_buffer.speech_stopped message
	if msg.RcvdMsgType() != RcvdMsgTypeAudioBufferSpeechStopped {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeAudioBufferSpeechStopped, msg.RcvdMsgType())
	}

	// Cast to AudioBufferSpeechStoppedMessage
	speechStoppedMsg, ok := msg.(*AudioBufferSpeechStoppedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to AudioBufferSpeechStoppedMessage")
	}

	// Verify the fields
	if speechStoppedMsg.EventID != "event_1718" {
		t.Errorf("Expected EventID to be %q, got %q", "event_1718", speechStoppedMsg.EventID)
	}

	if speechStoppedMsg.AudioEndMs != 2000 {
		t.Errorf("Expected AudioEndMs to be %d, got %d", 2000, speechStoppedMsg.AudioEndMs)
	}

	if speechStoppedMsg.ItemID != "msg_003" {
		t.Errorf("Expected ItemID to be %q, got %q", "msg_003", speechStoppedMsg.ItemID)
	}

	// Test marshaling back to JSON
	marshaled, err := json.Marshal(speechStoppedMsg)
	if err != nil {
		t.Fatalf("Failed to marshal input_audio_buffer.speech_stopped message: %v", err)
	}

	// Unmarshal the marshaled data to verify it's valid
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify fields in the marshaled JSON
	audioEndMs, ok := unmarshaled["audio_end_ms"].(float64)
	if !ok || int64(audioEndMs) != 2000 {
		t.Errorf("Expected audio_end_ms to be %d, got %v", 2000, unmarshaled["audio_end_ms"])
	}

	itemID, ok := unmarshaled["item_id"].(string)
	if !ok || itemID != "msg_003" {
		t.Errorf("Expected item_id to be %q, got %v", "msg_003", unmarshaled["item_id"])
	}
}
