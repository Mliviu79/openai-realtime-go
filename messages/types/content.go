package types

// MessageContentPart represents a single part of a message's content
type MessageContentPart struct {
	// Type specifies what kind of content this is
	Type MessageContentType `json:"type"`

	// Text contains the text content
	// Used for text, input_text content types
	Text string `json:"text,omitempty"`

	// ID references a previous conversation item
	// Used for item_reference content types
	ID string `json:"id,omitempty"`

	// Audio contains Base64-encoded audio data
	// Used for audio and input_audio types
	Audio string `json:"audio,omitempty"`

	// Transcript contains the text transcription of audio
	// Used for transcript and input_audio content types
	Transcript string `json:"transcript,omitempty"`
}

// TokenDetails contains information about token usage
type TokenDetails struct {
	// TextTokens is the number of tokens used for text content
	TextTokens int `json:"text_tokens"`

	// AudioTokens is the number of tokens used for audio content
	AudioTokens int `json:"audio_tokens"`
}

// InputTokenDetails contains information about input token usage
type InputTokenDetails struct {
	// CachedTokens is the number of tokens that were cached
	CachedTokens int `json:"cached_tokens"`

	// TextTokens is the number of tokens used for text content
	TextTokens int `json:"text_tokens"`

	// AudioTokens is the number of tokens used for audio content
	AudioTokens int `json:"audio_tokens"`

	// CachedTokensDetails contains detailed information about cached token usage
	CachedTokensDetails TokenDetails `json:"cached_tokens_details,omitempty"`
}

// OutputTokenDetails is an alias for TokenDetails
type OutputTokenDetails = TokenDetails

// CachedTokensDetails is an alias for TokenDetails
type CachedTokensDetails = TokenDetails

// Usage contains information about token usage
type Usage struct {
	// TotalTokens is the total number of tokens used
	TotalTokens int `json:"total_tokens"`

	// InputTokens is the number of tokens in the input
	InputTokens int `json:"input_tokens"`

	// OutputTokens is the number of tokens in the output
	OutputTokens int `json:"output_tokens"`

	// InputTokenDetails contains detailed information about input token usage
	InputTokenDetails InputTokenDetails `json:"input_token_details,omitempty"`

	// OutputTokenDetails contains detailed information about output token usage
	OutputTokenDetails OutputTokenDetails `json:"output_token_details,omitempty"`
}
