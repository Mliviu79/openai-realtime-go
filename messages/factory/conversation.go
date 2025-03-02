package factory

import (
	"github.com/Mliviu79/openai-realtime-go/messages/types"
)

// MessageItem creates a new message item
func MessageItem(role types.MessageRole, content []types.MessageContentPart) types.MessageItem {
	return types.MessageItem{
		Type:    types.MessageItemTypeMessage,
		Role:    role,
		Content: content,
	}
}

// FunctionCallItem creates a new function call item
func FunctionCallItem(name string, arguments string) types.MessageItem {
	return types.MessageItem{
		Type:      types.MessageItemTypeFunctionCall,
		Name:      name,
		Arguments: arguments,
	}
}

// FunctionResponseItem creates a new function response item
func FunctionResponseItem(callID string, output string) types.MessageItem {
	return types.MessageItem{
		Type:   types.MessageItemTypeFunctionCallOutput,
		CallID: callID,
		Output: output,
	}
}

// SystemMessage creates a new system message item
func SystemMessage(text string) types.MessageItem {
	return MessageItem(
		types.MessageRoleSystem,
		[]types.MessageContentPart{
			TextContent(text),
		},
	)
}

// UserMessage creates a new user message item
func UserMessage(content []types.MessageContentPart) types.MessageItem {
	return MessageItem(
		types.MessageRoleUser,
		content,
	)
}

// UserTextMessage creates a new user text message item
func UserTextMessage(text string) types.MessageItem {
	return UserMessage(
		[]types.MessageContentPart{
			InputTextContent(text),
		},
	)
}

// UserAudioMessage creates a new user audio message item
func UserAudioMessage(audio string, transcript string) types.MessageItem {
	return UserMessage(
		[]types.MessageContentPart{
			InputAudioContent(audio, transcript),
		},
	)
}

// AssistantMessage creates a new assistant message item
func AssistantMessage(content []types.MessageContentPart) types.MessageItem {
	return MessageItem(
		types.MessageRoleAssistant,
		content,
	)
}

// AssistantTextMessage creates a new assistant text message item
func AssistantTextMessage(text string) types.MessageItem {
	return AssistantMessage(
		[]types.MessageContentPart{
			TextContent(text),
		},
	)
}
