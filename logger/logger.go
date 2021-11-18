package logger

import "github.com/easeq/go-service/component"

const (
	LOGGER = "logger"
)

type Logger interface {
	component.Component
	// Debug logs a message
	Debug(args ...interface{})
	// Debugf logs a formatted message
	Debugf(template string, args ...interface{})
	// Debugw logs a message with variadic key-value pairs
	Debugw(message string, args ...interface{})
	// Debugp logs a message with variadic key-value pairs
	Debugp(message string, method string, pkg string, err error)
	// Info logs a message
	Info(args ...interface{})
	// Infof logs a formatted message
	Infof(template string, args ...interface{})
	// Infow logs a message with variadic key-value pairs
	Infow(message string, args ...interface{})
	// Warn logs a message
	Warn(args ...interface{})
	// Warnf logs a formatted message
	Warnf(template string, args ...interface{})
	// Warnw logs a message with variadic key-value pairs
	Warnw(message string, args ...interface{})
	// Error logs a message
	Error(args ...interface{})
	// Errorf logs a formatted message
	Errorf(template string, args ...interface{})
	// Errorw logs a message with variadic key-value pairs
	Errorw(message string, args ...interface{})
	// DPanic logs a message
	DPanic(args ...interface{})
	// DPanicf logs a formatted message
	DPanicf(template string, args ...interface{})
	// DPanicw logs a message with variadic key-value pairs
	DPanicw(message string, args ...interface{})
	// Fatal logs a message
	Fatal(args ...interface{})
	// Fatalf logs a formatted message
	Fatalf(template string, args ...interface{})
	// Fatalw logs a message with variadic key-value pairs
	Fatalw(message string, args ...interface{})
}
