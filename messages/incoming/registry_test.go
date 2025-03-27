package incoming

import (
	"testing"
)

func TestCreateMessage(t *testing.T) {
	// Test creating all registered message types
	for msgType := range MessageTypeRegistry {
		msg, exists := CreateMessage(msgType)
		if !exists {
			t.Errorf("CreateMessage(%q) returned exists=false, expected true", msgType)
			continue
		}

		if msg == nil {
			t.Errorf("CreateMessage(%q) returned nil message", msgType)
			continue
		}

		if msg.RcvdMsgType() != msgType {
			t.Errorf("Expected message type %q, got %q", msgType, msg.RcvdMsgType())
		}
	}

	// Test creating a non-existent message type
	msg, exists := CreateMessage("non.existent.type")
	if exists {
		t.Errorf("CreateMessage for non-existent type returned exists=true, expected false")
	}

	if msg != nil {
		t.Errorf("CreateMessage for non-existent type returned non-nil message: %T", msg)
	}
}

func TestMessageTypeRegistry(t *testing.T) {
	// Ensure all message types defined in types.go are registered
	expectedTypes := []RcvdMsgType{
		// Error message type
		RcvdMsgTypeError,

		// Session-related message types
		RcvdMsgTypeSessionCreated,
		RcvdMsgTypeSessionUpdated,

		// Transcription session-related message types
		RcvdMsgTypeTranscriptionSessionCreated,
		RcvdMsgTypeTranscriptionSessionUpdated,
		RcvdMsgTypeInputAudioTranscription,
		RcvdMsgTypeTranscriptionDone,

		// Conversation-related message types
		RcvdMsgTypeConversationCreated,
		RcvdMsgTypeConversationItemCreated,
		RcvdMsgTypeConversationItemInputAudioTranscriptionCompleted,
		RcvdMsgTypeConversationItemInputAudioTranscriptionFailed,
		RcvdMsgTypeConversationItemTruncated,
		RcvdMsgTypeConversationItemDeleted,

		// Audio buffer-related message types
		RcvdMsgTypeAudioBufferCommitted,
		RcvdMsgTypeAudioBufferCleared,
		RcvdMsgTypeAudioBufferSpeechStarted,
		RcvdMsgTypeAudioBufferSpeechStopped,

		// Response-related message types
		RcvdMsgTypeResponseCreated,
		RcvdMsgTypeResponseDone,
		RcvdMsgTypeResponseContentPartAdded,
		RcvdMsgTypeResponseContentPartDone,
		RcvdMsgTypeResponseTextDelta,
		RcvdMsgTypeResponseTextDone,
		RcvdMsgTypeResponseOutputItemAdded,
		RcvdMsgTypeResponseOutputItemDone,
		RcvdMsgTypeResponseAudioTranscriptDelta,
		RcvdMsgTypeResponseAudioTranscriptDone,
		RcvdMsgTypeResponseAudioDelta,
		RcvdMsgTypeResponseAudioDone,
		RcvdMsgTypeResponseFunctionCallArgumentsDelta,
		RcvdMsgTypeResponseFunctionCallArgumentsDone,

		// Rate limit-related message types
		RcvdMsgTypeRateLimitsUpdated,
	}

	for _, expectedType := range expectedTypes {
		if _, exists := MessageTypeRegistry[expectedType]; !exists {
			t.Errorf("Expected message type %q to be registered, but it is not", expectedType)
		}
	}

	// Check for equal lengths to ensure no extra types
	if len(MessageTypeRegistry) != len(expectedTypes) {
		t.Errorf("MessageTypeRegistry has %d types, but expected %d types", len(MessageTypeRegistry), len(expectedTypes))
	}
}
