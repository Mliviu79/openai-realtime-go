package types

import (
	"encoding/json"
	"testing"

	"github.com/Mliviu79/go-openai-realtime/session"
)

func TestMessageRoles(t *testing.T) {
	tests := []struct {
		name     string
		role     MessageRole
		expected string
	}{
		{
			name:     "System",
			role:     MessageRoleSystem,
			expected: "system",
		},
		{
			name:     "User",
			role:     MessageRoleUser,
			expected: "user",
		},
		{
			name:     "Assistant",
			role:     MessageRoleAssistant,
			expected: "assistant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.role) != tt.expected {
				t.Errorf("Expected role %s to be %q, got %q", tt.name, tt.expected, tt.role)
			}
		})
	}
}

func TestMessageSerialization(t *testing.T) {
	// Create a simple message with text content
	msg := Message{
		Role: MessageRoleUser,
		Content: []MessageContentPart{
			{
				Type: MessageContentTypeText,
				Text: "Hello, world!",
			},
		},
		EndTurn: true,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}

	var unmarshaled Message
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}

	// Verify role
	if unmarshaled.Role != MessageRoleUser {
		t.Errorf("Expected Role to be %v, got %v", MessageRoleUser, unmarshaled.Role)
	}

	// Verify content
	if len(unmarshaled.Content) != 1 {
		t.Fatalf("Expected 1 content part, got %d", len(unmarshaled.Content))
	}

	if unmarshaled.Content[0].Type != MessageContentTypeText {
		t.Errorf("Expected content Type to be %v, got %v", MessageContentTypeText, unmarshaled.Content[0].Type)
	}

	if unmarshaled.Content[0].Text != "Hello, world!" {
		t.Errorf("Expected content Text to be %q, got %q", "Hello, world!", unmarshaled.Content[0].Text)
	}

	// Verify end_turn
	if !unmarshaled.EndTurn {
		t.Errorf("Expected EndTurn to be true, got false")
	}

	// Create a more complex message with tools
	// First, create the parameters as JSON
	paramsJSON, err := json.Marshal(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "The city and state, e.g. San Francisco, CA",
			},
		},
		"required": []string{"location"},
	})
	if err != nil {
		t.Fatalf("Failed to marshal parameters: %v", err)
	}

	complexMsg := Message{
		Role: MessageRoleAssistant,
		Content: []MessageContentPart{
			{
				Type: MessageContentTypeText,
				Text: "I can help with that.",
			},
		},
		Tools: []session.Tool{
			{
				Type:        "function",
				Name:        "get_weather",
				Description: "Get the current weather",
				Parameters:  paramsJSON,
			},
		},
		Name:     "get_weather",
		Metadata: map[string]string{"user_id": "123"},
	}

	jsonData, err = json.Marshal(complexMsg)
	if err != nil {
		t.Fatalf("Failed to marshal complex message: %v", err)
	}

	// Check that JSON contains expected fields
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	if jsonMap["role"] != "assistant" {
		t.Errorf("Expected role to be %q, got %v", "assistant", jsonMap["role"])
	}

	if jsonMap["name"] != "get_weather" {
		t.Errorf("Expected name to be %q, got %v", "get_weather", jsonMap["name"])
	}

	content, ok := jsonMap["content"].([]interface{})
	if !ok || len(content) != 1 {
		t.Fatalf("Expected content to be array of length 1, got %v", jsonMap["content"])
	}

	tools, ok := jsonMap["tools"].([]interface{})
	if !ok || len(tools) != 1 {
		t.Fatalf("Expected tools to be array of length 1, got %v", jsonMap["tools"])
	}

	metadata, ok := jsonMap["metadata"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected metadata to be a map, got %v", jsonMap["metadata"])
	}

	if metadata["user_id"] != "123" {
		t.Errorf("Expected metadata.user_id to be %q, got %v", "123", metadata["user_id"])
	}
}
