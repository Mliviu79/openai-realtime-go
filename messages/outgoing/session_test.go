package outgoing

import (
	"encoding/json"
	"testing"

	"github.com/Mliviu79/go-openai-realtime/session"
)

func TestSessionUpdateMessage(t *testing.T) {
	// Create a sample model value
	model := session.GPT4oRealtimePreview

	// Create sample modalities
	modalities := []session.Modality{session.ModalityText, session.ModalityAudio}

	// Create sample instructions
	instructions := "You are a helpful assistant."

	// Create sample voice
	voice := session.VoiceSage

	// Create sample input audio format
	inputFormat := session.AudioFormatPCM16

	// Create sample output audio format
	outputFormat := session.AudioFormatPCM16

	// Create sample input audio transcription
	transcription := session.InputAudioTranscription{
		Model: session.TranscriptionModelWhisper1,
	}

	// Create sample turn detection values
	createResponse := true
	interruptResponse := true
	turnDetection := session.TurnDetection{
		Type:              session.TurnDetectionTypeServerVad,
		Threshold:         0.5,
		PrefixPaddingMs:   300,
		SilenceDurationMs: 500,
		CreateResponse:    &createResponse,
		InterruptResponse: &interruptResponse,
	}

	// Create sample tools
	tools := []session.Tool{
		{
			Type:        "function",
			Name:        "get_weather",
			Description: "Get the current weather...",
			Parameters:  json.RawMessage(`{"type":"object","properties":{"location":{"type":"string"}},"required":["location"]}`),
		},
	}

	// Create sample tool choice
	toolChoice := session.ToolChoiceObj{
		Type: session.ToolChoiceAuto,
	}

	// Create sample temperature
	temperature := 0.8

	// Create sample max response output tokens
	maxTokens := session.Inf

	// Create a session update message
	sessionReq := session.SessionRequest{
		Modalities:              &modalities,
		Model:                   &model,
		Instructions:            &instructions,
		Voice:                   &voice,
		InputAudioFormat:        &inputFormat,
		OutputAudioFormat:       &outputFormat,
		InputAudioTranscription: &transcription,
		TurnDetection:           &turnDetection,
		Tools:                   &tools,
		ToolChoice:              &toolChoice,
		Temperature:             &temperature,
		MaxResponseOutputTokens: &maxTokens,
	}

	// Create the session update message
	updateMsg := SessionUpdateMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeSessionUpdate,
			ID:   "event_123",
		},
		Session: sessionReq,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(updateMsg)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify the marshaled JSON matches what we expect
	// Here we are doing a basic structure check to ensure the right fields exist

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check top-level fields
	if result["event_id"] != "event_123" {
		t.Errorf("Expected event_id to be 'event_123', got %v", result["event_id"])
	}
	if result["type"] != "session.update" {
		t.Errorf("Expected type to be 'session.update', got %v", result["type"])
	}

	// Check session field exists
	if _, ok := result["session"]; !ok {
		t.Fatalf("Expected session field, but it's missing")
	}

	// This is a basic test. A more comprehensive test would check all fields.
	// The key point is that we're verifying our structure matches what OpenAI expects
	t.Logf("Session update message structure was validated successfully")
}
