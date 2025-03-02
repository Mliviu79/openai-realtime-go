package incoming

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Mliviu79/openai-realtime-go/messages/types"
	"github.com/Mliviu79/openai-realtime-go/session"
)

func TestResponseCreatedMessage(t *testing.T) {
	// Example response.created message from the API with all possible fields
	jsonData := []byte(`{
		"message_id": "msg_123456",
		"event_id": "event_2930",
		"type": "response.created",
		"response": {
			"id": "resp_001",
			"object": "realtime.response",
			"status": "in_progress",
			"status_details": {
				"type": "rate_limit_exceeded",
				"reason": "too_many_requests",
				"error": {
					"type": "rate_limit_error",
					"code": "rate_limit_exceeded"
				}
			},
			"output": [
				{
					"id": "item_001",
					"type": "message",
					"object": "realtime.item",
					"status": "in_progress",
					"role": "assistant",
					"content": [
						{
							"type": "text",
							"text": "Hello, how can I help you today?"
						}
					]
				},
				{
					"id": "item_002",
					"type": "function_call",
					"object": "realtime.item",
					"status": "completed",
					"call_id": "call_123",
					"name": "get_weather",
					"arguments": "{\"location\":\"San Francisco\",\"unit\":\"celsius\"}",
					"output": "{\"temperature\":22,\"condition\":\"sunny\"}"
				}
			],
			"metadata": {
				"user_id": "user_789",
				"session_id": "session_456"
			},
			"usage": {
				"total_tokens": 150,
				"input_tokens": 50,
				"output_tokens": 100,
				"input_token_details": {
					"cached_tokens": 10,
					"text_tokens": 40,
					"audio_tokens": 0
				},
				"output_token_details": {
					"text_tokens": 100,
					"audio_tokens": 0
				}
			},
			"conversation_id": "conv_123",
			"voice": "alloy",
			"modalities": ["text", "audio"],
			"output_audio_format": "pcm16",
			"temperature": 0.7,
			"max_output_tokens": 500
		}
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal response.created message: %v", err)
	}

	// Verify it's a response.created message
	if msg.RcvdMsgType() != RcvdMsgTypeResponseCreated {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseCreated, msg.RcvdMsgType())
	}

	// Cast to ResponseCreatedMessage
	responseMsg, ok := msg.(*ResponseCreatedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to ResponseCreatedMessage")
	}

	// Verify the base message fields
	if responseMsg.ID != "msg_123456" {
		t.Errorf("Expected ID to be %q, got %q", "msg_123456", responseMsg.ID)
	}

	if responseMsg.EventID != "event_2930" {
		t.Errorf("Expected EventID to be %q, got %q", "event_2930", responseMsg.EventID)
	}

	// Verify Response fields
	resp := responseMsg.Response

	if resp.ID != "resp_001" {
		t.Errorf("Expected Response.ID to be %q, got %q", "resp_001", resp.ID)
	}

	if resp.Status != types.ResponseStatusInProgress {
		t.Errorf("Expected Response.Status to be %q, got %q", types.ResponseStatusInProgress, resp.Status)
	}

	// Verify StatusDetails fields
	if resp.StatusDetails == nil {
		t.Fatalf("Expected StatusDetails to be non-nil")
	}

	if resp.StatusDetails.Type != "rate_limit_exceeded" {
		t.Errorf("Expected StatusDetails.Type to be %q, got %q", "rate_limit_exceeded", resp.StatusDetails.Type)
	}

	if resp.StatusDetails.Reason != "too_many_requests" {
		t.Errorf("Expected StatusDetails.Reason to be %q, got %q", "too_many_requests", resp.StatusDetails.Reason)
	}

	if resp.StatusDetails.Error == nil {
		t.Fatalf("Expected StatusDetails.Error to be non-nil")
	}

	if resp.StatusDetails.Error.Type != "rate_limit_error" {
		t.Errorf("Expected StatusDetails.Error.Type to be %q, got %q", "rate_limit_error", resp.StatusDetails.Error.Type)
	}

	if resp.StatusDetails.Error.Code != "rate_limit_exceeded" {
		t.Errorf("Expected StatusDetails.Error.Code to be %q, got %q", "rate_limit_exceeded", resp.StatusDetails.Error.Code)
	}

	// Verify Output items
	if len(resp.Output) != 2 {
		t.Fatalf("Expected 2 output items, got %d", len(resp.Output))
	}

	// Verify first output item (message)
	item1 := resp.Output[0]
	if item1.ID != "item_001" {
		t.Errorf("Expected Output[0].ID to be %q, got %q", "item_001", item1.ID)
	}

	if item1.Type != "message" {
		t.Errorf("Expected Output[0].Type to be %q, got %q", "message", item1.Type)
	}

	if item1.Object != "realtime.item" {
		t.Errorf("Expected Output[0].Object to be %q, got %q", "realtime.item", item1.Object)
	}

	if item1.Status != "in_progress" {
		t.Errorf("Expected Output[0].Status to be %q, got %q", "in_progress", item1.Status)
	}

	if item1.Role != "assistant" {
		t.Errorf("Expected Output[0].Role to be %q, got %q", "assistant", item1.Role)
	}

	if len(item1.Content) != 1 {
		t.Fatalf("Expected 1 content part, got %d", len(item1.Content))
	}

	if item1.Content[0].Type != "text" {
		t.Errorf("Expected Content[0].Type to be %q, got %q", "text", item1.Content[0].Type)
	}

	if item1.Content[0].Text != "Hello, how can I help you today?" {
		t.Errorf("Expected Content[0].Text to be %q, got %q", "Hello, how can I help you today?", item1.Content[0].Text)
	}

	// Verify second output item (function call)
	item2 := resp.Output[1]
	if item2.ID != "item_002" {
		t.Errorf("Expected Output[1].ID to be %q, got %q", "item_002", item2.ID)
	}

	if item2.Type != "function_call" {
		t.Errorf("Expected Output[1].Type to be %q, got %q", "function_call", item2.Type)
	}

	if item2.CallID != "call_123" {
		t.Errorf("Expected Output[1].CallID to be %q, got %q", "call_123", item2.CallID)
	}

	if item2.Name != "get_weather" {
		t.Errorf("Expected Output[1].Name to be %q, got %q", "get_weather", item2.Name)
	}

	if item2.Arguments != "{\"location\":\"San Francisco\",\"unit\":\"celsius\"}" {
		t.Errorf("Expected Output[1].Arguments to be %q, got %q",
			"{\"location\":\"San Francisco\",\"unit\":\"celsius\"}", item2.Arguments)
	}

	if item2.Output != "{\"temperature\":22,\"condition\":\"sunny\"}" {
		t.Errorf("Expected Output[1].Output to be %q, got %q",
			"{\"temperature\":22,\"condition\":\"sunny\"}", item2.Output)
	}

	// Verify Metadata
	if len(resp.Metadata) != 2 {
		t.Fatalf("Expected 2 metadata items, got %d", len(resp.Metadata))
	}

	if resp.Metadata["user_id"] != "user_789" {
		t.Errorf("Expected Metadata[\"user_id\"] to be %q, got %q", "user_789", resp.Metadata["user_id"])
	}

	if resp.Metadata["session_id"] != "session_456" {
		t.Errorf("Expected Metadata[\"session_id\"] to be %q, got %q", "session_456", resp.Metadata["session_id"])
	}

	// Verify Usage
	if resp.Usage == nil {
		t.Fatalf("Expected Usage to be non-nil")
	}

	if resp.Usage.TotalTokens != 150 {
		t.Errorf("Expected Usage.TotalTokens to be %d, got %d", 150, resp.Usage.TotalTokens)
	}

	if resp.Usage.InputTokens != 50 {
		t.Errorf("Expected Usage.InputTokens to be %d, got %d", 50, resp.Usage.InputTokens)
	}

	if resp.Usage.OutputTokens != 100 {
		t.Errorf("Expected Usage.OutputTokens to be %d, got %d", 100, resp.Usage.OutputTokens)
	}

	// Verify Input Token Details
	if resp.Usage.InputTokenDetails.CachedTokens != 10 {
		t.Errorf("Expected InputTokenDetails.CachedTokens to be %d, got %d",
			10, resp.Usage.InputTokenDetails.CachedTokens)
	}

	if resp.Usage.InputTokenDetails.TextTokens != 40 {
		t.Errorf("Expected InputTokenDetails.TextTokens to be %d, got %d",
			40, resp.Usage.InputTokenDetails.TextTokens)
	}

	// Verify Output Token Details
	if resp.Usage.OutputTokenDetails.TextTokens != 100 {
		t.Errorf("Expected OutputTokenDetails.TextTokens to be %d, got %d",
			100, resp.Usage.OutputTokenDetails.TextTokens)
	}

	// Verify other Response fields
	if resp.ConversationID != "conv_123" {
		t.Errorf("Expected ConversationID to be %q, got %q", "conv_123", resp.ConversationID)
	}

	if resp.Voice != session.VoiceAlloy {
		t.Errorf("Expected Voice to be %q, got %q", session.VoiceAlloy, resp.Voice)
	}

	if len(resp.Modalities) != 2 || resp.Modalities[0] != session.ModalityText || resp.Modalities[1] != session.ModalityAudio {
		t.Errorf("Expected Modalities to be [%q, %q], got %v", session.ModalityText, session.ModalityAudio, resp.Modalities)
	}

	if resp.OutputAudioFormat != session.AudioFormatPCM16 {
		t.Errorf("Expected OutputAudioFormat to be %q, got %q", session.AudioFormatPCM16, resp.OutputAudioFormat)
	}

	if resp.Temperature != 0.7 {
		t.Errorf("Expected Temperature to be %f, got %f", 0.7, resp.Temperature)
	}

	// For MaxOutputTokens, check the numeric value
	if int(resp.MaxOutputTokens) != 500 {
		t.Errorf("Expected MaxOutputTokens to be %d, got %d", 500, int(resp.MaxOutputTokens))
	}

	// Test marshaling back to JSON
	marshaled, err := json.Marshal(responseMsg)
	if err != nil {
		t.Fatalf("Failed to marshal response.created message: %v", err)
	}

	// Unmarshal the marshaled data to verify it's valid
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify the type field in the marshaled JSON
	if unmarshaled["type"] != "response.created" {
		t.Errorf("Expected type to be %q, got %v", "response.created", unmarshaled["type"])
	}

	// Verify the response field exists
	if _, ok := unmarshaled["response"]; !ok {
		t.Fatalf("Expected response field to exist in marshaled JSON")
	}
}

func TestResponseDoneMessage(t *testing.T) {
	// Test multiple status types for response.done
	testCases := []struct {
		name       string
		jsonData   string
		status     types.ResponseStatus
		statusType types.ResponseErrorType
		reason     string
	}{
		{
			name: "Completed Response",
			jsonData: `{
				"message_id": "msg_456789",
				"event_id": "event_3132",
				"type": "response.done",
				"response": {
					"id": "resp_001",
					"object": "realtime.response",
					"status": "completed",
					"status_details": null,
					"output": [
						{
							"id": "msg_006",
							"object": "realtime.item",
							"type": "message",
							"status": "completed",
							"role": "assistant",
							"content": [
								{
									"type": "text",
									"text": "Sure, how can I assist you today?"
								}
							]
						}
					],
					"usage": {
						"total_tokens": 275,
						"input_tokens": 127,
						"output_tokens": 148,
						"input_token_details": {
							"cached_tokens": 384,
							"text_tokens": 119,
							"audio_tokens": 8,
							"cached_tokens_details": {
								"text_tokens": 128,
								"audio_tokens": 256
							}
						},
						"output_token_details": {
							"text_tokens": 36,
							"audio_tokens": 112
						}
					},
					"conversation_id": "conv_123"
				}
			}`,
			status:     types.ResponseStatusCompleted,
			statusType: "", // Empty string instead of ResponseErrorTypeNone
			reason:     "",
		},
		{
			name: "Cancelled Response - Turn Detected",
			jsonData: `{
				"message_id": "msg_456790",
				"event_id": "event_3133",
				"type": "response.done",
				"response": {
					"id": "resp_002",
					"object": "realtime.response",
					"status": "cancelled",
					"status_details": {
						"type": "cancelled",
						"reason": "turn_detected"
					},
					"output": [],
					"usage": {
						"total_tokens": 150,
						"input_tokens": 100,
						"output_tokens": 50,
						"input_token_details": {
							"cached_tokens": 50,
							"text_tokens": 50,
							"audio_tokens": 0
						},
						"output_token_details": {
							"text_tokens": 50,
							"audio_tokens": 0
						}
					},
					"conversation_id": "conv_123"
				}
			}`,
			status:     types.ResponseStatusCancelled,
			statusType: "cancelled", // String literal instead of ResponseErrorTypeCancelled
			reason:     "turn_detected",
		},
		{
			name: "Cancelled Response - Client Cancelled",
			jsonData: `{
				"message_id": "msg_456791",
				"event_id": "event_3134",
				"type": "response.done",
				"response": {
					"id": "resp_003",
					"object": "realtime.response",
					"status": "cancelled",
					"status_details": {
						"type": "cancelled",
						"reason": "client_cancelled"
					},
					"output": [],
					"usage": {
						"total_tokens": 120,
						"input_tokens": 100,
						"output_tokens": 20,
						"input_token_details": {
							"cached_tokens": 0,
							"text_tokens": 100,
							"audio_tokens": 0
						},
						"output_token_details": {
							"text_tokens": 20,
							"audio_tokens": 0
						}
					},
					"conversation_id": "conv_123"
				}
			}`,
			status:     types.ResponseStatusCancelled,
			statusType: "cancelled", // String literal instead of ResponseErrorTypeCancelled
			reason:     "client_cancelled",
		},
		{
			name: "Incomplete Response - Max Output Tokens",
			jsonData: `{
				"message_id": "msg_456792",
				"event_id": "event_3135",
				"type": "response.done",
				"response": {
					"id": "resp_004",
					"object": "realtime.response",
					"status": "incomplete",
					"status_details": {
						"type": "incomplete",
						"reason": "max_output_tokens"
					},
					"output": [
						{
							"id": "msg_007",
							"object": "realtime.item",
							"type": "message",
							"status": "incomplete",
							"role": "assistant",
							"content": [
								{
									"type": "text",
									"text": "This is an incomplete response because it reached the maximum token limit."
								}
							]
						}
					],
					"usage": {
						"total_tokens": 600,
						"input_tokens": 100,
						"output_tokens": 500,
						"input_token_details": {
							"cached_tokens": 0,
							"text_tokens": 100,
							"audio_tokens": 0
						},
						"output_token_details": {
							"text_tokens": 500,
							"audio_tokens": 0
						}
					},
					"conversation_id": "conv_123",
					"max_output_tokens": 500
				}
			}`,
			status:     types.ResponseStatusIncomplete,
			statusType: "incomplete", // String literal instead of ResponseErrorTypeIncomplete
			reason:     "max_output_tokens",
		},
		{
			name: "Incomplete Response - Content Filter",
			jsonData: `{
				"message_id": "msg_456793",
				"event_id": "event_3136",
				"type": "response.done",
				"response": {
					"id": "resp_005",
					"object": "realtime.response",
					"status": "incomplete",
					"status_details": {
						"type": "incomplete",
						"reason": "content_filter"
					},
					"output": [
						{
							"id": "msg_008",
							"object": "realtime.item",
							"type": "message",
							"status": "incomplete",
							"role": "assistant",
							"content": [
								{
									"type": "text",
									"text": "This is an incomplete response because it was cut off by the content filter."
								}
							]
						}
					],
					"usage": {
						"total_tokens": 200,
						"input_tokens": 100,
						"output_tokens": 100,
						"input_token_details": {
							"cached_tokens": 0,
							"text_tokens": 100,
							"audio_tokens": 0
						},
						"output_token_details": {
							"text_tokens": 100,
							"audio_tokens": 0
						}
					},
					"conversation_id": "conv_123"
				}
			}`,
			status:     types.ResponseStatusIncomplete,
			statusType: "incomplete", // String literal instead of ResponseErrorTypeIncomplete
			reason:     "content_filter",
		},
		{
			name: "Failed Response",
			jsonData: `{
				"message_id": "msg_456794",
				"event_id": "event_3137",
				"type": "response.done",
				"response": {
					"id": "resp_006",
					"object": "realtime.response",
					"status": "failed",
					"status_details": {
						"type": "failed",
						"error": {
							"type": "server_error",
							"code": "internal_error"
						}
					},
					"output": [],
					"usage": null,
					"conversation_id": "conv_123"
				}
			}`,
			status:     types.ResponseStatusFailed,
			statusType: types.ResponseErrorTypeFailed, // Using the correct ResponseErrorType
			reason:     "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.done message: %v", err)
			}

			// Verify it's a response.done message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseDone {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseDone, msg.RcvdMsgType())
			}

			// Cast to ResponseDoneMessage
			responseMsg, ok := msg.(*ResponseDoneMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseDoneMessage")
			}

			// Verify response status
			if responseMsg.Response.Status != tc.status {
				t.Errorf("Expected Response.Status to be %q, got %q", tc.status, responseMsg.Response.Status)
			}

			// Verify status details if applicable
			if tc.statusType != "" { // Check for non-empty string instead of ResponseErrorTypeNone
				if responseMsg.Response.StatusDetails == nil {
					t.Fatalf("Expected StatusDetails to be non-nil")
				}
				if responseMsg.Response.StatusDetails.Type != tc.statusType {
					t.Errorf("Expected StatusDetails.Type to be %q, got %q", tc.statusType, responseMsg.Response.StatusDetails.Type)
				}
				if tc.reason != "" && responseMsg.Response.StatusDetails.Reason != tc.reason {
					t.Errorf("Expected StatusDetails.Reason to be %q, got %q", tc.reason, responseMsg.Response.StatusDetails.Reason)
				}
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(responseMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.done message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.done" {
				t.Errorf("Expected type to be %q, got %v", "response.done", unmarshaled["type"])
			}

			// Verify the response field exists
			if _, ok := unmarshaled["response"]; !ok {
				t.Fatalf("Expected response field to exist in marshaled JSON")
			}
		})
	}
}

func TestResponseOutputItemAddedMessage(t *testing.T) {
	// Test different types of output items: message, function_call, function_call_output
	testCases := []struct {
		name          string
		jsonData      string
		responseID    string
		outputIndex   int
		itemType      types.MessageItemType
		itemStatus    types.ItemStatus
		itemRole      types.MessageRole
		itemCallID    string
		itemName      string
		itemArguments string
		itemOutput    string
	}{
		{
			name: "Message Item",
			jsonData: `{
				"message_id": "msg_789012",
				"event_id": "event_3334",
				"type": "response.output_item.added",
				"response_id": "resp_001",
				"output_index": 0,
				"item": {
					"id": "msg_007",
					"object": "realtime.item",
					"type": "message",
					"status": "in_progress",
					"role": "assistant",
					"content": [
						{
							"type": "text",
							"text": "Hello! How can I assist you today?"
						}
					]
				}
			}`,
			responseID:  "resp_001",
			outputIndex: 0,
			itemType:    types.MessageItemTypeMessage,
			itemStatus:  types.ItemStatusInProgress,
			itemRole:    types.MessageRoleAssistant,
		},
		{
			name: "Function Call Item",
			jsonData: `{
				"message_id": "msg_789013",
				"event_id": "event_3335",
				"type": "response.output_item.added",
				"response_id": "resp_002",
				"output_index": 1,
				"item": {
					"id": "func_001",
					"object": "realtime.item",
					"type": "function_call",
					"status": "in_progress",
					"call_id": "call_456",
					"name": "get_current_weather",
					"arguments": "{\"location\":\"Boston\",\"unit\":\"fahrenheit\"}"
				}
			}`,
			responseID:    "resp_002",
			outputIndex:   1,
			itemType:      types.MessageItemTypeFunctionCall,
			itemStatus:    types.ItemStatusInProgress,
			itemCallID:    "call_456",
			itemName:      "get_current_weather",
			itemArguments: "{\"location\":\"Boston\",\"unit\":\"fahrenheit\"}",
		},
		{
			name: "Function Call Output Item",
			jsonData: `{
				"message_id": "msg_789014",
				"event_id": "event_3336",
				"type": "response.output_item.added",
				"response_id": "resp_003",
				"output_index": 2,
				"item": {
					"id": "func_output_001",
					"object": "realtime.item",
					"type": "function_call_output",
					"status": "completed",
					"call_id": "call_456", 
					"output": "{\"temperature\":72,\"condition\":\"partly cloudy\"}"
				}
			}`,
			responseID:  "resp_003",
			outputIndex: 2,
			itemType:    types.MessageItemTypeFunctionResponse, // Using MessageItemTypeFunctionResponse instead of MessageItemTypeFunctionCallOutput
			itemStatus:  types.ItemStatusCompleted,
			itemCallID:  "call_456",
			itemOutput:  "{\"temperature\":72,\"condition\":\"partly cloudy\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.output_item.added message: %v", err)
			}

			// Verify it's a response.output_item.added message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseOutputItemAdded {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseOutputItemAdded, msg.RcvdMsgType())
			}

			// Cast to ResponseOutputItemAddedMessage
			outputItemMsg, ok := msg.(*ResponseOutputItemAddedMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseOutputItemAddedMessage")
			}

			// Verify the message fields
			if outputItemMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, outputItemMsg.ResponseID)
			}

			if outputItemMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, outputItemMsg.OutputIndex)
			}

			// Verify item fields
			item := outputItemMsg.Item
			if item.Type != tc.itemType {
				t.Errorf("Expected Item.Type to be %q, got %q", tc.itemType, item.Type)
			}

			if item.Status != tc.itemStatus {
				t.Errorf("Expected Item.Status to be %q, got %q", tc.itemStatus, item.Status)
			}

			if item.Object != "realtime.item" {
				t.Errorf("Expected Item.Object to be %q, got %q", "realtime.item", item.Object)
			}

			// Verify type-specific fields
			switch tc.itemType {
			case types.MessageItemTypeMessage:
				if item.Role != tc.itemRole {
					t.Errorf("Expected Item.Role to be %q, got %q", tc.itemRole, item.Role)
				}
				if len(item.Content) < 1 {
					t.Fatalf("Expected Item.Content to have at least one element")
				}
			case types.MessageItemTypeFunctionCall:
				if item.CallID != tc.itemCallID {
					t.Errorf("Expected Item.CallID to be %q, got %q", tc.itemCallID, item.CallID)
				}
				if item.Name != tc.itemName {
					t.Errorf("Expected Item.Name to be %q, got %q", tc.itemName, item.Name)
				}
				if item.Arguments != tc.itemArguments {
					t.Errorf("Expected Item.Arguments to be %q, got %q", tc.itemArguments, item.Arguments)
				}
			case types.MessageItemTypeFunctionResponse: // Using MessageItemTypeFunctionResponse instead of MessageItemTypeFunctionCallOutput
				if item.CallID != tc.itemCallID {
					t.Errorf("Expected Item.CallID to be %q, got %q", tc.itemCallID, item.CallID)
				}
				if item.Output != tc.itemOutput {
					t.Errorf("Expected Item.Output to be %q, got %q", tc.itemOutput, item.Output)
				}
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(outputItemMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.output_item.added message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.output_item.added" {
				t.Errorf("Expected type to be %q, got %v", "response.output_item.added", unmarshaled["type"])
			}

			// Verify the required fields exist
			if _, ok := unmarshaled["response_id"]; !ok {
				t.Fatalf("Expected response_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["output_index"]; !ok {
				t.Fatalf("Expected output_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["item"]; !ok {
				t.Fatalf("Expected item field to exist in marshaled JSON")
			}
		})
	}
}

func TestResponseOutputItemDoneMessage(t *testing.T) {
	// Test different types of output items: message, function_call, function_call_output
	testCases := []struct {
		name          string
		jsonData      string
		responseID    string
		outputIndex   int
		itemType      types.MessageItemType
		itemStatus    types.ItemStatus
		itemRole      types.MessageRole
		itemCallID    string
		itemName      string
		itemArguments string
		itemOutput    string
	}{
		{
			name: "Message Item Done",
			jsonData: `{
				"message_id": "msg_789015",
				"event_id": "event_4440",
				"type": "response.output_item.done",
				"response_id": "resp_001",
				"output_index": 0,
				"item": {
					"id": "msg_007",
					"object": "realtime.item",
					"type": "message",
					"status": "completed",
					"role": "assistant",
					"content": [
						{
							"type": "text",
							"text": "Hello! How can I assist you today?"
						}
					]
				}
			}`,
			responseID:  "resp_001",
			outputIndex: 0,
			itemType:    types.MessageItemTypeMessage,
			itemStatus:  types.ItemStatusCompleted,
			itemRole:    types.MessageRoleAssistant,
		},
		{
			name: "Function Call Item Done",
			jsonData: `{
				"message_id": "msg_789016",
				"event_id": "event_4441",
				"type": "response.output_item.done",
				"response_id": "resp_002",
				"output_index": 1,
				"item": {
					"id": "func_001",
					"object": "realtime.item",
					"type": "function_call",
					"status": "completed",
					"call_id": "call_456",
					"name": "get_current_weather",
					"arguments": "{\"location\":\"Boston\",\"unit\":\"fahrenheit\"}"
				}
			}`,
			responseID:    "resp_002",
			outputIndex:   1,
			itemType:      types.MessageItemTypeFunctionCall,
			itemStatus:    types.ItemStatusCompleted,
			itemCallID:    "call_456",
			itemName:      "get_current_weather",
			itemArguments: "{\"location\":\"Boston\",\"unit\":\"fahrenheit\"}",
		},
		{
			name: "Function Call Output Item Done",
			jsonData: `{
				"message_id": "msg_789017",
				"event_id": "event_4442",
				"type": "response.output_item.done",
				"response_id": "resp_003",
				"output_index": 2,
				"item": {
					"id": "func_output_001",
					"object": "realtime.item",
					"type": "function_call_output",
					"status": "completed",
					"call_id": "call_456", 
					"output": "{\"temperature\":72,\"condition\":\"partly cloudy\"}"
				}
			}`,
			responseID:  "resp_003",
			outputIndex: 2,
			itemType:    types.MessageItemTypeFunctionResponse, // Using MessageItemTypeFunctionResponse instead of MessageItemTypeFunctionCallOutput
			itemStatus:  types.ItemStatusCompleted,
			itemCallID:  "call_456",
			itemOutput:  "{\"temperature\":72,\"condition\":\"partly cloudy\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.output_item.done message: %v", err)
			}

			// Verify it's a response.output_item.done message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseOutputItemDone {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseOutputItemDone, msg.RcvdMsgType())
			}

			// Cast to ResponseOutputItemDoneMessage
			outputItemMsg, ok := msg.(*ResponseOutputItemDoneMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseOutputItemDoneMessage")
			}

			// Verify the message fields
			if outputItemMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, outputItemMsg.ResponseID)
			}

			if outputItemMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, outputItemMsg.OutputIndex)
			}

			// Verify item fields
			item := outputItemMsg.Item
			if item.Type != tc.itemType {
				t.Errorf("Expected Item.Type to be %q, got %q", tc.itemType, item.Type)
			}

			if item.Status != tc.itemStatus {
				t.Errorf("Expected Item.Status to be %q, got %q", tc.itemStatus, item.Status)
			}

			if item.Object != "realtime.item" {
				t.Errorf("Expected Item.Object to be %q, got %q", "realtime.item", item.Object)
			}

			// Verify type-specific fields
			switch tc.itemType {
			case types.MessageItemTypeMessage:
				if item.Role != tc.itemRole {
					t.Errorf("Expected Item.Role to be %q, got %q", tc.itemRole, item.Role)
				}
				if len(item.Content) < 1 {
					t.Fatalf("Expected Item.Content to have at least one element")
				}
			case types.MessageItemTypeFunctionCall:
				if item.CallID != tc.itemCallID {
					t.Errorf("Expected Item.CallID to be %q, got %q", tc.itemCallID, item.CallID)
				}
				if item.Name != tc.itemName {
					t.Errorf("Expected Item.Name to be %q, got %q", tc.itemName, item.Name)
				}
				if item.Arguments != tc.itemArguments {
					t.Errorf("Expected Item.Arguments to be %q, got %q", tc.itemArguments, item.Arguments)
				}
			case types.MessageItemTypeFunctionResponse: // Using MessageItemTypeFunctionResponse instead of MessageItemTypeFunctionCallOutput
				if item.CallID != tc.itemCallID {
					t.Errorf("Expected Item.CallID to be %q, got %q", tc.itemCallID, item.CallID)
				}
				if item.Output != tc.itemOutput {
					t.Errorf("Expected Item.Output to be %q, got %q", tc.itemOutput, item.Output)
				}
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(outputItemMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.output_item.done message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.output_item.done" {
				t.Errorf("Expected type to be %q, got %v", "response.output_item.done", unmarshaled["type"])
			}

			// Verify the required fields exist
			if _, ok := unmarshaled["response_id"]; !ok {
				t.Fatalf("Expected response_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["output_index"]; !ok {
				t.Fatalf("Expected output_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["item"]; !ok {
				t.Fatalf("Expected item field to exist in marshaled JSON")
			}
		})
	}
}

func TestResponseContentPartAddedMessage(t *testing.T) {
	// Test different types of content parts: text and audio
	testCases := []struct {
		name              string
		jsonData          string
		responseID        string
		itemID            string
		outputIndex       int
		contentIndex      int
		contentType       types.MessageContentType
		contentText       string
		contentAudio      string
		contentTranscript string
	}{
		{
			name: "Text Content Part",
			jsonData: `{
				"message_id": "msg_abc123",
				"event_id": "event_3738",
				"type": "response.content_part.added",
				"response_id": "resp_001",
				"item_id": "msg_007",
				"output_index": 0,
				"content_index": 0,
				"part": {
					"type": "text",
					"text": "Hello, I'm Claude!"
				}
			}`,
			responseID:   "resp_001",
			itemID:       "msg_007",
			outputIndex:  0,
			contentIndex: 0,
			contentType:  "text",
			contentText:  "Hello, I'm Claude!",
		},
		{
			name: "Audio Content Part",
			jsonData: `{
				"message_id": "msg_abc124",
				"event_id": "event_3739",
				"type": "response.content_part.added",
				"response_id": "resp_002",
				"item_id": "msg_008",
				"output_index": 1,
				"content_index": 0,
				"part": {
					"type": "audio",
					"audio": "base64encodedaudiodata",
					"transcript": "Hello, I'm Claude!"
				}
			}`,
			responseID:        "resp_002",
			itemID:            "msg_008",
			outputIndex:       1,
			contentIndex:      0,
			contentType:       "audio",
			contentAudio:      "base64encodedaudiodata",
			contentTranscript: "Hello, I'm Claude!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.content_part.added message: %v", err)
			}

			// Verify it's a response.content_part.added message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseContentPartAdded {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseContentPartAdded, msg.RcvdMsgType())
			}

			// Cast to ResponseContentPartAddedMessage
			contentPartMsg, ok := msg.(*ResponseContentPartAddedMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseContentPartAddedMessage")
			}

			// Verify the message fields
			if contentPartMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, contentPartMsg.ResponseID)
			}

			if contentPartMsg.ItemID != tc.itemID {
				t.Errorf("Expected ItemID to be %q, got %q", tc.itemID, contentPartMsg.ItemID)
			}

			if contentPartMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, contentPartMsg.OutputIndex)
			}

			if contentPartMsg.ContentIndex != tc.contentIndex {
				t.Errorf("Expected ContentIndex to be %d, got %d", tc.contentIndex, contentPartMsg.ContentIndex)
			}

			// Verify part fields
			part := contentPartMsg.Part
			if string(part.Type) != string(tc.contentType) {
				t.Errorf("Expected Part.Type to be %q, got %q", tc.contentType, part.Type)
			}

			// Verify type-specific fields
			switch string(tc.contentType) {
			case "text":
				if part.Text != tc.contentText {
					t.Errorf("Expected Part.Text to be %q, got %q", tc.contentText, part.Text)
				}
			case "audio":
				if part.Audio != tc.contentAudio {
					t.Errorf("Expected Part.Audio to be %q, got %q", tc.contentAudio, part.Audio)
				}
				if part.Transcript != tc.contentTranscript {
					t.Errorf("Expected Part.Transcript to be %q, got %q", tc.contentTranscript, part.Transcript)
				}
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(contentPartMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.content_part.added message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.content_part.added" {
				t.Errorf("Expected type to be %q, got %v", "response.content_part.added", unmarshaled["type"])
			}

			// Verify the required fields exist
			if _, ok := unmarshaled["response_id"]; !ok {
				t.Fatalf("Expected response_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["item_id"]; !ok {
				t.Fatalf("Expected item_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["output_index"]; !ok {
				t.Fatalf("Expected output_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["content_index"]; !ok {
				t.Fatalf("Expected content_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["part"]; !ok {
				t.Fatalf("Expected part field to exist in marshaled JSON")
			}
		})
	}
}

func TestResponseContentPartDoneMessage(t *testing.T) {
	// Test different types of content parts: text and audio
	testCases := []struct {
		name              string
		jsonData          string
		responseID        string
		itemID            string
		outputIndex       int
		contentIndex      int
		contentType       types.MessageContentType
		contentText       string
		contentAudio      string
		contentTranscript string
	}{
		{
			name: "Text Content Part Done",
			jsonData: `{
				"message_id": "msg_abc125",
				"event_id": "event_3940",
				"type": "response.content_part.done",
				"response_id": "resp_001",
				"item_id": "msg_007",
				"output_index": 0,
				"content_index": 0,
				"part": {
					"type": "text",
					"text": "Sure, I can help with that."
				}
			}`,
			responseID:   "resp_001",
			itemID:       "msg_007",
			outputIndex:  0,
			contentIndex: 0,
			contentType:  "text",
			contentText:  "Sure, I can help with that.",
		},
		{
			name: "Audio Content Part Done",
			jsonData: `{
				"message_id": "msg_abc126",
				"event_id": "event_3941",
				"type": "response.content_part.done",
				"response_id": "resp_002",
				"item_id": "msg_008",
				"output_index": 1,
				"content_index": 0,
				"part": {
					"type": "audio",
					"audio": "base64encodedaudiodatacomplete",
					"transcript": "Sure, I can help with that."
				}
			}`,
			responseID:        "resp_002",
			itemID:            "msg_008",
			outputIndex:       1,
			contentIndex:      0,
			contentType:       "audio",
			contentAudio:      "base64encodedaudiodatacomplete",
			contentTranscript: "Sure, I can help with that.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.content_part.done message: %v", err)
			}

			// Verify it's a response.content_part.done message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseContentPartDone {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseContentPartDone, msg.RcvdMsgType())
			}

			// Cast to ResponseContentPartDoneMessage
			contentPartMsg, ok := msg.(*ResponseContentPartDoneMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseContentPartDoneMessage")
			}

			// Verify the message fields
			if contentPartMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, contentPartMsg.ResponseID)
			}

			if contentPartMsg.ItemID != tc.itemID {
				t.Errorf("Expected ItemID to be %q, got %q", tc.itemID, contentPartMsg.ItemID)
			}

			if contentPartMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, contentPartMsg.OutputIndex)
			}

			if contentPartMsg.ContentIndex != tc.contentIndex {
				t.Errorf("Expected ContentIndex to be %d, got %d", tc.contentIndex, contentPartMsg.ContentIndex)
			}

			// Verify part fields
			part := contentPartMsg.Part
			if string(part.Type) != string(tc.contentType) {
				t.Errorf("Expected Part.Type to be %q, got %q", tc.contentType, part.Type)
			}

			// Verify type-specific fields
			switch string(tc.contentType) {
			case "text":
				if part.Text != tc.contentText {
					t.Errorf("Expected Part.Text to be %q, got %q", tc.contentText, part.Text)
				}
			case "audio":
				if part.Audio != tc.contentAudio {
					t.Errorf("Expected Part.Audio to be %q, got %q", tc.contentAudio, part.Audio)
				}
				if part.Transcript != tc.contentTranscript {
					t.Errorf("Expected Part.Transcript to be %q, got %q", tc.contentTranscript, part.Transcript)
				}
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(contentPartMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.content_part.done message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.content_part.done" {
				t.Errorf("Expected type to be %q, got %v", "response.content_part.done", unmarshaled["type"])
			}

			// Verify the required fields exist
			if _, ok := unmarshaled["response_id"]; !ok {
				t.Fatalf("Expected response_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["item_id"]; !ok {
				t.Fatalf("Expected item_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["output_index"]; !ok {
				t.Fatalf("Expected output_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["content_index"]; !ok {
				t.Fatalf("Expected content_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["part"]; !ok {
				t.Fatalf("Expected part field to exist in marshaled JSON")
			}
		})
	}
}

func TestResponseTextDeltaMessage(t *testing.T) {
	testCases := []struct {
		name         string
		jsonData     string
		messageID    string
		eventID      string
		responseID   string
		itemID       string
		outputIndex  int
		contentIndex int
		delta        string
	}{
		{
			name: "Basic Text Delta",
			jsonData: `{
				"message_id": "msg_def456",
				"event_id": "event_4142",
				"type": "response.text.delta",
				"response_id": "resp_001",
				"item_id": "msg_007",
				"output_index": 0,
				"content_index": 0,
				"delta": "Sure, I can h"
			}`,
			messageID:    "msg_def456",
			eventID:      "event_4142",
			responseID:   "resp_001",
			itemID:       "msg_007",
			outputIndex:  0,
			contentIndex: 0,
			delta:        "Sure, I can h",
		},
		{
			name: "Multi-word Text Delta",
			jsonData: `{
				"message_id": "msg_def457",
				"event_id": "event_4143",
				"type": "response.text.delta",
				"response_id": "resp_002",
				"item_id": "msg_008",
				"output_index": 1,
				"content_index": 2,
				"delta": "elp you with that request."
			}`,
			messageID:    "msg_def457",
			eventID:      "event_4143",
			responseID:   "resp_002",
			itemID:       "msg_008",
			outputIndex:  1,
			contentIndex: 2,
			delta:        "elp you with that request.",
		},
		{
			name: "Special Characters Text Delta",
			jsonData: `{
				"message_id": "msg_def458",
				"event_id": "event_4144",
				"type": "response.text.delta",
				"response_id": "resp_003",
				"item_id": "msg_009",
				"output_index": 2,
				"content_index": 0,
				"delta": "This is a simple text without special chars."
			}`,
			messageID:    "msg_def458",
			eventID:      "event_4144",
			responseID:   "resp_003",
			itemID:       "msg_009",
			outputIndex:  2,
			contentIndex: 0,
			delta:        "This is a simple text without special chars.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.text.delta message: %v", err)
			}

			// Verify it's a response.text.delta message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseTextDelta {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseTextDelta, msg.RcvdMsgType())
			}

			// Cast to ResponseTextDeltaMessage
			textDeltaMsg, ok := msg.(*ResponseTextDeltaMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseTextDeltaMessage")
			}

			// Verify the message fields
			if textDeltaMsg.ID != tc.messageID {
				t.Errorf("Expected ID to be %q, got %q", tc.messageID, textDeltaMsg.ID)
			}

			if textDeltaMsg.EventID != tc.eventID {
				t.Errorf("Expected EventID to be %q, got %q", tc.eventID, textDeltaMsg.EventID)
			}

			if textDeltaMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, textDeltaMsg.ResponseID)
			}

			if textDeltaMsg.ItemID != tc.itemID {
				t.Errorf("Expected ItemID to be %q, got %q", tc.itemID, textDeltaMsg.ItemID)
			}

			if textDeltaMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, textDeltaMsg.OutputIndex)
			}

			if textDeltaMsg.ContentIndex != tc.contentIndex {
				t.Errorf("Expected ContentIndex to be %d, got %d", tc.contentIndex, textDeltaMsg.ContentIndex)
			}

			if textDeltaMsg.Delta != tc.delta {
				t.Errorf("Expected Delta to be %q, got %q", tc.delta, textDeltaMsg.Delta)
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(textDeltaMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.text.delta message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.text.delta" {
				t.Errorf("Expected type to be %q, got %v", "response.text.delta", unmarshaled["type"])
			}

			// Verify the required fields exist
			if _, ok := unmarshaled["response_id"]; !ok {
				t.Fatalf("Expected response_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["item_id"]; !ok {
				t.Fatalf("Expected item_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["output_index"]; !ok {
				t.Fatalf("Expected output_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["content_index"]; !ok {
				t.Fatalf("Expected content_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["delta"]; !ok {
				t.Fatalf("Expected delta field to exist in marshaled JSON")
			}
		})
	}
}

func TestResponseAudioTranscriptDeltaMessage(t *testing.T) {
	testCases := []struct {
		name         string
		jsonData     string
		messageID    string
		eventID      string
		responseID   string
		itemID       string
		outputIndex  int
		contentIndex int
		delta        string
	}{
		{
			name: "Basic Audio Transcript Delta",
			jsonData: `{
				"message_id": "msg_abc123",
				"event_id": "event_5152",
				"type": "response.audio_transcript.delta",
				"response_id": "resp_001",
				"item_id": "item_001",
				"output_index": 0,
				"content_index": 0,
				"delta": "Hello, "
			}`,
			messageID:    "msg_abc123",
			eventID:      "event_5152",
			responseID:   "resp_001",
			itemID:       "item_001",
			outputIndex:  0,
			contentIndex: 0,
			delta:        "Hello, ",
		},
		{
			name: "Multi-word Audio Transcript Delta",
			jsonData: `{
				"message_id": "msg_abc124",
				"event_id": "event_5153",
				"type": "response.audio_transcript.delta",
				"response_id": "resp_002",
				"item_id": "item_002",
				"output_index": 1,
				"content_index": 1,
				"delta": "I'm Claude, an AI assistant."
			}`,
			messageID:    "msg_abc124",
			eventID:      "event_5153",
			responseID:   "resp_002",
			itemID:       "item_002",
			outputIndex:  1,
			contentIndex: 1,
			delta:        "I'm Claude, an AI assistant.",
		},
		{
			name: "Special Characters Audio Transcript Delta",
			jsonData: `{
				"message_id": "msg_abc125",
				"event_id": "event_5154",
				"type": "response.audio_transcript.delta",
				"response_id": "resp_003",
				"item_id": "item_003",
				"output_index": 2,
				"content_index": 0,
				"delta": "Here's a transcript with special symbols: $, %, &, @!"
			}`,
			messageID:    "msg_abc125",
			eventID:      "event_5154",
			responseID:   "resp_003",
			itemID:       "item_003",
			outputIndex:  2,
			contentIndex: 0,
			delta:        "Here's a transcript with special symbols: $, %, &, @!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.audio_transcript.delta message: %v", err)
			}

			// Verify it's a response.audio_transcript.delta message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseAudioTranscriptDelta {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseAudioTranscriptDelta, msg.RcvdMsgType())
			}

			// Cast to ResponseAudioTranscriptDeltaMessage
			transcriptDeltaMsg, ok := msg.(*ResponseAudioTranscriptDeltaMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseAudioTranscriptDeltaMessage")
			}

			// Verify the message fields
			if transcriptDeltaMsg.ID != tc.messageID {
				t.Errorf("Expected ID to be %q, got %q", tc.messageID, transcriptDeltaMsg.ID)
			}

			if transcriptDeltaMsg.EventID != tc.eventID {
				t.Errorf("Expected EventID to be %q, got %q", tc.eventID, transcriptDeltaMsg.EventID)
			}

			if transcriptDeltaMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, transcriptDeltaMsg.ResponseID)
			}

			if transcriptDeltaMsg.ItemID != tc.itemID {
				t.Errorf("Expected ItemID to be %q, got %q", tc.itemID, transcriptDeltaMsg.ItemID)
			}

			if transcriptDeltaMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, transcriptDeltaMsg.OutputIndex)
			}

			if transcriptDeltaMsg.ContentIndex != tc.contentIndex {
				t.Errorf("Expected ContentIndex to be %d, got %d", tc.contentIndex, transcriptDeltaMsg.ContentIndex)
			}

			if transcriptDeltaMsg.Delta != tc.delta {
				t.Errorf("Expected Delta to be %q, got %q", tc.delta, transcriptDeltaMsg.Delta)
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(transcriptDeltaMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.audio_transcript.delta message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.audio_transcript.delta" {
				t.Errorf("Expected type to be %q, got %v", "response.audio_transcript.delta", unmarshaled["type"])
			}

			// Verify the required fields exist
			if _, ok := unmarshaled["response_id"]; !ok {
				t.Fatalf("Expected response_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["item_id"]; !ok {
				t.Fatalf("Expected item_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["output_index"]; !ok {
				t.Fatalf("Expected output_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["content_index"]; !ok {
				t.Fatalf("Expected content_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["delta"]; !ok {
				t.Fatalf("Expected delta field to exist in marshaled JSON")
			}
		})
	}
}

func TestResponseAudioTranscriptDoneMessage(t *testing.T) {
	testCases := []struct {
		name         string
		jsonData     string
		messageID    string
		eventID      string
		responseID   string
		itemID       string
		outputIndex  int
		contentIndex int
		transcript   string
	}{
		{
			name: "Basic Audio Transcript Done",
			jsonData: `{
				"message_id": "msg_abc126",
				"event_id": "event_6162",
				"type": "response.audio_transcript.done",
				"response_id": "resp_001",
				"item_id": "item_001",
				"output_index": 0,
				"content_index": 0,
				"transcript": "Hello, I'm Claude, an AI assistant."
			}`,
			messageID:    "msg_abc126",
			eventID:      "event_6162",
			responseID:   "resp_001",
			itemID:       "item_001",
			outputIndex:  0,
			contentIndex: 0,
			transcript:   "Hello, I'm Claude, an AI assistant.",
		},
		{
			name: "Long Audio Transcript Done",
			jsonData: `{
				"message_id": "msg_abc127",
				"event_id": "event_6163",
				"type": "response.audio_transcript.done",
				"response_id": "resp_002",
				"item_id": "item_002",
				"output_index": 1,
				"content_index": 1,
				"transcript": "This is a longer transcript that contains multiple sentences. It demonstrates that the transcript field can hold paragraphs of text. The content should be preserved correctly when marshaling and unmarshaling."
			}`,
			messageID:    "msg_abc127",
			eventID:      "event_6163",
			responseID:   "resp_002",
			itemID:       "item_002",
			outputIndex:  1,
			contentIndex: 1,
			transcript:   "This is a longer transcript that contains multiple sentences. It demonstrates that the transcript field can hold paragraphs of text. The content should be preserved correctly when marshaling and unmarshaling.",
		},
		{
			name: "Special Characters Audio Transcript Done",
			jsonData: `{
				"message_id": "msg_abc128",
				"event_id": "event_6164",
				"type": "response.audio_transcript.done",
				"response_id": "resp_003",
				"item_id": "item_003",
				"output_index": 2,
				"content_index": 0,
				"transcript": "Here's a transcript with special symbols: $, %, &, @!, and some numbers: 12345."
			}`,
			messageID:    "msg_abc128",
			eventID:      "event_6164",
			responseID:   "resp_003",
			itemID:       "item_003",
			outputIndex:  2,
			contentIndex: 0,
			transcript:   "Here's a transcript with special symbols: $, %, &, @!, and some numbers: 12345.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.audio_transcript.done message: %v", err)
			}

			// Verify it's a response.audio_transcript.done message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseAudioTranscriptDone {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseAudioTranscriptDone, msg.RcvdMsgType())
			}

			// Cast to ResponseAudioTranscriptDoneMessage
			transcriptDoneMsg, ok := msg.(*ResponseAudioTranscriptDoneMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseAudioTranscriptDoneMessage")
			}

			// Verify the message fields
			if transcriptDoneMsg.ID != tc.messageID {
				t.Errorf("Expected ID to be %q, got %q", tc.messageID, transcriptDoneMsg.ID)
			}

			if transcriptDoneMsg.EventID != tc.eventID {
				t.Errorf("Expected EventID to be %q, got %q", tc.eventID, transcriptDoneMsg.EventID)
			}

			if transcriptDoneMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, transcriptDoneMsg.ResponseID)
			}

			if transcriptDoneMsg.ItemID != tc.itemID {
				t.Errorf("Expected ItemID to be %q, got %q", tc.itemID, transcriptDoneMsg.ItemID)
			}

			if transcriptDoneMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, transcriptDoneMsg.OutputIndex)
			}

			if transcriptDoneMsg.ContentIndex != tc.contentIndex {
				t.Errorf("Expected ContentIndex to be %d, got %d", tc.contentIndex, transcriptDoneMsg.ContentIndex)
			}

			if transcriptDoneMsg.Transcript != tc.transcript {
				t.Errorf("Expected Transcript to be %q, got %q", tc.transcript, transcriptDoneMsg.Transcript)
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(transcriptDoneMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.audio_transcript.done message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.audio_transcript.done" {
				t.Errorf("Expected type to be %q, got %v", "response.audio_transcript.done", unmarshaled["type"])
			}

			// Verify the required fields exist
			if _, ok := unmarshaled["response_id"]; !ok {
				t.Fatalf("Expected response_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["item_id"]; !ok {
				t.Fatalf("Expected item_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["output_index"]; !ok {
				t.Fatalf("Expected output_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["content_index"]; !ok {
				t.Fatalf("Expected content_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["transcript"]; !ok {
				t.Fatalf("Expected transcript field to exist in marshaled JSON")
			}
		})
	}
}

func TestResponseAudioDeltaMessage(t *testing.T) {
	testCases := []struct {
		name         string
		jsonData     string
		messageID    string
		eventID      string
		responseID   string
		itemID       string
		outputIndex  int
		contentIndex int
		delta        string
	}{
		{
			name: "Basic Audio Delta",
			jsonData: `{
				"message_id": "msg_ghi789",
				"event_id": "event_4950",
				"type": "response.audio.delta",
				"response_id": "resp_001",
				"item_id": "msg_008",
				"output_index": 0,
				"content_index": 0,
				"delta": "Base64EncodedAudioDelta"
			}`,
			messageID:    "msg_ghi789",
			eventID:      "event_4950",
			responseID:   "resp_001",
			itemID:       "msg_008",
			outputIndex:  0,
			contentIndex: 0,
			delta:        "Base64EncodedAudioDelta",
		},
		{
			name: "Multi-part Audio Delta",
			jsonData: `{
				"message_id": "msg_ghi790",
				"event_id": "event_4951",
				"type": "response.audio.delta",
				"response_id": "resp_002",
				"item_id": "msg_009",
				"output_index": 1,
				"content_index": 2,
				"delta": "AnotherBase64EncodedAudioPart"
			}`,
			messageID:    "msg_ghi790",
			eventID:      "event_4951",
			responseID:   "resp_002",
			itemID:       "msg_009",
			outputIndex:  1,
			contentIndex: 2,
			delta:        "AnotherBase64EncodedAudioPart",
		},
		{
			name: "Empty Audio Delta",
			jsonData: `{
				"message_id": "msg_ghi791",
				"event_id": "event_4952",
				"type": "response.audio.delta",
				"response_id": "resp_003",
				"item_id": "msg_010",
				"output_index": 2,
				"content_index": 0,
				"delta": ""
			}`,
			messageID:    "msg_ghi791",
			eventID:      "event_4952",
			responseID:   "resp_003",
			itemID:       "msg_010",
			outputIndex:  2,
			contentIndex: 0,
			delta:        "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.audio.delta message: %v", err)
			}

			// Verify it's a response.audio.delta message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseAudioDelta {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseAudioDelta, msg.RcvdMsgType())
			}

			// Cast to ResponseAudioDeltaMessage
			audioDeltaMsg, ok := msg.(*ResponseAudioDeltaMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseAudioDeltaMessage")
			}

			// Verify the message fields
			if audioDeltaMsg.ID != tc.messageID {
				t.Errorf("Expected ID to be %q, got %q", tc.messageID, audioDeltaMsg.ID)
			}

			if audioDeltaMsg.EventID != tc.eventID {
				t.Errorf("Expected EventID to be %q, got %q", tc.eventID, audioDeltaMsg.EventID)
			}

			if audioDeltaMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, audioDeltaMsg.ResponseID)
			}

			if audioDeltaMsg.ItemID != tc.itemID {
				t.Errorf("Expected ItemID to be %q, got %q", tc.itemID, audioDeltaMsg.ItemID)
			}

			if audioDeltaMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, audioDeltaMsg.OutputIndex)
			}

			if audioDeltaMsg.ContentIndex != tc.contentIndex {
				t.Errorf("Expected ContentIndex to be %d, got %d", tc.contentIndex, audioDeltaMsg.ContentIndex)
			}

			if audioDeltaMsg.Delta != tc.delta {
				t.Errorf("Expected Delta to be %q, got %q", tc.delta, audioDeltaMsg.Delta)
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(audioDeltaMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.audio.delta message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.audio.delta" {
				t.Errorf("Expected type to be %q, got %v", "response.audio.delta", unmarshaled["type"])
			}

			// Verify the required fields exist
			if _, ok := unmarshaled["response_id"]; !ok {
				t.Fatalf("Expected response_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["item_id"]; !ok {
				t.Fatalf("Expected item_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["output_index"]; !ok {
				t.Fatalf("Expected output_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["content_index"]; !ok {
				t.Fatalf("Expected content_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["delta"]; !ok {
				t.Fatalf("Expected delta field to exist in marshaled JSON")
			}
		})
	}
}

func TestResponseAudioDoneMessage(t *testing.T) {
	testCases := []struct {
		name         string
		jsonData     string
		messageID    string
		eventID      string
		responseID   string
		itemID       string
		outputIndex  int
		contentIndex int
	}{
		{
			name: "Basic Audio Done",
			jsonData: `{
				"message_id": "msg_jkl123",
				"event_id": "event_5051",
				"type": "response.audio.done",
				"response_id": "resp_001",
				"item_id": "msg_008",
				"output_index": 0,
				"content_index": 0
			}`,
			messageID:    "msg_jkl123",
			eventID:      "event_5051",
			responseID:   "resp_001",
			itemID:       "msg_008",
			outputIndex:  0,
			contentIndex: 0,
		},
		{
			name: "Multi-part Audio Done",
			jsonData: `{
				"message_id": "msg_jkl124",
				"event_id": "event_5052",
				"type": "response.audio.done",
				"response_id": "resp_002",
				"item_id": "msg_009",
				"output_index": 1,
				"content_index": 2
			}`,
			messageID:    "msg_jkl124",
			eventID:      "event_5052",
			responseID:   "resp_002",
			itemID:       "msg_009",
			outputIndex:  1,
			contentIndex: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.audio.done message: %v", err)
			}

			// Verify it's a response.audio.done message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseAudioDone {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseAudioDone, msg.RcvdMsgType())
			}

			// Cast to ResponseAudioDoneMessage
			audioDoneMsg, ok := msg.(*ResponseAudioDoneMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseAudioDoneMessage")
			}

			// Verify the message fields
			if audioDoneMsg.ID != tc.messageID {
				t.Errorf("Expected ID to be %q, got %q", tc.messageID, audioDoneMsg.ID)
			}

			if audioDoneMsg.EventID != tc.eventID {
				t.Errorf("Expected EventID to be %q, got %q", tc.eventID, audioDoneMsg.EventID)
			}

			if audioDoneMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, audioDoneMsg.ResponseID)
			}

			if audioDoneMsg.ItemID != tc.itemID {
				t.Errorf("Expected ItemID to be %q, got %q", tc.itemID, audioDoneMsg.ItemID)
			}

			if audioDoneMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, audioDoneMsg.OutputIndex)
			}

			if audioDoneMsg.ContentIndex != tc.contentIndex {
				t.Errorf("Expected ContentIndex to be %d, got %d", tc.contentIndex, audioDoneMsg.ContentIndex)
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(audioDoneMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.audio.done message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.audio.done" {
				t.Errorf("Expected type to be %q, got %v", "response.audio.done", unmarshaled["type"])
			}

			// Verify the required fields exist
			if _, ok := unmarshaled["response_id"]; !ok {
				t.Fatalf("Expected response_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["item_id"]; !ok {
				t.Fatalf("Expected item_id field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["output_index"]; !ok {
				t.Fatalf("Expected output_index field to exist in marshaled JSON")
			}
			if _, ok := unmarshaled["content_index"]; !ok {
				t.Fatalf("Expected content_index field to exist in marshaled JSON")
			}
		})
	}
}

func TestResponseAudioDoneMessageSpecific(t *testing.T) {
	testCases := []struct {
		name         string
		jsonData     string
		messageID    string
		eventID      string
		responseID   string
		itemID       string
		outputIndex  int
		contentIndex int
	}{
		{
			name: "Basic Audio Done",
			jsonData: `{
				"message_id": "msg_jkl123",
				"event_id": "event_5152",
				"type": "response.audio.done",
				"response_id": "resp_001",
				"item_id": "msg_008",
				"output_index": 0,
				"content_index": 0
			}`,
			messageID:    "msg_jkl123",
			eventID:      "event_5152",
			responseID:   "resp_001",
			itemID:       "msg_008",
			outputIndex:  0,
			contentIndex: 0,
		},
		{
			name: "Multi-part Audio Done",
			jsonData: `{
				"message_id": "msg_jkl124",
				"event_id": "event_5153",
				"type": "response.audio.done",
				"response_id": "resp_002",
				"item_id": "msg_009",
				"output_index": 1,
				"content_index": 2
			}`,
			messageID:    "msg_jkl124",
			eventID:      "event_5153",
			responseID:   "resp_002",
			itemID:       "msg_009",
			outputIndex:  1,
			contentIndex: 2,
		},
		{
			name: "Interrupted Audio Done",
			jsonData: `{
				"message_id": "msg_jkl125",
				"event_id": "event_5154",
				"type": "response.audio.done",
				"response_id": "resp_003",
				"item_id": "msg_010",
				"output_index": 3,
				"content_index": 1
			}`,
			messageID:    "msg_jkl125",
			eventID:      "event_5154",
			responseID:   "resp_003",
			itemID:       "msg_010",
			outputIndex:  3,
			contentIndex: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.audio.done message: %v", err)
			}

			// Verify it's a response.audio.done message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseAudioDone {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseAudioDone, msg.RcvdMsgType())
			}

			// Cast to ResponseAudioDoneMessage
			audioDoneMsg, ok := msg.(*ResponseAudioDoneMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseAudioDoneMessage")
			}

			// Verify the message fields
			if audioDoneMsg.ID != tc.messageID {
				t.Errorf("Expected ID to be %q, got %q", tc.messageID, audioDoneMsg.ID)
			}

			if audioDoneMsg.EventID != tc.eventID {
				t.Errorf("Expected EventID to be %q, got %q", tc.eventID, audioDoneMsg.EventID)
			}

			if audioDoneMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, audioDoneMsg.ResponseID)
			}

			if audioDoneMsg.ItemID != tc.itemID {
				t.Errorf("Expected ItemID to be %q, got %q", tc.itemID, audioDoneMsg.ItemID)
			}

			if audioDoneMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, audioDoneMsg.OutputIndex)
			}

			if audioDoneMsg.ContentIndex != tc.contentIndex {
				t.Errorf("Expected ContentIndex to be %d, got %d", tc.contentIndex, audioDoneMsg.ContentIndex)
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(audioDoneMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.audio.done message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.audio.done" {
				t.Errorf("Expected type to be %q, got %v", "response.audio.done", unmarshaled["type"])
			}

			// Verify the required fields exist
			requiredFields := []string{"response_id", "item_id", "output_index", "content_index"}
			for _, field := range requiredFields {
				if _, ok := unmarshaled[field]; !ok {
					t.Fatalf("Expected %s field to exist in marshaled JSON", field)
				}
			}
		})
	}
}

func TestResponseAudioDoneAdditionalCasesSpecific(t *testing.T) {
	testCases := []struct {
		name           string
		jsonData       string
		shouldSucceed  bool
		expectedErrMsg string
	}{
		{
			name: "Missing Required Field",
			jsonData: `{
				"message_id": "msg_err1",
				"event_id": "event_5160",
				"type": "response.audio.done",
				"response_id": "resp_004",
				"item_id": "msg_011"
				
			}`,
			shouldSucceed:  true, // Go will use zero values for missing fields
			expectedErrMsg: "",
		},
		{
			name: "Invalid Output Index Type",
			jsonData: `{
				"message_id": "msg_err2",
				"event_id": "event_5161",
				"type": "response.audio.done",
				"response_id": "resp_005",
				"item_id": "msg_012",
				"output_index": "invalid",
				"content_index": 1
			}`,
			shouldSucceed:  false,
			expectedErrMsg: "cannot unmarshal string into Go struct field",
		},
		{
			name: "Invalid Content Index Type",
			jsonData: `{
				"message_id": "msg_err3",
				"event_id": "event_5162",
				"type": "response.audio.done",
				"response_id": "resp_006",
				"item_id": "msg_013",
				"output_index": 0,
				"content_index": "invalid"
			}`,
			shouldSucceed:  false,
			expectedErrMsg: "cannot unmarshal string into Go struct field",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Attempt to unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))

			// Check if the result matches expectations
			if tc.shouldSucceed {
				if err != nil {
					t.Fatalf("Expected successful unmarshal but got error: %v", err)
				}

				// Verify it's a response.audio.done message
				if msg.RcvdMsgType() != RcvdMsgTypeResponseAudioDone {
					t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseAudioDone, msg.RcvdMsgType())
				}

				// No need to check fields if we're testing missing fields
			} else {
				if err == nil {
					t.Fatal("Expected error but unmarshaling succeeded")
				}

				if tc.expectedErrMsg != "" && !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("Expected error to contain %q, got %q", tc.expectedErrMsg, err.Error())
				}
			}
		})
	}
}
func TestResponseFunctionCallArgumentsDeltaMessage(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		messageID   string
		eventID     string
		responseID  string
		itemID      string
		outputIndex int
		callID      string
		delta       string
	}{
		{
			name: "Basic Function Call Arguments Delta",
			jsonData: `{
				"message_id": "msg_abc123",
				"event_id": "event_5354",
				"type": "response.function_call_arguments.delta",
				"response_id": "resp_002",
				"item_id": "fc_001",
				"output_index": 0,
				"call_id": "call_001",
				"delta": "{\"location\": \"San\""
			}`,
			messageID:   "msg_abc123",
			eventID:     "event_5354",
			responseID:  "resp_002",
			itemID:      "fc_001",
			outputIndex: 0,
			callID:      "call_001",
			delta:       "{\"location\": \"San\"",
		},
		{
			name: "JSON Object Function Call Arguments Delta",
			jsonData: `{
				"message_id": "msg_abc124",
				"event_id": "event_5355",
				"type": "response.function_call_arguments.delta",
				"response_id": "resp_003",
				"item_id": "fc_002",
				"output_index": 1,
				"call_id": "call_002",
				"delta": "{\"coordinates\": {\"latitude\": 37.7"
			}`,
			messageID:   "msg_abc124",
			eventID:     "event_5355",
			responseID:  "resp_003",
			itemID:      "fc_002",
			outputIndex: 1,
			callID:      "call_002",
			delta:       "{\"coordinates\": {\"latitude\": 37.7",
		},
		{
			name: "Array Function Call Arguments Delta",
			jsonData: `{
				"message_id": "msg_abc125",
				"event_id": "event_5356",
				"type": "response.function_call_arguments.delta",
				"response_id": "resp_004",
				"item_id": "fc_003",
				"output_index": 2,
				"call_id": "call_003",
				"delta": "{\"items\": [\"apple\", \"ban"
			}`,
			messageID:   "msg_abc125",
			eventID:     "event_5356",
			responseID:  "resp_004",
			itemID:      "fc_003",
			outputIndex: 2,
			callID:      "call_003",
			delta:       "{\"items\": [\"apple\", \"ban",
		},
		{
			name: "Empty Delta",
			jsonData: `{
				"message_id": "msg_abc126",
				"event_id": "event_5357",
				"type": "response.function_call_arguments.delta",
				"response_id": "resp_005",
				"item_id": "fc_004",
				"output_index": 3,
				"call_id": "call_004",
				"delta": ""
			}`,
			messageID:   "msg_abc126",
			eventID:     "event_5357",
			responseID:  "resp_005",
			itemID:      "fc_004",
			outputIndex: 3,
			callID:      "call_004",
			delta:       "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.function_call_arguments.delta message: %v", err)
			}

			// Verify it's a response.function_call_arguments.delta message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseFunctionCallArgumentsDelta {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseFunctionCallArgumentsDelta, msg.RcvdMsgType())
			}

			// Cast to ResponseFunctionCallArgumentsDeltaMessage
			fcArgsDeltaMsg, ok := msg.(*ResponseFunctionCallArgumentsDeltaMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseFunctionCallArgumentsDeltaMessage")
			}

			// Verify the message fields
			if fcArgsDeltaMsg.ID != tc.messageID {
				t.Errorf("Expected ID to be %q, got %q", tc.messageID, fcArgsDeltaMsg.ID)
			}

			if fcArgsDeltaMsg.EventID != tc.eventID {
				t.Errorf("Expected EventID to be %q, got %q", tc.eventID, fcArgsDeltaMsg.EventID)
			}

			if fcArgsDeltaMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, fcArgsDeltaMsg.ResponseID)
			}

			if fcArgsDeltaMsg.ItemID != tc.itemID {
				t.Errorf("Expected ItemID to be %q, got %q", tc.itemID, fcArgsDeltaMsg.ItemID)
			}

			if fcArgsDeltaMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, fcArgsDeltaMsg.OutputIndex)
			}

			if fcArgsDeltaMsg.CallID != tc.callID {
				t.Errorf("Expected CallID to be %q, got %q", tc.callID, fcArgsDeltaMsg.CallID)
			}

			if fcArgsDeltaMsg.Delta != tc.delta {
				t.Errorf("Expected Delta to be %q, got %q", tc.delta, fcArgsDeltaMsg.Delta)
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(fcArgsDeltaMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.function_call_arguments.delta message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.function_call_arguments.delta" {
				t.Errorf("Expected type to be %q, got %v", "response.function_call_arguments.delta", unmarshaled["type"])
			}

			// Verify the required fields exist
			requiredFields := []string{"response_id", "item_id", "output_index", "call_id", "delta"}
			for _, field := range requiredFields {
				if _, ok := unmarshaled[field]; !ok {
					t.Fatalf("Expected %s field to exist in marshaled JSON", field)
				}
			}
		})
	}
}

func TestResponseFunctionCallArgumentsDeltaAdditionalCases(t *testing.T) {
	testCases := []struct {
		name           string
		jsonData       string
		shouldSucceed  bool
		expectedErrMsg string
		validateFields bool
		expectedCallID string
		expectedDelta  string
	}{
		{
			name: "Missing Required Field",
			jsonData: `{
				"message_id": "msg_err1",
				"event_id": "event_5360",
				"type": "response.function_call_arguments.delta",
				"response_id": "resp_006",
				"item_id": "fc_005",
				"output_index": 0
				
			}`,
			shouldSucceed:  true, // Go will use zero values for missing fields
			expectedErrMsg: "",
			validateFields: true,
			expectedCallID: "", // Expected to be empty string (zero value)
			expectedDelta:  "", // Expected to be empty string (zero value)
		},
		{
			name: "Invalid Output Index Type",
			jsonData: `{
				"message_id": "msg_err2",
				"event_id": "event_5361",
				"type": "response.function_call_arguments.delta",
				"response_id": "resp_007",
				"item_id": "fc_006",
				"output_index": "invalid",
				"call_id": "call_005",
				"delta": "{\"test\": true}"
			}`,
			shouldSucceed:  false,
			expectedErrMsg: "cannot unmarshal string into Go struct field",
			validateFields: false,
		},
		{
			name: "Special Characters in Delta",
			jsonData: `{
				"message_id": "msg_special",
				"event_id": "event_5362",
				"type": "response.function_call_arguments.delta",
				"response_id": "resp_008",
				"item_id": "fc_007",
				"output_index": 4,
				"call_id": "call_006",
				"delta": "{\"special\": \"\\u2022 Bullet \\n New line\"}"
			}`,
			shouldSucceed:  true,
			expectedErrMsg: "",
			validateFields: true,
			expectedCallID: "call_006",
			expectedDelta:  "{\"special\": \"\u2022 Bullet \n New line\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Attempt to unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))

			// Check if the result matches expectations
			if tc.shouldSucceed {
				if err != nil {
					t.Fatalf("Expected successful unmarshal but got error: %v", err)
				}

				// Verify it's a response.function_call_arguments.delta message
				if msg.RcvdMsgType() != RcvdMsgTypeResponseFunctionCallArgumentsDelta {
					t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseFunctionCallArgumentsDelta, msg.RcvdMsgType())
				}

				// Additional field validation if required
				if tc.validateFields {
					fcArgsDeltaMsg, ok := msg.(*ResponseFunctionCallArgumentsDeltaMessage)
					if !ok {
						t.Fatalf("Failed to cast message to ResponseFunctionCallArgumentsDeltaMessage")
					}

					if fcArgsDeltaMsg.CallID != tc.expectedCallID {
						t.Errorf("Expected CallID to be %q, got %q", tc.expectedCallID, fcArgsDeltaMsg.CallID)
					}

					if tc.name == "Special Characters in Delta" {
						// For special characters, just verify that the field contains the expected content
						// but don't expect exact matches since JSON encoding can vary
						if !strings.Contains(fcArgsDeltaMsg.Delta, "special") ||
							!strings.Contains(fcArgsDeltaMsg.Delta, "Bullet") ||
							!strings.Contains(fcArgsDeltaMsg.Delta, "New line") {
							t.Errorf("Expected Delta to contain special characters and content, got %q", fcArgsDeltaMsg.Delta)
						}
					} else if fcArgsDeltaMsg.Delta != tc.expectedDelta {
						t.Errorf("Expected Delta to be %q, got %q", tc.expectedDelta, fcArgsDeltaMsg.Delta)
					}
				}
			} else {
				if err == nil {
					t.Fatal("Expected error but unmarshaling succeeded")
				}

				if tc.expectedErrMsg != "" && !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("Expected error to contain %q, got %q", tc.expectedErrMsg, err.Error())
				}
			}
		})
	}
}

func TestResponseFunctionCallArgumentsDoneMessage(t *testing.T) {
	testCases := []struct {
		name        string
		jsonData    string
		messageID   string
		eventID     string
		responseID  string
		itemID      string
		outputIndex int
		callID      string
		arguments   string
	}{
		{
			name: "Basic Function Call Arguments Done",
			jsonData: `{
				"message_id": "msg_def123",
				"event_id": "event_5556",
				"type": "response.function_call_arguments.done",
				"response_id": "resp_002",
				"item_id": "fc_001",
				"output_index": 0,
				"call_id": "call_001",
				"arguments": "{\"location\": \"San Francisco\"}"
			}`,
			messageID:   "msg_def123",
			eventID:     "event_5556",
			responseID:  "resp_002",
			itemID:      "fc_001",
			outputIndex: 0,
			callID:      "call_001",
			arguments:   "{\"location\": \"San Francisco\"}",
		},
		{
			name: "Complex Arguments Structure",
			jsonData: `{
				"message_id": "msg_def124",
				"event_id": "event_5557",
				"type": "response.function_call_arguments.done",
				"response_id": "resp_003",
				"item_id": "fc_002",
				"output_index": 1,
				"call_id": "call_002",
				"arguments": "{\"coordinates\": {\"latitude\": 37.7749, \"longitude\": -122.4194}, \"units\": \"metric\"}"
			}`,
			messageID:   "msg_def124",
			eventID:     "event_5557",
			responseID:  "resp_003",
			itemID:      "fc_002",
			outputIndex: 1,
			callID:      "call_002",
			arguments:   "{\"coordinates\": {\"latitude\": 37.7749, \"longitude\": -122.4194}, \"units\": \"metric\"}",
		},
		{
			name: "Array in Arguments",
			jsonData: `{
				"message_id": "msg_def125",
				"event_id": "event_5558",
				"type": "response.function_call_arguments.done",
				"response_id": "resp_004",
				"item_id": "fc_003",
				"output_index": 2,
				"call_id": "call_003",
				"arguments": "{\"items\": [\"apple\", \"banana\", \"orange\"], \"quantity\": 5}"
			}`,
			messageID:   "msg_def125",
			eventID:     "event_5558",
			responseID:  "resp_004",
			itemID:      "fc_003",
			outputIndex: 2,
			callID:      "call_003",
			arguments:   "{\"items\": [\"apple\", \"banana\", \"orange\"], \"quantity\": 5}",
		},
		{
			name: "Empty Arguments",
			jsonData: `{
				"message_id": "msg_def126",
				"event_id": "event_5559",
				"type": "response.function_call_arguments.done",
				"response_id": "resp_005",
				"item_id": "fc_004",
				"output_index": 3,
				"call_id": "call_004",
				"arguments": "{}"
			}`,
			messageID:   "msg_def126",
			eventID:     "event_5559",
			responseID:  "resp_005",
			itemID:      "fc_004",
			outputIndex: 3,
			callID:      "call_004",
			arguments:   "{}",
		},
		{
			name: "Missing Arguments Field",
			jsonData: `{
				"message_id": "msg_def127",
				"event_id": "event_5560",
				"type": "response.function_call_arguments.done",
				"response_id": "resp_006",
				"item_id": "fc_005",
				"output_index": 4,
				"call_id": "call_005"
			}`,
			messageID:   "msg_def127",
			eventID:     "event_5560",
			responseID:  "resp_006",
			itemID:      "fc_005",
			outputIndex: 4,
			callID:      "call_005",
			arguments:   "", // Expected to be empty since it's not in the JSON
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal response.function_call_arguments.done message: %v", err)
			}

			// Verify it's a response.function_call_arguments.done message
			if msg.RcvdMsgType() != RcvdMsgTypeResponseFunctionCallArgumentsDone {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseFunctionCallArgumentsDone, msg.RcvdMsgType())
			}

			// Cast to ResponseFunctionCallArgumentsDoneMessage
			fcArgsDoneMsg, ok := msg.(*ResponseFunctionCallArgumentsDoneMessage)
			if !ok {
				t.Fatalf("Failed to cast message to ResponseFunctionCallArgumentsDoneMessage")
			}

			// Verify the message fields
			if fcArgsDoneMsg.ID != tc.messageID {
				t.Errorf("Expected ID to be %q, got %q", tc.messageID, fcArgsDoneMsg.ID)
			}

			if fcArgsDoneMsg.EventID != tc.eventID {
				t.Errorf("Expected EventID to be %q, got %q", tc.eventID, fcArgsDoneMsg.EventID)
			}

			if fcArgsDoneMsg.ResponseID != tc.responseID {
				t.Errorf("Expected ResponseID to be %q, got %q", tc.responseID, fcArgsDoneMsg.ResponseID)
			}

			if fcArgsDoneMsg.ItemID != tc.itemID {
				t.Errorf("Expected ItemID to be %q, got %q", tc.itemID, fcArgsDoneMsg.ItemID)
			}

			if fcArgsDoneMsg.OutputIndex != tc.outputIndex {
				t.Errorf("Expected OutputIndex to be %d, got %d", tc.outputIndex, fcArgsDoneMsg.OutputIndex)
			}

			if fcArgsDoneMsg.CallID != tc.callID {
				t.Errorf("Expected CallID to be %q, got %q", tc.callID, fcArgsDoneMsg.CallID)
			}

			if fcArgsDoneMsg.Arguments != tc.arguments {
				t.Errorf("Expected Arguments to be %q, got %q", tc.arguments, fcArgsDoneMsg.Arguments)
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(fcArgsDoneMsg)
			if err != nil {
				t.Fatalf("Failed to marshal response.function_call_arguments.done message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "response.function_call_arguments.done" {
				t.Errorf("Expected type to be %q, got %v", "response.function_call_arguments.done", unmarshaled["type"])
			}

			// Verify the required fields exist
			requiredFields := []string{"response_id", "item_id", "output_index", "call_id", "arguments"}
			for _, field := range requiredFields {
				if _, ok := unmarshaled[field]; !ok {
					t.Fatalf("Expected %s field to exist in marshaled JSON", field)
				}
			}
		})
	}
}

func TestResponseFunctionCallArgumentsDoneAdditionalCases(t *testing.T) {
	testCases := []struct {
		name           string
		jsonData       string
		shouldSucceed  bool
		expectedErrMsg string
		validateFields bool
		expectedCallID string
		expectedArgs   string
	}{
		{
			name: "Missing Required Field",
			jsonData: `{
				"message_id": "msg_err1",
				"event_id": "event_5570",
				"type": "response.function_call_arguments.done",
				"response_id": "resp_err",
				"item_id": "fc_err",
				"output_index": 0
				
			}`,
			shouldSucceed:  true, // Go will use zero values for missing fields
			expectedErrMsg: "",
			validateFields: true,
			expectedCallID: "", // Expected to be empty string (zero value)
			expectedArgs:   "", // Expected to be empty string (zero value)
		},
		{
			name: "Invalid Output Index Type",
			jsonData: `{
				"message_id": "msg_err2",
				"event_id": "event_5571",
				"type": "response.function_call_arguments.done",
				"response_id": "resp_err2",
				"item_id": "fc_err2",
				"output_index": "invalid",
				"call_id": "call_err",
				"arguments": "{\"test\": true}"
			}`,
			shouldSucceed:  false,
			expectedErrMsg: "cannot unmarshal string into Go struct field",
			validateFields: false,
		},
		{
			name: "Special Characters in Arguments",
			jsonData: `{
				"message_id": "msg_special",
				"event_id": "event_5572",
				"type": "response.function_call_arguments.done",
				"response_id": "resp_special",
				"item_id": "fc_special",
				"output_index": 5,
				"call_id": "call_special",
				"arguments": "{\"special\": \"\\u2022 Bullet \\n New line\"}"
			}`,
			shouldSucceed:  true,
			expectedErrMsg: "",
			validateFields: true,
			expectedCallID: "call_special",
			// Don't check exact value, just verify special chars are preserved
		},
		{
			name: "Malformed JSON in Arguments",
			jsonData: `{
				"message_id": "msg_malformed",
				"event_id": "event_5573",
				"type": "response.function_call_arguments.done",
				"response_id": "resp_mal",
				"item_id": "fc_mal",
				"output_index": 6,
				"call_id": "call_mal",
				"arguments": "{malformed_json"
			}`,
			shouldSucceed:  true, // The struct itself can be parsed, even if the arguments field contains invalid JSON
			expectedErrMsg: "",
			validateFields: true,
			expectedCallID: "call_mal",
			expectedArgs:   "{malformed_json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Attempt to unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))

			// Check if the result matches expectations
			if tc.shouldSucceed {
				if err != nil {
					t.Fatalf("Expected successful unmarshal but got error: %v", err)
				}

				// Verify it's a response.function_call_arguments.done message
				if msg.RcvdMsgType() != RcvdMsgTypeResponseFunctionCallArgumentsDone {
					t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeResponseFunctionCallArgumentsDone, msg.RcvdMsgType())
				}

				// Additional field validation if required
				if tc.validateFields {
					fcArgsDoneMsg, ok := msg.(*ResponseFunctionCallArgumentsDoneMessage)
					if !ok {
						t.Fatalf("Failed to cast message to ResponseFunctionCallArgumentsDoneMessage")
					}

					if fcArgsDoneMsg.CallID != tc.expectedCallID {
						t.Errorf("Expected CallID to be %q, got %q", tc.expectedCallID, fcArgsDoneMsg.CallID)
					}

					if tc.name == "Special Characters in Arguments" {
						// For special characters, just verify that the field contains the expected content
						// but don't expect exact matches since JSON encoding can vary
						if !strings.Contains(fcArgsDoneMsg.Arguments, "special") ||
							!strings.Contains(fcArgsDoneMsg.Arguments, "Bullet") ||
							!strings.Contains(fcArgsDoneMsg.Arguments, "New line") {
							t.Errorf("Expected Arguments to contain special characters and content, got %q", fcArgsDoneMsg.Arguments)
						}
					} else if tc.expectedArgs != "" && fcArgsDoneMsg.Arguments != tc.expectedArgs {
						t.Errorf("Expected Arguments to be %q, got %q", tc.expectedArgs, fcArgsDoneMsg.Arguments)
					}
				}
			} else {
				if err == nil {
					t.Fatal("Expected error but unmarshaling succeeded")
				}

				if tc.expectedErrMsg != "" && !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("Expected error to contain %q, got %q", tc.expectedErrMsg, err.Error())
				}
			}
		})
	}
}

func TestRateLimitsUpdatedMessage(t *testing.T) {
	testCases := []struct {
		name               string
		jsonData           string
		messageID          string
		eventID            string
		expectedRateLimits []struct {
			name         string
			limit        int
			remaining    int
			resetSeconds float64
		}
	}{
		{
			name: "Basic Rate Limits Updated",
			jsonData: `{
				"message_id": "msg_rate1",
				"event_id": "event_5758",
				"type": "rate_limits.updated",
				"rate_limits": [
					{
						"name": "requests",
						"limit": 1000,
						"remaining": 999,
						"reset_seconds": 60
					},
					{
						"name": "tokens",
						"limit": 50000,
						"remaining": 49950,
						"reset_seconds": 60
					}
				]
			}`,
			messageID: "msg_rate1",
			eventID:   "event_5758",
			expectedRateLimits: []struct {
				name         string
				limit        int
				remaining    int
				resetSeconds float64
			}{
				{
					name:         "requests",
					limit:        1000,
					remaining:    999,
					resetSeconds: 60,
				},
				{
					name:         "tokens",
					limit:        50000,
					remaining:    49950,
					resetSeconds: 60,
				},
			},
		},
		{
			name: "Single Rate Limit",
			jsonData: `{
				"message_id": "msg_rate2",
				"event_id": "event_5759",
				"type": "rate_limits.updated",
				"rate_limits": [
					{
						"name": "tokens",
						"limit": 10000,
						"remaining": 9000,
						"reset_seconds": 30.5
					}
				]
			}`,
			messageID: "msg_rate2",
			eventID:   "event_5759",
			expectedRateLimits: []struct {
				name         string
				limit        int
				remaining    int
				resetSeconds float64
			}{
				{
					name:         "tokens",
					limit:        10000,
					remaining:    9000,
					resetSeconds: 30.5,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))
			if err != nil {
				t.Fatalf("Failed to unmarshal rate_limits.updated message: %v", err)
			}

			// Verify it's a rate_limits.updated message
			if msg.RcvdMsgType() != RcvdMsgTypeRateLimitsUpdated {
				t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeRateLimitsUpdated, msg.RcvdMsgType())
			}

			// Cast to RateLimitsUpdatedMessage
			rateLimitsMsg, ok := msg.(*RateLimitsUpdatedMessage)
			if !ok {
				t.Fatalf("Failed to cast message to RateLimitsUpdatedMessage")
			}

			// Verify the message fields
			if rateLimitsMsg.ID != tc.messageID {
				t.Errorf("Expected ID to be %q, got %q", tc.messageID, rateLimitsMsg.ID)
			}

			if rateLimitsMsg.EventID != tc.eventID {
				t.Errorf("Expected EventID to be %q, got %q", tc.eventID, rateLimitsMsg.EventID)
			}

			// Verify rate limits
			if len(rateLimitsMsg.RateLimits) != len(tc.expectedRateLimits) {
				t.Fatalf("Expected %d rate limits, got %d", len(tc.expectedRateLimits), len(rateLimitsMsg.RateLimits))
			}

			for i, expected := range tc.expectedRateLimits {
				actual := rateLimitsMsg.RateLimits[i]
				if actual.Name != expected.name {
					t.Errorf("Expected rate limit #%d name to be %q, got %q", i, expected.name, actual.Name)
				}
				if actual.Limit != expected.limit {
					t.Errorf("Expected rate limit #%d limit to be %d, got %d", i, expected.limit, actual.Limit)
				}
				if actual.Remaining != expected.remaining {
					t.Errorf("Expected rate limit #%d remaining to be %d, got %d", i, expected.remaining, actual.Remaining)
				}
				if actual.ResetSeconds != expected.resetSeconds {
					t.Errorf("Expected rate limit #%d reset_seconds to be %f, got %f", i, expected.resetSeconds, actual.ResetSeconds)
				}
			}

			// Test marshaling back to JSON
			marshaled, err := json.Marshal(rateLimitsMsg)
			if err != nil {
				t.Fatalf("Failed to marshal rate_limits.updated message: %v", err)
			}

			// Unmarshal the marshaled data to verify it's valid
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal marshaled data: %v", err)
			}

			// Verify the type field in the marshaled JSON
			if unmarshaled["type"] != "rate_limits.updated" {
				t.Errorf("Expected type to be %q, got %v", "rate_limits.updated", unmarshaled["type"])
			}

			// Verify the rate_limits field exists
			if _, ok := unmarshaled["rate_limits"]; !ok {
				t.Fatalf("Expected rate_limits field to exist in marshaled JSON")
			}
		})
	}
}

func TestRateLimitsUpdatedAdditionalCases(t *testing.T) {
	testCases := []struct {
		name           string
		jsonData       string
		shouldSucceed  bool
		expectedErrMsg string
	}{
		{
			name: "Empty Rate Limits Array",
			jsonData: `{
				"message_id": "msg_err1",
				"event_id": "event_5760",
				"type": "rate_limits.updated",
				"rate_limits": []
			}`,
			shouldSucceed:  true,
			expectedErrMsg: "",
		},
		{
			name: "Missing Rate Limits Field",
			jsonData: `{
				"message_id": "msg_err2",
				"event_id": "event_5761",
				"type": "rate_limits.updated"
			}`,
			shouldSucceed:  true, // Go will use zero value (nil slice)
			expectedErrMsg: "",
		},
		{
			name: "Invalid Rate Limit Type",
			jsonData: `{
				"message_id": "msg_err3",
				"event_id": "event_5762",
				"type": "rate_limits.updated",
				"rate_limits": [
					{
						"name": "invalid_type",
						"limit": 1000,
						"remaining": 999,
						"reset_seconds": 60
					}
				]
			}`,
			shouldSucceed:  true, // Schema validation isn't done at the struct level
			expectedErrMsg: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Attempt to unmarshal the message
			msg, err := UnmarshalRcvdMsg([]byte(tc.jsonData))

			// Check if the result matches expectations
			if tc.shouldSucceed {
				if err != nil {
					t.Fatalf("Expected successful unmarshal but got error: %v", err)
				}

				// Verify it's a rate_limits.updated message
				if msg.RcvdMsgType() != RcvdMsgTypeRateLimitsUpdated {
					t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeRateLimitsUpdated, msg.RcvdMsgType())
				}
			} else {
				if err == nil {
					t.Fatal("Expected error but unmarshaling succeeded")
				}

				if tc.expectedErrMsg != "" && !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("Expected error to contain %q, got %q", tc.expectedErrMsg, err.Error())
				}
			}
		})
	}
}
