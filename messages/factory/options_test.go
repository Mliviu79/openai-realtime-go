package factory

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Mliviu79/go-openai-realtime/messages/types"
	"github.com/Mliviu79/go-openai-realtime/session"
)

func TestWithEndTurn(t *testing.T) {
	tests := []struct {
		name    string
		endTurn bool
		want    bool
	}{
		{
			name:    "set end turn true",
			endTurn: true,
			want:    true,
		},
		{
			name:    "set end turn false",
			endTurn: false,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &types.Message{}
			opt := WithEndTurn(tt.endTurn)
			opt(msg)

			if msg.EndTurn != tt.want {
				t.Errorf("WithEndTurn(%v) = %v, want %v", tt.endTurn, msg.EndTurn, tt.want)
			}
		})
	}
}

func TestWithMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata interface{}
	}{
		{
			name:     "set string metadata",
			metadata: "test-metadata",
		},
		{
			name:     "set map metadata",
			metadata: map[string]string{"key": "value"},
		},
		{
			name:     "set struct metadata",
			metadata: struct{ Value string }{"test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &types.Message{}
			opt := WithMetadata(tt.metadata)
			opt(msg)

			if !reflect.DeepEqual(msg.Metadata, tt.metadata) {
				t.Errorf("WithMetadata() = %v, want %v", msg.Metadata, tt.metadata)
			}
		})
	}
}

func TestWithRole(t *testing.T) {
	tests := []struct {
		name string
		role types.MessageRole
		want types.MessageRole
	}{
		{
			name: "set role to system",
			role: types.MessageRoleSystem,
			want: types.MessageRoleSystem,
		},
		{
			name: "set role to user",
			role: types.MessageRoleUser,
			want: types.MessageRoleUser,
		},
		{
			name: "set role to assistant",
			role: types.MessageRoleAssistant,
			want: types.MessageRoleAssistant,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &types.Message{}
			opt := WithRole(tt.role)
			opt(msg)

			if msg.Role != tt.want {
				t.Errorf("WithRole(%v) = %v, want %v", tt.role, msg.Role, tt.want)
			}
		})
	}
}

func TestWithContent(t *testing.T) {
	content := []types.MessageContentPart{
		{
			Type: types.MessageContentTypeText,
			Text: "Hello, world!",
		},
	}

	msg := &types.Message{}
	opt := WithContent(content)
	opt(msg)

	if !reflect.DeepEqual(msg.Content, content) {
		t.Errorf("WithContent() = %v, want %v", msg.Content, content)
	}
}

func TestWithName(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		want     string
	}{
		{
			name:     "set name",
			funcName: "testFunction",
			want:     "testFunction",
		},
		{
			name:     "set empty name",
			funcName: "",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &types.Message{}
			opt := WithName(tt.funcName)
			opt(msg)

			if msg.Name != tt.want {
				t.Errorf("WithName(%v) = %v, want %v", tt.funcName, msg.Name, tt.want)
			}
		})
	}
}

func TestWithTools(t *testing.T) {
	// Create parameters as a JSON object
	parameters, _ := json.Marshal(map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"test": map[string]interface{}{
				"type": "string",
			},
		},
		"required": []string{"test"},
	})

	tools := []session.Tool{
		{
			Type:        "function",
			Name:        "test_function",
			Description: "A test function",
			Parameters:  json.RawMessage(parameters),
		},
	}

	msg := &types.Message{}
	opt := WithTools(tools)
	opt(msg)

	if !reflect.DeepEqual(msg.Tools, tools) {
		t.Errorf("WithTools() = %v, want %v", msg.Tools, tools)
	}
}

func TestNewMessage(t *testing.T) {
	t.Run("create message with options", func(t *testing.T) {
		role := types.MessageRoleUser
		name := "testFunction"
		content := []types.MessageContentPart{
			{
				Type: types.MessageContentTypeText,
				Text: "Hello, world!",
			},
		}
		endTurn := true
		metadata := map[string]string{"key": "value"}

		msg := NewMessage(
			WithRole(role),
			WithName(name),
			WithContent(content),
			WithEndTurn(endTurn),
			WithMetadata(metadata),
		)

		if msg.Role != role {
			t.Errorf("NewMessage().Role = %v, want %v", msg.Role, role)
		}

		if msg.Name != name {
			t.Errorf("NewMessage().Name = %v, want %v", msg.Name, name)
		}

		if !reflect.DeepEqual(msg.Content, content) {
			t.Errorf("NewMessage().Content = %v, want %v", msg.Content, content)
		}

		if msg.EndTurn != endTurn {
			t.Errorf("NewMessage().EndTurn = %v, want %v", msg.EndTurn, endTurn)
		}

		if !reflect.DeepEqual(msg.Metadata, metadata) {
			t.Errorf("NewMessage().Metadata = %v, want %v", msg.Metadata, metadata)
		}
	})

	t.Run("create message with no options", func(t *testing.T) {
		msg := NewMessage()

		// First check if msg is nil before dereferencing
		if msg == nil {
			t.Error("NewMessage() returned nil")
			return // Return early to avoid nil dereference
		}

		// Empty message should have zero values for all fields
		if msg.Role != "" {
			t.Errorf("NewMessage().Role = %v, want %v", msg.Role, "")
		}

		if msg.Name != "" {
			t.Errorf("NewMessage().Name = %v, want %v", msg.Name, "")
		}

		if msg.Content != nil {
			t.Errorf("NewMessage().Content = %v, want nil", msg.Content)
		}

		if msg.EndTurn != false {
			t.Errorf("NewMessage().EndTurn = %v, want false", msg.EndTurn)
		}

		if msg.Metadata != nil {
			t.Errorf("NewMessage().Metadata = %v, want nil", msg.Metadata)
		}
	})
}

func TestOptionWithNilMessage(t *testing.T) {
	// Test that options don't panic when given a nil message
	var msg *types.Message = nil

	// Each option should safely handle a nil message
	WithEndTurn(true)(msg)
	WithMetadata("test")(msg)
	WithRole(types.MessageRoleSystem)(msg)
	WithContent([]types.MessageContentPart{})(msg)
	WithName("test")(msg)
	WithTools([]session.Tool{})(msg)
}
