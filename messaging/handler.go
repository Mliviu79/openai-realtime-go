package messaging

import (
	"context"

	"github.com/Mliviu79/go-openai-realtime/logger"
	"github.com/Mliviu79/go-openai-realtime/messages/incoming"
	"github.com/Mliviu79/go-openai-realtime/ws"
)

// MessageHandler is a function that processes an incoming OpenAI message
type MessageHandler func(ctx context.Context, event incoming.RcvdMsg)

// Handler handles incoming OpenAI messages from a WebSocket connection.
// It reads messages in a standalone goroutine and calls the registered handlers.
// It is the responsibility of the caller to call Start and Stop.
type Handler struct {
	ctx       context.Context
	cancel    context.CancelFunc
	client    *Client
	wsHandler *ws.ConnHandler
	handlers  []MessageHandler
	logger    logger.Logger
	errCh     chan error
}

// NewHandler creates a new Handler for the OpenAI Realtime API.
func NewHandler(parentCtx context.Context, client *Client, handlers ...MessageHandler) *Handler {
	if client == nil {
		panic("client cannot be nil")
	}

	ctx, cancel := context.WithCancel(parentCtx)

	h := &Handler{
		ctx:      ctx,
		cancel:   cancel,
		client:   client,
		handlers: handlers,
		errCh:    make(chan error, 1),
	}

	// Create a WebSocket handler that will decode raw messages into OpenAI messages
	wsHandler := ws.NewConnHandler(ctx, client.conn, h.handleRawMessage)
	h.wsHandler = wsHandler

	return h
}

// SetLogger sets the logger for the handler
func (h *Handler) SetLogger(logger logger.Logger) {
	h.logger = logger
}

// Start starts the handler.
func (h *Handler) Start() {
	if h.logger != nil {
		h.logger.Debugf("Starting message handler")
	}
	h.wsHandler.Start()
}

// Err returns a channel that receives errors from the handler.
func (h *Handler) Err() <-chan error {
	return h.errCh
}

// AddHandler adds a message handler.
// This is safe to call before Start() but not after.
func (h *Handler) AddHandler(handler MessageHandler) {
	if handler == nil {
		if h.logger != nil {
			h.logger.Warnf("Attempted to add nil handler, ignoring")
		}
		return
	}
	h.handlers = append(h.handlers, handler)
}

// Stop gracefully stops the handler by canceling its context.
func (h *Handler) Stop() {
	if h.logger != nil {
		h.logger.Debugf("Stopping message handler")
	}
	h.wsHandler.Stop()
	if h.cancel != nil {
		h.cancel()
	}
}

// handleRawMessage is called by the WebSocket handler when a raw message is received.
// It decodes the raw message into an OpenAI message and calls the handlers.
func (h *Handler) handleRawMessage(ctx context.Context, messageType ws.MessageType, data []byte) {
	// We only handle text messages
	if messageType != ws.MessageText {
		if h.logger != nil {
			h.logger.Warnf("Received non-text message: %s", messageType.String())
		}
		return
	}

	// Decode the message
	msg, err := incoming.UnmarshalRcvdMsg(data)
	if err != nil {
		if h.logger != nil {
			h.logger.Errorf("Failed to unmarshal message: %v", err)
		}
		return
	}

	if h.logger != nil {
		h.logger.Debugf("Received message of type: %s", msg.RcvdMsgType())
	}

	// Call the handlers
	for i, handler := range h.handlers {
		if handler == nil {
			if h.logger != nil {
				h.logger.Warnf("Skipping nil handler at index %d", i)
			}
			continue
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					if h.logger != nil {
						h.logger.Errorf("Handler %d panicked: %v", i, r)
					}
				}
			}()
			handler(ctx, msg)
		}()
	}
}
