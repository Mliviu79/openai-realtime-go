package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Mliviu79/openai-realtime-go/logger"
	"github.com/Mliviu79/openai-realtime-go/messages/incoming"
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Create a client
	client := openaiClient.NewClient(apiKey)

	// Configure transcription session
	inputFormat := session.AudioFormatPCM16
	transcriptionModel := session.TranscriptionModelGPT4oTranscribe
	includes := []session.TranscriptionSessionInclude{
		session.TranscriptionSessionIncludeLogprobs,
	}

	// Create a transcription session
	includeSlice := make([]session.TranscriptionSessionInclude, len(includes))
	copy(includeSlice, includes)

	createReq := &session.CreateTranscriptionSessionRequest{
		TranscriptionSessionRequest: session.TranscriptionSessionRequest{
			InputAudioFormat: &inputFormat,
			InputAudioTranscription: &session.InputAudioTranscription{
				Model: transcriptionModel,
				// Optionally set language and prompt
				Language: "en",
				Prompt:   "Technical vocabulary related to programming",
			},
			Include: &includeSlice,
		},
	}

	fmt.Println("Creating transcription session...")
	sessionResp, err := client.CreateTranscriptionSession(ctx, createReq)
	if err != nil {
		log.Fatalf("Failed to create transcription session: %v", err)
	}

	fmt.Printf("Created transcription session with ID: %s\n", sessionResp.ID)
	fmt.Printf("Client secret expires at: %d\n", sessionResp.ClientSecret.ExpiresAt)

	// Connect to the transcription session
	fmt.Println("Connecting to transcription session...")
	fmt.Printf("Session ID: %s\n", sessionResp.ID)

	// Enable debug logging for more details
	debugLoggerOpts := logger.LoggerOptions{
		Level:  zerolog.DebugLevel,
		Output: os.Stdout,
	}
	debugLogger := logger.NewZeroLogger(debugLoggerOpts)

	conn, err := client.ConnectTranscription(ctx,
		openaiClient.WithTranscriptionSessionID(sessionResp.ID),
		openaiClient.WithTranscriptionLogger(debugLogger))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Create messaging client to handle the protocol
	msgClient := messaging.NewClient(conn)

	// Update the transcription session with new settings
	fmt.Println("Updating transcription session...")

	// Create noise reduction configuration
	noiseReduction := &session.InputAudioNoiseReduction{
		Type: session.NoiseReductionTypeNearField,
	}

	// Update turn detection to use semantic VAD
	turnDetectionType := session.TurnDetectionTypeSemanticVad
	eagerness := session.EagernessLevelHigh
	turnDetection := &session.TurnDetection{
		Type:      turnDetectionType,
		Eagerness: eagerness,
	}

	// Create update request
	updateReq := session.TranscriptionSessionRequest{
		InputAudioNoiseReduction: noiseReduction,
		TurnDetection:            turnDetection,
	}

	// Send the update with a custom ID so we can track the response
	err = msgClient.SendTranscriptionSessionUpdateWithID(ctx, "update-1", updateReq)
	if err != nil {
		log.Fatalf("Failed to update transcription session: %v", err)
	}

	// Wait for confirmation of update
	fmt.Println("Waiting for update confirmation...")
	updateConfirmed := false
	timeout := time.After(10 * time.Second)

updateLoop:
	for !updateConfirmed {
		select {
		case <-timeout:
			log.Fatalf("Timed out waiting for update confirmation")
		default:
			msg, err := msgClient.ReadMessage(ctx)
			if err != nil {
				log.Fatalf("Error reading message: %v", err)
			}

			switch msg.RcvdMsgType() {
			case incoming.RcvdMsgTypeTranscriptionSessionUpdated:
				if updated, ok := msg.(*incoming.TranscriptionSessionUpdatedMessage); ok {
					fmt.Println("Session updated successfully")
					fmt.Printf("  - Noise reduction type: %v\n", updated.Session.InputAudioNoiseReduction.Type)
					fmt.Printf("  - Turn detection type: %v\n", updated.Session.TurnDetection.Type)
					updateConfirmed = true
					break updateLoop
				}
			case incoming.RcvdMsgTypeError:
				if errMsg, ok := msg.(*incoming.ErrorMessage); ok {
					fmt.Printf("Error updating session: %s\n", errMsg.Error.Message)
					return
				}
			}
		}
	}

	// Send audio data (this is a placeholder - you would read actual audio data from a file or microphone)
	// In a real application, you would stream audio chunks continuously
	fmt.Println("Sending audio data...")
	audioChunk := make([]byte, 1024) // Placeholder for audio data
	audioBase64 := base64.StdEncoding.EncodeToString(audioChunk)
	err = msgClient.SendAudio(ctx, audioBase64, "")
	if err != nil {
		log.Fatalf("Failed to send audio: %v", err)
	}

	// Alternative method: Use audio buffer append and commit
	// err = msgClient.SendAudioBufferAppend(ctx, audioBase64)
	// if err != nil {
	//     log.Fatalf("Failed to append audio: %v", err)
	// }
	// err = msgClient.SendAudioBufferCommit(ctx, "")
	// if err != nil {
	//     log.Fatalf("Failed to commit audio: %v", err)
	// }

	// Read and process messages
	fmt.Println("Waiting for transcriptions...")
	for {
		msg, err := msgClient.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("Error reading message: %v", err)
		}

		// Handle different message types
		switch msg.RcvdMsgType() {
		case incoming.RcvdMsgTypeInputAudioTranscription:
			transcription, ok := msg.(*incoming.InputAudioTranscriptionMessage)
			if ok {
				fmt.Printf("Transcription: %s\n", transcription.Text)
				if len(transcription.Logprobs) > 0 {
					fmt.Printf("Log probabilities available: %d items\n", len(transcription.Logprobs))
				}
			}
		case incoming.RcvdMsgTypeTranscriptionDone:
			fmt.Println("Transcription complete")
			return
		case incoming.RcvdMsgTypeError:
			if errMsg, ok := msg.(*incoming.ErrorMessage); ok {
				fmt.Printf("Error: %s\n", errMsg.Error.Message)
			}
			return
		}
	}
}
