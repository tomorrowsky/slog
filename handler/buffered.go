package handler

import (
	"bufio"

	"github.com/gookit/slog"
)

const defaultFlushInterval = 1000

// BufferedHandler definition
type BufferedHandler struct {
	BaseHandler
	number int
	buffer *bufio.Writer
	handler slog.Handler
	// options
	FlushInterval int
}

// NewBufferedHandler create new BufferedHandler
func NewBufferedHandler(handler slog.WriterHandler) *BufferedHandler {
	return &BufferedHandler{
		buffer: bufio.NewWriter(handler.Writer()),
		handler: handler,
		FlushInterval: defaultFlushInterval,
	}
}

// Flush all buffers to the `h.handler.Writer()`
func (h *BufferedHandler) Flush() error {
	return h.buffer.Flush()
}

// Handle log record
func (h *BufferedHandler) Handle(record *slog.Record)  error {
	bts, err := h.Formatter().Format(record)
	if err != nil {
		return err
	}

	h.number++
	_, err = h.buffer.Write(bts)

	if h.number >= h.FlushInterval {
		return h.Flush()
	}

	return err
}

