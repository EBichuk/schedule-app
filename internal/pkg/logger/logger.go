package logger

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"schedule-app/internal/pkg/contexts"
)

func GetLogger(LOG_FILE string) *slog.Logger {
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Not open log file", err.Error())
	}

	handler := slog.Handler(slog.NewJSONHandler(io.MultiWriter(os.Stdout, file), &slog.HandlerOptions{Level: slog.LevelDebug}))

	handler = NewHandlerMiddleware(handler)

	log := slog.New(handler)
	slog.SetDefault(log)

	return log
}

type HandlerMiddlware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddlware {
	return &HandlerMiddlware{next: next}
}

func (h *HandlerMiddlware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddlware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := contexts.RequestIDFromContext(ctx); ok == nil {
		rec.Add("X-Request-Id", c)
	}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithAttrs(attrs)}
}

func (h *HandlerMiddlware) WithGroup(name string) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithGroup(name)}
}
