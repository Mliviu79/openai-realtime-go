package session

//-----------------------------------------------------------------------------
// Turn Detection Types
//-----------------------------------------------------------------------------

// TurnDetectionType represents the type of turn detection
type TurnDetectionType string

const (
	// TurnDetectionTypeServerVad use server-side VAD to detect turn.
	// This is default value for newly created session.
	TurnDetectionTypeServerVad TurnDetectionType = "server_vad"
)

// TurnDetection represents turn detection configuration
type TurnDetection struct {
	// Type specifies the turn detection method (currently only server_vad is supported)
	Type TurnDetectionType `json:"type,omitempty"`

	// Threshold is the activation threshold for VAD (0.0 to 1.0), defaults to 0.5
	// A higher threshold will require louder audio to activate the model
	Threshold float64 `json:"threshold,omitempty"`

	// PrefixPaddingMs is the amount of audio to include before VAD detected speech (in milliseconds)
	// Defaults to 300ms
	PrefixPaddingMs int `json:"prefix_padding_ms,omitempty"`

	// SilenceDurationMs is the duration of silence to detect speech stop (in milliseconds)
	// Defaults to 500ms
	SilenceDurationMs int `json:"silence_duration_ms,omitempty"`

	// CreateResponse determines whether to automatically generate a response when a VAD stop event occurs
	// Defaults to true
	CreateResponse *bool `json:"create_response,omitempty"`

	// InterruptResponse determines whether to automatically interrupt any ongoing response
	// when a VAD start event occurs. Defaults to true
	InterruptResponse *bool `json:"interrupt_response,omitempty"`
}
