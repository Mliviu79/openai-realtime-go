package incoming

import (
	"github.com/Mliviu79/openai-realtime-go/messages/types"
)

// ResponseCreatedMessage is sent when a new response is created
type ResponseCreatedMessage struct {
	RcvdMsgBase
	// Response contains the details of the newly created response
	Response types.Response `json:"response"`
}

// ResponseDoneMessage is sent when a response is completed
type ResponseDoneMessage struct {
	RcvdMsgBase
	// Response contains the final state of the completed response
	Response types.Response `json:"response"`
}

// ResponseContentPartAddedMessage is sent when a content part is added to a response
type ResponseContentPartAddedMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this content belongs to
	ResponseID string `json:"response_id"`
	// ItemID identifies which item within the response this content belongs to
	ItemID string `json:"item_id"`
	// OutputIndex specifies which output within the item this content belongs to
	OutputIndex int `json:"output_index"`
	// ContentIndex specifies the position of this content within the output
	ContentIndex int `json:"content_index"`
	// Part contains the actual content that was added
	Part types.MessageContentPart `json:"part"`
}

// ResponseContentPartDoneMessage is sent when a content part is completed
type ResponseContentPartDoneMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this content belongs to
	ResponseID string `json:"response_id"`
	// ItemID identifies which item within the response this content belongs to
	ItemID string `json:"item_id"`
	// OutputIndex specifies which output within the item this content belongs to
	OutputIndex int `json:"output_index"`
	// ContentIndex specifies the position of this content within the output
	ContentIndex int `json:"content_index"`
	// Part contains the final state of the completed content part
	Part types.MessageContentPart `json:"part"`
}

// ResponseTextDeltaMessage is sent when new text is added to a response
type ResponseTextDeltaMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this text belongs to
	ResponseID string `json:"response_id"`
	// ItemID identifies which item within the response this text belongs to
	ItemID string `json:"item_id"`
	// OutputIndex specifies which output within the item this text belongs to
	OutputIndex int `json:"output_index"`
	// ContentIndex specifies the position of this text within the output
	ContentIndex int `json:"content_index"`
	// Delta contains the new text fragment to append
	Delta string `json:"delta"`
}

// ResponseTextDoneMessage is sent when text generation is completed
type ResponseTextDoneMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this text belongs to
	ResponseID string `json:"response_id"`
	// ItemID identifies which item within the response this text belongs to
	ItemID string `json:"item_id"`
	// OutputIndex specifies which output within the item this text belongs to
	OutputIndex int `json:"output_index"`
	// ContentIndex specifies the position of this text within the output
	ContentIndex int `json:"content_index"`
	// Text contains the final text content
	Text string `json:"text"`
}

// ResponseOutputItemAddedMessage is sent when an output item is added to a response
type ResponseOutputItemAddedMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this output belongs to
	ResponseID string `json:"response_id"`
	// OutputIndex specifies the position of this output within the response
	OutputIndex int `json:"output_index"`
	// Item contains the details of the newly added output item
	Item types.OutputItem `json:"item"`
}

// ResponseOutputItemDoneMessage is sent when an output item is completed
type ResponseOutputItemDoneMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this output belongs to
	ResponseID string `json:"response_id"`
	// OutputIndex specifies the position of this output within the response
	OutputIndex int `json:"output_index"`
	// Item contains the final state of the completed output item
	Item types.OutputItem `json:"item"`
}

// ResponseAudioTranscriptDeltaMessage is sent when new transcript text is added
type ResponseAudioTranscriptDeltaMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this transcript belongs to
	ResponseID string `json:"response_id"`
	// ItemID identifies which item within the response this transcript belongs to
	ItemID string `json:"item_id"`
	// OutputIndex specifies which output within the item this transcript belongs to
	OutputIndex int `json:"output_index"`
	// ContentIndex specifies the position of this transcript within the output
	ContentIndex int `json:"content_index"`
	// Delta contains the new transcript fragment to append
	Delta string `json:"delta"`
}

// ResponseAudioTranscriptDoneMessage is sent when transcript generation is completed
type ResponseAudioTranscriptDoneMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this transcript belongs to
	ResponseID string `json:"response_id"`
	// ItemID identifies which item within the response this transcript belongs to
	ItemID string `json:"item_id"`
	// OutputIndex specifies which output within the item this transcript belongs to
	OutputIndex int `json:"output_index"`
	// ContentIndex specifies the position of this transcript within the output
	ContentIndex int `json:"content_index"`
	// Transcript contains the complete transcript text
	Transcript string `json:"transcript"`
}

// ResponseAudioDeltaMessage is sent when new audio data is added
type ResponseAudioDeltaMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this audio belongs to
	ResponseID string `json:"response_id"`
	// ItemID identifies which item within the response this audio belongs to
	ItemID string `json:"item_id"`
	// OutputIndex specifies which output within the item this audio belongs to
	OutputIndex int `json:"output_index"`
	// ContentIndex specifies the position of this audio within the output
	ContentIndex int `json:"content_index"`
	// Delta contains the new audio data fragment as base64-encoded string
	Delta string `json:"delta"`
}

// ResponseAudioDoneMessage is sent when audio generation is completed
type ResponseAudioDoneMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this audio belongs to
	ResponseID string `json:"response_id"`
	// ItemID identifies which item within the response this audio belongs to
	ItemID string `json:"item_id"`
	// OutputIndex specifies which output within the item this audio belongs to
	OutputIndex int `json:"output_index"`
	// ContentIndex specifies the position of this audio within the output
	ContentIndex int `json:"content_index"`
}

// ResponseFunctionCallArgumentsDeltaMessage is sent when new function call arguments are added
type ResponseFunctionCallArgumentsDeltaMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this function call belongs to
	ResponseID string `json:"response_id"`
	// ItemID identifies which item within the response this function call belongs to
	ItemID string `json:"item_id"`
	// OutputIndex specifies which output within the item this function call belongs to
	OutputIndex int `json:"output_index"`
	// CallID uniquely identifies this function call
	CallID string `json:"call_id"`
	// Delta contains the new arguments fragment as a JSON string
	Delta string `json:"delta"`
}

// ResponseFunctionCallArgumentsDoneMessage is sent when function call arguments are completed
type ResponseFunctionCallArgumentsDoneMessage struct {
	RcvdMsgBase
	// ResponseID identifies which response this function call belongs to
	ResponseID string `json:"response_id"`
	// ItemID identifies which item within the response this function call belongs to
	ItemID string `json:"item_id"`
	// OutputIndex specifies which output within the item this function call belongs to
	OutputIndex int `json:"output_index"`
	// CallID uniquely identifies this function call
	CallID string `json:"call_id"`
	// Arguments contains the complete function arguments as a JSON string
	Arguments string `json:"arguments"`
}

// RateLimitsUpdatedMessage is sent when rate limits are updated
type RateLimitsUpdatedMessage struct {
	RcvdMsgBase
	// RateLimits contains the updated rate limit information
	RateLimits []types.RateLimit `json:"rate_limits"`
}
