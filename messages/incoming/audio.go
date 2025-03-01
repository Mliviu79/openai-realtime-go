package incoming

// AudioBufferCommittedMessage is sent when an audio buffer is committed
type AudioBufferCommittedMessage struct {
	RcvdMsgBase
	// PreviousItemID references the item that comes before this one, if any
	PreviousItemID string `json:"previous_item_id,omitempty"`
	// ItemID identifies the item this audio buffer belongs to
	ItemID string `json:"item_id"`
}

// AudioBufferClearedMessage is sent when an audio buffer is cleared
type AudioBufferClearedMessage struct {
	RcvdMsgBase
}

// AudioBufferSpeechStartedMessage is sent when speech is detected in the audio buffer
type AudioBufferSpeechStartedMessage struct {
	RcvdMsgBase
	// AudioStartMs indicates when speech was first detected, in milliseconds
	AudioStartMs int64 `json:"audio_start_ms"`
	// ItemID identifies the item this audio buffer belongs to
	ItemID string `json:"item_id"`
}

// AudioBufferSpeechStoppedMessage is sent when speech stops in the audio buffer
type AudioBufferSpeechStoppedMessage struct {
	RcvdMsgBase
	// AudioEndMs indicates when speech ended, in milliseconds
	AudioEndMs int64 `json:"audio_end_ms"`
	// ItemID identifies the item this audio buffer belongs to
	ItemID string `json:"item_id"`
}
