package incoming

// This file contains a registry of message types and factory functions
// for creating new instances of each message type.

// MessageTypeRegistry maps message types to factory functions
var MessageTypeRegistry = map[RcvdMsgType]func() RcvdMsg{
	// Error message
	RcvdMsgTypeError: func() RcvdMsg { return &ErrorMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeError}} },

	// Session-related messages
	RcvdMsgTypeSessionCreated: func() RcvdMsg {
		return &SessionCreatedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeSessionCreated}}
	},
	RcvdMsgTypeSessionUpdated: func() RcvdMsg {
		return &SessionUpdatedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeSessionUpdated}}
	},

	// Transcription session-related messages
	RcvdMsgTypeTranscriptionSessionCreated: func() RcvdMsg {
		return &TranscriptionSessionCreatedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeTranscriptionSessionCreated}}
	},
	RcvdMsgTypeTranscriptionSessionUpdated: func() RcvdMsg {
		return &TranscriptionSessionUpdatedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeTranscriptionSessionUpdated}}
	},
	RcvdMsgTypeInputAudioTranscription: func() RcvdMsg {
		return &InputAudioTranscriptionMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeInputAudioTranscription}}
	},
	RcvdMsgTypeTranscriptionDone: func() RcvdMsg {
		return &TranscriptionDoneMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeTranscriptionDone}}
	},

	// Conversation-related messages
	RcvdMsgTypeConversationCreated: func() RcvdMsg {
		return &ConversationCreatedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeConversationCreated}}
	},
	RcvdMsgTypeConversationItemCreated: func() RcvdMsg {
		return &ConversationItemCreatedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeConversationItemCreated}}
	},
	RcvdMsgTypeConversationItemInputAudioTranscriptionCompleted: func() RcvdMsg {
		return &ConversationItemTranscriptionCompletedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeConversationItemInputAudioTranscriptionCompleted}}
	},
	RcvdMsgTypeConversationItemInputAudioTranscriptionFailed: func() RcvdMsg {
		return &ConversationItemTranscriptionFailedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeConversationItemInputAudioTranscriptionFailed}}
	},
	RcvdMsgTypeConversationItemTruncated: func() RcvdMsg {
		return &ConversationItemTruncatedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeConversationItemTruncated}}
	},
	RcvdMsgTypeConversationItemDeleted: func() RcvdMsg {
		return &ConversationItemDeletedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeConversationItemDeleted}}
	},

	// Audio buffer-related messages
	RcvdMsgTypeAudioBufferCommitted: func() RcvdMsg {
		return &AudioBufferCommittedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeAudioBufferCommitted}}
	},
	RcvdMsgTypeAudioBufferCleared: func() RcvdMsg {
		return &AudioBufferClearedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeAudioBufferCleared}}
	},
	RcvdMsgTypeAudioBufferSpeechStarted: func() RcvdMsg {
		return &AudioBufferSpeechStartedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeAudioBufferSpeechStarted}}
	},
	RcvdMsgTypeAudioBufferSpeechStopped: func() RcvdMsg {
		return &AudioBufferSpeechStoppedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeAudioBufferSpeechStopped}}
	},

	// Response-related messages
	RcvdMsgTypeResponseCreated: func() RcvdMsg {
		return &ResponseCreatedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseCreated}}
	},
	RcvdMsgTypeResponseDone: func() RcvdMsg { return &ResponseDoneMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseDone}} },
	RcvdMsgTypeResponseContentPartAdded: func() RcvdMsg {
		return &ResponseContentPartAddedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseContentPartAdded}}
	},
	RcvdMsgTypeResponseContentPartDone: func() RcvdMsg {
		return &ResponseContentPartDoneMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseContentPartDone}}
	},
	RcvdMsgTypeResponseTextDelta: func() RcvdMsg {
		return &ResponseTextDeltaMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseTextDelta}}
	},
	RcvdMsgTypeResponseTextDone: func() RcvdMsg {
		return &ResponseTextDoneMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseTextDone}}
	},
	RcvdMsgTypeResponseOutputItemAdded: func() RcvdMsg {
		return &ResponseOutputItemAddedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseOutputItemAdded}}
	},
	RcvdMsgTypeResponseOutputItemDone: func() RcvdMsg {
		return &ResponseOutputItemDoneMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseOutputItemDone}}
	},
	RcvdMsgTypeResponseAudioTranscriptDelta: func() RcvdMsg {
		return &ResponseAudioTranscriptDeltaMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseAudioTranscriptDelta}}
	},
	RcvdMsgTypeResponseAudioTranscriptDone: func() RcvdMsg {
		return &ResponseAudioTranscriptDoneMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseAudioTranscriptDone}}
	},
	RcvdMsgTypeResponseAudioDelta: func() RcvdMsg {
		return &ResponseAudioDeltaMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseAudioDelta}}
	},
	RcvdMsgTypeResponseAudioDone: func() RcvdMsg {
		return &ResponseAudioDoneMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseAudioDone}}
	},
	RcvdMsgTypeResponseFunctionCallArgumentsDelta: func() RcvdMsg {
		return &ResponseFunctionCallArgumentsDeltaMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseFunctionCallArgumentsDelta}}
	},
	RcvdMsgTypeResponseFunctionCallArgumentsDone: func() RcvdMsg {
		return &ResponseFunctionCallArgumentsDoneMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeResponseFunctionCallArgumentsDone}}
	},

	// Rate limit-related messages
	RcvdMsgTypeRateLimitsUpdated: func() RcvdMsg {
		return &RateLimitsUpdatedMessage{RcvdMsgBase: RcvdMsgBase{Type: RcvdMsgTypeRateLimitsUpdated}}
	},
}

// CreateMessage creates a new instance of the specified message type
func CreateMessage(msgType RcvdMsgType) (RcvdMsg, bool) {
	factory, exists := MessageTypeRegistry[msgType]
	if !exists {
		return nil, false
	}
	return factory(), true
}
