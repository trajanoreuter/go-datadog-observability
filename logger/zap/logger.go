package zap

import (
	"context"
	"os"

	ginCtx "github.com/trajanoreuter/go-datadog-observability/context/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ILogger is the interface that wraps the basic logging methods and automatically adds the datadog fields to the log message.
// Info, Warn, Debug, Error, Panic, Fatal
// The base logger is zap logger.
type ILogger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Panic(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
}

// Config is the configuration struct for the logger.
type Config struct {
	Level   zap.AtomicLevel
	Datadog struct {
		Service     string
		version     string
		Environment string
	}
}

type Logger struct {
	logger *zap.Logger
	cfg    *Config
}

// setContextFields is a helper function that adds the datadog trace and span id to the log message.
func setContextFields(ctx context.Context, fields ...zap.Field) []zap.Field {
	var ok bool
	traceID, ok := ginCtx.TraceID(ctx)

	if ok {
		fields = append(fields, zap.Uint64("dd.trace_id", traceID))
	}

	spanID, ok := ginCtx.SpanID(ctx)

	if ok {
		fields = append(fields, zap.Uint64("dd.span_id", spanID))
	}

	return fields
}

// Info logs a message at the info level with Datadog trace and span_id.
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = setContextFields(ctx, fields...)
	l.logger.Info(msg, fields...)
}

// Warn logs a message at the warn level with Datadog trace and span_id.
func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = setContextFields(ctx, fields...)
	l.logger.Warn(msg, fields...)
}

// Debug logs a message at the debug level with Datadog trace and span_id.
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = setContextFields(ctx, fields...)
	l.logger.Debug(msg, fields...)
}

// Error logs a message at the error level with Datadog trace and span_id.
func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = setContextFields(ctx, fields...)
	l.logger.Error(msg, fields...)
}

// Panic logs a message at the panic level with Datadog trace and span_id.
func (l *Logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	fields = setContextFields(ctx, fields...)
	l.logger.Panic(msg, fields...)
}

// Fatal logs a message at the fatal level with Datadog trace and span_id.
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = setContextFields(ctx, fields...)
	l.logger.Fatal(msg, fields...)
}

// configureLogger is a helper function that configures the zap logger with the desired settings.
func configureLogger(cfg *Config) zap.Config {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "message",
		LevelKey:    "severity",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "time",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	}

	hostname, _ := os.Hostname()

	config := zap.Config{
		Encoding:         "json",
		Level:            cfg.Level,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Development:      false,
		InitialFields: map[string]interface{}{
			"host.name":        hostname,
			"service.language": "golang",
			"dd.service":       cfg.Datadog.Service,
			"dd.version":       cfg.Datadog.version,
			"dd.env":           cfg.Datadog.Environment,
		},
		EncoderConfig: encoderConfig,
	}

	return config
}

// NewLogger creates a new instance of the Logger.
func NewLogger(cfg *Config) ILogger {
	config := configureLogger(cfg)

	l, err := config.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{
		logger: l,
		cfg:    cfg,
	}
}
