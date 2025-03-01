package session

// ConfigOption is a function that configures a SessionRequest
type ConfigOption func(*SessionRequest)

// NewSessionRequest creates a new session request with the given options
func NewSessionRequest(opts ...ConfigOption) *SessionRequest {
	req := &SessionRequest{}
	for _, opt := range opts {
		opt(req)
	}
	return req
}

// WithModalities sets the modalities for the session
func WithModalities(modalities []Modality) ConfigOption {
	return func(c *SessionRequest) {
		c.Modalities = &modalities
	}
}

// WithModel sets the model for the session
func WithModel(model Model) ConfigOption {
	return func(c *SessionRequest) {
		c.Model = &model
	}
}

// WithInstructions sets the instructions for the session
func WithInstructions(instructions string) ConfigOption {
	return func(c *SessionRequest) {
		c.Instructions = &instructions
	}
}

// WithVoice sets the voice for the session
func WithVoice(voice Voice) ConfigOption {
	return func(c *SessionRequest) {
		c.Voice = &voice
	}
}

// WithInputAudioFormat sets the input audio format for the session
func WithInputAudioFormat(format AudioFormat) ConfigOption {
	return func(c *SessionRequest) {
		c.InputAudioFormat = &format
	}
}

// WithOutputAudioFormat sets the output audio format for the session
func WithOutputAudioFormat(format AudioFormat) ConfigOption {
	return func(c *SessionRequest) {
		c.OutputAudioFormat = &format
	}
}

// WithInputAudioTranscription sets the input audio transcription configuration for the session
func WithInputAudioTranscription(transcription InputAudioTranscription) ConfigOption {
	return func(c *SessionRequest) {
		c.InputAudioTranscription = &transcription
	}
}

// WithTurnDetection sets the turn detection configuration for the session
func WithTurnDetection(turnDetection TurnDetection) ConfigOption {
	return func(c *SessionRequest) {
		c.TurnDetection = &turnDetection
	}
}

// WithTools sets the tools for the session
func WithTools(tools []Tool) ConfigOption {
	return func(c *SessionRequest) {
		c.Tools = &tools
	}
}

// WithToolChoice sets the tool choice for the session
func WithToolChoice(toolChoice ToolChoice) ConfigOption {
	return func(c *SessionRequest) {
		c.ToolChoice = &ToolChoiceObj{Type: toolChoice}
	}
}

// WithTemperature sets the temperature for the session
func WithTemperature(temperature float64) ConfigOption {
	return func(c *SessionRequest) {
		c.Temperature = &temperature
	}
}

// DefaultMaxResponseOutputTokens is the default value for max_response_output_tokens if specified
const DefaultMaxResponseOutputTokens = 4000

// WithMaxResponseOutputTokens sets the maximum number of tokens for the response output.
// If tokens is -1, it will set the value to "inf".
// If tokens is 0, it will use the default value of 4000.
// Otherwise, it will use the specified integer value.
func WithMaxResponseOutputTokens(tokens int) ConfigOption {
	return func(c *SessionRequest) {
		if tokens == -1 {
			value := Inf
			c.MaxResponseOutputTokens = &value
		} else if tokens == 0 {
			value := IntOrInf(DefaultMaxResponseOutputTokens)
			c.MaxResponseOutputTokens = &value
		} else {
			value := IntOrInf(tokens)
			c.MaxResponseOutputTokens = &value
		}
	}
}
