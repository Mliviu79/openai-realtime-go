package ws

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/Mliviu79/openai-realtime-go/apierrs"
)

// RawMessageHandler is a function that processes raw WebSocket messages
type RawMessageHandler func(ctx context.Context, messageType MessageType, data []byte)

// ConnHandler is a handler for a WebSocket connection.
// It reads messages from the server in a standalone goroutine and calls the registered handlers.
// It is the responsibility of the caller to call Start and Stop.
// The handlers are called in the order they are registered.
// Users should not call ReadRaw directly when using ConnHandler.
type ConnHandler struct {
	ctx      context.Context
	cancel   context.CancelFunc
	conn     *Conn
	handlers []RawMessageHandler
	errCh    chan error
}

// NewConnHandler creates a new ConnHandler with the given context and connection.
func NewConnHandler(parentCtx context.Context, conn *Conn, handlers ...RawMessageHandler) *ConnHandler {
	if conn == nil {
		panic("conn cannot be nil")
	}

	ctx, cancel := context.WithCancel(parentCtx)

	return &ConnHandler{
		ctx:      ctx,
		cancel:   cancel,
		conn:     conn,
		handlers: handlers,
		errCh:    make(chan error, 1),
	}
}

// Start starts the ConnHandler.
func (c *ConnHandler) Start() {
	if c.conn.logger != nil {
		c.conn.logger.Debugf("Starting connection handler")
	}
	go func() {
		err := c.run()
		if err != nil {
			if c.conn.logger != nil {
				c.conn.logger.Errorf("Connection handler exited with error: %v", err)
			}
			c.errCh <- err
		} else {
			if c.conn.logger != nil {
				c.conn.logger.Debugf("Connection handler exited without error")
			}
		}
		close(c.errCh)
	}()
}

// Err returns a channel that receives errors from the ConnHandler.
// This could be used to wait for the goroutine to exit.
// If you don't need to wait for the goroutine to exit, there's no need to call this.
// This must be called after the connection is closed, otherwise it will block indefinitely.
func (c *ConnHandler) Err() <-chan error {
	return c.errCh
}

// AddHandler adds a message handler to the ConnHandler.
// This is safe to call before Start() but not after.
func (c *ConnHandler) AddHandler(handler RawMessageHandler) {
	if handler == nil {
		if c.conn.logger != nil {
			c.conn.logger.Warnf("Attempted to add nil handler, ignoring")
		}
		return
	}
	c.handlers = append(c.handlers, handler)
}

// Stop gracefully stops the ConnHandler by canceling its context.
func (c *ConnHandler) Stop() {
	if c.conn.logger != nil {
		c.conn.logger.Debugf("Stopping connection handler")
	}
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *ConnHandler) run() error {
	if c.conn.logger != nil {
		c.conn.logger.Debugf("Connection handler running")
	}

	for {
		select {
		case <-c.ctx.Done():
			if c.conn.logger != nil {
				c.conn.logger.Debugf("Context done, exiting connection handler: %v", c.ctx.Err())
			}
			return c.ctx.Err()
		default:
		}

		messageType, data, err := c.conn.ReadRaw(c.ctx)
		if err != nil {
			// Check for existing wrapped API errors
			var apiErr *apierrs.APIError
			var permanentErr *apierrs.PermanentError

			// First, check if this is already an API error
			if errors.As(err, &apiErr) {
				if c.conn.logger != nil {
					c.conn.logger.Errorf("API error reading message: %v", apiErr)
				}
				return apiErr
			}

			// Then check if it's a permanent error
			if errors.As(err, &permanentErr) {
				if c.conn.logger != nil {
					c.conn.logger.Errorf("Permanent error reading message: %v", permanentErr.Err)
				}
				return permanentErr.Err
			}

			// Special cases for network errors
			var netErr net.Error
			if errors.As(err, &netErr) {
				if netErr.Timeout() {
					// This is a timeout error (temporary)
					if c.conn.logger != nil {
						c.conn.logger.Warnf("Network timeout error: %v", err)
					}
					continue
				}
			}

			// Handle connection closed errors
			if strings.Contains(err.Error(), "use of closed network connection") ||
				strings.Contains(err.Error(), "connection reset by peer") {
				if c.conn.logger != nil {
					c.conn.logger.Infof("Connection closed: %v", err)
				}
				return apierrs.NewServerError("The connection was closed")
			}

			// For all other errors, treat as temporary and continue
			if c.conn.logger != nil {
				c.conn.logger.Warnf("Temporary error reading message: %v", err)
			}
			continue
		}

		if c.conn.logger != nil {
			c.conn.logger.Debugf("Received message of type: %s", messageType.String())
		}

		for i, handler := range c.handlers {
			if handler == nil {
				if c.conn.logger != nil {
					c.conn.logger.Warnf("Skipping nil handler at index %d", i)
				}
				continue
			}

			func() {
				defer func() {
					if r := recover(); r != nil {
						if c.conn.logger != nil {
							c.conn.logger.Errorf("Handler %d panicked: %v", i, r)
						}
					}
				}()
				handler(c.ctx, messageType, data)
			}()
		}
	}
}
