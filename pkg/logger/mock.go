package logger

import (
	"context"
	"log/slog"
)

func NewMock() *slog.Logger {
	return slog.New(MockHandler{})
}

type MockHandler struct {
	slog.Handler
}

func (MockHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (MockHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h MockHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h MockHandler) WithGroup(_ string) slog.Handler {
	return h
}
