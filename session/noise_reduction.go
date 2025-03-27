package session

//-----------------------------------------------------------------------------
// Noise Reduction Types
//-----------------------------------------------------------------------------

// NoiseReductionType represents the type of noise reduction
type NoiseReductionType string

const (
	// NoiseReductionTypeNearField is for close-talking microphones such as headphones
	NoiseReductionTypeNearField NoiseReductionType = "near_field"

	// NoiseReductionTypeFarField is for far-field microphones such as laptop or conference room microphones
	NoiseReductionTypeFarField NoiseReductionType = "far_field"
)

// InputAudioNoiseReduction represents configuration for noise reduction on input audio
type InputAudioNoiseReduction struct {
	// Type specifies the type of noise reduction
	// NearField is for close-talking microphones such as headphones
	// FarField is for far-field microphones such as laptop or conference room microphones
	Type NoiseReductionType `json:"type,omitempty"`
}
