package types

import (
	"encoding/json"
	"testing"
)

func TestMessageContentTypes(t *testing.T) {
	tests := []struct {
		name        string
		contentType MessageContentType
		expected    string
	}{
		{
			name:        "Text",
			contentType: MessageContentTypeText,
			expected:    "text",
		},
		{
			name:        "InputText",
			contentType: MessageContentTypeInputText,
			expected:    "input_text",
		},
		{
			name:        "InputAudio",
			contentType: MessageContentTypeInputAudio,
			expected:    "input_audio",
		},
		{
			name:        "Audio",
			contentType: MessageContentTypeAudio,
			expected:    "audio",
		},
		{
			name:        "Transcript",
			contentType: MessageContentTypeTranscript,
			expected:    "transcript",
		},
		{
			name:        "ItemReference",
			contentType: MessageContentTypeItemReference,
			expected:    "item_reference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.contentType) != tt.expected {
				t.Errorf("Expected content type %s to be %q, got %q", tt.name, tt.expected, tt.contentType)
			}
		})
	}
}

func TestMessageContentPartSerialization(t *testing.T) {
	// Test text content
	textContent := MessageContentPart{
		Type: MessageContentTypeText,
		Text: "Hello, world!",
	}

	jsonData, err := json.Marshal(textContent)
	if err != nil {
		t.Fatalf("Failed to marshal text content: %v", err)
	}

	var unmarshaled MessageContentPart
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal text content: %v", err)
	}

	if unmarshaled.Type != MessageContentTypeText {
		t.Errorf("Expected Type to be %v, got %v", MessageContentTypeText, unmarshaled.Type)
	}

	if unmarshaled.Text != "Hello, world!" {
		t.Errorf("Expected Text to be %q, got %q", "Hello, world!", unmarshaled.Text)
	}

	// Test audio content
	audioContent := MessageContentPart{
		Type:  MessageContentTypeAudio,
		Audio: "audio-data",
	}

	jsonData, err = json.Marshal(audioContent)
	if err != nil {
		t.Fatalf("Failed to marshal audio content: %v", err)
	}

	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal audio content: %v", err)
	}

	if unmarshaled.Type != MessageContentTypeAudio {
		t.Errorf("Expected Type to be %v, got %v", MessageContentTypeAudio, unmarshaled.Type)
	}

	if unmarshaled.Audio != "audio-data" {
		t.Errorf("Expected Audio to be %q, got %q", "audio-data", unmarshaled.Audio)
	}

	// Test transcript content
	transcriptContent := MessageContentPart{
		Type:       MessageContentTypeTranscript,
		Transcript: "transcript text",
	}

	jsonData, err = json.Marshal(transcriptContent)
	if err != nil {
		t.Fatalf("Failed to marshal transcript content: %v", err)
	}

	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal transcript content: %v", err)
	}

	if unmarshaled.Type != MessageContentTypeTranscript {
		t.Errorf("Expected Type to be %v, got %v", MessageContentTypeTranscript, unmarshaled.Type)
	}

	if unmarshaled.Transcript != "transcript text" {
		t.Errorf("Expected Transcript to be %q, got %q", "transcript text", unmarshaled.Transcript)
	}

	// Test item reference content
	referenceContent := MessageContentPart{
		Type: MessageContentTypeItemReference,
		ID:   "ref-123",
	}

	jsonData, err = json.Marshal(referenceContent)
	if err != nil {
		t.Fatalf("Failed to marshal reference content: %v", err)
	}

	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal reference content: %v", err)
	}

	if unmarshaled.Type != MessageContentTypeItemReference {
		t.Errorf("Expected Type to be %v, got %v", MessageContentTypeItemReference, unmarshaled.Type)
	}

	if unmarshaled.ID != "ref-123" {
		t.Errorf("Expected ID to be %q, got %q", "ref-123", unmarshaled.ID)
	}
}

func TestInputAudioTranscription(t *testing.T) {
	// Test with all fields set
	transcription := InputAudioTranscription{
		Model:    TranscriptionModelWhisper1,
		Language: "fr",
		Prompt:   "Testing prompt",
	}

	jsonData, err := json.Marshal(transcription)
	if err != nil {
		t.Fatalf("Failed to marshal transcription: %v", err)
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	// Check each field
	if jsonMap["model"] != string(TranscriptionModelWhisper1) {
		t.Errorf("Expected model to be %q, got %v", TranscriptionModelWhisper1, jsonMap["model"])
	}

	if jsonMap["language"] != "fr" {
		t.Errorf("Expected language to be %q, got %v", "fr", jsonMap["language"])
	}

	if jsonMap["prompt"] != "Testing prompt" {
		t.Errorf("Expected prompt to be %q, got %v", "Testing prompt", jsonMap["prompt"])
	}
}
