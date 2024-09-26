package otel_logger

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
)

// Logger is a wrapper around zap.Logger that adds OpenTelemetry support
type Logger struct {
	logger *zap.Logger
	name   string
}

// New creates a new Logger
func New(name string) *Logger {
	l, _ := zap.NewDevelopment()
	return &Logger{logger: l, name: name}
}

// NewWithZap creates a new Logger with a custom zap.Logger
func NewWithZap(l *zap.Logger) *Logger {
	return &Logger{logger: l}
}

// Info logs a message at the info level
func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
	span, attr := l.getSpanAndAttributes(ctx, fields...)
	l.addEventToSpan(span, msg, attr)
	// Maybe future option to log info to span
	// _, childSpan := otel.Tracer(l.name).Start(ctx, "Log Info")
	// childSpan.AddEvent(msg, trace.WithAttributes(attr...))
	// defer childSpan.End()
}

// Error logs a message at the error level and records an error on the span
func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
	span, attr := l.getSpanAndAttributes(ctx, fields...)
	if span != nil {
		_, childSpan := otel.Tracer(l.name).Start(ctx, "Error")
		defer childSpan.End()
		childSpan.RecordError(errors.New(msg), trace.WithAttributes(attr...))
		childSpan.SetStatus(codes.Error, msg)
		// span.SetStatus(codes.Error, msg) // set the outer span as an error
	}
}

// Debug logs a message at the debug level
func (l *Logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
	span, attr := l.getSpanAndAttributes(ctx, fields...)
	l.addEventToSpan(span, msg, attr)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
	span, attr := l.getSpanAndAttributes(ctx, fields...)
	l.addEventToSpan(span, msg, attr)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
	span, attr := l.getSpanAndAttributes(ctx, fields...)
	l.addEventToSpan(span, msg, attr)
}

func (l *Logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
	span, attr := l.getSpanAndAttributes(ctx, fields...)
	l.addEventToSpan(span, msg, attr)
}

func (l *Logger) zapFieldsToAttributes(fields ...zap.Field) []attribute.KeyValue {
	attr := make([]attribute.KeyValue, 0, len(fields))
	for _, f := range fields {
		attr = append(attr, attribute.String(f.Key, f.String))
	}
	return attr
}

func (l *Logger) getSpanAndAttributes(ctx context.Context, fields ...zap.Field) (trace.Span, []attribute.KeyValue) {
	if ctx != nil {
		span := trace.SpanFromContext(ctx)
		attr := l.zapFieldsToAttributes(fields...)
		return span, attr
	}
	return nil, nil
}

func (l *Logger) addEventToSpan(span trace.Span, msg string, attrributes []attribute.KeyValue) {
	if span == nil {
		return
	}
	span.AddEvent(msg, trace.WithAttributes(attrributes...))
}
