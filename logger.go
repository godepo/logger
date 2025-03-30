//go:generate go tool mockery
package logger

import (
	"context"
	"log/slog"
)

type Logger interface {
	Error(ctx context.Context, message string, tags ...slog.Attr)
	Info(ctx context.Context, message string, tags ...slog.Attr)
	Debug(ctx context.Context, message string, tags ...slog.Attr)
	Warn(ctx context.Context, message string, tags ...slog.Attr)
	With(tags ...slog.Attr) Logger
}

type loggerTags int8

const (
	logInstance loggerTags = iota
)

func With(ctx context.Context, attrs ...slog.Attr) context.Context {
	return context.WithValue(ctx, logInstance, From(ctx).With(attrs...))
}

func Wrap(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, logInstance, logger)
}

func From(ctx context.Context) Logger {
	log, ok := ctx.Value(logInstance).(Logger)
	if !ok {
		return wrapSLog()
	}

	return log
}

type slogWrapper struct {
	slog *slog.Logger
}

func (wrapper slogWrapper) Error(ctx context.Context, message string, tags ...slog.Attr) {
	wrapper.slog.LogAttrs(ctx, slog.LevelError, message, tags...)
}

func (wrapper slogWrapper) Info(ctx context.Context, message string, tags ...slog.Attr) {
	wrapper.slog.LogAttrs(ctx, slog.LevelInfo, message, tags...)
}

func (wrapper slogWrapper) Debug(ctx context.Context, message string, tags ...slog.Attr) {
	wrapper.slog.LogAttrs(ctx, slog.LevelDebug, message, tags...)
}

func (wrapper slogWrapper) Warn(ctx context.Context, message string, tags ...slog.Attr) {
	wrapper.slog.LogAttrs(ctx, slog.LevelWarn, message, tags...)
}

func (wrapper slogWrapper) With(tags ...slog.Attr) Logger {
	var out slogWrapper
	fork := make([]any, 0, len(tags))

	for _, tag := range tags {
		fork = append(fork, tag)
	}

	out.slog = wrapper.slog.With(fork...)

	return out
}

func wrapSLog() slogWrapper {
	return slogWrapper{
		slog: slog.Default(),
	}
}

func Error(ctx context.Context, message string, tags ...slog.Attr) {
	From(ctx).Error(ctx, message, tags...)
}

func Info(ctx context.Context, message string, tags ...slog.Attr) {
	From(ctx).Info(ctx, message, tags...)
}

func Debug(ctx context.Context, message string, tags ...slog.Attr) {
	From(ctx).Debug(ctx, message, tags...)
}

func Warn(ctx context.Context, message string, tags ...slog.Attr) {
	From(ctx).Warn(ctx, message, tags...)
}
