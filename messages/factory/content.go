// Package factory provides factory functions for creating message components for the OpenAI Realtime API.
// It simplifies the creation of complex message structures by providing helper functions
// that ensure proper construction according to the API specification.
package factory

import (
	"github.com/Mliviu79/openai-realtime-go/messages/types"
)

// TextContent creates a new text content part
func TextContent(text string) types.MessageContentPart {
	return types.MessageContentPart{
		Type: types.MessageContentTypeText,
		Text: text,
	}
}

// InputTextContent creates a new input text content part
func InputTextContent(text string) types.MessageContentPart {
	return types.MessageContentPart{
		Type: types.MessageContentTypeInputText,
		Text: text,
	}
}

// InputAudioContent creates a new input audio content part
func InputAudioContent(audio string, transcript string) types.MessageContentPart {
	return types.MessageContentPart{
		Type:       types.MessageContentTypeInputAudio,
		Audio:      audio,
		Transcript: transcript,
	}
}

// ItemReferenceContent creates a new item reference content part
func ItemReferenceContent(id string) types.MessageContentPart {
	return types.MessageContentPart{
		Type: types.MessageContentTypeItemReference,
		ID:   id,
	}
}

// AudioContent creates a new audio content part
func AudioContent(audio string) types.MessageContentPart {
	return types.MessageContentPart{
		Type:  types.MessageContentTypeAudio,
		Audio: audio,
	}
}

// TranscriptContent creates a new transcript content part
func TranscriptContent(transcript string) types.MessageContentPart {
	return types.MessageContentPart{
		Type:       types.MessageContentTypeTranscript,
		Transcript: transcript,
	}
}
