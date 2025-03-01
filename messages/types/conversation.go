package types

//-----------------------------------------------------------------------------
// Message Item Types
//-----------------------------------------------------------------------------

// MessageItemType represents the type of a message item
type MessageItemType string

const (
	// MessageItemTypeMessage represents a standard message
	MessageItemTypeMessage MessageItemType = "message"

	// MessageItemTypeFunctionCall represents a function call
	MessageItemTypeFunctionCall MessageItemType = "function_call"

	// MessageItemTypeFunctionCallOutput represents a function call output
	MessageItemTypeFunctionCallOutput MessageItemType = "function_call_output"

	// MessageItemTypeFunctionResponse is an alias for MessageItemTypeFunctionCallOutput
	// Deprecated: Use MessageItemTypeFunctionCallOutput instead
	MessageItemTypeFunctionResponse = MessageItemTypeFunctionCallOutput
)

//-----------------------------------------------------------------------------
// Item Status Types
//-----------------------------------------------------------------------------

// ItemStatus represents the status of an item
type ItemStatus string

const (
	// ItemStatusInProgress indicates the item is still being processed
	// Deprecated: Used only in tests, not in the API
	ItemStatusInProgress ItemStatus = "in_progress"

	// ItemStatusCompleted indicates the item has been fully processed
	ItemStatusCompleted ItemStatus = "completed"

	// ItemStatusIncomplete indicates the item is incomplete
	ItemStatusIncomplete ItemStatus = "incomplete"
)

// MessageItem represents an item in a message
type MessageItem struct {
	// ID is an optional identifier for this item
	ID string `json:"id,omitempty"`

	// Object is always "realtime.item" when present
	Object string `json:"object,omitempty"`

	// Type specifies what kind of item this is
	Type MessageItemType `json:"type"`

	// Status indicates the current status of the item
	Status ItemStatus `json:"status,omitempty"`

	// Role specifies who created this message
	Role MessageRole `json:"role,omitempty"`

	// Content contains the actual content of the message
	Content []MessageContentPart `json:"content,omitempty"`

	// CallID identifies the function call this item relates to
	CallID string `json:"call_id,omitempty"`

	// Name specifies the function being called
	Name string `json:"name,omitempty"`

	// Arguments contains the function call arguments as a JSON string
	Arguments string `json:"arguments,omitempty"`

	// Output contains the result of the function call
	Output string `json:"output,omitempty"`
}

// NewMessageItem creates a new MessageItem with default values
func NewMessageItem() *MessageItem {
	return &MessageItem{
		Object: "realtime.item",
	}
}

// ConversationItem represents an item in a conversation
type ConversationItem struct {
	// ID is an optional identifier for this item
	ID string `json:"id,omitempty"`

	// Type specifies what kind of item this is
	Type MessageItemType `json:"type"`

	// Status indicates the current status of the item
	Status ItemStatus `json:"status,omitempty"`

	// Role specifies who created this message
	Role *MessageRole `json:"role,omitempty"`

	// Content contains the actual content of the message
	Content []MessageContentPart `json:"content,omitempty"`

	// CallID identifies the function call this item relates to
	CallID string `json:"call_id,omitempty"`

	// Name specifies the function being called
	Name string `json:"name,omitempty"`

	// Arguments contains the function call arguments as a JSON string
	Arguments string `json:"arguments,omitempty"`

	// Output contains the result of the function call
	Output string `json:"output,omitempty"`
}

// Conversation represents a complete conversation
type Conversation struct {
	// ID is the unique identifier for this conversation
	ID string `json:"id"`

	// Object is always "realtime.conversation" when present in conversation.created
	Object string `json:"object,omitempty"`

	// Items contains the messages and other items in this conversation
	Items []MessageItem `json:"items,omitempty"`
}
