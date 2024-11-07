package zap

import (
	"context"

	"github.com/zahidhasanpapon/go-clean-architecture/pkg/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger implements the logger.Logger interface using Zap
type Logger struct {
	logger *zap.Logger
}

// Config holds the configuration for the Zap logger
type Config struct {
	Level      string
	OutputPath string
	Encoding   string // "json" or "console"
	DevMode    bool
}

// NewZapLogger New creates a new Logger instance
func NewZapLogger(cfg *Config) (*Logger, error) {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	zapCfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{cfg.OutputPath},
		ErrorOutputPaths: []string{cfg.OutputPath},
		Encoding:         cfg.Encoding,
		EncoderConfig:    getEncoderConfig(cfg.DevMode),
		Development:      cfg.DevMode,
	}

	l, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{logger: l}, nil
}

func getEncoderConfig(devMode bool) zapcore.EncoderConfig {
	if devMode {
		return zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}

	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// Implementation of logger.Logger interface

func (l *Logger) Debug(ctx context.Context, msg string, fields ...logger.Fields) {
	l.logger.Debug(msg, l.getZapFields(ctx, fields...)...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...logger.Fields) {
	l.logger.Info(msg, l.getZapFields(ctx, fields...)...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...logger.Fields) {
	l.logger.Warn(msg, l.getZapFields(ctx, fields...)...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...logger.Fields) {
	l.logger.Error(msg, l.getZapFields(ctx, fields...)...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...logger.Fields) {
	l.logger.Fatal(msg, l.getZapFields(ctx, fields...)...)
}

func (l *Logger) WithFields(fields logger.Fields) logger.Logger {
	return &Logger{
		logger: l.logger.With(l.convertToZapFields(fields)...),
	}
}

func (l *Logger) WithField(key string, value interface{}) logger.Logger {
	return &Logger{
		logger: l.logger.With(zap.Any(key, value)),
	}
}

func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// Helper methods

func (l *Logger) getZapFields(ctx context.Context, fields ...logger.Fields) []zap.Field {
	if len(fields) == 0 {
		return []zap.Field{}
	}

	// Merge all fields into a single map
	mergedFields := make(logger.Fields)
	for _, f := range fields {
		for k, v := range f {
			mergedFields[k] = v
		}
	}

	// Convert to zap fields
	return l.convertToZapFields(mergedFields)
}

func (l *Logger) convertToZapFields(fields logger.Fields) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}
