package outgoing

import (
	"encoding/json"
	"testing"

	"github.com/Mliviu79/openai-realtime-go/messages/types"
)

func TestAudioBufferAppendMessageStructure(t *testing.T) {
	// Create a sample message based on the OpenAI API reference
	audioData := "Base64EncodedAudioData"

	// Create the message
	message := AudioBufferAppendMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeAudioBufferAppend,
			ID:   "event_456",
		},
		Audio: audioData,
		// Intentionally not setting Transcription
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal AudioBufferAppendMessage to JSON: %v", err)
	}

	// Verify the JSON structure matches the OpenAI API reference
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required fields
	if result["type"] != "input_audio_buffer.append" {
		t.Errorf("Expected type to be 'input_audio_buffer.append', got %v", result["type"])
	}

	if result["audio"] != audioData {
		t.Errorf("Expected audio to be %q, got %v", audioData, result["audio"])
	}

	// Check optional event_id field
	if result["event_id"] != "event_456" {
		t.Errorf("Expected event_id to be 'event_456', got %v", result["event_id"])
	}

	// The transcription field should not be included when nil
	if _, exists := result["transcription"]; exists {
		t.Errorf("Expected transcription field to be omitted, but it was included")
	}

	// Create another message with a mock transcription field (which is not part of OpenAI API)
	messageWithTranscription := NewAudioBufferAppendMessage(audioData, nil)
	messageWithTranscription.ID = "event_789"

	jsonDataWithTranscription, err := json.Marshal(messageWithTranscription)
	if err != nil {
		t.Fatalf("Failed to marshal message with transcription: %v", err)
	}

	var resultWithTranscription map[string]interface{}
	if err := json.Unmarshal(jsonDataWithTranscription, &resultWithTranscription); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify that transcription is not included (should be omitted when nil)
	if _, exists := resultWithTranscription["transcription"]; exists {
		t.Errorf("The transcription field should be omitted when nil")
	}

	// Compare structure with the OpenAI API example
	expectedJSON := `{
		"event_id": "event_456",
		"type": "input_audio_buffer.append",
		"audio": "Base64EncodedAudioData"
	}`

	var expectedResult map[string]interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expectedResult); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	// Check if our structure has only the fields in the expected structure
	for key := range result {
		if _, exists := expectedResult[key]; !exists {
			t.Errorf("Unexpected field %q in the JSON output", key)
		}
	}

	t.Logf("AudioBufferAppendMessage JSON structure matches OpenAI API reference")
}

func TestAudioBufferCommitMessageStructure(t *testing.T) {
	// Create the message
	message := AudioBufferCommitMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeAudioBufferCommit,
			ID:   "event_789",
		},
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal AudioBufferCommitMessage to JSON: %v", err)
	}

	// Verify the JSON structure matches the OpenAI API reference
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required field
	if result["type"] != "input_audio_buffer.commit" {
		t.Errorf("Expected type to be 'input_audio_buffer.commit', got %v", result["type"])
	}

	// Check optional event_id field
	if result["event_id"] != "event_789" {
		t.Errorf("Expected event_id to be 'event_789', got %v", result["event_id"])
	}

	// The previous_item_id field should not be included according to OpenAI API
	if _, exists := result["previous_item_id"]; exists {
		t.Errorf("Expected previous_item_id field to be omitted, but it was included")
	}

	// Compare structure with the OpenAI API example
	expectedJSON := `{
		"event_id": "event_789",
		"type": "input_audio_buffer.commit"
	}`

	var expectedResult map[string]interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expectedResult); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	// Check if our structure has only the fields in the expected structure
	for key := range result {
		if _, exists := expectedResult[key]; !exists {
			t.Errorf("Unexpected field %q in the JSON output", key)
		}
	}

	t.Logf("AudioBufferCommitMessage JSON structure matches OpenAI API reference")
}

func TestAudioBufferClearMessageStructure(t *testing.T) {
	// Create a message using the constructor function
	message := NewAudioBufferClearMessage()
	message.ID = "event_012"

	// Marshal to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal AudioBufferClearMessage to JSON: %v", err)
	}

	// Verify the JSON structure matches the OpenAI API reference
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required field
	if result["type"] != "input_audio_buffer.clear" {
		t.Errorf("Expected type to be 'input_audio_buffer.clear', got %v", result["type"])
	}

	// Check optional event_id field
	if result["event_id"] != "event_012" {
		t.Errorf("Expected event_id to be 'event_012', got %v", result["event_id"])
	}

	// Compare structure with the OpenAI API example
	expectedJSON := `{
		"event_id": "event_012",
		"type": "input_audio_buffer.clear"
	}`

	var expectedResult map[string]interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expectedResult); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	// Check if our structure has only the fields in the expected structure
	for key := range result {
		if _, exists := expectedResult[key]; !exists {
			t.Errorf("Unexpected field %q in the JSON output", key)
		}
	}

	t.Logf("AudioBufferClearMessage JSON structure matches OpenAI API reference")
}

func TestNewAudioBufferAppendMessage(t *testing.T) {
	// Test with audio only
	audio := "base64-encoded-audio-data"
	// Transcription parameter is kept for backward compatibility but is no longer used
	transcription := &types.InputAudioTranscription{
		Language: "en",
	}

	message := NewAudioBufferAppendMessage(audio, transcription)

	// Verify the message type
	if message.Type != OutMsgTypeAudioBufferAppend {
		t.Errorf("Expected Type to be %v, got %v", OutMsgTypeAudioBufferAppend, message.Type)
	}

	// Verify the audio data
	if message.Audio != audio {
		t.Errorf("Expected Audio to be %q, got %q", audio, message.Audio)
	}

	// Test without transcription (should work the same)
	message = NewAudioBufferAppendMessage(audio, nil)

	// Verify the message type
	if message.Type != OutMsgTypeAudioBufferAppend {
		t.Errorf("Expected Type to be %v, got %v", OutMsgTypeAudioBufferAppend, message.Type)
	}

	// Verify the audio data
	if message.Audio != audio {
		t.Errorf("Expected Audio to be %q, got %q", audio, message.Audio)
	}
}

func TestNewAudioBufferCommitMessage(t *testing.T) {
	// previousItemID parameter is kept for backward compatibility but is no longer used
	previousItemID := "item-123"
	message := NewAudioBufferCommitMessage(previousItemID)

	// Verify the message type
	if message.Type != OutMsgTypeAudioBufferCommit {
		t.Errorf("Expected Type to be %v, got %v", OutMsgTypeAudioBufferCommit, message.Type)
	}

	// Test with empty previous item ID (should work the same)
	message = NewAudioBufferCommitMessage("")

	// Verify the message type
	if message.Type != OutMsgTypeAudioBufferCommit {
		t.Errorf("Expected Type to be %v, got %v", OutMsgTypeAudioBufferCommit, message.Type)
	}
}

func TestNewAudioBufferClearMessage(t *testing.T) {
	message := NewAudioBufferClearMessage()

	// Verify the message type
	if message.Type != OutMsgTypeAudioBufferClear {
		t.Errorf("Expected Type to be %v, got %v", OutMsgTypeAudioBufferClear, message.Type)
	}
}

func TestAudioMessageSerialization(t *testing.T) {
	// Test serializing AudioBufferAppendMessage to JSON
	audio := "base64-encoded-audio-data"

	// Transcription parameter is kept for backward compatibility but is no longer used
	transcription := &types.InputAudioTranscription{
		Language: "en",
	}

	appendMessage := NewAudioBufferAppendMessage(audio, transcription)
	appendMessage.ID = "append-event-123"

	jsonData, err := json.Marshal(appendMessage)
	if err != nil {
		t.Fatalf("Failed to marshal AudioBufferAppendMessage to JSON: %v", err)
	}

	// Check that JSON contains expected fields
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if jsonMap["type"] != string(OutMsgTypeAudioBufferAppend) {
		t.Errorf("Expected type field to be %q, got %q", OutMsgTypeAudioBufferAppend, jsonMap["type"])
	}

	if jsonMap["event_id"] != "append-event-123" {
		t.Errorf("Expected event_id field to be %q, got %q", "append-event-123", jsonMap["event_id"])
	}

	if jsonMap["audio"] != audio {
		t.Errorf("Expected audio field to be %q, got %q", audio, jsonMap["audio"])
	}

	// Verify transcription field is not included
	if _, hasTranscription := jsonMap["transcription"]; hasTranscription {
		t.Errorf("Transcription field should not be included in the JSON")
	}
}

func TestOutMsgInterface(t *testing.T) {
	// Verify that all message types implement the OutMsg interface
	var _ OutMsg = NewAudioBufferAppendMessage("", nil)
	var _ OutMsg = NewAudioBufferCommitMessage("")
	var _ OutMsg = NewAudioBufferClearMessage()

	// Test the interface methods
	message := NewAudioBufferAppendMessage("test", nil)
	message.ID = "test-id"

	if message.OutMsgType() != string(OutMsgTypeAudioBufferAppend) {
		t.Errorf("Expected OutMsgType() to be %q, got %q", OutMsgTypeAudioBufferAppend, message.OutMsgType())
	}

	if message.OutMsgID() != "test-id" {
		t.Errorf("Expected OutMsgID() to be %q, got %q", "test-id", message.OutMsgID())
	}
}
