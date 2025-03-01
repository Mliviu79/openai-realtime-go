package outgoing

import (
	"github.com/Mliviu79/go-openai-realtime/messages/types"
)

// ConversationCreateMessage is used to create a new conversation item
type ConversationCreateMessage struct {
	OutMsgBase
	// PreviousItemID references the item that comes before this one, if any
	PreviousItemID string `json:"previous_item_id,omitempty"`
	// Item contains the details of the conversation item to create
	Item types.MessageItem `json:"item"`
}

// NewConversationCreateMessage creates a new conversation create message
func NewConversationCreateMessage(previousItemID string, item types.MessageItem) ConversationCreateMessage {
	return ConversationCreateMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeConversationCreate,
		},
		PreviousItemID: previousItemID,
		Item:           item,
	}
}

// ConversationTruncateMessage is used to truncate a conversation item
type ConversationTruncateMessage struct {
	OutMsgBase
	// ItemID identifies the conversation item to truncate
	ItemID string `json:"item_id"`
	// ContentIndex specifies which content part within the item to truncate
	ContentIndex int `json:"content_index"`
	// AudioEndMs indicates the new end time of the audio in milliseconds
	AudioEndMs int `json:"audio_end_ms"`
}

// NewConversationTruncateMessage creates a new conversation truncate message
func NewConversationTruncateMessage(itemID string, contentIndex int, audioEndMs int) ConversationTruncateMessage {
	return ConversationTruncateMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeConversationTruncate,
		},
		ItemID:       itemID,
		ContentIndex: contentIndex,
		AudioEndMs:   audioEndMs,
	}
}

// ConversationDeleteMessage is used to delete a conversation item
type ConversationDeleteMessage struct {
	OutMsgBase
	// ItemID identifies the conversation item to delete
	ItemID string `json:"item_id"`
}

// NewConversationDeleteMessage creates a new conversation delete message
func NewConversationDeleteMessage(itemID string) ConversationDeleteMessage {
	return ConversationDeleteMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeConversationDelete,
		},
		ItemID: itemID,
	}
}
