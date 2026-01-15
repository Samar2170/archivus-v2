package logging

import (
	"archivus-v2/config"
	"context"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

var Errorlogger zerolog.Logger
var AuditLogger zerolog.Logger

func SetupLogging() {
	auditLogFilePath := filepath.Join(config.Config.LogsDir, "audit.log")
	errorLogFilePath := filepath.Join(config.Config.LogsDir, "error.log")
	auditLogFile := &lumberjack.Logger{
		Filename:   auditLogFilePath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   false,
	}
	logFile := &lumberjack.Logger{
		Filename:   errorLogFilePath,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   false,
	}

	Errorlogger = zerolog.New(logFile).With().Timestamp().Logger()
	AuditLogger = zerolog.New(auditLogFile).With().Timestamp().Logger()
}

func HandleError(err error) error {
	if err != nil {
		Errorlogger.Error().Err(err).Msg("An error occurred")
	}
	return err
}

// Log returns a logger enriched with trace information from the context.
// If the context has no span, it returns the standard Errorlogger.
func Log(ctx context.Context) *zerolog.Event {
	span := trace.SpanFromContext(ctx)
	logger := Errorlogger.With()
	if span.IsRecording() {
		logger = logger.
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("span_id", span.SpanContext().SpanID().String())
	}
	l := logger.Logger()
	return l.Info()
}

// LogError returns an error event enriched with trace information
func LogError(ctx context.Context, err error) *zerolog.Event {
	span := trace.SpanFromContext(ctx)
	logger := Errorlogger.Error().Err(err)
	if span.IsRecording() {
		logger = logger.
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("span_id", span.SpanContext().SpanID().String())
	}
	return logger
}

// LogWith returns a logger event from a specific logger, enriched with trace information from the context.
func LogWith(ctx context.Context, logger zerolog.Logger) *zerolog.Event {
	span := trace.SpanFromContext(ctx)
	l := logger.With()
	if span.IsRecording() {
		l = l.
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("span_id", span.SpanContext().SpanID().String())
	}
	lg := l.Logger()
	return lg.Info()
}
