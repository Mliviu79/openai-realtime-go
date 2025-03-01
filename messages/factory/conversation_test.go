package factory

import (
	"reflect"
	"testing"

	"github.com/Mliviu79/go-openai-realtime/messages/types"
)

func TestMessageItem(t *testing.T) {
	// Create a message item with user role and text content
	role := types.MessageRoleUser
	content := []types.MessageContentPart{
		{
			Type: types.MessageContentTypeText,
			Text: "Hello, world!",
		},
	}

	item := MessageItem(role, content)

	// Check type
	if item.Type != types.MessageItemTypeMessage {
		t.Errorf("MessageItem().Type = %v, want %v", item.Type, types.MessageItemTypeMessage)
	}

	// Check role
	if item.Role != role {
		t.Errorf("MessageItem().Role = %v, want %v", item.Role, role)
	}

	// Check content
	if !reflect.DeepEqual(item.Content, content) {
		t.Errorf("MessageItem().Content = %v, want %v", item.Content, content)
	}

	// Other fields should be empty
	if item.CallID != "" {
		t.Errorf("MessageItem().CallID = %v, want empty string", item.CallID)
	}

	if item.Name != "" {
		t.Errorf("MessageItem().Name = %v, want empty string", item.Name)
	}

	if item.Arguments != "" {
		t.Errorf("MessageItem().Arguments = %v, want empty string", item.Arguments)
	}

	if item.Output != "" {
		t.Errorf("MessageItem().Output = %v, want empty string", item.Output)
	}
}

func TestFunctionCallItem(t *testing.T) {
	name := "calculate"
	arguments := `{"x": 1, "y": 2}`

	item := FunctionCallItem(name, arguments)

	// Check type
	if item.Type != types.MessageItemTypeFunctionCall {
		t.Errorf("FunctionCallItem().Type = %v, want %v", item.Type, types.MessageItemTypeFunctionCall)
	}

	// Check name
	if item.Name != name {
		t.Errorf("FunctionCallItem().Name = %v, want %v", item.Name, name)
	}

	// Check arguments
	if item.Arguments != arguments {
		t.Errorf("FunctionCallItem().Arguments = %v, want %v", item.Arguments, arguments)
	}

	// Other fields should be empty or nil
	if item.Role != "" {
		t.Errorf("FunctionCallItem().Role = %v, want empty string", item.Role)
	}

	if item.Content != nil {
		t.Errorf("FunctionCallItem().Content = %v, want nil", item.Content)
	}

	if item.CallID != "" {
		t.Errorf("FunctionCallItem().CallID = %v, want empty string", item.CallID)
	}

	if item.Output != "" {
		t.Errorf("FunctionCallItem().Output = %v, want empty string", item.Output)
	}
}

func TestFunctionResponseItem(t *testing.T) {
	callID := "call-123456"
	output := `{"result": 3}`

	item := FunctionResponseItem(callID, output)

	// Check type
	if item.Type != types.MessageItemTypeFunctionCallOutput {
		t.Errorf("FunctionResponseItem().Type = %v, want %v", item.Type, types.MessageItemTypeFunctionCallOutput)
	}

	// Check callID
	if item.CallID != callID {
		t.Errorf("FunctionResponseItem().CallID = %v, want %v", item.CallID, callID)
	}

	// Check output
	if item.Output != output {
		t.Errorf("FunctionResponseItem().Output = %v, want %v", item.Output, output)
	}

	// Other fields should be empty or nil
	if item.Role != "" {
		t.Errorf("FunctionResponseItem().Role = %v, want empty string", item.Role)
	}

	if item.Content != nil {
		t.Errorf("FunctionResponseItem().Content = %v, want nil", item.Content)
	}

	if item.Name != "" {
		t.Errorf("FunctionResponseItem().Name = %v, want empty string", item.Name)
	}

	if item.Arguments != "" {
		t.Errorf("FunctionResponseItem().Arguments = %v, want empty string", item.Arguments)
	}
}

func TestSystemMessage(t *testing.T) {
	text := "You are a helpful assistant."

	item := SystemMessage(text)

	// Check type
	if item.Type != types.MessageItemTypeMessage {
		t.Errorf("SystemMessage().Type = %v, want %v", item.Type, types.MessageItemTypeMessage)
	}

	// Check role
	if item.Role != types.MessageRoleSystem {
		t.Errorf("SystemMessage().Role = %v, want %v", item.Role, types.MessageRoleSystem)
	}

	// Check content
	if len(item.Content) != 1 {
		t.Errorf("SystemMessage().Content has length %v, want 1", len(item.Content))
	} else {
		content := item.Content[0]
		if content.Type != types.MessageContentTypeText {
			t.Errorf("SystemMessage().Content[0].Type = %v, want %v", content.Type, types.MessageContentTypeText)
		}

		if content.Text != text {
			t.Errorf("SystemMessage().Content[0].Text = %v, want %v", content.Text, text)
		}
	}
}

func TestUserMessage(t *testing.T) {
	content := []types.MessageContentPart{
		{
			Type: types.MessageContentTypeInputText,
			Text: "Hello, assistant!",
		},
	}

	item := UserMessage(content)

	// Check type
	if item.Type != types.MessageItemTypeMessage {
		t.Errorf("UserMessage().Type = %v, want %v", item.Type, types.MessageItemTypeMessage)
	}

	// Check role
	if item.Role != types.MessageRoleUser {
		t.Errorf("UserMessage().Role = %v, want %v", item.Role, types.MessageRoleUser)
	}

	// Check content
	if !reflect.DeepEqual(item.Content, content) {
		t.Errorf("UserMessage().Content = %v, want %v", item.Content, content)
	}
}

func TestUserTextMessage(t *testing.T) {
	text := "Hello, assistant!"

	item := UserTextMessage(text)

	// Check type
	if item.Type != types.MessageItemTypeMessage {
		t.Errorf("UserTextMessage().Type = %v, want %v", item.Type, types.MessageItemTypeMessage)
	}

	// Check role
	if item.Role != types.MessageRoleUser {
		t.Errorf("UserTextMessage().Role = %v, want %v", item.Role, types.MessageRoleUser)
	}

	// Check content
	if len(item.Content) != 1 {
		t.Errorf("UserTextMessage().Content has length %v, want 1", len(item.Content))
	} else {
		content := item.Content[0]
		if content.Type != types.MessageContentTypeInputText {
			t.Errorf("UserTextMessage().Content[0].Type = %v, want %v",
				content.Type, types.MessageContentTypeInputText)
		}

		if content.Text != text {
			t.Errorf("UserTextMessage().Content[0].Text = %v, want %v", content.Text, text)
		}
	}
}

func TestUserAudioMessage(t *testing.T) {
	audio := "base64-encoded-audio"
	transcript := "Hello, assistant!"

	item := UserAudioMessage(audio, transcript)

	// Check type
	if item.Type != types.MessageItemTypeMessage {
		t.Errorf("UserAudioMessage().Type = %v, want %v", item.Type, types.MessageItemTypeMessage)
	}

	// Check role
	if item.Role != types.MessageRoleUser {
		t.Errorf("UserAudioMessage().Role = %v, want %v", item.Role, types.MessageRoleUser)
	}

	// Check content
	if len(item.Content) != 1 {
		t.Errorf("UserAudioMessage().Content has length %v, want 1", len(item.Content))
	} else {
		content := item.Content[0]
		if content.Type != types.MessageContentTypeInputAudio {
			t.Errorf("UserAudioMessage().Content[0].Type = %v, want %v",
				content.Type, types.MessageContentTypeInputAudio)
		}

		if content.Audio != audio {
			t.Errorf("UserAudioMessage().Content[0].Audio = %v, want %v", content.Audio, audio)
		}

		if content.Transcript != transcript {
			t.Errorf("UserAudioMessage().Content[0].Transcript = %v, want %v",
				content.Transcript, transcript)
		}
	}
}

func TestAssistantMessage(t *testing.T) {
	content := []types.MessageContentPart{
		{
			Type: types.MessageContentTypeText,
			Text: "I'm here to help!",
		},
	}

	item := AssistantMessage(content)

	// Check type
	if item.Type != types.MessageItemTypeMessage {
		t.Errorf("AssistantMessage().Type = %v, want %v", item.Type, types.MessageItemTypeMessage)
	}

	// Check role
	if item.Role != types.MessageRoleAssistant {
		t.Errorf("AssistantMessage().Role = %v, want %v", item.Role, types.MessageRoleAssistant)
	}

	// Check content
	if !reflect.DeepEqual(item.Content, content) {
		t.Errorf("AssistantMessage().Content = %v, want %v", item.Content, content)
	}
}

func TestAssistantTextMessage(t *testing.T) {
	text := "I'm here to help!"

	item := AssistantTextMessage(text)

	// Check type
	if item.Type != types.MessageItemTypeMessage {
		t.Errorf("AssistantTextMessage().Type = %v, want %v", item.Type, types.MessageItemTypeMessage)
	}

	// Check role
	if item.Role != types.MessageRoleAssistant {
		t.Errorf("AssistantTextMessage().Role = %v, want %v", item.Role, types.MessageRoleAssistant)
	}

	// Check content
	if len(item.Content) != 1 {
		t.Errorf("AssistantTextMessage().Content has length %v, want 1", len(item.Content))
	} else {
		content := item.Content[0]
		if content.Type != types.MessageContentTypeText {
			t.Errorf("AssistantTextMessage().Content[0].Type = %v, want %v",
				content.Type, types.MessageContentTypeText)
		}

		if content.Text != text {
			t.Errorf("AssistantTextMessage().Content[0].Text = %v, want %v", content.Text, text)
		}
	}
}
