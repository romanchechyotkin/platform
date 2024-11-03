package logger

import (
	"log/slog"
	"os"
	"strings"
)

// Logger is a concrete implementation of the Logger interface using slog.Logger.
type Logger struct {
	l *slog.Logger
}

// convertAttrsToAny converts []slog.Attr to []any for slog.Logger methods
func convertAttrsToAny(attrs []slog.Attr) []any {
	anyAttr := make([]any, 0, len(attrs))

	for _, attr := range attrs {
		anyAttr = append(anyAttr, attr)
	}

	return anyAttr
}

func (l *Logger) Info(msg string, attrs ...slog.Attr) {
	l.l.Info(msg, convertAttrsToAny(attrs)...)
}

func (l *Logger) Warn(msg string, attrs ...slog.Attr) {
	l.l.Warn(msg, convertAttrsToAny(attrs)...)
}

func (l *Logger) Debug(msg string, attrs ...slog.Attr) {
	l.l.Debug(msg, convertAttrsToAny(attrs)...)
}

func (l *Logger) Error(msg string, err error, attrs ...slog.Attr) {
	attrs = append(attrs, slog.Any("error", err))
	l.l.Error(msg, convertAttrsToAny(attrs)...)
}

func (l *Logger) With(attrs ...slog.Attr) *Logger {
	c := l.clone()
	c.l = c.l.With(convertAttrsToAny(attrs)...)
	return c
}

const defaultLevel = slog.LevelDebug

// New creates a new instance of logger based on environment settings.
func New() *Logger {
	var handler slog.Handler

	if env := os.Getenv("APP_ENV"); env == "prod" {
		handler = prodHandler()
	} else if env == "dev" {
		handler = devHandler()
	} else {
		return &Logger{
			l: slog.Default(),
		}
	}

	return &Logger{l: slog.New(handler)}
}

// prodHandler configures the handler for production environments
func prodHandler() slog.Handler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: configLevel(),
	})
}

// devHandler configures the handler for development environments
func devHandler() slog.Handler {
	return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     configLevel(),
		AddSource: true,
	})
}

// configLevel sets the logging level based on the LOG_LEVEL environment variable
func configLevel() slog.Level {
	var logLevel slog.Level

	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = defaultLevel
	}

	return logLevel
}
func (l *Logger) clone() *Logger {
	c := *l
	return &c
}
