# Go OpenAI Realtime API Client

[![Go Reference](https://pkg.go.dev/badge/github.com/Mliviu79/openai-realtime-go.svg)](https://pkg.go.dev/github.com/Mliviu79/openai-realtime-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Mliviu79/openai-realtime-go)](https://goreportcard.com/report/github.com/Mliviu79/openai-realtime-go)
[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/Mliviu79/openai-realtime-go/blob/main/LICENSE)

A fully-featured Go client for the OpenAI Realtime API, supporting multi-modal conversations with text and audio. This project is a heavily refactored and restructured fork of [WqyJh/go-openai-realtime](https://github.com/WqyJh/go-openai-realtime).

## Latest Release: v0.1.3

This release includes several important improvements:

- **Fixed API Compatibility**: Updated the `SendText` method to use `input_text` instead of `text` content type to properly match OpenAI's API expectations
- **New Examples**: Added comprehensive examples demonstrating practical use cases, including:
  - **Dual Session Example**: Shows how to run both transcription and conversation sessions simultaneously
  - **Simulated Transcription Example**: Demonstrates how to process text as if it were transcribed audio
- **Better Documentation**: Added detailed READMEs for all examples
- **Bug Fixes**: Fixed various issues and linter errors

For a complete list of changes, see the [CHANGES.md](CHANGES.md) file.

## Overview

This client allows you to integrate with OpenAI's Realtime API to build applications that can have natural, streaming conversations with OpenAI models like GPT-4o. The Realtime API enables:

- **Multi-modal conversations** - Support for both text and audio modalities
- **Bidirectional streaming** - Real-time streaming of both user inputs and model outputs
- **Voice interactions** - Voice input/output with different voice options and formats
- **Function calling** - Define and call tools/functions from the model
- **Turn detection** - Voice Activity Detection (VAD) for natural conversation turns

This client library provides a complete implementation of all 28 incoming and 9 outgoing message types supported by the Realtime API.

## Key Differences from Original

This fork has been extensively refactored to provide:

- More modular package structure with clear separation of concerns
- Comprehensive godoc documentation for every package and type
- Type-safe API with well-defined interfaces
- Improved error handling and logging
- Expanded examples for various use cases

## Supported Models

- `gpt-4o-realtime-preview`
- `gpt-4o-realtime-preview-2024-10-01`
- `gpt-4o-realtime-preview-2024-12-17`
- `gpt-4o-mini-realtime-preview`
- `gpt-4o-mini-realtime-preview-2024-12-17`

## Installation

```bash
go get github.com/Mliviu79/openai-realtime-go
```

This library requires Go version 1.23 or greater.

## Quick Start

Here's a simple example to get started:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Mliviu79/openai-realtime-go/messages/incoming"
	"github.com/Mliviu79/openai-realtime-go/messages/types"
	"github.com/Mliviu79/openai-realtime-go/messaging"
	"github.com/Mliviu79/openai-realtime-go/openaiClient"
	"github.com/Mliviu79/openai-realtime-go/session"
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

	// Connect to the API with the specified model
	conn, err := client.Connect(ctx, 
		openaiClient.WithModel(session.GPT4oRealtimePreview))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// Create messaging client to handle the protocol
	msgClient := messaging.NewClient(conn)

	// Send a text message
	err = msgClient.SendText(ctx, "Tell me about the OpenAI Realtime API")
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	
	// Request the model to generate a response
	// This step is required - without it, no response will be generated
	responseConfig := &types.ResponseConfig{
		Modalities: []session.Modality{session.ModalityText},
	}
	err = msgClient.SendResponseCreate(ctx, responseConfig)
	if err != nil {
		log.Fatalf("Failed to request response: %v", err)
	}

	// Read and process messages
	fmt.Println("Response:")
	for {
		msg, err := msgClient.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("Error reading message: %v", err)
		}

		// Handle different message types
		switch m := msg.(type) {
		case *incoming.ResponseTextDeltaMessage:
			// Print text deltas as they arrive
			fmt.Print(m.Delta)
		case *incoming.ResponseDoneMessage:
			// End of the response
			fmt.Println("\nResponse complete")
			return
		}
	}
}
```

## Examples

The library includes several example applications to demonstrate different use cases:

### Text Messaging Example

A simple example showing how to send text messages and handle responses.

```bash
cd examples/text_message
go run text_message_example.go
```

### Transcription Example

Demonstrates how to use the transcription API to convert audio to text.

```bash
cd examples/transcription_example
go run main.go
```

### Simulated Transcription Example

Shows how to send simulated audio transcriptions to a conversation session and get responses.

```bash
cd examples/simulated_transcription
go run simulated_transcription_example.go
```

### Dual Session Example

Advanced example showing how to run both transcription and conversation sessions simultaneously, mimicking a complete voice assistant pipeline.

```bash
cd examples/dual_session
go run dual_session_example.go
```

## Important: Two-Step Process for Getting Responses

The OpenAI Realtime API uses a two-step process to get a response from the model:

1. **Send your message(s)** using methods like `SendText`, `SendAudio`, etc.
2. **Request a response** using `SendResponseCreate` with a configuration object

This design enables powerful capabilities:
- You can send multiple messages before requesting a response
- You can configure exactly how the model should respond (modalities, temperature, etc.)
- You can control precisely when the model should start generating a response

**Without the second step, the model will not generate any response to your messages.**

The `SendResponseCreate` method takes a `ResponseConfig` object that can include:
- Which modalities to use for the response (text, audio)
- Voice options for audio responses
- Temperature and other generation parameters
- Tools the model can use
- Optional system instructions

## Session Management

You can create and manage sessions either through the REST API or WebSocket messages:

```go
// Create a session via REST API
model := session.GPT4oRealtimePreview
modalities := []session.Modality{session.ModalityText, session.ModalityAudio}

createReq := &session.CreateRequest{
    SessionRequest: session.SessionRequest{
        Model:      &model,
        Modalities: &modalities,
    },
}

sessionResp, err := client.CreateSession(ctx, createReq)
if err != nil {
    log.Fatalf("Failed to create session: %v", err)
}

// Connect using the session ID
conn, err := client.Connect(ctx,
    openaiClient.WithModel(model),
    openaiClient.WithSessionID(sessionResp.ID))
```

## Advanced Features

### Audio Input/Output

The client supports sending and receiving audio in various formats:

```go
// Configure audio formats in the session
inputFormat := session.AudioFormatPCM16
outputFormat := session.AudioFormatPCM16
voice := session.VoiceAlloy

createReq := &session.CreateRequest{
    SessionRequest: session.SessionRequest{
        Model:             &model,
        Modalities:        &[]session.Modality{session.ModalityText, session.ModalityAudio},
        InputAudioFormat:  &inputFormat,
        OutputAudioFormat: &outputFormat,
        Voice:             &voice,
    },
}

// Send audio data
audioData := []byte{...} // Your PCM16 audio data
err = msgClient.SendAudioMessage(ctx, audioData, nil)
```

### Function Calling

You can define functions for the model to call:

```go
// Define a tool/function
getWeatherParams := map[string]interface{}{
    "type": "object",
    "properties": map[string]interface{}{
        "location": map[string]interface{}{
            "type": "string",
            "description": "The city and state, e.g. San Francisco, CA",
        },
    },
    "required": []string{"location"},
}

tools := []session.Tool{
    {
        Type: "function",
        Function: &session.Function{
            Name:        "get_weather",
            Description: "Get the current weather in a given location",
            Parameters:  getWeatherParams,
        },
    },
}

// Add tools to the session
createReq := &session.CreateRequest{
    SessionRequest: session.SessionRequest{
        Model: &model,
        Tools: &tools,
    },
}
```

### Turn Detection

Enable Voice Activity Detection (VAD) for natural conversation turns:

```go
// Configure turn detection
turnDetectionType := "server_vad"
threshold := 0.5
prefixPaddingMs := 300
silenceDurationMs := 500
createResponse := true
interruptResponse := true

turnDetection := &session.TurnDetection{
    Type:              &turnDetectionType,
    Threshold:         &threshold,
    PrefixPaddingMs:   &prefixPaddingMs,
    SilenceDurationMs: &silenceDurationMs,
    CreateResponse:    &createResponse,
    InterruptResponse: &interruptResponse,
}

createReq := &session.CreateRequest{
    SessionRequest: session.SessionRequest{
        Model:         &model,
        TurnDetection: turnDetection,
    },
}
```

### Transcription Sessions

The library supports OpenAI's dedicated transcription sessions for real-time speech-to-text:

```go
// Configure transcription session
inputFormat := session.AudioFormatPCM16
transcriptionModel := session.TranscriptionModelGPT4oTranscribe

// Optional: request log probabilities 
includes := []session.TranscriptionSessionInclude{
    session.TranscriptionSessionIncludeLogprobs,
}
includeSlice := make([]session.TranscriptionSessionInclude, len(includes))
copy(includeSlice, includes)

// Create a transcription session
createReq := &session.CreateTranscriptionSessionRequest{
    TranscriptionSessionRequest: session.TranscriptionSessionRequest{
        InputAudioFormat: &inputFormat,
        InputAudioTranscription: &session.InputAudioTranscription{
            Model: transcriptionModel,
            Language: "en",  // Optional language hint
            Prompt: "Technical vocabulary",  // Optional domain hint
        },
        Include: &includeSlice,  // Optional: include log probabilities
    },
}

// Create the session via the API
sessionResp, err := client.CreateTranscriptionSession(ctx, createReq)
if err != nil {
    log.Fatalf("Failed to create transcription session: %v", err)
}

// Connect to the transcription session
conn, err := client.Connect(ctx,
    openaiClient.WithModel(session.GPT4oRealtimePreview),
    openaiClient.WithSessionID(sessionResp.ID),
    openaiClient.WithTranscriptionSession())  // Special flag for transcription sessions

// Update the transcription session with new settings while connected
noiseReduction := &session.InputAudioNoiseReduction{
    Type: session.NoiseReductionTypeNearField,
}

turnDetection := &session.TurnDetection{
    Type:      session.TurnDetectionTypeSemanticVad,
    Eagerness: session.EagernessLevelHigh,
}

updateReq := session.TranscriptionSessionRequest{
    InputAudioNoiseReduction: noiseReduction,
    TurnDetection:            turnDetection,
}

// Send the update
err = msgClient.SendTranscriptionSessionUpdate(ctx, updateReq)
if err != nil {
    log.Fatalf("Failed to update transcription session: %v", err)
}

// Now send audio and receive transcriptions
// ...
```

See the [transcription example](examples/transcription_example/main.go) for a complete demonstration.

## Comprehensive Example

For a complete demonstration of all message types and features, see the [comprehensive example](examples/comprehensive_example.go).

This example shows:

- All 29 incoming message types
- All 10 outgoing message types
- Audio handling
- Function calling
- Turn detection
- Error handling

## Package Structure

Our library is organized into several packages with clear separation of concerns:

- `openaiClient` - Main client package for API connection
- `session` - Session management and configuration
- `messages` - Message types and handling
  - `incoming` - Incoming message types
  - `outgoing` - Outgoing message types
  - `types` - Shared message type definitions
  - `factory` - Message factory functions
- `messaging` - High-level messaging interface
- `ws` - WebSocket connection management
- `httpClient` - HTTP client for REST API endpoints
- `logger` - Logging utilities
- `apierrs` - Error handling

## Documentation

- [GoDoc](https://pkg.go.dev/github.com/Mliviu79/openai-realtime-go)
- [OpenAI Realtime API Documentation](https://platform.openai.com/docs/api-reference/realtime)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 
