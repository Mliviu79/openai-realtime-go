// Package main provides an example demonstrating how to send pretend transcriptions
// directly to a conversation session, without using a real transcription service.
//
// This example simulates the second half of a voice assistant application where
// you already have transcribed text and just need to send it to a conversation model.
//
// To run this example:
//  1. Set the OPENAI_API_KEY environment variable
//  2. Run: go run examples/simulated_transcription/simulated_transcription_example.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Mliviu79/openai-realtime-go/logger"
	"github.com/Mliviu79/openai-realtime-go/messages/incoming"
	"github.com/Mliviu79/openai-realtime-go/messages/types"
	"github.com/Mliviu79/openai-realtime-go/messaging"
	"github.com/Mliviu79/openai-realtime-go/openaiClient"
	"github.com/Mliviu79/openai-realtime-go/session"
	"github.com/rs/zerolog"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create a client with API key
	client := openaiClient.NewClient(apiKey)

	// Create logger
	loggerOpts := logger.LoggerOptions{
		Level:  zerolog.InfoLevel,
		Output: os.Stdout,
	}
	debugLogger := logger.NewZeroLogger(loggerOpts)

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
	fmt.Println("Creating conversation session...")
	sessionResp, err := client.CreateSession(ctx, createReq)
	if err != nil {
		log.Fatalf("Failed to create conversation session: %v", err)
	}

	fmt.Printf("Created conversation session with ID: %s\n", sessionResp.ID)

	// Connect to the conversation session
	fmt.Println("Connecting to session...")
	conn, err := client.Connect(ctx,
		openaiClient.WithModel(model),
		openaiClient.WithSessionID(sessionResp.ID),
		openaiClient.WithLogger(debugLogger))
	if err != nil {
		log.Fatalf("Failed to connect to conversation session: %v", err)
	}

	// Create messaging client
	msgClient := messaging.NewClient(conn)

	// Send initial message
	initialMsg := "Hello, I'm going to send you some audio transcriptions. Please respond to them as they come in."
	if err := msgClient.SendText(ctx, initialMsg); err != nil {
		log.Fatalf("Failed to send initial message: %v", err)
	}

	fmt.Printf("Sent: %s\n", initialMsg)

	// Request the model to generate a response to the initial message
	fmt.Println("Requesting initial response...")
	responseConfig := &types.ResponseConfig{
		Modalities: modalities,
	}
	if err := msgClient.SendResponseCreate(ctx, responseConfig); err != nil {
		log.Fatalf("Failed to request initial response: %v", err)
	}

	// Wait for and process the initial response
	printResponse(ctx, msgClient)

	// These are simulated transcriptions - pretend they came from a transcription service
	simulatedTranscriptions := []string{
		"This is a test of the transcription system.",
		"Can you explain how machine learning works?",
		"What are the main differences between supervised and unsupervised learning?",
	}

	// Send each simulated transcription and get responses
	for i, transcription := range simulatedTranscriptions {
		// Wait a bit between messages
		time.Sleep(2 * time.Second)

		fmt.Printf("\nSending simulated transcription %d: '%s'\n", i+1, transcription)

		// Format the message as if it came from a transcription service
		prompt := fmt.Sprintf("I just transcribed this audio: \"%s\". Please respond to it.", transcription)

		if err := msgClient.SendText(ctx, prompt); err != nil {
			log.Fatalf("Failed to send transcription: %v", err)
		}

		// Request the model to generate a response
		fmt.Println("Requesting response...")
		if err := msgClient.SendResponseCreate(ctx, responseConfig); err != nil {
			log.Fatalf("Failed to request response: %v", err)
		}

		// Wait for and process the response
		printResponse(ctx, msgClient)
	}

	fmt.Println("\nExample completed successfully")
}

// printResponse reads and prints the model's response
func printResponse(ctx context.Context, msgClient *messaging.Client) {
	var messageBuffer string
	responseDone := false

	fmt.Println("Waiting for AI response...")

	// Loop until we get a response.done message
	for !responseDone {
		msg, err := msgClient.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("Error reading message: %v", err)
		}

		switch msg.RcvdMsgType() {
		case incoming.RcvdMsgTypeResponseCreated:
			fmt.Println("AI is generating a response...")

		case incoming.RcvdMsgTypeResponseOutputTextDelta:
			if delta, ok := msg.(*incoming.ResponseOutputTextDeltaMessage); ok {
				messageBuffer += delta.Delta
				fmt.Print(delta.Delta)
			}

		case incoming.RcvdMsgTypeResponseDone:
			fmt.Println("\nResponse complete")
			fmt.Printf("Full response: %s\n", messageBuffer)
			messageBuffer = ""
			responseDone = true

		case incoming.RcvdMsgTypeError:
			if errMsg, ok := msg.(*incoming.ErrorMessage); ok {
				log.Fatalf("Error: %s", errMsg.Error.Message)
			}

		// These message types are normal and can be ignored for this simple example
		case incoming.RcvdMsgTypeSessionCreated:
		case incoming.RcvdMsgTypeRateLimitsUpdated:
		case incoming.RcvdMsgTypeConversationItemCreated:
		case incoming.RcvdMsgTypeResponseOutputItemAdded:
		case incoming.RcvdMsgTypeResponseOutputItemDone:
		case incoming.RcvdMsgTypeResponseContentPartAdded:
		case incoming.RcvdMsgTypeResponseContentPartDone:
		case incoming.RcvdMsgTypeResponseOutputTextDone:

		default:
			fmt.Printf("Unhandled message type: %s\n", msg.RcvdMsgType())
		}
	}
}
