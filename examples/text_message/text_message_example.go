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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create a client
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
	fmt.Println("Creating session...")
	sessionResp, err := client.CreateSession(ctx, createReq)
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	fmt.Printf("Created session with ID: %s\n", sessionResp.ID)

	// Connect to the conversation session
	fmt.Println("Connecting to session...")
	conn, err := client.Connect(ctx,
		openaiClient.WithModel(model),
		openaiClient.WithSessionID(sessionResp.ID),
		openaiClient.WithLogger(debugLogger))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Create messaging client
	msgClient := messaging.NewClient(conn)

	// Send a simple text message
	fmt.Println("Sending text message...")
	if err := msgClient.SendText(ctx, "Hello, world!"); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	// Request the model to generate a response
	fmt.Println("Requesting model response...")
	responseConfig := &types.ResponseConfig{
		Modalities: []session.Modality{session.ModalityText},
	}
	if err := msgClient.SendResponseCreate(ctx, responseConfig); err != nil {
		log.Fatalf("Failed to request response: %v", err)
	}

	// Process responses
	fmt.Println("Waiting for response...")
	for {
		msg, err := msgClient.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("Error reading message: %v", err)
		}

		switch msg.RcvdMsgType() {
		case incoming.RcvdMsgTypeResponseTextDelta:
			if delta, ok := msg.(*incoming.ResponseTextDeltaMessage); ok {
				fmt.Print(delta.Delta)
			}
		case incoming.RcvdMsgTypeResponseDone:
			fmt.Println("\nResponse complete")
			return
		case incoming.RcvdMsgTypeError:
			if errMsg, ok := msg.(*incoming.ErrorMessage); ok {
				log.Fatalf("Error: %s", errMsg.Error.Message)
			}
		default:
			log.Printf("Unhandled message type: %s", msg.RcvdMsgType())
		}
	}
}
