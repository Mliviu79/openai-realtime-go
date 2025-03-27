// Package openaiClient provides a client for the OpenAI Realtime API.
//
// The Realtime API allows developers to build streaming multi-modal AI applications
// with voice and text conversations. This client handles WebSocket connections,
// session management, and message handling for the OpenAI Realtime API.
//
// # Overview
//
// The OpenAI Realtime API allows for bidirectional streaming communication with
// OpenAI models like GPT-4o. It supports multiple modalities including text and
// audio, enabling voice conversations, function calling, and other interactive features.
//
// This client library provides:
//   - Session management (creation, updating, retrieval)
//   - WebSocket connection handling
//   - Message parsing and serialization
//   - Support for all message types (28 incoming and 9 outgoing)
//   - Audio format handling
//   - Tool/function calling
//
// # Getting Started
//
// To use this client, you need an OpenAI API key with access to the Realtime API.
// Here's a basic example of creating a client and connecting to the API:
//
//	// Create a client with your API key
//	client := openaiClient.NewClient("your-api-key")
//
//	// Establish a WebSocket connection
//	conn, err := client.Connect(context.Background(),
//		openaiClient.WithModel(session.GPT4oRealtimePreview))
//	if err != nil {
//		log.Fatalf("Failed to connect: %v", err)
//	}
//
// # Session Management
//
// You can create and manage sessions either through the REST API or WebSocket messages:
//
//	// Create a session via REST API
//	model := session.GPT4oRealtimePreview
//	modalities := []session.Modality{session.ModalityText, session.ModalityAudio}
//
//	createReq := &session.CreateRequest{
//		SessionRequest: session.SessionRequest{
//			Model:      &model,
//			Modalities: &modalities,
//		},
//	}
//
//	sessionResp, err := client.CreateSession(context.Background(), createReq)
//	if err != nil {
//		log.Fatalf("Failed to create session: %v", err)
//	}
//
//	// Connect using the session ID
//	conn, err := client.Connect(context.Background(),
//		openaiClient.WithModel(model),
//		openaiClient.WithSessionID(sessionResp.ID))
//
// # Message Handling
//
// After establishing a connection, you can create a messaging client to handle the
// communication protocol:
//
//	// Create a messaging client
//	msgClient := messaging.NewClient(conn)
//
//	// Send a message
//	err = msgClient.SendTextMessage(context.Background(), "Hello, how are you?", nil)
//	if err != nil {
//		log.Fatalf("Failed to send message: %v", err)
//	}
//
//	// Read messages in a loop
//	for {
//		msg, err := msgClient.ReadMessage(context.Background())
//		if err != nil {
//			log.Fatalf("Error reading message: %v", err)
//			break
//		}
//
//		// Handle different message types
//		switch m := msg.(type) {
//		case *incoming.ResponseTextDeltaMessage:
//			fmt.Print(m.Delta.Text)
//		case *incoming.ResponseDoneMessage:
//			fmt.Println("\nResponse complete")
//		// Handle other message types...
//		}
//	}
//
// # Advanced Features
//
// This client supports advanced features like:
//   - Audio input/output in various formats
//   - Turn detection (VAD) for natural conversations
//   - Function calling through tools
//   - Audio transcription
//   - Temperature and token limit controls
//
// For a comprehensive example demonstrating all message types and features,
// see the examples/comprehensive_run.go file.
//
// # Documentation
//
// For detailed API documentation, refer to the OpenAI Realtime API documentation.
// For detailed usage examples, refer to the examples directory in this repository.
package openaiClient

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Mliviu79/openai-realtime-go/httpClient"
	logger "github.com/Mliviu79/openai-realtime-go/logger"
	"github.com/Mliviu79/openai-realtime-go/session"
	"github.com/Mliviu79/openai-realtime-go/ws"
)

// ConnectOption is a function that configures connection options
type ConnectOption func(*connectOptions)

// connectOptions holds the options for establishing a connection
type connectOptions struct {
	model     string        // The model to use for the connection
	logger    logger.Logger // Logger for the connection
	sessionID string        // Session ID for the connection
	readLimit int64         // Maximum size of a WebSocket message in bytes
}

// WithModel sets the model for the connection
//
// Parameters:
//   - model: The model to use (e.g., "gpt-4o")
func WithModel(model session.Model) ConnectOption {
	return func(o *connectOptions) {
		o.model = string(model)
	}
}

// WithLogger sets the logger for the connection
//
// Parameters:
//   - logger: The logger to use for the connection
func WithLogger(logger logger.Logger) ConnectOption {
	return func(o *connectOptions) {
		o.logger = logger
	}
}

// WithSessionID sets the session ID for the connection
//
// Parameters:
//   - sessionID: The session ID to use for the connection
func WithSessionID(sessionID string) ConnectOption {
	return func(o *connectOptions) {
		o.sessionID = sessionID
	}
}

// WithReadLimit sets the maximum size of a WebSocket message in bytes
//
// Parameters:
//   - readLimit: The maximum size in bytes (0 or negative means no limit)
func WithReadLimit(readLimit int64) ConnectOption {
	return func(o *connectOptions) {
		o.readLimit = readLimit
	}
}

// TranscriptionConnectOption is a function that configures transcription connection options
type TranscriptionConnectOption func(*transcriptionConnectOptions)

// transcriptionConnectOptions holds the options for establishing a transcription connection
type transcriptionConnectOptions struct {
	logger    logger.Logger // Logger for the connection
	sessionID string        // Session ID for the connection
	readLimit int64         // Maximum size of a WebSocket message in bytes
}

// WithTranscriptionLogger sets the logger for the transcription connection
//
// Parameters:
//   - logger: The logger to use for the connection
func WithTranscriptionLogger(logger logger.Logger) TranscriptionConnectOption {
	return func(o *transcriptionConnectOptions) {
		o.logger = logger
	}
}

// WithTranscriptionSessionID sets the session ID for the transcription connection
//
// Parameters:
//   - sessionID: The session ID to use for the connection
func WithTranscriptionSessionID(sessionID string) TranscriptionConnectOption {
	return func(o *transcriptionConnectOptions) {
		o.sessionID = sessionID
	}
}

// WithTranscriptionReadLimit sets the maximum size of a WebSocket message in bytes
//
// Parameters:
//   - readLimit: The maximum size in bytes (0 or negative means no limit)
func WithTranscriptionReadLimit(readLimit int64) TranscriptionConnectOption {
	return func(o *transcriptionConnectOptions) {
		o.readLimit = readLimit
	}
}

// Client is OpenAI Realtime API client
type Client struct {
	config httpClient.ClientConfig
}

// NewClient creates new OpenAI Realtime API client with the given auth token
//
// Parameters:
//   - authToken: The authentication token for the OpenAI API
//
// Returns:
//   - *Client: A new OpenAI Realtime API client
func NewClient(authToken string) *Client {
	config := httpClient.DefaultConfig(authToken)
	return NewClientWithConfig(config)
}

// NewClientWithConfig creates new OpenAI Realtime API client with specified config
//
// Parameters:
//   - config: The client configuration
//
// Returns:
//   - *Client: A new OpenAI Realtime API client
func NewClientWithConfig(config httpClient.ClientConfig) *Client {
	return &Client{
		config: config,
	}
}

// CreateSession creates a new session
//
// Parameters:
//   - ctx: The context for the request
//   - req: The session creation request
//
// Returns:
//   - *session.CreateResponse: The session creation response
//   - error: An error if the request failed
func (c *Client) CreateSession(ctx context.Context, req *session.CreateRequest) (*session.CreateResponse, error) {
	return httpClient.Do[session.CreateRequest, session.CreateResponse](
		ctx,
		c.config.APIBaseURL+"/realtime/sessions",
		req,
		httpClient.WithHeaders(httpClient.GetHeaders(c.config)),
		httpClient.WithClient(c.config.HTTPClient),
	)
}

// CreateTranscriptionSession creates a new transcription session
//
// Parameters:
//   - ctx: The context for the request
//   - req: The transcription session creation request
//
// Returns:
//   - *session.CreateTranscriptionSessionResponse: The transcription session creation response
//   - error: An error if the request failed
func (c *Client) CreateTranscriptionSession(ctx context.Context, req *session.CreateTranscriptionSessionRequest) (*session.CreateTranscriptionSessionResponse, error) {
	return httpClient.Do[session.CreateTranscriptionSessionRequest, session.CreateTranscriptionSessionResponse](
		ctx,
		c.config.APIBaseURL+"/realtime/transcription_sessions",
		req,
		httpClient.WithHeaders(httpClient.GetHeaders(c.config)),
		httpClient.WithClient(c.config.HTTPClient),
	)
}

// Connect establishes a WebSocket connection to the OpenAI Realtime API for model-based conversations
//
// Parameters:
//   - ctx: The context for the connection
//   - opts: Options for the connection
//
// Returns:
//   - *ws.Conn: The WebSocket connection
//   - error: An error if the connection failed
func (c *Client) Connect(ctx context.Context, opts ...ConnectOption) (*ws.Conn, error) {
	options := &connectOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if options.model == "" {
		return nil, fmt.Errorf("model is required")
	}

	// Create dialer with custom read limit if specified
	dialer := ws.DirectDialer(ws.DialerOptions{
		ReadLimit: options.readLimit,
	})

	// Construct URL with query parameters
	query := url.Values{}
	query.Set("model", options.model)
	if options.sessionID != "" {
		query.Set("session_id", options.sessionID)
	}

	// Set the base URL
	baseURL := c.config.BaseURL
	url := baseURL + "?" + query.Encode()

	headers := httpClient.GetHeaders(c.config)

	wsConn, err := dialer.Dial(ctx, url, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OpenAI: %w", err)
	}

	conn := ws.NewConn(wsConn)
	if options.logger != nil {
		conn.SetLogger(options.logger)
	}

	return conn, nil
}

// ConnectTranscription establishes a WebSocket connection to the OpenAI Realtime API for transcription
//
// Parameters:
//   - ctx: The context for the connection
//   - opts: Options for the connection
//
// Returns:
//   - *ws.Conn: The WebSocket connection
//   - error: An error if the connection failed
func (c *Client) ConnectTranscription(ctx context.Context, opts ...TranscriptionConnectOption) (*ws.Conn, error) {
	options := &transcriptionConnectOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Create dialer with custom read limit if specified
	dialer := ws.DirectDialer(ws.DialerOptions{
		ReadLimit: options.readLimit,
	})

	// Construct URL with query parameters
	query := url.Values{}
	query.Set("intent", "transcription")
	if options.sessionID != "" {
		query.Set("session_id", options.sessionID)
	}

	// Set the base URL
	baseURL := c.config.BaseURL
	url := baseURL + "?" + query.Encode()

	headers := httpClient.GetHeaders(c.config)

	wsConn, err := dialer.Dial(ctx, url, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OpenAI transcription service: %w", err)
	}

	conn := ws.NewConn(wsConn)
	if options.logger != nil {
		conn.SetLogger(options.logger)
	}

	return conn, nil
}
