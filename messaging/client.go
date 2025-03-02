// Package messaging provides high-level messaging functionality for the OpenAI Realtime API.
// It wraps the WebSocket connection and provides methods for sending and receiving
// various types of messages, including conversation items, audio data, and responses.
//
// The messaging package is the primary interface for interacting with the OpenAI Realtime API.
// It provides:
//   - Simple methods for sending different types of messages (text, audio, functions)
//   - Message serialization and deserialization
//   - Conversation management
//   - Response handling
//   - Audio buffer operations
//
// This package abstracts away the low-level details of the WebSocket protocol and
// JSON message formats, providing a clean, type-safe API for OpenAI Realtime operations.
//
// Most users will create a Client using the openai package's Connect method, which handles
// authentication and connection establishment. The resulting Client can then be used to
// send and receive messages.
//
// Example:
//
//	// Create a client with your OpenAI API key
//	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
//
//	// Connect to the OpenAI Realtime API
//	conn, err := client.Connect(ctx, openai.WithModel("gpt-4o"))
//	if err != nil {
//		log.Fatalf("Failed to connect: %v", err)
//	}
//	defer conn.Close()
//
//	// Create a messaging client
//	msgClient := messaging.NewClient(conn)
//
//	// Send a text message
//	err = msgClient.SendText(ctx, "Hello, how are you?")
//
//	// Read the response
//	for {
//		msg, err := msgClient.ReadMessage(ctx)
//		// Process different message types...
//	}
//
// The messaging package works with the types defined in the messages/types package
// and uses the factory functions from messages/factory to create properly formatted messages.
package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Mliviu79/openai-realtime-go/logger"
	"github.com/Mliviu79/openai-realtime-go/messages/factory"
	"github.com/Mliviu79/openai-realtime-go/messages/incoming"
	"github.com/Mliviu79/openai-realtime-go/messages/outgoing"
	"github.com/Mliviu79/openai-realtime-go/messages/types"
	"github.com/Mliviu79/openai-realtime-go/session"
	"github.com/Mliviu79/openai-realtime-go/ws"
)

// Client is a client for the OpenAI Realtime API that handles message serialization/deserialization.
// It provides high-level methods for sending different types of messages and processing responses.
// All methods are thread-safe and can be called from multiple goroutines.
type Client struct {
	mu     sync.RWMutex
	conn   *ws.Conn
	logger logger.Logger
}

// NewClient creates a new messaging client that wraps a WebSocket connection.
// The client provides high-level methods for sending and receiving messages.
//
// Parameters:
//   - conn: A WebSocket connection wrapper (usually obtained from openai.Connect)
//
// Returns:
//   - A new Client instance that can be used to send and receive messages
func NewClient(conn *ws.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

// SetLogger sets the logger for the client.
// The logger is used to log message operations for debugging purposes.
// If nil, no logging is performed.
func (c *Client) SetLogger(logger logger.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
	// Also set the logger on the underlying connection
	c.conn.SetLogger(logger)
}

// Close closes the underlying connection.
// After closing, no more messages can be sent or received.
// This method is thread-safe and can be called from any goroutine.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Ping sends a ping to the server to keep the connection alive.
// This can be useful for long-lived connections to prevent timeouts.
// This method is thread-safe and can be called from any goroutine.
func (c *Client) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

// SendMessage sends a message to the server.
// This is a low-level method that takes any message implementing the OutMsg interface.
// Most users should use higher-level methods like SendText, SendAudio, etc.
//
// Parameters:
//   - ctx: A context for cancellation and timeouts
//   - msg: The message to send, must implement outgoing.OutMsg
//
// Returns:
//   - An error if the message could not be sent
func (c *Client) SendMessage(ctx context.Context, msg outgoing.OutMsg) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if c.logger != nil {
		c.logger.Debugf("sending message: type=%s data=%s", msg.OutMsgType(), string(data))
	}

	return c.conn.SendRaw(ctx, ws.MessageText, data)
}

// ReadMessage reads a message from the server.
// This method blocks until a message is received, the context is canceled, or an error occurs.
// The returned message is automatically deserialized into the appropriate Go type.
//
// Parameters:
//   - ctx: A context for cancellation and timeouts
//
// Returns:
//   - A message implementing the incoming.RcvdMsg interface
//   - An error if the message could not be read or deserialized
func (c *Client) ReadMessage(ctx context.Context) (incoming.RcvdMsg, error) {
	messageType, data, err := c.conn.ReadRaw(ctx)
	if err != nil {
		return nil, err
	}

	if messageType != ws.MessageText {
		return nil, fmt.Errorf("expected text message, got %s", messageType.String())
	}

	msg, err := incoming.UnmarshalRcvdMsg(data)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Convenience methods for sending specific types of messages

// SendSessionUpdate sends a session update message.
func (c *Client) SendSessionUpdate(ctx context.Context, sessionReq session.SessionRequest) error {
	msg := outgoing.NewSessionUpdateMessage(sessionReq)
	return c.SendMessage(ctx, msg)
}

// SendAudioBufferAppend sends an audio buffer append message.
func (c *Client) SendAudioBufferAppend(ctx context.Context, audioData string) error {
	msg := outgoing.NewAudioBufferAppendMessage(audioData, nil)
	return c.SendMessage(ctx, msg)
}

// SendAudioBufferCommit sends an audio buffer commit message.
func (c *Client) SendAudioBufferCommit(ctx context.Context, previousItemID string) error {
	msg := outgoing.NewAudioBufferCommitMessage(previousItemID)
	return c.SendMessage(ctx, msg)
}

// SendAudioBufferClear sends an audio buffer clear message.
func (c *Client) SendAudioBufferClear(ctx context.Context) error {
	msg := outgoing.NewAudioBufferClearMessage()
	return c.SendMessage(ctx, msg)
}

// SendConversationItemCreate sends a conversation item create message.
func (c *Client) SendConversationItemCreate(ctx context.Context, item *types.MessageItem, previousItemID *string) error {
	prevID := ""
	if previousItemID != nil {
		prevID = *previousItemID
	}
	msg := outgoing.NewConversationCreateMessage(prevID, *item)
	return c.SendMessage(ctx, msg)
}

// SendResponseCreate sends a response create message.
func (c *Client) SendResponseCreate(ctx context.Context, config *types.ResponseConfig) error {
	if config == nil {
		return fmt.Errorf("response config cannot be nil")
	}
	msg := outgoing.NewResponseCreateMessage(*config)
	return c.SendMessage(ctx, msg)
}

// SendResponseCancel sends a response cancel message.
func (c *Client) SendResponseCancel(ctx context.Context, responseID string) error {
	msg := outgoing.NewResponseCancelMessage(responseID)
	return c.SendMessage(ctx, msg)
}

// SendText sends a text message from the user.
func (c *Client) SendText(ctx context.Context, text string) error {
	content := []types.MessageContentPart{
		factory.TextContent(text),
	}
	item := factory.MessageItem(types.MessageRoleUser, content)
	return c.SendConversationItemCreate(ctx, &item, nil)
}

// SendAudio sends an audio message from the user.
func (c *Client) SendAudio(ctx context.Context, audioBase64 string, transcript string) error {
	content := []types.MessageContentPart{
		factory.InputAudioContent(audioBase64, transcript),
	}
	item := factory.MessageItem(types.MessageRoleUser, content)
	return c.SendConversationItemCreate(ctx, &item, nil)
}

// SendSystemMessage sends a system message.
func (c *Client) SendSystemMessage(ctx context.Context, text string) error {
	content := []types.MessageContentPart{
		factory.TextContent(text),
	}
	item := factory.MessageItem(types.MessageRoleSystem, content)
	return c.SendConversationItemCreate(ctx, &item, nil)
}

// SendConversationItemTruncate sends a conversation item truncate message.
// This truncates the conversation history to the specified index.
func (c *Client) SendConversationItemTruncate(ctx context.Context, itemID string, contentIndex int, audioEndMs int) error {
	msg := outgoing.NewConversationTruncateMessage(itemID, contentIndex, audioEndMs)
	return c.SendMessage(ctx, msg)
}

// SendConversationItemDelete sends a conversation item delete message.
// This deletes the conversation item with the specified ID.
func (c *Client) SendConversationItemDelete(ctx context.Context, itemID string) error {
	msg := outgoing.NewConversationDeleteMessage(itemID)
	return c.SendMessage(ctx, msg)
}
