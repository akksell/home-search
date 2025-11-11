package logger

import (
	"context"
	"log/slog"
	"os"
)

type Level = slog.Level

const (
	LevelDebug Level = slog.LevelDebug
	LevelInfo  Level = slog.LevelInfo
	LevelWarn  Level = slog.LevelWarn
	LevelError Level = slog.LevelError
)

type Attr = slog.Attr

var (
	String   = slog.String
	Int      = slog.Int
	Int64    = slog.Int64
	Float64  = slog.Float64
	Bool     = slog.Bool
	Duration = slog.Duration
	Time     = slog.Time
	Any      = slog.Any
	Group    = slog.Group
)

var instance *ServiceLogger

type ServiceLogger struct {
	logger *slog.Logger
}

func Init(serviceName string) *ServiceLogger {
	if instance == nil {
		instance = &ServiceLogger{
			logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", serviceName),
		}
	}
	return instance
}

func get() *ServiceLogger {
	if instance == nil {
		panic("Failed to log with attributes - instance must be created first with logger.Init()")
	}
	return instance
}

func (l *ServiceLogger) LogAttrs(ctx context.Context, level Level, message string, attrs ...Attr) {
	l.logger.LogAttrs(ctx, level, message, attrs...)
}

func LogAttrs(ctx context.Context, level Level, message string, attrs ...Attr) {
	get().logger.LogAttrs(ctx, level, message, attrs...)
}
