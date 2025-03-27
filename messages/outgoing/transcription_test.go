package outgoing

import (
	"encoding/json"
	"testing"

	"github.com/Mliviu79/openai-realtime-go/session"
	"github.com/stretchr/testify/assert"
)

func TestTranscriptionSessionUpdateMessage(t *testing.T) {
	// Set up test data
	format := session.AudioFormatPCM16
	model := session.TranscriptionModelGPT4oTranscribe
	language := "en"
	prompt := "technical vocabulary"
	includes := []session.TranscriptionSessionInclude{
		session.TranscriptionSessionIncludeLogprobs,
	}

	// Create a TranscriptionSessionRequest with all fields populated
	req := session.TranscriptionSessionRequest{
		InputAudioFormat: &format,
		InputAudioTranscription: &session.InputAudioTranscription{
			Model:    model,
			Language: language,
			Prompt:   prompt,
		},
		Include: &includes,
	}

	// Create message
	msg := NewTranscriptionSessionUpdateMessage(req)

	// Check the message type
	assert.Equal(t, "transcription_session.update", msg.OutMsgType())
	assert.Equal(t, "", msg.OutMsgID())

	// Create message with ID
	msgWithID := NewTranscriptionSessionUpdateMessageWithID("test-id-123", req)
	assert.Equal(t, "test-id-123", msgWithID.OutMsgID())

	// Test JSON serialization
	serialized, err := json.Marshal(msg)
	assert.NoError(t, err)

	// Test that the serialized format contains expected fields
	var deserialized map[string]interface{}
	err = json.Unmarshal(serialized, &deserialized)
	assert.NoError(t, err)

	// Verify basic structure
	assert.Equal(t, "transcription_session.update", deserialized["type"])
	sessionData, ok := deserialized["session"].(map[string]interface{})
	assert.True(t, ok)

	// Check session fields
	assert.Equal(t, "pcm16", sessionData["input_audio_format"])

	transcription, ok := sessionData["input_audio_transcription"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "gpt-4o-transcribe", transcription["model"])
	assert.Equal(t, "en", transcription["language"])
	assert.Equal(t, "technical vocabulary", transcription["prompt"])

	includeArray, ok := sessionData["include"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 1, len(includeArray))
	assert.Equal(t, "item.input_audio_transcription.logprobs", includeArray[0])
}
