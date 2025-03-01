package incoming

import (
	"encoding/json"
	"testing"
)

func TestConversationCreatedMessage(t *testing.T) {
	// Example conversation.created message from the API
	jsonData := []byte(`{
		"event_id": "event_9101",
		"type": "conversation.created",
		"conversation": {
			"id": "conv_001",
			"object": "realtime.conversation"
		}
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal conversation.created message: %v", err)
	}

	// Verify it's a conversation.created message
	if msg.RcvdMsgType() != RcvdMsgTypeConversationCreated {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeConversationCreated, msg.RcvdMsgType())
	}

	// Cast to ConversationCreatedMessage
	convMsg, ok := msg.(*ConversationCreatedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to ConversationCreatedMessage")
	}

	// Verify the fields
	if convMsg.EventID != "event_9101" {
		t.Errorf("Expected EventID to be %q, got %q", "event_9101", convMsg.EventID)
	}

	if convMsg.Conversation.ID != "conv_001" {
		t.Errorf("Expected Conversation.ID to be %q, got %q", "conv_001", convMsg.Conversation.ID)
	}

	if convMsg.Conversation.Object != "realtime.conversation" {
		t.Errorf("Expected Conversation.Object to be %q, got %q", "realtime.conversation", convMsg.Conversation)
	}

	// Verify Items is empty
	if len(convMsg.Conversation.Items) != 0 {
		t.Errorf("Expected Conversation.Items to be empty, got %d items", len(convMsg.Conversation.Items))
	}

	// Test marshaling back to JSON
	marshaled, err := json.Marshal(convMsg)
	if err != nil {
		t.Fatalf("Failed to marshal conversation.created message: %v", err)
	}

	// Unmarshal the marshaled data to verify it's valid
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify the conversation is included
	conv, ok := unmarshaled["conversation"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected conversation field to be an object")
	}

	// Verify conversation fields
	id, ok := conv["id"].(string)
	if !ok || id != "conv_001" {
		t.Errorf("Expected conversation.id to be %q, got %v", "conv_001", conv["id"])
	}

	obj, ok := conv["object"].(string)
	if !ok || obj != "realtime.conversation" {
		t.Errorf("Expected conversation.object to be %q, got %v", "realtime.conversation", conv["object"])
	}

	// Verify 'items' field isn't included
	if _, exists := conv["items"]; exists {
		t.Errorf("Expected 'items' field to be omitted, but it was included")
	}
}

func TestConversationItemCreatedMessage(t *testing.T) {
	// Example conversation.item.created message from the API
	jsonData := []byte(`{
		"event_id": "event_1920",
		"type": "conversation.item.created",
		"previous_item_id": "msg_002",
		"item": {
			"id": "msg_003",
			"object": "realtime.item",
			"type": "message",
			"status": "completed",
			"role": "user",
			"content": [
				{
					"type": "input_audio",
					"transcript": "hello how are you",
					"audio": "base64encodedaudio=="
				}
			]
		}
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal conversation.item.created message: %v", err)
	}

	// Verify it's a conversation.item.created message
	if msg.RcvdMsgType() != RcvdMsgTypeConversationItemCreated {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeConversationItemCreated, msg.RcvdMsgType())
	}

	// Cast to ConversationItemCreatedMessage
	itemMsg, ok := msg.(*ConversationItemCreatedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to ConversationItemCreatedMessage")
	}

	// Verify the fields
	if itemMsg.EventID != "event_1920" {
		t.Errorf("Expected EventID to be %q, got %q", "event_1920", itemMsg.EventID)
	}

	if itemMsg.PreviousItemID != "msg_002" {
		t.Errorf("Expected PreviousItemID to be %q, got %q", "msg_002", itemMsg.PreviousItemID)
	}

	// Verify item fields
	if itemMsg.Item.ID != "msg_003" {
		t.Errorf("Expected Item.ID to be %q, got %q", "msg_003", itemMsg.Item.ID)
	}

	if itemMsg.Item.Object != "realtime.item" {
		t.Errorf("Expected Item.Object to be %q, got %q", "realtime.item", itemMsg.Item.Object)
	}

	if itemMsg.Item.Type != "message" {
		t.Errorf("Expected Item.Type to be %q, got %q", "message", itemMsg.Item.Type)
	}

	if itemMsg.Item.Status != "completed" {
		t.Errorf("Expected Item.Status to be %q, got %q", "completed", itemMsg.Item.Status)
	}

	if itemMsg.Item.Role != "user" {
		t.Errorf("Expected Item.Role to be %q, got %q", "user", itemMsg.Item.Role)
	}

	// Verify content
	if len(itemMsg.Item.Content) != 1 {
		t.Fatalf("Expected 1 content part, got %d", len(itemMsg.Item.Content))
	}

	contentPart := itemMsg.Item.Content[0]
	if contentPart.Type != "input_audio" {
		t.Errorf("Expected content type to be %q, got %q", "input_audio", contentPart.Type)
	}

	if contentPart.Transcript != "hello how are you" {
		t.Errorf("Expected transcript to be %q, got %q", "hello how are you", contentPart.Transcript)
	}

	if contentPart.Audio != "base64encodedaudio==" {
		t.Errorf("Expected audio to be %q, got %q", "base64encodedaudio==", contentPart.Audio)
	}

	// Test marshaling back to JSON
	marshaled, err := json.Marshal(itemMsg)
	if err != nil {
		t.Fatalf("Failed to marshal conversation.item.created message: %v", err)
	}

	// Unmarshal the marshaled data to verify it's valid
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify the previous_item_id is included
	prevItemID, ok := unmarshaled["previous_item_id"].(string)
	if !ok || prevItemID != "msg_002" {
		t.Errorf("Expected previous_item_id to be %q, got %v", "msg_002", unmarshaled["previous_item_id"])
	}

	// Verify the item is included
	item, ok := unmarshaled["item"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected item field to be an object")
	}

	// Verify item fields in the marshaled JSON
	id, ok := item["id"].(string)
	if !ok || id != "msg_003" {
		t.Errorf("Expected item.id to be %q, got %v", "msg_003", item["id"])
	}

	object, ok := item["object"].(string)
	if !ok || object != "realtime.item" {
		t.Errorf("Expected item.object to be %q, got %v", "realtime.item", item["object"])
	}
}

func TestConversationItemTranscriptionCompletedMessage(t *testing.T) {
	// Example conversation.item.input_audio_transcription.completed message from the API
	jsonData := []byte(`{
		"event_id": "event_2122",
		"type": "conversation.item.input_audio_transcription.completed",
		"item_id": "msg_003",
		"content_index": 0,
		"transcript": "Hello, how are you?"
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal conversation.item.input_audio_transcription.completed message: %v", err)
	}

	// Verify it's a conversation.item.input_audio_transcription.completed message
	if msg.RcvdMsgType() != RcvdMsgTypeConversationItemInputAudioTranscriptionCompleted {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeConversationItemInputAudioTranscriptionCompleted, msg.RcvdMsgType())
	}

	// Cast to ConversationItemTranscriptionCompletedMessage
	transcriptMsg, ok := msg.(*ConversationItemTranscriptionCompletedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to ConversationItemTranscriptionCompletedMessage")
	}

	// Verify the fields
	if transcriptMsg.EventID != "event_2122" {
		t.Errorf("Expected EventID to be %q, got %q", "event_2122", transcriptMsg.EventID)
	}

	if transcriptMsg.ItemID != "msg_003" {
		t.Errorf("Expected ItemID to be %q, got %q", "msg_003", transcriptMsg.ItemID)
	}

	if transcriptMsg.ContentIndex != 0 {
		t.Errorf("Expected ContentIndex to be %d, got %d", 0, transcriptMsg.ContentIndex)
	}

	if transcriptMsg.Transcript != "Hello, how are you?" {
		t.Errorf("Expected Transcript to be %q, got %q", "Hello, how are you?", transcriptMsg.Transcript)
	}

	// Test marshaling back to JSON
	marshaled, err := json.Marshal(transcriptMsg)
	if err != nil {
		t.Fatalf("Failed to marshal conversation.item.input_audio_transcription.completed message: %v", err)
	}

	// Unmarshal the marshaled data to verify it's valid
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify fields in the marshaled JSON
	itemID, ok := unmarshaled["item_id"].(string)
	if !ok || itemID != "msg_003" {
		t.Errorf("Expected item_id to be %q, got %v", "msg_003", unmarshaled["item_id"])
	}

	contentIndex, ok := unmarshaled["content_index"].(float64)
	if !ok || int(contentIndex) != 0 {
		t.Errorf("Expected content_index to be %d, got %v", 0, unmarshaled["content_index"])
	}

	transcript, ok := unmarshaled["transcript"].(string)
	if !ok || transcript != "Hello, how are you?" {
		t.Errorf("Expected transcript to be %q, got %v", "Hello, how are you?", unmarshaled["transcript"])
	}
}

func TestConversationItemTranscriptionFailedMessage(t *testing.T) {
	// Example conversation.item.input_audio_transcription.failed message from the API
	jsonData := []byte(`{
		"event_id": "event_2324",
		"type": "conversation.item.input_audio_transcription.failed",
		"item_id": "msg_003",
		"content_index": 0,
		"error": {
			"type": "transcription_error",
			"code": "audio_unintelligible",
			"message": "The audio could not be transcribed.",
			"param": null
		}
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal conversation.item.input_audio_transcription.failed message: %v", err)
	}

	// Verify it's a conversation.item.input_audio_transcription.failed message
	if msg.RcvdMsgType() != RcvdMsgTypeConversationItemInputAudioTranscriptionFailed {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeConversationItemInputAudioTranscriptionFailed, msg.RcvdMsgType())
	}

	// Cast to ConversationItemTranscriptionFailedMessage
	failedMsg, ok := msg.(*ConversationItemTranscriptionFailedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to ConversationItemTranscriptionFailedMessage")
	}

	// Verify the fields
	if failedMsg.EventID != "event_2324" {
		t.Errorf("Expected EventID to be %q, got %q", "event_2324", failedMsg.EventID)
	}

	if failedMsg.ItemID != "msg_003" {
		t.Errorf("Expected ItemID to be %q, got %q", "msg_003", failedMsg.ItemID)
	}

	if failedMsg.ContentIndex != 0 {
		t.Errorf("Expected ContentIndex to be %d, got %d", 0, failedMsg.ContentIndex)
	}

	// Verify error details
	if failedMsg.Error.Type != "transcription_error" {
		t.Errorf("Expected Error.Type to be %q, got %q", "transcription_error", failedMsg.Error.Type)
	}

	if failedMsg.Error.Code != "audio_unintelligible" {
		t.Errorf("Expected Error.Code to be %q, got %q", "audio_unintelligible", failedMsg.Error.Code)
	}

	if failedMsg.Error.Message != "The audio could not be transcribed." {
		t.Errorf("Expected Error.Message to be %q, got %q", "The audio could not be transcribed.", failedMsg.Error.Message)
	}

	if failedMsg.Error.Param != nil {
		t.Errorf("Expected Error.Param to be nil, got %v", *failedMsg.Error.Param)
	}

	// Test marshaling back to JSON
	marshaled, err := json.Marshal(failedMsg)
	if err != nil {
		t.Fatalf("Failed to marshal conversation.item.input_audio_transcription.failed message: %v", err)
	}

	// Unmarshal the marshaled data to verify it's valid
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify fields in the marshaled JSON
	itemID, ok := unmarshaled["item_id"].(string)
	if !ok || itemID != "msg_003" {
		t.Errorf("Expected item_id to be %q, got %v", "msg_003", unmarshaled["item_id"])
	}

	contentIndex, ok := unmarshaled["content_index"].(float64)
	if !ok || int(contentIndex) != 0 {
		t.Errorf("Expected content_index to be %d, got %v", 0, unmarshaled["content_index"])
	}

	// Verify the error object in marshaled JSON
	errorObj, ok := unmarshaled["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected error field to be an object")
	}

	errorType, ok := errorObj["type"].(string)
	if !ok || errorType != "transcription_error" {
		t.Errorf("Expected error.type to be %q, got %v", "transcription_error", errorObj["type"])
	}

	errorCode, ok := errorObj["code"].(string)
	if !ok || errorCode != "audio_unintelligible" {
		t.Errorf("Expected error.code to be %q, got %v", "audio_unintelligible", errorObj["code"])
	}

	errorMessage, ok := errorObj["message"].(string)
	if !ok || errorMessage != "The audio could not be transcribed." {
		t.Errorf("Expected error.message to be %q, got %v", "The audio could not be transcribed.", errorObj["message"])
	}
}

func TestConversationItemTruncatedMessage(t *testing.T) {
	// Example conversation.item.truncated message from the API
	jsonData := []byte(`{
		"event_id": "event_2526",
		"type": "conversation.item.truncated",
		"item_id": "msg_004",
		"content_index": 0,
		"audio_end_ms": 1500
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal conversation.item.truncated message: %v", err)
	}

	// Verify it's a conversation.item.truncated message
	if msg.RcvdMsgType() != RcvdMsgTypeConversationItemTruncated {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeConversationItemTruncated, msg.RcvdMsgType())
	}

	// Cast to ConversationItemTruncatedMessage
	truncatedMsg, ok := msg.(*ConversationItemTruncatedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to ConversationItemTruncatedMessage")
	}

	// Verify the fields
	if truncatedMsg.EventID != "event_2526" {
		t.Errorf("Expected EventID to be %q, got %q", "event_2526", truncatedMsg.EventID)
	}

	if truncatedMsg.ItemID != "msg_004" {
		t.Errorf("Expected ItemID to be %q, got %q", "msg_004", truncatedMsg.ItemID)
	}

	if truncatedMsg.ContentIndex != 0 {
		t.Errorf("Expected ContentIndex to be %d, got %d", 0, truncatedMsg.ContentIndex)
	}

	if truncatedMsg.AudioEndMs != 1500 {
		t.Errorf("Expected AudioEndMs to be %d, got %d", 1500, truncatedMsg.AudioEndMs)
	}

	// Test marshaling back to JSON
	marshaled, err := json.Marshal(truncatedMsg)
	if err != nil {
		t.Fatalf("Failed to marshal conversation.item.truncated message: %v", err)
	}

	// Unmarshal the marshaled data to verify it's valid
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify fields in the marshaled JSON
	itemID, ok := unmarshaled["item_id"].(string)
	if !ok || itemID != "msg_004" {
		t.Errorf("Expected item_id to be %q, got %v", "msg_004", unmarshaled["item_id"])
	}

	contentIndex, ok := unmarshaled["content_index"].(float64)
	if !ok || int(contentIndex) != 0 {
		t.Errorf("Expected content_index to be %d, got %v", 0, unmarshaled["content_index"])
	}

	audioEndMs, ok := unmarshaled["audio_end_ms"].(float64)
	if !ok || int(audioEndMs) != 1500 {
		t.Errorf("Expected audio_end_ms to be %d, got %v", 1500, unmarshaled["audio_end_ms"])
	}
}

func TestConversationItemDeletedMessage(t *testing.T) {
	// Example conversation.item.deleted message from the API
	jsonData := []byte(`{
		"event_id": "event_2728",
		"type": "conversation.item.deleted",
		"item_id": "msg_005"
	}`)

	// Unmarshal the message
	msg, err := UnmarshalRcvdMsg(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal conversation.item.deleted message: %v", err)
	}

	// Verify it's a conversation.item.deleted message
	if msg.RcvdMsgType() != RcvdMsgTypeConversationItemDeleted {
		t.Fatalf("Expected message type to be %q, got %q", RcvdMsgTypeConversationItemDeleted, msg.RcvdMsgType())
	}

	// Cast to ConversationItemDeletedMessage
	deletedMsg, ok := msg.(*ConversationItemDeletedMessage)
	if !ok {
		t.Fatalf("Failed to cast message to ConversationItemDeletedMessage")
	}

	// Verify the fields
	if deletedMsg.EventID != "event_2728" {
		t.Errorf("Expected EventID to be %q, got %q", "event_2728", deletedMsg.EventID)
	}

	if deletedMsg.ItemID != "msg_005" {
		t.Errorf("Expected ItemID to be %q, got %q", "msg_005", deletedMsg.ItemID)
	}

	// Test marshaling back to JSON
	marshaled, err := json.Marshal(deletedMsg)
	if err != nil {
		t.Fatalf("Failed to marshal conversation.item.deleted message: %v", err)
	}

	// Unmarshal the marshaled data to verify it's valid
	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(marshaled, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	// Verify the item_id is included in the marshaled data
	itemID, ok := unmarshaled["item_id"].(string)
	if !ok || itemID != "msg_005" {
		t.Errorf("Expected item_id to be %q, got %v", "msg_005", unmarshaled["item_id"])
	}
}
