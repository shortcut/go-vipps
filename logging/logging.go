package logging

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

// Logger provides a basic leveled logging interface for printing informational, and error messages.
// The context will usually be the http.Request's context.
// It's implemented by a a no-op logger, and a Stdout-logger, and can also be used with a thin shim for use with other logging libraries.
type Logger interface {
	Info(ctx context.Context, message string, arguments ...LogArgument)
	Error(ctx context.Context, message string, arguments ...LogArgument)
}

// LogArgument is a simple key-value pair for logging contextual information.
type LogArgument struct {
	key   string
	value interface{}
}

// NewArg creates a new LogArgument.
func NewArg(key string, value interface{}) LogArgument { return LogArgument{key: key, value: value} }

type nopLogger struct{}

// NewNopLogger returns a logger that doesn't do anything.
func NewNopLogger() Logger { return nopLogger{} }

// Info logs an info-message.
func (nopLogger) Info(_ context.Context, _ string, _ ...LogArgument) {}

// Error logs an error-message.
func (nopLogger) Error(_ context.Context, _ string, _ ...LogArgument) {}

type stdOutLogger struct {
	l *log.Logger
}

// NewStdOutLogger returns a logger that logs everything to StdOut
func NewStdOutLogger() Logger { return &stdOutLogger{l: log.New(os.Stdout, "", log.LstdFlags)} }

// Info logs an info-message.
func (l *stdOutLogger) Info(_ context.Context, message string, arguments ...LogArgument) {
	msg := "info: " + message
	l.log(msg, arguments...)
}

// Error logs an error-message.
func (l *stdOutLogger) Error(_ context.Context, message string, arguments ...LogArgument) {
	msg := "error: " + message
	l.log(msg, arguments...)
}

func (l *stdOutLogger) log(message string, arguments ...LogArgument) {
	argStrings := make([]string, len(arguments))
	for i, a := range arguments {
		argStrings[i] = fmt.Sprintf("%s: %s", a.key, a.value)
	}
	if len(argStrings) > 0 {
		message = message + ":"
	}
	l.l.Println(message, strings.Join(argStrings, ", "))
}
