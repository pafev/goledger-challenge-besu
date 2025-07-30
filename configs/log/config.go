package logConfig

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/fatih/color"
)

type HandlerOptions struct {
	slog.HandlerOptions
}

type Handler struct {
	slog.Handler
	logger *log.Logger
}

func (handler *Handler) Handle(ctx context.Context, record slog.Record) error {
	level := record.Level.String() + ":"
	switch record.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	timeStr := record.Time.Format("[15:05:05.000]")
	msg := color.CyanString(record.Message)

	if record.NumAttrs() > 0 {
		fields := make(map[string]any, record.NumAttrs())
		record.Attrs(func(a slog.Attr) bool {
			fields[a.Key] = a.Value.Any()
			return true
		})
		attrs, err := json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
		handler.logger.Println(timeStr, level, msg, "\n"+string(attrs))
	} else {
		handler.logger.Println(timeStr, level, msg)
	}

	return nil
}

func NewHandler(
	out io.Writer,
	opts HandlerOptions,
) *Handler {
	return &Handler{
		Handler: slog.NewJSONHandler(out, &opts.HandlerOptions),
		logger:  log.New(out, "", 0),
	}
}

func Config() {
	opts := HandlerOptions{
		slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := NewHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}
