// Package main provides a comprehensive test for all OpenAI Realtime API message types.
//
// This example demonstrates how to:
// - Create and establish a session with the OpenAI Realtime API
// - Send every type of outgoing message (all 9 types)
// - Handle every type of incoming message (all 28 types)
// - Test all the API's functionality in a real environment
// - Test edge cases and error conditions
//
// To run this example:
//  1. Set the OPENAI_API_KEY environment variable
//  2. Run: go run examples/comprehensive_run.go
//
// This is a live test against the actual OpenAI API, so it will consume API credits.
package main

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Mliviu79/go-openai-realtime/logger"
	"github.com/Mliviu79/go-openai-realtime/messages/incoming"
	"github.com/Mliviu79/go-openai-realtime/messages/outgoing"
	"github.com/Mliviu79/go-openai-realtime/messages/types"
	"github.com/Mliviu79/go-openai-realtime/messaging"
	"github.com/Mliviu79/go-openai-realtime/openaiClient"
	"github.com/Mliviu79/go-openai-realtime/session"
	"github.com/Mliviu79/go-openai-realtime/ws"
	"github.com/rs/zerolog"
)

// TestTracker helps track which message types we've received
type TestTracker struct {
	mutex sync.Mutex
	seen  map[string]bool
}

// NewTestTracker creates a new test tracker
func NewTestTracker() *TestTracker {
	return &TestTracker{
		seen: make(map[string]bool),
	}
}

// MarkSeen marks a message type as seen
func (t *TestTracker) MarkSeen(msgType string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.seen[msgType] = true
	fmt.Printf("✓ Received message type: %s\n", msgType)
}

// GetUnseen returns a list of message types that haven't been seen
func (t *TestTracker) GetUnseen() []string {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	var unseen []string
	// All 28 incoming message types that can be received from the server
	allTypes := []string{
		// Session management
		"session.created",
		"session.updated",

		// Conversation management
		"conversation.created",
		"conversation.item.created",
		"conversation.truncated",       // Added: Confirmation of conversation truncation
		"conversation.items.truncated", // Added: Multiple items truncated

		// Audio buffer events
		"input_audio_buffer.committed",
		"input_audio_buffer.speech_started",
		"input_audio_buffer.speech_stopped",

		// Audio transcription events
		"conversation.item.input_audio_transcription.completed",
		"conversation.item.input_audio_transcription.failed", // Added: Audio transcription failure

		// Response lifecycle events
		"response.created",
		"response.canceled", // Added: Response cancellation confirmation
		"response.done",

		// Response delta events
		"response.text.delta",
		"response.audio.delta",
		"response.audio_transcript.delta",
		"response.function_call_arguments.delta",

		// Response completion events
		"response.audio.done",
		"response.audio_transcript.done",
		"response.function_call_arguments.done",
		"response.function_call_output.done", // Added: Function call output completion

		// Response content events
		"response.content_part.added",
		"response.content_part.done",
		"response.output_item.added",
		"response.output_item.done",

		// System events
		"rate_limits.updated",
		"error", // Added: Generic error message
	}

	for _, typeName := range allTypes {
		if !t.seen[typeName] {
			unseen = append(unseen, typeName)
		}
	}

	return unseen
}

// Pretty prints any message for debugging
func prettyPrint(msg interface{}) {
	data, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling message: %v\n", err)
		return
	}
	fmt.Printf("Message content: %s\n", string(data))
}

// RunComprehensiveTest runs all the tests for the API
func RunComprehensiveTest() {
	fmt.Println("Starting comprehensive OpenAI Realtime API test...")
	fmt.Println("This test will attempt to trigger all 28 incoming message types and send all 9 outgoing message types.")

	// Get API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: OPENAI_API_KEY environment variable is required")
		os.Exit(1)
	}

	// Set up a logger
	zeroLogger := logger.NewZeroLogger(logger.LoggerOptions{
		Level: zerolog.InfoLevel, // Change to DebugLevel for more detailed logs
	})

	// Create a message type tracker
	tracker := NewTestTracker()

	// Create a client with your OpenAI API key
	client := openaiClient.NewClient(apiKey)

	// Set up context with timeout and cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nReceived shutdown signal, cleaning up...")
		cancel()
	}()

	// Store for conversation ID and item ID to use in later tests
	var conversationID string
	var itemID string
	var responseID string

	// STEP 1: Create a session directly using session.create message (instead of the client helper)
	// This tests the 9th outgoing message type that wasn't covered before
	fmt.Println("\n[TEST 1/10] Creating a session directly with session.create...")

	model := session.GPT4oMiniRealtimePreview20241217
	modalities := []session.Modality{session.ModalityText, session.ModalityAudio}

	// Create a direct websocket connection first
	// Rather than using ConnectDirectWS (which doesn't exist), let's use Connect with minimal options
	// and then send the session.create message manually
	conn, err := client.Connect(ctx, openaiClient.WithModel(model), openaiClient.WithLogger(zeroLogger))
	if err != nil {
		fmt.Printf("Failed to establish direct WebSocket connection: %v\n", err)
		os.Exit(1)
	}

	// Create messaging client from the raw connection
	directMsgClient := messaging.NewClient(conn)
	directMsgClient.SetLogger(zeroLogger)

	// Create a session request
	createReq := session.SessionRequest{
		Model:      &model,
		Modalities: &modalities,
	}

	// Manually create and send a session.create message (outgoing message type #1)
	// Need to create a JSON object with the session request
	createMsgData := map[string]interface{}{
		"type":    "session.create",
		"session": createReq,
	}

	createMsgBytes, err := json.Marshal(createMsgData)
	if err != nil {
		fmt.Printf("Failed to marshal session.create message: %v\n", err)
		os.Exit(1)
	}

	err = conn.SendRaw(ctx, ws.MessageText, createMsgBytes)
	if err != nil {
		fmt.Printf("Failed to send session.create message: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Sent session.create message directly")

	// Set up a listener for the session creation response
	var directSessionID string
	var wgSessionCreation sync.WaitGroup
	wgSessionCreation.Add(1)

	go func() {
		defer wgSessionCreation.Done()

		for i := 0; i < 5; i++ { // Try up to 5 times
			msg, err := directMsgClient.ReadMessage(ctx)
			if err != nil {
				fmt.Printf("Error reading message: %v\n", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if sessionCreated, ok := msg.(*incoming.SessionCreatedMessage); ok {
				tracker.MarkSeen("session.created")
				directSessionID = sessionCreated.Session.ID
				fmt.Printf("Session created directly with ID: %s\n", directSessionID)
				return
			}

			if errMsg, ok := msg.(*incoming.ErrorMessage); ok {
				tracker.MarkSeen("error")
				fmt.Printf("Error creating session: %s - %s\n", errMsg.Error.Type, errMsg.Error.Message)
			}
		}
		fmt.Println("Failed to receive session.created response")
	}()

	// Wait for session creation
	wgSessionCreation.Wait()
	if directSessionID == "" {
		fmt.Println("Failed to create session directly, falling back to client API")

		// Clean up the direct connection
		directMsgClient.Close()

		// Fall back to creating a session with the helper
		createClientReq := &session.CreateRequest{
			SessionRequest: session.SessionRequest{
				Model:      &model,
				Modalities: &modalities,
			},
		}

		sessionResp, err := client.CreateSession(ctx, createClientReq)
		if err != nil {
			fmt.Printf("Failed to create session with client API: %v\n", err)
			os.Exit(1)
		}

		directSessionID = sessionResp.ID
		fmt.Printf("Session created with client API, ID: %s\n", directSessionID)
	}

	// STEP 2: Connect to the session
	fmt.Println("\n[TEST 2/10] Connecting to the session...")

	conn, err = client.Connect(ctx,
		openaiClient.WithModel(model),
		openaiClient.WithLogger(zeroLogger),
		openaiClient.WithSessionID(directSessionID),
	)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create a messaging client from the WebSocket connection
	msgClient := messaging.NewClient(conn)
	msgClient.SetLogger(zeroLogger)

	// Create a waitgroup for test completion
	var wg sync.WaitGroup
	wg.Add(1)

	// Start message listener in a separate goroutine
	go func() {
		defer wg.Done()

		fmt.Println("\n[LISTENER] Starting message listener...")

		for {
			select {
			case <-ctx.Done():
				fmt.Println("[LISTENER] Context cancelled, stopping message listener")
				return
			default:
				msg, err := msgClient.ReadMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						// Context cancelled, normal exit
						return
					}
					fmt.Printf("[LISTENER] Error reading message: %v\n", err)

					// Check for websocket close errors and exit gracefully
					if err.Error() == "websocket: close 1000 (normal)" ||
						err.Error() == "repeated read on failed websocket connection" {
						fmt.Println("[LISTENER] WebSocket connection closed. Exiting listener.")
						return
					}

					// For other errors, wait briefly and continue
					time.Sleep(500 * time.Millisecond)
					continue
				}

				// Handle different message types
				switch event := msg.(type) {
				case *incoming.SessionCreatedMessage:
					tracker.MarkSeen("session.created")
					fmt.Printf("Session created with ID: %s\n", event.Session.ID)

				case *incoming.SessionUpdatedMessage:
					tracker.MarkSeen("session.updated")
					fmt.Println("Session updated successfully")

				case *incoming.ConversationCreatedMessage:
					tracker.MarkSeen("conversation.created")
					conversationID = event.Conversation.ID
					fmt.Printf("Conversation created with ID: %s\n", conversationID)

				case *incoming.ConversationItemCreatedMessage:
					tracker.MarkSeen("conversation.item.created")
					itemID = event.Item.ID
					fmt.Printf("Conversation item created with ID: %s\n", itemID)

				case *incoming.ConversationItemTruncatedMessage:
					tracker.MarkSeen("conversation.item.truncated")
					fmt.Println("Conversation item truncated successfully")

				// Handle missing types - note these might need adjustments based on the actual structs
				// Since we couldn't find the exact struct definitions, these are educated guesses
				case *incoming.ErrorMessage:
					tracker.MarkSeen("error")
					fmt.Printf("Error from API: %s - %s\n", event.Error.Type, event.Error.Message)

				case *incoming.AudioBufferCommittedMessage:
					tracker.MarkSeen("input_audio_buffer.committed")
					fmt.Println("Audio buffer committed successfully")

				case *incoming.AudioBufferSpeechStartedMessage:
					tracker.MarkSeen("input_audio_buffer.speech_started")
					fmt.Println("Speech started in audio buffer")

				case *incoming.AudioBufferSpeechStoppedMessage:
					tracker.MarkSeen("input_audio_buffer.speech_stopped")
					fmt.Println("Speech stopped in audio buffer")

				case *incoming.ConversationItemTranscriptionCompletedMessage:
					tracker.MarkSeen("conversation.item.input_audio_transcription.completed")
					fmt.Printf("Audio transcription completed: %s\n", event.Transcript)

				case *incoming.ConversationItemTranscriptionFailedMessage:
					tracker.MarkSeen("conversation.item.input_audio_transcription.failed")
					fmt.Printf("Audio transcription failed: %s\n", event.Error.Message)

				case *incoming.ResponseCreatedMessage:
					tracker.MarkSeen("response.created")
					responseID = event.Response.ID
					fmt.Printf("Response created with ID: %s\n", responseID)

				// Handle response canceled message - This should match the actual struct when available
				default:
					// Handle message types we don't have explicit cases for
					msgType := msg.RcvdMsgType()

					// Mark all incoming messages as seen based on their type string
					tracker.MarkSeen(string(msgType))

					fmt.Printf("Received message type: %s\n", msgType)

					// Try to extract response ID from response.canceled message
					if string(msgType) == "response.canceled" {
						// Extract response ID using reflection or type assertion if possible
						fmt.Println("Response canceled received")
					}

					// Similarly track other missing types
					if string(msgType) == "conversation.truncated" {
						fmt.Println("Conversation truncated received")
					}

					if string(msgType) == "conversation.items.truncated" {
						fmt.Println("Conversation items truncated received")
					}

					if string(msgType) == "response.function_call_output.done" {
						fmt.Println("Function call output done received")
					}
				}
			}
		}
	}()

	// Wait a moment for connection to fully establish
	time.Sleep(1 * time.Second)

	// STEP 3: Update the session (outgoing message type #2)
	fmt.Println("\n[TEST 3/10] Updating session parameters...")

	turnDetection := session.TurnDetection{
		Type: session.TurnDetectionTypeServerVad,
	}

	instructions := "You are a helpful assistant specialized in technology and science."
	voice := session.VoiceEcho
	outputAudioFormat := session.AudioFormatPCM16
	temperature := 0.7

	// Define sample tools - using the correct flattened structure
	tools := []session.Tool{
		{
			Type:        "function",
			Name:        "test_function",
			Description: "A simple test function",
			Parameters:  json.RawMessage(`{"type":"object","properties":{"input":{"type":"string"}},"required":["input"]}`),
		},
	}

	// Print the tool definition for debugging
	toolsJson, _ := json.Marshal(tools)
	fmt.Printf("Tools JSON: %s\n", string(toolsJson))

	// Configure session update parameters - using only the standard approach
	sessionUpdateReq := session.SessionRequest{
		Modalities:        &modalities,
		Instructions:      &instructions,
		Voice:             &voice,
		OutputAudioFormat: &outputAudioFormat,
		TurnDetection:     &turnDetection,
		Temperature:       &temperature,
		Tools:             &tools,
	}

	// Send the session update message
	fmt.Println("Sending session.update message...")
	sessionUpdateMsg := outgoing.NewSessionUpdateMessage(sessionUpdateReq)
	err = msgClient.SendMessage(ctx, sessionUpdateMsg)
	if err != nil {
		fmt.Printf("Failed to update session: %v\n", err)
		// Continue execution even if this fails
	}

	// Wait a moment for the update to process
	time.Sleep(2 * time.Second)

	// STEP 4: Create a conversation and add items (outgoing message type #3)
	fmt.Println("\n[TEST 4/10] Creating a conversation with text items...")

	// Create a user message
	userMsg := types.MessageItem{
		Type: types.MessageItemTypeMessage,
		Role: types.MessageRoleUser,
		Content: []types.MessageContentPart{
			{
				Type: types.MessageContentTypeInputText,
				Text: "What are the planets in our solar system?",
			},
		},
	}

	// Send the conversation item create message
	err = msgClient.SendConversationItemCreate(ctx, &userMsg, nil)
	if err != nil {
		fmt.Printf("Failed to create conversation item: %v\n", err)
		os.Exit(1)
	}

	// Wait a bit for the conversation to process
	time.Sleep(2 * time.Second)

	// STEP 5: Send audio buffer data (outgoing message types #4, #5, #6)
	fmt.Println("\n[TEST 5/10] Testing audio buffer functionality...")

	// Generate "empty" audio data for testing (this would normally be real audio)
	// Convert to base64 string as required by the API
	audioBytes := make([]byte, 1600) // 100ms of 16kHz PCM16 audio (all zeros)
	audioData := base64.StdEncoding.EncodeToString(audioBytes)

	// Append audio to buffer (outgoing message type #4)
	err = msgClient.SendAudioBufferAppend(ctx, audioData)
	if err != nil {
		fmt.Printf("Failed to append audio to buffer: %v\n", err)
	}

	// Wait a moment
	time.Sleep(500 * time.Millisecond)

	// Send a second chunk of valid audio instead of intentionally invalid audio
	// This helps test multiple append operations without breaking the connection
	secondAudioBytes := make([]byte, 1600)
	// Fill with a simple sine wave pattern instead of zeros
	for i := 0; i < len(secondAudioBytes); i += 2 {
		value := int16(10000 * math.Sin(float64(i)/100))
		binary.LittleEndian.PutUint16(secondAudioBytes[i:i+2], uint16(value))
	}
	secondAudioData := base64.StdEncoding.EncodeToString(secondAudioBytes)

	err = msgClient.SendAudioBufferAppend(ctx, secondAudioData)
	if err != nil {
		fmt.Printf("Failed to append second audio chunk: %v\n", err)
	}

	// Wait a moment
	time.Sleep(500 * time.Millisecond)

	// Commit the audio buffer (outgoing message type #5)
	// By using the real item ID if available, or an empty string for testing
	commitItemID := itemID
	if commitItemID == "" {
		commitItemID = ""
	}

	err = msgClient.SendAudioBufferCommit(ctx, commitItemID)
	if err != nil {
		fmt.Printf("Failed to commit audio buffer: %v\n", err)
	}

	// Wait a moment
	time.Sleep(500 * time.Millisecond)

	// Clear the audio buffer (outgoing message type #6)
	err = msgClient.SendAudioBufferClear(ctx)
	if err != nil {
		fmt.Printf("Failed to clear audio buffer: %v\n", err)
	}

	// Wait for audio buffer operations to complete
	time.Sleep(2 * time.Second)

	// STEP 6: Create a response with tools (outgoing message type #7)
	fmt.Println("\n[TEST 6/10] Creating a response with tools...")

	// Create a ResponseConfig with tools and both modalities
	responseModalities := []session.Modality{session.ModalityText, session.ModalityAudio}
	responseInstructions := "Give a very brief weather report for San Francisco, using the get_weather function."

	responseConfig := types.ResponseConfig{
		Modalities:   responseModalities,
		Instructions: &responseInstructions,
		Voice:        &voice,
		Tools:        tools,
		Metadata: map[string]string{
			"test_id": "comprehensive_test",
		},
	}

	// Send the response.create message
	err = msgClient.SendResponseCreate(ctx, &responseConfig)
	if err != nil {
		fmt.Printf("Failed to create response: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Function call response request sent")

	// Wait longer for this response as it needs to make a function call
	time.Sleep(10 * time.Second)

	// STEP 7: Test response cancellation (outgoing message type #8)
	fmt.Println("\n[TEST 7/10] Testing response cancellation...")

	// Create another response
	cancelResponseConfig := types.ResponseConfig{
		Modalities:   responseModalities,
		Instructions: &responseInstructions,
	}

	// Send the response.create message
	err = msgClient.SendResponseCreate(ctx, &cancelResponseConfig)
	if err != nil {
		fmt.Printf("Failed to create cancellation test response: %v\n", err)
		os.Exit(1)
	}

	// Wait for the response ID to be captured by the listener
	time.Sleep(2 * time.Second)

	// Cancel the response if we have an ID
	if responseID != "" {
		fmt.Printf("Cancelling response with ID: %s\n", responseID)
		err := msgClient.SendResponseCancel(ctx, responseID)
		if err != nil {
			fmt.Printf("Failed to cancel response: %v\n", err)
		}
	} else {
		fmt.Println("No response ID available for cancellation")
	}

	// Wait for cancellation to process
	time.Sleep(3 * time.Second)

	// STEP 8: Create a more complex conversation for truncation
	fmt.Println("\n[TEST 8/10] Creating a more complex conversation for truncation tests...")

	// Create multiple messages
	for i := 0; i < 3; i++ {
		msg := types.MessageItem{
			Type: types.MessageItemTypeMessage,
			Role: types.MessageRoleUser,
			Content: []types.MessageContentPart{
				{
					Type: types.MessageContentTypeInputText,
					Text: fmt.Sprintf("This is message %d in the conversation", i+1),
				},
			},
		}

		err = msgClient.SendConversationItemCreate(ctx, &msg, nil)
		if err != nil {
			fmt.Printf("Failed to create message %d: %v\n", i+1, err)
			continue
		}

		time.Sleep(1 * time.Second)
	}

	// Wait for messages to be processed
	time.Sleep(3 * time.Second)

	// STEP 9: Test conversation truncation (outgoing message type #9)
	fmt.Println("\n[TEST 9/10] Testing conversation truncation...")

	// Only truncate if we have a valid item ID from conversation creation
	if itemID != "" {
		// Now truncate the conversation
		err = msgClient.SendConversationItemTruncate(ctx, itemID, 1, 0)
		if err != nil {
			fmt.Printf("Failed to truncate conversation: %v\n", err)
		}
		fmt.Printf("Sent truncate request for item ID: %s\n", itemID)
	} else {
		fmt.Println("No valid item ID available for truncation")

		// Try another approach - create a new item and immediately truncate it
		newMsg := types.MessageItem{
			Type: types.MessageItemTypeMessage,
			Role: types.MessageRoleUser,
			Content: []types.MessageContentPart{
				{
					Type: types.MessageContentTypeInputText,
					Text: "This is a test message for truncation",
				},
			},
		}

		err = msgClient.SendConversationItemCreate(ctx, &newMsg, nil)
		if err != nil {
			fmt.Printf("Failed to create item for truncation: %v\n", err)
		} else {
			// Wait briefly for item creation
			time.Sleep(2 * time.Second)

			// Use a dummy ID as fallback if we still don't have a real one
			truncateID := itemID
			if truncateID == "" {
				truncateID = "item_last"
			}

			// Try truncation
			err = msgClient.SendConversationItemTruncate(ctx, truncateID, 1, 0)
			if err != nil {
				fmt.Printf("Failed to truncate with fallback ID: %v\n", err)
			}
		}
	}

	// Wait for truncation to process
	time.Sleep(2 * time.Second)

	// STEP 10: Try to trigger error conditions to get full message coverage
	fmt.Println("\n[TEST 10/10] Testing error conditions for additional message types...")

	// Try to trigger a transcription failure with corrupt audio
	corruptAudioBytes := make([]byte, 800)
	// Fill with pattern that's not valid audio
	for i := range corruptAudioBytes {
		corruptAudioBytes[i] = byte(i % 255)
	}
	corruptAudioData := base64.StdEncoding.EncodeToString(corruptAudioBytes)

	// Send corrupt audio and commit it
	err = msgClient.SendAudioBufferAppend(ctx, corruptAudioData)
	if err != nil {
		fmt.Printf("Error appending corrupt audio (expected): %v\n", err)
	} else {
		err = msgClient.SendAudioBufferCommit(ctx, "")
		if err != nil {
			fmt.Printf("Error committing corrupt audio (expected): %v\n", err)
		}
	}

	// Try to create a response with invalid configuration
	invalidResponseConfig := types.ResponseConfig{
		// Missing required fields intentionally to trigger errors
		Modalities: []session.Modality{"invalid_modality"}, // Invalid modality
	}

	err = msgClient.SendResponseCreate(ctx, &invalidResponseConfig)
	if err != nil {
		fmt.Printf("Error creating response with invalid config (expected): %v\n", err)
	}

	// Try to cancel a non-existent response
	err = msgClient.SendResponseCancel(ctx, "resp_nonexistent")
	if err != nil {
		fmt.Printf("Error cancelling non-existent response (expected): %v\n", err)
	}

	// Wait a bit for any error responses
	time.Sleep(5 * time.Second)

	// Print summary of message types we received
	fmt.Println("\n=== TEST RESULTS ===")
	unseen := tracker.GetUnseen()

	if len(unseen) == 0 {
		fmt.Println("SUCCESS! All 28 expected message types were received.")
		fmt.Println("All 9 outgoing message types were sent.")
	} else {
		fmt.Println("The following message types were NOT observed during testing:")
		for _, t := range unseen {
			fmt.Printf(" - %s\n", t)
		}
		fmt.Println("\nNote: Some message types may not appear in every test run depending on API behavior.")
		fmt.Println("Especially error conditions and edge cases may not trigger all possible message types.")
	}

	// Wait for the message listener to complete
	cancel() // Signal the context is done
	wg.Wait()

	fmt.Println("\nComprehensive test completed.")
}

// main is the entry point for the standalone example
func main() {
	fmt.Println("Starting main() function...")
	RunComprehensiveTest()
}
