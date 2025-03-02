# Go OpenAI Realtime API Client

[![Go Reference](https://pkg.go.dev/badge/github.com/Mliviu79/openai-realtime-go.svg)](https://pkg.go.dev/github.com/Mliviu79/openai-realtime-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/Mliviu79/openai-realtime-go)](https://goreportcard.com/report/github.com/Mliviu79/openai-realtime-go)
[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/Mliviu79/openai-realtime-go/blob/main/LICENSE)

A fully-featured Go client for the OpenAI Realtime API, supporting multi-modal conversations with text and audio. This project is a heavily refactored and restructured fork of [WqyJh/go-openai-realtime](https://github.com/WqyJh/go-openai-realtime).

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
	err = msgClient.SendTextMessage(ctx, "Tell me about the OpenAI Realtime API", nil)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
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
			fmt.Print(m.Delta.Text)
		case *incoming.ResponseDoneMessage:
			// End of the response
			fmt.Println("\nResponse complete")
			return
		}
	}
}
```

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

## Comprehensive Example

For a complete demonstration of all message types and features, see the [comprehensive example](examples/comprehensive_example.go).

This example shows:

- All 28 incoming message types
- All 9 outgoing message types
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
