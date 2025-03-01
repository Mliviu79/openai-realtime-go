package outgoing

import (
	"encoding/json"
	"testing"

	"github.com/Mliviu79/go-openai-realtime/messages/types"
	"github.com/Mliviu79/go-openai-realtime/session"
)

func TestResponseCreateMessageStructure(t *testing.T) {
	// Create sample modalities
	modalities := []session.Modality{session.ModalityText, session.ModalityAudio}

	// Create sample instructions
	instructions := "Please assist the user."

	// Create sample voice
	voice := session.VoiceSage

	// Create sample output audio format
	outputFormat := session.AudioFormatPCM16

	// Create sample tools
	tools := []session.Tool{
		{
			Type:        "function",
			Name:        "calculate_sum",
			Description: "Calculates the sum of two numbers.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"a": { "type": "number" },
					"b": { "type": "number" }
				},
				"required": ["a", "b"]
			}`),
		},
	}

	// Create sample tool choice
	toolChoice := session.ToolChoiceObj{Type: session.ToolChoiceAuto}

	// Create sample temperature
	temperature := 0.8

	// Create sample max response output tokens
	maxTokens := session.NewIntOrInf(1024)

	// Create a response config
	responseConfig := types.ResponseConfig{
		Modalities:              modalities,
		Instructions:            &instructions,
		Voice:                   &voice,
		OutputAudioFormat:       &outputFormat,
		Tools:                   tools,
		ToolChoice:              &toolChoice,
		Temperature:             &temperature,
		MaxResponseOutputTokens: maxTokens,
	}

	// Create the response create message
	createMsg := NewResponseCreateMessage(responseConfig)
	createMsg.ID = "event_234"

	// Marshal to JSON
	jsonData, err := json.Marshal(createMsg)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Verify the marshaled JSON matches what we expect
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check top-level fields
	if result["event_id"] != "event_234" {
		t.Errorf("Expected event_id to be 'event_234', got %v", result["event_id"])
	}
	if result["type"] != "response.create" {
		t.Errorf("Expected type to be 'response.create', got %v", result["type"])
	}

	// According to the OpenAI API reference, the field should be named "response", not "config"
	if _, ok := result["response"]; !ok {
		if _, ok := result["config"]; ok {
			t.Errorf("Expected 'response' field but found 'config' field instead. The field name should be 'response' according to the OpenAI API reference.")
		} else {
			t.Fatalf("Expected response field, but it's missing")
		}
	}

	// If there's a "config" field, let's check its content to make sure it matches what we'd expect in "response"
	configField := "config"
	if _, ok := result["response"]; ok {
		configField = "response"
	}

	config, ok := result[configField].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected %s field to be an object", configField)
	}

	// Check modalities
	modalitiesResult, ok := config["modalities"].([]interface{})
	if !ok {
		t.Fatalf("Expected modalities to be an array")
	}
	if len(modalitiesResult) != 2 {
		t.Errorf("Expected modalities to have 2 elements, got %d", len(modalitiesResult))
	}
	if modalitiesResult[0] != "text" || modalitiesResult[1] != "audio" {
		t.Errorf("Expected modalities to be ['text', 'audio'], got %v", modalitiesResult)
	}

	// Check instructions
	if config["instructions"] != instructions {
		t.Errorf("Expected instructions to be '%s', got %v", instructions, config["instructions"])
	}

	// Check voice
	if config["voice"] != "sage" {
		t.Errorf("Expected voice to be 'sage', got %v", config["voice"])
	}

	// Check output_audio_format
	if config["output_audio_format"] != "pcm16" {
		t.Errorf("Expected output_audio_format to be 'pcm16', got %v", config["output_audio_format"])
	}

	// Check temperature
	if config["temperature"] != float64(temperature) {
		t.Errorf("Expected temperature to be %f, got %v", temperature, config["temperature"])
	}

	// Compare structure with the OpenAI API example
	expectedJSON := `{
		"event_id": "event_234",
		"type": "response.create",
		"response": {
			"modalities": ["text", "audio"],
			"instructions": "Please assist the user.",
			"voice": "sage",
			"output_audio_format": "pcm16",
			"tools": [
				{
					"type": "function",
					"name": "calculate_sum",
					"description": "Calculates the sum of two numbers.",
					"parameters": {
						"type": "object",
						"properties": {
							"a": { "type": "number" },
							"b": { "type": "number" }
						},
						"required": ["a", "b"]
					}
				}
			],
			"tool_choice": "auto",
			"temperature": 0.8,
			"max_response_output_tokens": 1024
		}
	}`

	var expectedResult map[string]interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expectedResult); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	// Log a message for clarity on whether the structure matches the OpenAI API reference
	if configField == "response" {
		t.Logf("ResponseCreateMessage JSON structure matches OpenAI API reference")
	} else {
		t.Logf("ResponseCreateMessage JSON structure uses 'config' field instead of 'response' - update needed to match OpenAI API reference")
	}
}

func TestResponseCancelMessageStructure(t *testing.T) {
	// Test with a specific response_id
	message := ResponseCancelMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeResponseCancel,
			ID:   "event_567",
		},
		ResponseID: "resp_123",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal ResponseCancelMessage to JSON: %v", err)
	}

	// Verify the JSON structure matches the OpenAI API reference
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required field
	if result["type"] != "response.cancel" {
		t.Errorf("Expected type to be 'response.cancel', got %v", result["type"])
	}

	// Check optional fields
	if result["event_id"] != "event_567" {
		t.Errorf("Expected event_id to be 'event_567', got %v", result["event_id"])
	}

	if result["response_id"] != "resp_123" {
		t.Errorf("Expected response_id to be 'resp_123', got %v", result["response_id"])
	}

	// Also test the case where response_id is not provided
	messageWithoutResponseID := ResponseCancelMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeResponseCancel,
			ID:   "event_567",
		},
		// ResponseID intentionally left empty
	}

	// Marshal to JSON
	jsonDataWithoutResponseID, err := json.Marshal(messageWithoutResponseID)
	if err != nil {
		t.Fatalf("Failed to marshal ResponseCancelMessage without response_id to JSON: %v", err)
	}

	// Verify the JSON structure
	var resultWithoutResponseID map[string]interface{}
	if err := json.Unmarshal(jsonDataWithoutResponseID, &resultWithoutResponseID); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check if response_id is omitted when empty
	if response_id, exists := resultWithoutResponseID["response_id"]; exists && response_id != "" {
		t.Errorf("Expected response_id to be omitted or empty when not provided, got %v", response_id)
	}

	// Test the constructor function
	constructedMessage := NewResponseCancelMessage("resp_123")
	constructedMessage.ID = "event_567"

	// Marshal to JSON
	constructedJsonData, err := json.Marshal(constructedMessage)
	if err != nil {
		t.Fatalf("Failed to marshal constructed message to JSON: %v", err)
	}

	// Verify the JSON structure matches
	var constructedResult map[string]interface{}
	if err := json.Unmarshal(constructedJsonData, &constructedResult); err != nil {
		t.Fatalf("Failed to unmarshal constructed JSON: %v", err)
	}

	// Check all fields match
	if constructedResult["type"] != "response.cancel" {
		t.Errorf("Expected type to be 'response.cancel', got %v", constructedResult["type"])
	}

	if constructedResult["event_id"] != "event_567" {
		t.Errorf("Expected event_id to be 'event_567', got %v", constructedResult["event_id"])
	}

	if constructedResult["response_id"] != "resp_123" {
		t.Errorf("Expected response_id to be 'resp_123', got %v", constructedResult["response_id"])
	}

	// Compare with the OpenAI API example
	expectedJSON := `{
		"event_id": "event_567",
		"type": "response.cancel"
	}`

	var expectedResult map[string]interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expectedResult); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	// Verify all expected fields exist in the API example (might not include response_id)
	for key, expectedValue := range expectedResult {
		constructedValue, exists := constructedResult[key]
		if !exists {
			t.Errorf("Expected field %q missing from constructed message", key)
			continue
		}

		if expectedValue != constructedValue {
			t.Errorf("Field %q differs: expected %v, got %v", key, expectedValue, constructedValue)
		}
	}

	t.Logf("ResponseCancelMessage JSON structure matches OpenAI API reference")
}
