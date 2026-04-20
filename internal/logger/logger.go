package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Info(ctx context.Context, tag, method, msg string, fields ...zap.Field)
	Error(ctx context.Context, tag, method, msg string, fields ...zap.Field)
	Warn(ctx context.Context, tag, method, msg string, fields ...zap.Field)
}

type zapLogger struct {
	zap *zap.Logger
}

func New(level string) Logger {
	lvl := zap.InfoLevel
	if level == "debug" {
		lvl = zap.DebugLevel
	}
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(lvl),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	z, _ := cfg.Build()
	return &zapLogger{zap: z}
}

func (l *zapLogger) Info(ctx context.Context, tag, method, msg string, fields ...zap.Field) {
	l.zap.Info(msg, append([]zap.Field{zap.String("tag", tag), zap.String("method", method)}, fields...)...)
}

func (l *zapLogger) Error(ctx context.Context, tag, method, msg string, fields ...zap.Field) {
	l.zap.Error(msg, append([]zap.Field{zap.String("tag", tag), zap.String("method", method)}, fields...)...)
}

func (l *zapLogger) Warn(ctx context.Context, tag, method, msg string, fields ...zap.Field) {
	l.zap.Warn(msg, append([]zap.Field{zap.String("tag", tag), zap.String("method", method)}, fields...)...)
}

func String(key, val string) zap.Field  { return zap.String(key, val) }
func Int(key string, val int) zap.Field { return zap.Int(key, val) }
func Err(err error) zap.Field           { return zap.Error(err) }
