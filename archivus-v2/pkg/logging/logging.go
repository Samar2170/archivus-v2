package logging

import (
	"archivus-v2/config"
	"context"
	"path/filepath"
	"runtime/debug"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

var Errorlogger zerolog.Logger
var AuditLogger zerolog.Logger
var DebugLogger zerolog.Logger

func SetupLogging() {
	auditLogFilePath := filepath.Join(config.Config.LogsDir, "audit.log")
	errorLogFilePath := filepath.Join(config.Config.LogsDir, "error.log")
	debugLogFilePath := filepath.Join(config.Config.LogsDir, "debug.log")

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
	debugLogFile := &lumberjack.Logger{
		Filename:   debugLogFilePath,
		MaxSize:    20,
		MaxBackups: 2,
		MaxAge:     7,
		Compress:   false,
	}

	Errorlogger = zerolog.New(logFile).With().Timestamp().Logger()
	AuditLogger = zerolog.New(auditLogFile).With().Timestamp().Logger()
	DebugLogger = zerolog.New(debugLogFile).Level(zerolog.DebugLevel).With().Timestamp().Logger()
}

func HandleError(err error) error {
	if err != nil {
		Errorlogger.Error().
			Err(err).
			Str("stack", string(debug.Stack())).
			Msg("an error occurred")
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

// LogError returns an error event enriched with trace information.
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

// LogErrorWithStack logs an error with a full stack trace and optional trace context.
func LogErrorWithStack(ctx context.Context, err error, msg string) {
	span := trace.SpanFromContext(ctx)
	ev := Errorlogger.Error().
		Err(err).
		Str("stack", string(debug.Stack()))
	if span.IsRecording() {
		ev = ev.
			Str("trace_id", span.SpanContext().TraceID().String()).
			Str("span_id", span.SpanContext().SpanID().String())
	}
	ev.Msg(msg)
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
