package zaplog

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/suphawit-wt/linebk-assignment/linebk-api/pkg/go-libs/config"
)

type ctxKey string

const (
	loggerCtxKey ctxKey = "app_logger"
	// TraceIDKey is the field key used for logging trace IDs.
	TraceIDKey string = "trace_id"
	// RequestIDKey is the field key used for logging request IDs.
	RequestIDKey string = "request_id"
)

// Logger wraps the zap logger and atomic level control.
type Logger struct {
	log         *zap.Logger
	atomicLevel zap.AtomicLevel
}

// NewLogger creates and configures a new Logger instance based on the application mode and configuration.
func NewLogger(appMode string, cfg *config.LoggerConfig) (*Logger, error) {
	appMode = strings.ToLower(appMode)

	var zapCfg zap.Config

	switch appMode {
	case "prod", "production", "staging":
		zapCfg = zap.NewProductionConfig()
		zapCfg.EncoderConfig.TimeKey = "timestamp"
		zapCfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		zapCfg.EncoderConfig.StacktraceKey = "stacktrace"
		zapCfg.EncoderConfig.LevelKey = "level"
		zapCfg.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	case "dev", "development", "local", "":
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapCfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	default:
		return nil, fmt.Errorf(
			"logger setup error: unknown env %q (must be dev, staging, or prod)",
			appMode,
		)
	}

	atomicLevel := zap.NewAtomicLevelAt(zapCfg.Level.Level())

	if cfg != nil && cfg.Level != "" {
		parsedLevel, err := zapcore.ParseLevel(cfg.Level)
		if err != nil {
			return nil, fmt.Errorf("invalid log level %q: %w", cfg.Level, err)
		}
		atomicLevel.SetLevel(parsedLevel)
	}
	zapCfg.Level = atomicLevel

	zapCfg.OutputPaths = []string{"stdout"}
	zapCfg.ErrorOutputPaths = []string{"stderr"}

	logger, err := zapCfg.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	return &Logger{log: logger, atomicLevel: atomicLevel}, nil
}

// NewNop creates a no-operation Logger that produces no output.
func NewNop() *Logger {
	return &Logger{
		log:         zap.NewNop(),
		atomicLevel: zap.NewAtomicLevelAt(zapcore.InvalidLevel),
	}
}

// Info logs a message at InfoLevel.
func (l *Logger) Info(msg string, fields ...zap.Field) { l.log.Info(msg, fields...) }

// Debug logs a message at DebugLevel.
func (l *Logger) Debug(msg string, fields ...zap.Field) { l.log.Debug(msg, fields...) }

// Warn logs a message at WarnLevel.
func (l *Logger) Warn(msg string, fields ...zap.Field) { l.log.Warn(msg, fields...) }

// Error logs a message at ErrorLevel.
func (l *Logger) Error(msg string, fields ...zap.Field) { l.log.Error(msg, fields...) }

// Fatal logs a message at FatalLevel and terminates execution.
func (l *Logger) Fatal(msg string, fields ...zap.Field) { l.log.Fatal(msg, fields...) }

// SetLevel dynamically updates the logger's log level.
func (l *Logger) SetLevel(level string) error {
	parsedLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("failed to parse level %q: %w", level, err)
	}
	l.atomicLevel.SetLevel(parsedLevel)
	return nil
}

// Named adds a sub-scope to the logger's name.
func (l *Logger) Named(name string) *Logger {
	return &Logger{
		log:         l.log.Named(name),
		atomicLevel: l.atomicLevel,
	}
}

// With adds structured context fields to the logger.
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		log:         l.log.With(fields...),
		atomicLevel: l.atomicLevel,
	}
}

// IntoContext injects the Logger into the given context.
func (l *Logger) IntoContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

// FromContext retrieves the Logger from the context, or returns a nop Logger if not found.
func FromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(loggerCtxKey).(*Logger); ok {
		return l
	}
	return NewNop()
}

// WithTrace attaches trace and request IDs to the logger fields.
func (l *Logger) WithTrace(traceID, requestID string) *Logger {
	fields := make([]zap.Field, 0, 2)
	if traceID != "" {
		fields = append(fields, zap.String(TraceIDKey, traceID))
	}
	if requestID != "" {
		fields = append(fields, zap.String(RequestIDKey, requestID))
	}
	return l.With(fields...)
}

// Sync flushes buffered log entries.
func (l *Logger) Sync() error {
	err := l.log.Sync()
	if err != nil &&
		(os.IsPermission(err) || strings.Contains(err.Error(), "inappropriate ioctl for device") || strings.Contains(err.Error(), "bad file descriptor")) {
		return nil
	}
	return err
}
