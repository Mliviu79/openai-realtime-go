package outgoing

import (
	"github.com/Mliviu79/go-openai-realtime/messages/types"
)

// AudioBufferAppendMessage is used to append audio data to the buffer
type AudioBufferAppendMessage struct {
	OutMsgBase
	// Audio contains the Base64-encoded audio data to append
	Audio string `json:"audio"`
}

// NewAudioBufferAppendMessage creates a new audio buffer append message
func NewAudioBufferAppendMessage(audio string, transcription *types.InputAudioTranscription) AudioBufferAppendMessage {
	return AudioBufferAppendMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeAudioBufferAppend,
		},
		Audio: audio,
	}
}

// AudioBufferCommitMessage is used to commit the audio buffer
type AudioBufferCommitMessage struct {
	OutMsgBase
}

// NewAudioBufferCommitMessage creates a new audio buffer commit message
// Note: previousItemID parameter is kept for backward compatibility but is no longer used
func NewAudioBufferCommitMessage(previousItemID string) AudioBufferCommitMessage {
	return AudioBufferCommitMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeAudioBufferCommit,
		},
	}
}

// AudioBufferClearMessage is used to clear the audio buffer
type AudioBufferClearMessage struct {
	OutMsgBase
}

// NewAudioBufferClearMessage creates a new audio buffer clear message
func NewAudioBufferClearMessage() AudioBufferClearMessage {
	return AudioBufferClearMessage{
		OutMsgBase: OutMsgBase{
			Type: OutMsgTypeAudioBufferClear,
		},
	}
}
