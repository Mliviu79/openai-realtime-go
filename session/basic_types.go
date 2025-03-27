package session

//-----------------------------------------------------------------------------
// Basic Types and Constants
//-----------------------------------------------------------------------------

// AudioFormat represents the supported audio formats
type AudioFormat string

const (
	// AudioFormatPCM16 is the PCM 16-bit audio format
	AudioFormatPCM16 AudioFormat = "pcm16"

	// AudioFormatG711ULaw is the G.711 μ-law audio format
	AudioFormatG711ULaw AudioFormat = "g711_ulaw"

	// AudioFormatG711ALaw is the G.711 A-law audio format
	AudioFormatG711ALaw AudioFormat = "g711_alaw"
)

// Modality represents the supported modalities
type Modality string

const (
	// ModalityAudio represents audio input/output capability
	ModalityAudio Modality = "audio"

	// ModalityText represents text input/output capability
	ModalityText Modality = "text"
)

// Model represents available OpenAI models
type Model string

const (
	// GPT4oRealtimePreview is the base GPT-4o realtime preview model
	GPT4oRealtimePreview Model = "gpt-4o-realtime-preview"

	// GPT4oRealtimePreview20241001 is the October 2024 version of GPT-4o realtime
	GPT4oRealtimePreview20241001 Model = "gpt-4o-realtime-preview-2024-10-01"

	// GPT4oRealtimePreview20241217 is the December 2024 version of GPT-4o realtime
	GPT4oRealtimePreview20241217 Model = "gpt-4o-realtime-preview-2024-12-17"

	// GPT4oMiniRealtimePreview is the base GPT-4o mini realtime preview model
	GPT4oMiniRealtimePreview Model = "gpt-4o-mini-realtime-preview"

	// GPT4oMiniRealtimePreview20241217 is the December 2024 version of GPT-4o mini realtime
	GPT4oMiniRealtimePreview20241217 Model = "gpt-4o-mini-realtime-preview-2024-12-17"
)

type Intent string

const (
	transcribe Intent = "transcribe"
)

// Voice represents the available voice options
type Voice string

const (
	// VoiceAlloy is a neutral voice with a slight British accent
	VoiceAlloy Voice = "alloy"

	// VoiceAsh is a warm, clear voice with an American accent
	VoiceAsh Voice = "ash"

	// VoiceBallad is a soft, melodic voice with an American accent
	VoiceBallad Voice = "ballad"

	// VoiceCoral is a warm, friendly voice with an American accent
	VoiceCoral Voice = "coral"

	// VoiceEcho is a baritone voice with an American accent
	VoiceEcho Voice = "echo"

	// VoiceFable is a narration-focused voice with an American accent
	VoiceFable Voice = "fable"

	// VoiceNova is a gender-neutral voice with an American accent
	VoiceNova Voice = "nova"

	// VoiceOnyx is a deep, authoritative voice with an American accent
	VoiceOnyx Voice = "onyx"

	// VoiceSage is a gentle, thoughtful voice with an American accent
	VoiceSage Voice = "sage"

	// VoiceShimmer is a bright, enthusiastic voice with an American accent
	VoiceShimmer Voice = "shimmer"

	// VoiceVerse is a deep, resonant voice with an American accent
	VoiceVerse Voice = "verse"
)
