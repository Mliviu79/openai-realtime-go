package incoming

import (
	"github.com/Mliviu79/openai-realtime-go/messages/types"
)

// ConversationCreatedMessage is sent when a new conversation is created
type ConversationCreatedMessage struct {
	RcvdMsgBase
	// Conversation contains the details of the newly created conversation
	Conversation types.Conversation `json:"conversation"`
}

// ConversationItemCreatedMessage is sent when a new item is added to a conversation
type ConversationItemCreatedMessage struct {
	RcvdMsgBase
	// PreviousItemID references the item that comes before this one, if any
	PreviousItemID string `json:"previous_item_id,omitempty"`
	// Item contains the details of the newly created conversation item
	Item types.ResponseMessageItem `json:"item"`
}

// ConversationItemTranscriptionCompletedMessage is sent when audio transcription completes
type ConversationItemTranscriptionCompletedMessage struct {
	RcvdMsgBase
	// ItemID identifies the conversation item this transcription belongs to
	ItemID string `json:"item_id"`
	// ContentIndex specifies which content part within the item was transcribed
	ContentIndex int `json:"content_index"`
	// Transcript contains the text transcribed from audio
	Transcript string `json:"transcript"`
	// Logprobs contains the log probabilities of the transcription
	Logprobs []logprob `json:"logprobs,omitempty"`
}

type logprob struct {
	//The bytes that were used to generate the log probability.
	Bytes []byte `json:"bytes"`
	//The log probability of the token.
	Logprob float64 `json:"logprob"`
	//The token that was used to generate the log probability.
	Token string `json:"token"`
}

// ConversationItemTranscriptionFailedMessage is sent when audio transcription fails
type ConversationItemTranscriptionFailedMessage struct {
	RcvdMsgBase
	// ItemID identifies the conversation item where transcription failed
	ItemID string `json:"item_id"`
	// ContentIndex specifies which content part within the item failed transcription
	ContentIndex int `json:"content_index"`
	// Error contains details about why the transcription failed
	Error ErrorInfo `json:"error"`
}

// ConversationItemTruncatedMessage is sent when an item's audio is truncated
type ConversationItemTruncatedMessage struct {
	RcvdMsgBase
	// ItemID identifies the conversation item that was truncated
	ItemID string `json:"item_id"`
	// ContentIndex specifies which content part within the item was truncated
	ContentIndex int `json:"content_index"`
	// AudioEndMs indicates the new end time of the audio in milliseconds
	AudioEndMs int `json:"audio_end_ms"`
}

// ConversationItemDeletedMessage is sent when an item is deleted from a conversation
type ConversationItemDeletedMessage struct {
	RcvdMsgBase
	// ItemID identifies the conversation item that was deleted
	ItemID string `json:"item_id"`
}

// ConversationItemTranscriptionDeltaMessage is sent as an incremental update during real-time transcription
type ConversationItemTranscriptionDeltaMessage struct {
	RcvdMsgBase
	// ItemID identifies the conversation item this transcription belongs to
	ItemID string `json:"item_id"`
	// ContentIndex specifies which content part within the item is being transcribed
	ContentIndex int `json:"content_index"`
	// Delta contains the incremental text transcribed from audio
	Delta string `json:"delta"`
}
