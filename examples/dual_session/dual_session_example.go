// Package main provides an example demonstrating how to run both a realtime conversation
// session and a transcription session simultaneously.
//
// This example showcases how to:
// - Create and connect to both session types at the same time
// - Handle messages from both sessions concurrently using goroutines
// - Send audio to the transcription session and use the results in the conversation
// - Manage separate WebSocket connections with different configurations
//
// To run this example:
//  1. Set the OPENAI_API_KEY environment variable
//  2. Run: go run examples/dual_session_example.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Mliviu79/openai-realtime-go/logger"
	"github.com/Mliviu79/openai-realtime-go/messages/incoming"
	"github.com/Mliviu79/openai-realtime-go/messages/types"
	"github.com/Mliviu79/openai-realtime-go/messaging"
	"github.com/Mliviu79/openai-realtime-go/openaiClient"
	"github.com/Mliviu79/openai-realtime-go/session"
	"github.com/rs/zerolog"
)

// Simulated audio phrases we pretend are being spoken
var simulatedPhrases = []string{
	"This is a test of the transcription system.",
	"Can you explain how machine learning works?",
	"What are the main differences between supervised and unsupervised learning?",
}

// Transcription simulation results, to track what we've "received" to make the simulation more realistic
var receivedTranscriptions = make(map[int]string)
var transcriptionMutex sync.Mutex

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create a context with cancellation for clean shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling for graceful shutdown
	setupSignalHandler(cancel)

	// Create a client with API key
	client := openaiClient.NewClient(apiKey)

	// Create logger
	loggerOpts := logger.LoggerOptions{
		Level:  zerolog.InfoLevel,
		Output: os.Stdout,
	}
	debugLogger := logger.NewZeroLogger(loggerOpts)

	// Use waitgroups to ensure all goroutines complete before exiting
	var wg sync.WaitGroup

	// Channel for passing transcribed text from transcription to conversation
	transcriptionCh := make(chan string, 10)

	// Start transcription session
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runTranscriptionSession(ctx, client, debugLogger, transcriptionCh); err != nil {
			log.Printf("Transcription session error: %v", err)
		}
	}()

	// Start conversation session (with slight delay to ensure transcription is ready)
	time.Sleep(500 * time.Millisecond)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := runConversationSession(ctx, client, debugLogger, transcriptionCh); err != nil {
			log.Printf("Conversation session error: %v", err)
		}
	}()

	// Wait for both sessions to complete
	wg.Wait()
	fmt.Println("Both sessions completed successfully")
}

// setupSignalHandler configures handlers for OS signals to ensure graceful shutdown
func setupSignalHandler(cancel context.CancelFunc) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		fmt.Println("\nReceived shutdown signal, cleaning up...")
		cancel()
	}()
}

// runTranscriptionSession creates and manages a transcription session
func runTranscriptionSession(ctx context.Context, client *openaiClient.Client, logger logger.Logger, transcriptionCh chan<- string) error {
	fmt.Println("Starting transcription session...")

	// Configure transcription session
	inputFormat := session.AudioFormatPCM16
	transcriptionModel := session.TranscriptionModelGPT4oTranscribe

	// Create transcription session
	createReq := &session.CreateTranscriptionSessionRequest{
		TranscriptionSessionRequest: session.TranscriptionSessionRequest{
			InputAudioFormat: &inputFormat,
			InputAudioTranscription: &session.InputAudioTranscription{
				Model:    transcriptionModel,
				Language: "en",
			},
		},
	}

	sessionResp, err := client.CreateTranscriptionSession(ctx, createReq)
	if err != nil {
		return fmt.Errorf("failed to create transcription session: %w", err)
	}

	fmt.Printf("Created transcription session with ID: %s\n", sessionResp.ID)

	// Connect to the transcription session
	conn, err := client.ConnectTranscription(ctx,
		openaiClient.WithTranscriptionSessionID(sessionResp.ID),
		openaiClient.WithTranscriptionLogger(logger))
	if err != nil {
		return fmt.Errorf("failed to connect to transcription session: %w", err)
	}

	// Create messaging client
	msgClient := messaging.NewClient(conn)

	// Create noise reduction configuration
	noiseReduction := &session.InputAudioNoiseReduction{
		Type: session.NoiseReductionTypeNearField,
	}

	// Configure turn detection
	turnDetection := &session.TurnDetection{
		Type:      session.TurnDetectionTypeServerVad,
		Threshold: 0.6,
	}

	// Update session settings
	updateReq := session.TranscriptionSessionRequest{
		InputAudioNoiseReduction: noiseReduction,
		TurnDetection:            turnDetection,
	}

	if err := msgClient.SendTranscriptionSessionUpdate(ctx, updateReq); err != nil {
		return fmt.Errorf("failed to update transcription session: %w", err)
	}

	// --------------------------------------
	// SIMULATION: Directly inject simulated transcriptions
	// --------------------------------------
	fmt.Println("[IMPORTANT] This example uses direct simulation instead of real audio processing")
	fmt.Println("[SIMULATION] In a real app, you would send actual audio and process real transcriptions")
	fmt.Println("[SIMULATION] For demonstration purposes, we're simulating the transcription process")

	// Start a goroutine to simulate transcription results
	go func() {
		// Wait a bit for initialization
		time.Sleep(2 * time.Second)

		// Send three simulated transcriptions with pauses
		for i := 0; i < len(simulatedPhrases); i++ {
			select {
			case <-ctx.Done():
				return
			default:
				// The phrase we pretend is being spoken in this audio chunk
				phrase := simulatedPhrases[i]

				fmt.Printf("\n====== SIMULATED AUDIO CHUNK %d ======\n", i+1)
				fmt.Printf("[SIMULATION] Pretending user said: '%s'\n", phrase)

				// We would normally send audio and wait for a transcription.
				// For demonstration, we'll simulate both the sending and the receiving.
				fmt.Println("[SIMULATION] Sending audio chunk (simulated)")
				fmt.Println("[SIMULATION] Audio processing started (simulated)")
				time.Sleep(1 * time.Second) // Simulate processing time

				// Instead of actually sending audio and waiting for the API to transcribe it,
				// we'll directly simulate receiving a transcription
				fmt.Printf("[SIMULATION] Received transcription: %s\n", phrase)

				// Send our simulated transcription to the conversation session
				select {
				case transcriptionCh <- phrase:
					fmt.Println("[SIMULATION] Sent transcription to conversation session")
				default:
					// Don't block if channel is full
					fmt.Println("[SIMULATION] Couldn't send transcription - channel full")
				}

				// Pause between simulated audio chunks
				time.Sleep(8 * time.Second) // Longer pause to allow for AI response
			}
		}

		// Signal that we're done simulating audio
		fmt.Println("\n[SIMULATION] Finished simulating all audio chunks")
	}()

	// In a real app, we would process actual message events from the API.
	// For this simulation, we'll just keep the connection alive and wait for the context to be done
	done := ctx.Done()
	for {
		select {
		case <-done:
			fmt.Println("[TRANSCRIPTION] Session ending due to context cancellation")
			return nil
		default:
			// Just wait for a message, but we don't expect any meaningful ones in our simulation
			msg, err := msgClient.ReadMessage(ctx)
			if err != nil {
				// Just log errors but continue the simulation
				fmt.Printf("[TRANSCRIPTION] Error reading message (expected in simulation): %v\n", err)
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// Log the received message type but don't take any special action
			// as our simulation doesn't depend on the real messages
			switch msg.RcvdMsgType() {
			case incoming.RcvdMsgTypeError:
				if errMsg, ok := msg.(*incoming.ErrorMessage); ok {
					fmt.Printf("[TRANSCRIPTION] Received error (expected in simulation): %s\n", errMsg.Error.Message)
				}
			default:
				// For simulation purposes, we'll ignore most message types silently
				// to avoid cluttering the output
			}
		}
	}
}

// runConversationSession creates and manages a conversation session
func runConversationSession(ctx context.Context, client *openaiClient.Client, logger logger.Logger, transcriptionCh <-chan string) error {
	fmt.Println("Starting conversation session...")

	// Configure the conversation session
	model := session.GPT4oRealtimePreview
	modalities := []session.Modality{session.ModalityText}

	// Create a session request
	createReq := &session.CreateRequest{
		SessionRequest: session.SessionRequest{
			Model:      &model,
			Modalities: &modalities,
		},
	}

	// Create the session
	sessionResp, err := client.CreateSession(ctx, createReq)
	if err != nil {
		return fmt.Errorf("failed to create conversation session: %w", err)
	}

	fmt.Printf("Created conversation session with ID: %s\n", sessionResp.ID)

	// Connect to the conversation session
	conn, err := client.Connect(ctx,
		openaiClient.WithModel(model),
		openaiClient.WithSessionID(sessionResp.ID),
		openaiClient.WithLogger(logger))
	if err != nil {
		return fmt.Errorf("failed to connect to conversation session: %w", err)
	}

	// Create messaging client
	msgClient := messaging.NewClient(conn)

	// Send initial message
	initialMsg := "Hello, I'm going to send you some audio transcriptions. Please respond to them as they come in."
	if err := msgClient.SendText(ctx, initialMsg); err != nil {
		return fmt.Errorf("failed to send initial message: %w", err)
	}

	fmt.Printf("[CONVERSATION] Sent: %s\n", initialMsg)

	// Request the model to generate a response to the initial message
	responseConfig := &types.ResponseConfig{
		Modalities: []session.Modality{session.ModalityText},
	}
	if err := msgClient.SendResponseCreate(ctx, responseConfig); err != nil {
		return fmt.Errorf("failed to request initial response: %w", err)
	}

	// Create channels for concurrency control
	doneCh := make(chan struct{})
	errorCh := make(chan error)

	// Handle receiving transcriptions and sending them to the conversation
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-doneCh:
				return
			case transcription := <-transcriptionCh:
				fmt.Printf("\n====== PROCESSING TRANSCRIPTION ======\n")
				// Prompt with the transcription
				prompt := fmt.Sprintf("I just transcribed this audio: \"%s\". Please respond to it.", transcription)

				if err := msgClient.SendText(ctx, prompt); err != nil {
					errorCh <- fmt.Errorf("failed to send transcription to conversation: %w", err)
					return
				}

				fmt.Printf("[CONVERSATION] Sent transcription: %s\n", transcription)

				// Request the model to generate a response to the transcription
				responseConfig := &types.ResponseConfig{
					Modalities: []session.Modality{session.ModalityText},
				}
				if err := msgClient.SendResponseCreate(ctx, responseConfig); err != nil {
					errorCh <- fmt.Errorf("failed to request response for transcription: %w", err)
					return
				}

				fmt.Println("[CONVERSATION] Waiting for AI response...")
			}
		}
	}()

	// Process incoming messages from the conversation
	go func() {
		var messageBuffer string
		var isCollectingResponse bool

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := msgClient.ReadMessage(ctx)
				if err != nil {
					errorCh <- fmt.Errorf("error reading conversation message: %w", err)
					return
				}

				switch msg.RcvdMsgType() {
				case incoming.RcvdMsgTypeResponseTextDelta:
					if delta, ok := msg.(*incoming.ResponseTextDeltaMessage); ok {
						if !isCollectingResponse {
							fmt.Print("[CONVERSATION] AI: ")
							isCollectingResponse = true
						}
						messageBuffer += delta.Delta
						fmt.Print(delta.Delta)
					}
				case incoming.RcvdMsgTypeResponseDone:
					if isCollectingResponse {
						fmt.Println("\n[CONVERSATION] Response complete")
						fmt.Printf("[CONVERSATION] Full response: %s\n", messageBuffer)
						messageBuffer = ""
						isCollectingResponse = false
					}
				case incoming.RcvdMsgTypeResponseCreated:
					fmt.Println("[CONVERSATION] AI is generating a response...")
				case incoming.RcvdMsgTypeError:
					if errMsg, ok := msg.(*incoming.ErrorMessage); ok {
						errorCh <- fmt.Errorf("conversation error: %s", errMsg.Error.Message)
						return
					}
				// Handle other common message types without logging them as "unhandled"
				case incoming.RcvdMsgTypeSessionCreated:
				case incoming.RcvdMsgTypeRateLimitsUpdated:
				case incoming.RcvdMsgTypeConversationItemCreated:
				case incoming.RcvdMsgTypeResponseOutputItemAdded:
				case incoming.RcvdMsgTypeResponseOutputItemDone:
				case incoming.RcvdMsgTypeResponseContentPartAdded:
				case incoming.RcvdMsgTypeResponseContentPartDone:
				case incoming.RcvdMsgTypeResponseTextDone:
					// These message types are expected and can be safely ignored
				default:
					fmt.Printf("[CONVERSATION] Unhandled message type: %s\n", msg.RcvdMsgType())
				}
			}
		}
	}()

	// Wait for completion or error
	select {
	case <-ctx.Done():
		return nil
	case err := <-errorCh:
		return err
	case <-time.After(60 * time.Second):
		// Time limit for the example
		close(doneCh)
		fmt.Println("[CONVERSATION] Time limit reached, completing example")
		return nil
	}
}
