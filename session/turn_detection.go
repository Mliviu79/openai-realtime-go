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

	// TurnDetectionTypeSemanticVad uses semantic turn detection with VAD.
	// This mode uses a turn detection model to semantically estimate whether the user has finished speaking.
	TurnDetectionTypeSemanticVad TurnDetectionType = "semantic_vad"
)

// EagernessLevel represents the eagerness of the model to respond in semantic VAD mode
type EagernessLevel string

const (
	// EagernessLevelLow means the model will wait longer for the user to continue speaking
	EagernessLevelLow EagernessLevel = "low"

	// EagernessLevelMedium is a balanced eagerness level
	EagernessLevelMedium EagernessLevel = "medium"

	// EagernessLevelHigh means the model will respond more quickly
	EagernessLevelHigh EagernessLevel = "high"

	// EagernessLevelAuto is the default eagerness level (equivalent to medium)
	EagernessLevelAuto EagernessLevel = "auto"
)

// TurnDetection represents turn detection configuration
type TurnDetection struct {
	// Type specifies the turn detection method (server_vad or semantic_vad)
	Type TurnDetectionType `json:"type,omitempty"`

	// Eagerness is used only for semantic_vad mode and controls how eager the model is to respond
	// Low will wait longer for the user to continue speaking, high will respond more quickly
	// Auto is the default and is equivalent to medium
	Eagerness EagernessLevel `json:"eagerness,omitempty"`

	// Threshold is the activation threshold for VAD (0.0 to 1.0), defaults to 0.5
	// A higher threshold will require louder audio to activate the model
	// Only used for server_vad mode
	Threshold float64 `json:"threshold,omitempty"`

	// PrefixPaddingMs is the amount of audio to include before VAD detected speech (in milliseconds)
	// Defaults to 300ms
	// Only used for server_vad mode
	PrefixPaddingMs int `json:"prefix_padding_ms,omitempty"`

	// SilenceDurationMs is the duration of silence to detect speech stop (in milliseconds)
	// Defaults to 500ms
	// Only used for server_vad mode
	SilenceDurationMs int `json:"silence_duration_ms,omitempty"`

	// CreateResponse determines whether to automatically generate a response when a VAD stop event occurs
	// Defaults to true
	CreateResponse *bool `json:"create_response,omitempty"`

	// InterruptResponse determines whether to automatically interrupt any ongoing response
	// when a VAD start event occurs. Defaults to true
	InterruptResponse *bool `json:"interrupt_response,omitempty"`
}
