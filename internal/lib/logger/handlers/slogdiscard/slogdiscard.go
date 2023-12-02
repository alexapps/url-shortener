package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscardHandler struct{}

// Enabled implements slog.Handler.
func (*DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

// Handle implements slog.Handler.
func (*DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

// WithAttrs implements slog.Handler.
func (h *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// WithGroup implements slog.Handler.
func (h *DiscardHandler) WithGroup(_ string) slog.Handler {
	return h
}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}
