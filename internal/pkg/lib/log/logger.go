package log

import (
	"log/slog"
	"os"
)

// Level - logging level.
type Level int

// Level - base levels.
const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

// Opts - options to create new logger instance.
type Opts struct {
	Level  Level
	JSON   bool
	Writer *os.File
}

func New(opts Opts) *slog.Logger {
	writer := opts.Writer
	if writer == nil {
		writer = os.Stderr
	}

	level := new(slog.LevelVar)
	level.Set(slog.Level(opts.Level))

	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if opts.JSON {
		handler = slog.NewJSONHandler(writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(writer, handlerOpts)
	}

	return slog.New(handler)
}
