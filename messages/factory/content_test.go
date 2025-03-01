package factory

import (
	"reflect"
	"testing"

	"github.com/Mliviu79/go-openai-realtime/messages/types"
)

func TestTextContent(t *testing.T) {
	text := "Hello, world!"
	content := TextContent(text)

	if content.Type != types.MessageContentTypeText {
		t.Errorf("TextContent().Type = %v, want %v", content.Type, types.MessageContentTypeText)
	}

	if content.Text != text {
		t.Errorf("TextContent().Text = %v, want %v", content.Text, text)
	}

	// Other fields should be empty
	if content.Audio != "" {
		t.Errorf("TextContent().Audio = %v, want empty string", content.Audio)
	}

	if content.Transcript != "" {
		t.Errorf("TextContent().Transcript = %v, want empty string", content.Transcript)
	}

	if content.ID != "" {
		t.Errorf("TextContent().ID = %v, want empty string", content.ID)
	}
}

func TestInputTextContent(t *testing.T) {
	text := "User input"
	content := InputTextContent(text)

	if content.Type != types.MessageContentTypeInputText {
		t.Errorf("InputTextContent().Type = %v, want %v", content.Type, types.MessageContentTypeInputText)
	}

	if content.Text != text {
		t.Errorf("InputTextContent().Text = %v, want %v", content.Text, text)
	}

	// Other fields should be empty
	if content.Audio != "" {
		t.Errorf("InputTextContent().Audio = %v, want empty string", content.Audio)
	}

	if content.Transcript != "" {
		t.Errorf("InputTextContent().Transcript = %v, want empty string", content.Transcript)
	}

	if content.ID != "" {
		t.Errorf("InputTextContent().ID = %v, want empty string", content.ID)
	}
}

func TestInputAudioContent(t *testing.T) {
	audio := "base64-encoded-audio"
	transcript := "Transcribed audio text"
	content := InputAudioContent(audio, transcript)

	if content.Type != types.MessageContentTypeInputAudio {
		t.Errorf("InputAudioContent().Type = %v, want %v", content.Type, types.MessageContentTypeInputAudio)
	}

	if content.Audio != audio {
		t.Errorf("InputAudioContent().Audio = %v, want %v", content.Audio, audio)
	}

	if content.Transcript != transcript {
		t.Errorf("InputAudioContent().Transcript = %v, want %v", content.Transcript, transcript)
	}

	// Other fields should be empty
	if content.Text != "" {
		t.Errorf("InputAudioContent().Text = %v, want empty string", content.Text)
	}

	if content.ID != "" {
		t.Errorf("InputAudioContent().ID = %v, want empty string", content.ID)
	}
}

func TestItemReferenceContent(t *testing.T) {
	id := "item-123456"
	content := ItemReferenceContent(id)

	if content.Type != types.MessageContentTypeItemReference {
		t.Errorf("ItemReferenceContent().Type = %v, want %v", content.Type, types.MessageContentTypeItemReference)
	}

	if content.ID != id {
		t.Errorf("ItemReferenceContent().ID = %v, want %v", content.ID, id)
	}

	// Other fields should be empty
	if content.Text != "" {
		t.Errorf("ItemReferenceContent().Text = %v, want empty string", content.Text)
	}

	if content.Audio != "" {
		t.Errorf("ItemReferenceContent().Audio = %v, want empty string", content.Audio)
	}

	if content.Transcript != "" {
		t.Errorf("ItemReferenceContent().Transcript = %v, want empty string", content.Transcript)
	}
}

func TestAudioContent(t *testing.T) {
	audio := "base64-encoded-audio"
	content := AudioContent(audio)

	if content.Type != types.MessageContentTypeAudio {
		t.Errorf("AudioContent().Type = %v, want %v", content.Type, types.MessageContentTypeAudio)
	}

	if content.Audio != audio {
		t.Errorf("AudioContent().Audio = %v, want %v", content.Audio, audio)
	}

	// Other fields should be empty
	if content.Text != "" {
		t.Errorf("AudioContent().Text = %v, want empty string", content.Text)
	}

	if content.Transcript != "" {
		t.Errorf("AudioContent().Transcript = %v, want empty string", content.Transcript)
	}

	if content.ID != "" {
		t.Errorf("AudioContent().ID = %v, want empty string", content.ID)
	}
}

func TestTranscriptContent(t *testing.T) {
	transcript := "Transcribed audio text"
	content := TranscriptContent(transcript)

	if content.Type != types.MessageContentTypeTranscript {
		t.Errorf("TranscriptContent().Type = %v, want %v", content.Type, types.MessageContentTypeTranscript)
	}

	if content.Transcript != transcript {
		t.Errorf("TranscriptContent().Transcript = %v, want %v", content.Transcript, transcript)
	}

	// Other fields should be empty
	if content.Text != "" {
		t.Errorf("TranscriptContent().Text = %v, want empty string", content.Text)
	}

	if content.Audio != "" {
		t.Errorf("TranscriptContent().Audio = %v, want empty string", content.Audio)
	}

	if content.ID != "" {
		t.Errorf("TranscriptContent().ID = %v, want empty string", content.ID)
	}
}

func TestContentEquality(t *testing.T) {
	// Test that content creators create distinct objects
	text1 := TextContent("Hello")
	text2 := TextContent("Hello")

	if !reflect.DeepEqual(text1, text2) {
		t.Errorf("TextContent equality failed: %v != %v", text1, text2)
	}

	// Test that different content types with same text are not equal
	text := "Same text"
	textContent := TextContent(text)
	inputTextContent := InputTextContent(text)

	if reflect.DeepEqual(textContent, inputTextContent) {
		t.Errorf("TextContent and InputTextContent should not be equal: %v == %v", textContent, inputTextContent)
	}
}
