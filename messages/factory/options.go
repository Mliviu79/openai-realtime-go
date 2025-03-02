package factory

import (
	"github.com/Mliviu79/openai-realtime-go/messages/types"
	"github.com/Mliviu79/openai-realtime-go/session"
)

// WithEndTurn sets whether the message ends the turn
func WithEndTurn(endTurn bool) types.MessageOption {
	return func(m *types.Message) {
		if m != nil {
			m.EndTurn = endTurn
		}
	}
}

// WithMetadata sets the metadata for the message
func WithMetadata(metadata any) types.MessageOption {
	return func(m *types.Message) {
		if m != nil {
			m.Metadata = metadata
		}
	}
}

// WithRole sets the role for the message
func WithRole(role types.MessageRole) types.MessageOption {
	return func(m *types.Message) {
		if m != nil {
			m.Role = role
		}
	}
}

// WithContent sets the content for the message
func WithContent(content []types.MessageContentPart) types.MessageOption {
	return func(m *types.Message) {
		if m != nil {
			m.Content = content
		}
	}
}

// WithName sets the name for the message
func WithName(name string) types.MessageOption {
	return func(m *types.Message) {
		if m != nil {
			m.Name = name
		}
	}
}

// WithTools sets the tools for the message
func WithTools(tools []session.Tool) types.MessageOption {
	return func(m *types.Message) {
		if m != nil {
			m.Tools = tools
		}
	}
}

// NewMessage creates a new message with the given options
func NewMessage(options ...types.MessageOption) *types.Message {
	msg := &types.Message{}
	for _, option := range options {
		option(msg)
	}
	return msg
}
