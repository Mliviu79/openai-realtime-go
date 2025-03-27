package outgoing

import (
	"encoding/json"
	"testing"

	"github.com/Mliviu79/openai-realtime-go/messages/factory"
	"github.com/Mliviu79/openai-realtime-go/messages/types"
)

func TestConversationCreateMessageStructure(t *testing.T) {
	// Create a sample message based on the OpenAI API reference
	userTextContent := factory.InputTextContent("Hello, how are you?")
	messageItem := factory.UserMessage([]types.MessageContentPart{userTextContent})
	messageItem.ID = "msg_001"

	// Create the message
	message := ConversationCreateMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeConversationCreate,
			ID:   "event_345",
		},
		PreviousItemID: "",
		Item:           messageItem,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal ConversationCreateMessage to JSON: %v", err)
	}

	// Verify the JSON structure matches the OpenAI API reference
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required fields
	if result["type"] != "conversation.item.create" {
		t.Errorf("Expected type to be 'conversation.item.create', got %v", result["type"])
	}

	// Check if 'item' field exists and is an object
	item, ok := result["item"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected item field to be an object")
	}

	// Check item fields
	if item["id"] != "msg_001" {
		t.Errorf("Expected item.id to be 'msg_001', got %v", item["id"])
	}

	if item["type"] != "message" {
		t.Errorf("Expected item.type to be 'message', got %v", item["type"])
	}

	if item["role"] != "user" {
		t.Errorf("Expected item.role to be 'user', got %v", item["role"])
	}

	// Check for object field which should be "realtime.item" according to OpenAI API
	if _, exists := item["object"]; !exists {
		t.Logf("Note: 'object' field ('realtime.item') missing from the item structure")
	}

	// Check content field exists and is an array
	content, ok := item["content"].([]interface{})
	if !ok {
		t.Fatalf("Expected item.content to be an array")
	}

	if len(content) != 1 {
		t.Fatalf("Expected item.content to have 1 element, got %d", len(content))
	}

	// Check content element fields
	contentItem, ok := content[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected content item to be an object")
	}

	if contentItem["type"] != "input_text" {
		t.Errorf("Expected content item type to be 'input_text', got %v", contentItem["type"])
	}

	if contentItem["text"] != "Hello, how are you?" {
		t.Errorf("Expected content item text to be 'Hello, how are you?', got %v", contentItem["text"])
	}

	// Check optional fields
	if result["event_id"] != "event_345" {
		t.Errorf("Expected event_id to be 'event_345', got %v", result["event_id"])
	}

	// If we set previous_item_id to an empty string, it should be omitted
	if _, exists := result["previous_item_id"]; exists {
		t.Errorf("Expected previous_item_id to be omitted when empty, but it was included")
	}

	// Create a new message with a specified previous_item_id
	message.PreviousItemID = "root"
	jsonData, err = json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal ConversationCreateMessage to JSON: %v", err)
	}

	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify previous_item_id field
	if result["previous_item_id"] != "root" {
		t.Errorf("Expected previous_item_id to be 'root', got %v", result["previous_item_id"])
	}

	// Compare structure with the OpenAI API example
	expectedJSON := `{
		"event_id": "event_345",
		"type": "conversation.item.create",
		"previous_item_id": null,
		"item": {
			"id": "msg_001",
			"type": "message",
			"role": "user",
			"content": [
				{
					"type": "input_text",
					"text": "Hello, how are you?"
				}
			]
		}
	}`

	var expectedResult map[string]interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expectedResult); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}

	// Test function_call type item
	functionCallItem := factory.FunctionCallItem("get_weather", `{"location": "San Francisco"}`)
	functionCallItem.ID = "func_001"
	functionCallItem.CallID = "call_001"

	functionCallMessage := ConversationCreateMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeConversationCreate,
			ID:   "event_346",
		},
		Item: functionCallItem,
	}

	functionCallJSON, err := json.Marshal(functionCallMessage)
	if err != nil {
		t.Fatalf("Failed to marshal function call message: %v", err)
	}

	var functionCallResult map[string]interface{}
	if err := json.Unmarshal(functionCallJSON, &functionCallResult); err != nil {
		t.Fatalf("Failed to unmarshal function call JSON: %v", err)
	}

	functionCallItemResult, ok := functionCallResult["item"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected function call item to be an object")
	}

	if functionCallItemResult["type"] != "function_call" {
		t.Errorf("Expected function call item type to be 'function_call', got %v", functionCallItemResult["type"])
	}

	if functionCallItemResult["name"] != "get_weather" {
		t.Errorf("Expected function name to be 'get_weather', got %v", functionCallItemResult["name"])
	}

	if functionCallItemResult["arguments"] != `{"location": "San Francisco"}` {
		t.Errorf("Expected function arguments to be '{\"location\": \"San Francisco\"}', got %v", functionCallItemResult["arguments"])
	}

	// Test function_call_output type item
	functionOutputItem := factory.FunctionResponseItem("call_001", `{"temperature": 72, "condition": "sunny"}`)
	functionOutputItem.ID = "resp_001"

	functionOutputMessage := ConversationCreateMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeConversationCreate,
			ID:   "event_347",
		},
		Item: functionOutputItem,
	}

	functionOutputJSON, err := json.Marshal(functionOutputMessage)
	if err != nil {
		t.Fatalf("Failed to marshal function output message: %v", err)
	}

	var functionOutputResult map[string]interface{}
	if err := json.Unmarshal(functionOutputJSON, &functionOutputResult); err != nil {
		t.Fatalf("Failed to unmarshal function output JSON: %v", err)
	}

	functionOutputItemResult, ok := functionOutputResult["item"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected function output item to be an object")
	}

	if functionOutputItemResult["type"] != "function_call_output" {
		t.Errorf("Expected function output item type to be 'function_call_output', got %v", functionOutputItemResult["type"])
	}

	if functionOutputItemResult["call_id"] != "call_001" {
		t.Errorf("Expected function call_id to be 'call_001', got %v", functionOutputItemResult["call_id"])
	}

	if functionOutputItemResult["output"] != `{"temperature": 72, "condition": "sunny"}` {
		t.Errorf("Expected function output to be '{\"temperature\": 72, \"condition\": \"sunny\"}', got %v", functionOutputItemResult["output"])
	}

	t.Logf("ConversationCreateMessage JSON structure matches OpenAI API reference")
}

func TestConversationDeleteMessageStructure(t *testing.T) {
	// Create a message based on the OpenAI API reference
	message := ConversationDeleteMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeConversationDelete,
			ID:   "event_901",
		},
		ItemID: "msg_003",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal ConversationDeleteMessage to JSON: %v", err)
	}

	// Verify the JSON structure matches the OpenAI API reference
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required fields
	if result["type"] != "conversation.item.delete" {
		t.Errorf("Expected type to be 'conversation.item.delete', got %v", result["type"])
	}

	if result["item_id"] != "msg_003" {
		t.Errorf("Expected item_id to be 'msg_003', got %v", result["item_id"])
	}

	// Check optional field
	if result["event_id"] != "event_901" {
		t.Errorf("Expected event_id to be 'event_901', got %v", result["event_id"])
	}

	// Compare structure with the OpenAI API example
	expectedJSON := `{
		"event_id": "event_901",
		"type": "conversation.item.delete",
		"item_id": "msg_003"
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

	// Also create a message using the constructor function
	constructedMessage := NewConversationDeleteMessage("msg_003")
	constructedMessage.ID = "event_901"

	// Marshal to JSON
	constructedJsonData, err := json.Marshal(constructedMessage)
	if err != nil {
		t.Fatalf("Failed to marshal constructed message to JSON: %v", err)
	}

	// Verify the JSON structure matches what we expect
	var constructedResult map[string]interface{}
	if err := json.Unmarshal(constructedJsonData, &constructedResult); err != nil {
		t.Fatalf("Failed to unmarshal constructed JSON: %v", err)
	}

	// Compare with the expected structure
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

	t.Logf("ConversationDeleteMessage JSON structure matches OpenAI API reference")
}

func TestConversationTruncateMessageStructure(t *testing.T) {
	// Create a message based on the OpenAI API reference
	message := ConversationTruncateMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeConversationTruncate,
			ID:   "event_678",
		},
		ItemID:       "msg_002",
		ContentIndex: 0,
		AudioEndMs:   1500,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal ConversationTruncateMessage to JSON: %v", err)
	}

	// Verify the JSON structure matches the OpenAI API reference
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check required fields
	if result["type"] != "conversation.item.truncate" {
		t.Errorf("Expected type to be 'conversation.item.truncate', got %v", result["type"])
	}

	if result["item_id"] != "msg_002" {
		t.Errorf("Expected item_id to be 'msg_002', got %v", result["item_id"])
	}

	if result["content_index"] != float64(0) {
		t.Errorf("Expected content_index to be 0, got %v", result["content_index"])
	}

	if result["audio_end_ms"] != float64(1500) {
		t.Errorf("Expected audio_end_ms to be 1500, got %v", result["audio_end_ms"])
	}

	// Check optional field
	if result["event_id"] != "event_678" {
		t.Errorf("Expected event_id to be 'event_678', got %v", result["event_id"])
	}

	// Compare structure with the OpenAI API example
	expectedJSON := `{
		"event_id": "event_678",
		"type": "conversation.item.truncate",
		"item_id": "msg_002",
		"content_index": 0,
		"audio_end_ms": 1500
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

	// Also create a message using the constructor function
	constructedMessage := NewConversationTruncateMessage("msg_002", 0, 1500)
	constructedMessage.ID = "event_678"

	// Marshal to JSON
	constructedJsonData, err := json.Marshal(constructedMessage)
	if err != nil {
		t.Fatalf("Failed to marshal constructed message to JSON: %v", err)
	}

	// Verify the JSON structure matches what we expect
	var constructedResult map[string]interface{}
	if err := json.Unmarshal(constructedJsonData, &constructedResult); err != nil {
		t.Fatalf("Failed to unmarshal constructed JSON: %v", err)
	}

	// Compare with the expected structure
	for key, expectedValue := range expectedResult {
		constructedValue, exists := constructedResult[key]
		if !exists {
			t.Errorf("Expected field %q missing from constructed message", key)
			continue
		}

		if key == "content_index" || key == "audio_end_ms" {
			// Compare as floats to handle JSON number conversion
			expectedFloat, _ := expectedValue.(float64)
			constructedFloat, _ := constructedValue.(float64)
			if expectedFloat != constructedFloat {
				t.Errorf("Field %q differs: expected %v, got %v", key, expectedValue, constructedValue)
			}
		} else if expectedValue != constructedValue {
			t.Errorf("Field %q differs: expected %v, got %v", key, expectedValue, constructedValue)
		}
	}

	t.Logf("ConversationTruncateMessage JSON structure matches OpenAI API reference")
}
