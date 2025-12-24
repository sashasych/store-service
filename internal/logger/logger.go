package logger

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

// New builds a zap logger with the given level (info, debug, warn, error).
func New(level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"

	lvl := zapcore.InfoLevel
	if err := lvl.Set(strings.ToLower(level)); err == nil {
		cfg.Level = zap.NewAtomicLevelAt(lvl)
	}

	log, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	zap.ReplaceGlobals(log)
	return log, nil
}

// WithContext stores logger in the context.
func WithContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

// FromContext returns logger from context or global fallback.
func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return zap.L()
	}

	if log, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok && log != nil {
		return log
	}
	return zap.L()
}
