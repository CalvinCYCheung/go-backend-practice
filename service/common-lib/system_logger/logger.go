package systemLogger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type CustomHandler struct {
	writer io.Writer
	level  slog.Level
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *CustomHandler) Handle(ctx context.Context, record slog.Record) error {
	// Format: LEVEL MESSAGE TIME
	record.Attrs(func(a slog.Attr) bool {
		fmt.Println(a.Key, a.Value)
		return true
	})
	timeStr := record.Time.Format("2006-01-02 15:04:05")
	levelStr := record.Level.String()
	message := record.Message

	output := fmt.Sprintf("%s %s %s \n", levelStr, message, timeStr)
	_, err := h.writer.Write([]byte(output))
	return err
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fmt.Println("WithAttrs", attrs)
	return h
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	fmt.Println("WithGroup", name)
	return h
}

func InitLogger(level slog.Level) *slog.Logger {
	handler := &CustomHandler{
		writer: os.Stdout,
		level:  level,
	}
	// slog.TextHandler()
	return slog.New(handler)
}
