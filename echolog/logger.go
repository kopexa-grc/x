// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package echolog provides a high-performance structured logging wrapper for the Echo web framework.
//
// This package bridges the gap between Echo's logging interface and zerolog's powerful JSON logging
// capabilities, providing structured logging with minimal performance overhead.
//
// # Key Features
//
//   - Full compatibility with Echo's Logger interface
//   - Zero-allocation JSON logging through zerolog
//   - Configurable middleware for request logging
//   - Context-aware logging with request ID tracking
//   - Flexible enrichment and customization options
//   - Smart error handling and latency-based log levels
//
// # Quick Start
//
// Basic usage with Echo:
//
//	e := echo.New()
//	e.Logger = echolog.New(os.Stdout)
//
// With configuration options:
//
//	logger := echolog.New(
//		os.Stdout,
//		echolog.WithLevel(log.DEBUG),
//		echolog.WithTimestamp(),
//		echolog.WithCaller(),
//	)
//	e.Logger = logger
//
// # Middleware Usage
//
// Add request logging middleware:
//
//	e.Use(echolog.LoggingMiddleware(echolog.Config{
//		Logger: logger,
//		RequestLatencyLimit: 500 * time.Millisecond,
//		RequestLatencyLevel: zerolog.WarnLevel,
//	}))
//
// # Performance
//
// This package is built on zerolog, which provides:
//   - Zero allocation JSON logger
//   - Lazy evaluation of expensive operations
//   - Structured logging without reflection
//   - Minimal overhead over raw zerolog
//
// For more information, see the package documentation at:
// https://pkg.go.dev/github.com/kopexa-grc/x/echolog
package echolog

import (
	"fmt"
	"io"

	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
)

// Logger provides a high-performance implementation of Echo's Logger interface using zerolog.
//
// Logger wraps a zerolog.Logger instance and maintains compatibility with Echo's logging
// interface while providing structured JSON logging capabilities. It supports all Echo
// logging methods and extends them with zerolog's performance optimizations.
//
// The Logger maintains internal state for configuration options and can be dynamically
// reconfigured using the SetLevel and SetPrefix methods.
//
// Example usage:
//
//	logger := echolog.New(os.Stdout, echolog.WithTimestamp())
//	logger.Info("application started")
//	logger.Debugf("user %s logged in", userID)
type Logger struct {
	log     zerolog.Logger // underlying zerolog instance
	out     io.Writer      // output writer (may be nil)
	level   log.Lvl        // current log level in Echo format
	prefix  string         // log prefix
	setters []ConfigSetter // configuration setters for reconfiguration
}

// New creates a new Logger instance with the specified output writer and configuration options.
//
// The function accepts an io.Writer for log output and optional ConfigSetter functions
// for customizing the logger behavior. If the output writer is already a zerolog.Logger,
// it will be used directly; otherwise, a new zerolog.Logger will be created.
//
// Parameters:
//   - out: The output destination for log messages. Can be any io.Writer or a zerolog.Logger.
//   - setters: Optional configuration functions to customize logger behavior.
//
// Returns a configured Logger instance ready for use with Echo.
//
// Example:
//
//	// Basic logger writing to stdout
//	logger := echolog.New(os.Stdout)
//
//	// Logger with timestamp and debug level
//	logger := echolog.New(os.Stdout,
//		echolog.WithTimestamp(),
//		echolog.WithLevel(log.DEBUG),
//	)
//
//	// Logger with multiple configuration options
//	logger := echolog.New(os.Stdout,
//		echolog.WithTimestamp(),
//		echolog.WithCaller(),
//		echolog.WithFields(map[string]interface{}{
//			"service": "api",
//			"version": "1.0.0",
//		}),
//	)
func New(out io.Writer, setters ...ConfigSetter) *Logger {
	switch l := out.(type) {
	case zerolog.Logger:
		return newLogger(l, setters)
	default:
		return newLogger(zerolog.New(out), setters)
	}
}

// From creates a new Logger instance using an existing zerolog.Logger.
//
// This function is useful when you already have a configured zerolog.Logger instance
// and want to wrap it with Echo's Logger interface. The existing logger's configuration
// is preserved, and additional configuration can be applied via setters.
//
// Parameters:
//   - log: An existing zerolog.Logger instance to wrap.
//   - setters: Optional configuration functions to apply additional customization.
//
// Returns a Logger instance that wraps the provided zerolog.Logger.
//
// Example:
//
//	// Create a zerolog logger with custom configuration
//	zlog := zerolog.New(os.Stdout).With().
//		Timestamp().
//		Str("service", "api").
//		Logger()
//
//	// Wrap it with echolog
//	logger := echolog.From(zlog)
//
//	// Or with additional configuration
//	logger := echolog.From(zlog, echolog.WithLevel(log.DEBUG))
func From(log zerolog.Logger, setters ...ConfigSetter) *Logger {
	return newLogger(log, setters)
}

// newLogger creates a new Logger instance with the provided zerolog.Logger and configuration setters.
// This is an internal constructor function used by New and From.
func newLogger(log zerolog.Logger, setters []ConfigSetter) *Logger {
	opts := newOptions(log, setters)

	return &Logger{
		log:     opts.zcontext.Logger(),
		out:     nil,
		level:   opts.level,
		prefix:  opts.prefix,
		setters: setters,
	}
}

// Write implements io.Writer interface, allowing the Logger to be used as a writer.
//
// This method delegates to the underlying zerolog.Logger's Write method, enabling
// the Logger to be used anywhere an io.Writer is expected.
//
// Parameters:
//   - p: byte slice containing the data to write
//
// Returns:
//   - n: number of bytes written
//   - err: any error that occurred during writing
//
// Example:
//
//	logger := echolog.New(os.Stdout)
//	fmt.Fprintf(logger, "This will be written to the logger")
func (l *Logger) Write(p []byte) (n int, err error) {
	return l.log.Write(p)
}

// Debug logs a message at debug level using Sprint-style formatting.
//
// Debug level is typically used for detailed information that is valuable for
// debugging and development but should not be logged in production.
//
// Parameters:
//   - i: variadic arguments that will be formatted using fmt.Sprint
//
// Example:
//
//	logger.Debug("processing request", requestID, "for user", userID)
func (l *Logger) Debug(i ...any) {
	l.log.Debug().Msg(fmt.Sprint(i...))
}

// Debugf logs a formatted message at debug level using Printf-style formatting.
//
// This method provides formatted string logging with placeholder substitution,
// similar to fmt.Printf.
//
// Parameters:
//   - format: format string with placeholders
//   - i: variadic arguments to substitute into the format string
//
// Example:
//
//	logger.Debugf("user %s accessed endpoint %s at %v", userID, endpoint, time.Now())
func (l *Logger) Debugf(format string, i ...any) {
	l.log.Debug().Msgf(format, i...)
}

// Debugj logs a JSON object at debug level.
//
// This method accepts a map of key-value pairs and logs them as structured
// JSON fields. This is useful for logging complex data structures.
//
// Parameters:
//   - j: map containing key-value pairs to log as JSON fields
//
// Example:
//
//	logger.Debugj(log.JSON{
//		"user_id": "123",
//		"action": "login",
//		"ip": "192.168.1.1",
//	})
func (l *Logger) Debugj(j log.JSON) {
	l.logJSON(l.log.Debug(), j)
}

// Info logs a message at info level using Sprint-style formatting.
//
// Info level is used for general information about application operation.
// This is typically the default logging level for production systems.
//
// Parameters:
//   - i: variadic arguments that will be formatted using fmt.Sprint
//
// Example:
//
//	logger.Info("server started on port", 8080)
func (l *Logger) Info(i ...any) {
	l.log.Info().Msg(fmt.Sprint(i...))
}

// Infof logs a formatted message at info level using Printf-style formatting.
//
// Parameters:
//   - format: format string with placeholders
//   - i: variadic arguments to substitute into the format string
//
// Example:
//
//	logger.Infof("server started on port %d", 8080)
func (l *Logger) Infof(format string, i ...any) {
	l.log.Info().Msgf(format, i...)
}

// Infoj logs a JSON object at info level.
//
// Parameters:
//   - j: map containing key-value pairs to log as JSON fields
//
// Example:
//
//	logger.Infoj(log.JSON{"event": "startup", "port": 8080})
func (l *Logger) Infoj(j log.JSON) {
	l.logJSON(l.log.Info(), j)
}

// Warn logs a message at warning level using Sprint-style formatting.
//
// Warning level indicates potentially harmful situations or unusual conditions
// that should be noted but don't prevent the application from continuing.
//
// Parameters:
//   - i: variadic arguments that will be formatted using fmt.Sprint
//
// Example:
//
//	logger.Warn("deprecated API endpoint accessed", endpoint)
func (l *Logger) Warn(i ...any) {
	l.log.Warn().Msg(fmt.Sprint(i...))
}

// Warnf logs a formatted message at warning level using Printf-style formatting.
//
// Parameters:
//   - format: format string with placeholders
//   - i: variadic arguments to substitute into the format string
//
// Example:
//
//	logger.Warnf("slow query detected: %s took %v", query, duration)
func (l *Logger) Warnf(format string, i ...any) {
	l.log.Warn().Msgf(format, i...)
}

// Warnj logs a JSON object at warning level.
//
// Parameters:
//   - j: map containing key-value pairs to log as JSON fields
//
// Example:
//
//	logger.Warnj(log.JSON{"event": "slow_query", "duration_ms": 5000})
func (l *Logger) Warnj(j log.JSON) {
	l.logJSON(l.log.Warn(), j)
}

// Error logs a message at error level using Sprint-style formatting.
//
// Error level indicates error conditions that should be investigated.
// Errors represent failures that don't necessarily require immediate attention
// but should be addressed.
//
// Parameters:
//   - e: variadic arguments that will be formatted using fmt.Sprint
//
// Example:
//
//	logger.Error("failed to connect to database", err)
func (l *Logger) Error(e ...interface{}) {
	l.log.Error().Msg(fmt.Sprint(e...))
}

// Errorf logs a formatted message at error level using Printf-style formatting.
//
// Parameters:
//   - format: format string with placeholders
//   - i: variadic arguments to substitute into the format string
//
// Example:
//
//	logger.Errorf("database connection failed after %d attempts: %v", attempts, err)
func (l *Logger) Errorf(format string, i ...any) {
	l.log.Error().Msgf(format, i...)
}

// Errorj logs a JSON object at error level.
//
// Parameters:
//   - j: map containing key-value pairs to log as JSON fields
//
// Example:
//
//	logger.Errorj(log.JSON{"error": err.Error(), "operation": "db_connect"})
func (l *Logger) Errorj(j log.JSON) {
	l.logJSON(l.log.Error(), j)
}

// Fatal logs a message at fatal level and terminates the program.
//
// Fatal level indicates critical errors that require immediate program termination.
// After logging the message, the program will exit with status code 1.
// Use with caution as this will stop program execution.
//
// Parameters:
//   - i: variadic arguments that will be formatted using fmt.Sprint
//
// Example:
//
//	logger.Fatal("cannot start server", err)
func (l *Logger) Fatal(i ...any) {
	l.log.Fatal().Msg(fmt.Sprint(i...))
}

// Fatalf logs a formatted message at fatal level and terminates the program.
//
// After logging the message, the program will exit with status code 1.
//
// Parameters:
//   - format: format string with placeholders
//   - i: variadic arguments to substitute into the format string
//
// Example:
//
//	logger.Fatalf("cannot bind to port %d: %v", port, err)
func (l *Logger) Fatalf(format string, i ...any) {
	l.log.Fatal().Msgf(format, i...)
}

// Fatalj logs a JSON object at fatal level and terminates the program.
//
// After logging the message, the program will exit with status code 1.
//
// Parameters:
//   - j: map containing key-value pairs to log as JSON fields
//
// Example:
//
//	logger.Fatalj(log.JSON{"error": "startup_failed", "port": 8080})
func (l *Logger) Fatalj(j log.JSON) {
	l.logJSON(l.log.Fatal(), j)
}

// Panic logs a message at panic level and panics.
//
// Panic level indicates critical errors that require immediate program termination
// with a stack trace. After logging the message, the function calls panic().
// Use only for unrecoverable errors.
//
// Parameters:
//   - i: variadic arguments that will be formatted using fmt.Sprint
//
// Example:
//
//	logger.Panic("critical system failure", err)
func (l *Logger) Panic(i ...any) {
	l.log.Panic().Msg(fmt.Sprint(i...))
}

// Panicf logs a formatted message at panic level and panics.
//
// After logging the message, the function calls panic().
//
// Parameters:
//   - format: format string with placeholders
//   - i: variadic arguments to substitute into the format string
//
// Example:
//
//	logger.Panicf("invalid state: %s", state)
func (l *Logger) Panicf(format string, i ...any) {
	l.log.Panic().Msgf(format, i...)
}

// Panicj logs a JSON object at panic level and panics.
//
// After logging the message, the function calls panic().
//
// Parameters:
//   - j: map containing key-value pairs to log as JSON fields
//
// Example:
//
//	logger.Panicj(log.JSON{"error": "system_failure", "component": "auth"})
func (l *Logger) Panicj(j log.JSON) {
	l.logJSON(l.log.Panic(), j)
}

// Print logs a message without a specific level using Sprint-style formatting.
//
// Print outputs messages without a log level classification, useful for
// general output that doesn't fit into standard logging levels.
//
// Parameters:
//   - i: variadic arguments that will be formatted using fmt.Sprint
//
// Example:
//
//	logger.Print("application banner or status message")
func (l *Logger) Print(i ...any) {
	l.log.WithLevel(zerolog.NoLevel).Str("level", "-").Msg(fmt.Sprint(i...))
}

// Printf logs a formatted message without a specific level.
//
// Parameters:
//   - format: format string with placeholders
//   - i: variadic arguments to substitute into the format string
//
// Example:
//
//	logger.Printf("Application %s v%s starting", name, version)
func (l *Logger) Printf(format string, i ...any) {
	l.log.WithLevel(zerolog.NoLevel).Str("level", "-").Msgf(format, i...)
}

// Printj logs a JSON object without a specific level.
//
// Parameters:
//   - j: map containing key-value pairs to log as JSON fields
//
// Example:
//
//	logger.Printj(log.JSON{"app": "myapp", "version": "1.0.0"})
func (l *Logger) Printj(j log.JSON) {
	l.logJSON(l.log.WithLevel(zerolog.NoLevel).Str("level", "-"), j)
}

// Output returns the current output writer for the logger.
//
// This method implements the Echo Logger interface and returns the io.Writer
// that the logger is currently writing to. May return nil if no explicit
// output writer was set.
//
// Returns:
//   - io.Writer: the current output destination, or nil if not set
//
// Example:
//
//	writer := logger.Output()
//	if writer != nil {
//		fmt.Fprintf(writer, "Direct write to logger output")
//	}
func (l *Logger) Output() io.Writer {
	return l.out
}

// SetOutput changes the output destination for the logger.
//
// This method updates both the internal output writer reference and
// reconfigures the underlying zerolog.Logger to write to the new destination.
//
// Parameters:
//   - newOut: the new io.Writer to use for log output
//
// Example:
//
//	// Change logger output to a file
//	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//	if err == nil {
//		logger.SetOutput(file)
//	}
func (l *Logger) SetOutput(newOut io.Writer) {
	l.out = newOut
	l.log = l.log.Output(newOut)
}

// Level returns the current logging level.
//
// This method returns the current log level in Echo's log.Lvl format.
// The level determines which log messages will be output.
//
// Returns:
//   - log.Lvl: current logging level (DEBUG, INFO, WARN, ERROR, OFF)
//
// Example:
//
//	currentLevel := logger.Level()
//	if currentLevel == log.DEBUG {
//		logger.Debug("Debug logging is enabled")
//	}
func (l *Logger) Level() log.Lvl {
	return l.level
}

// SetLevel changes the current logging level.
//
// This method updates the logging level for the logger, affecting which
// messages will be output. Messages below the set level will be discarded.
//
// Parameters:
//   - level: new logging level to set (log.DEBUG, log.INFO, log.WARN, log.ERROR, log.OFF)
//
// Example:
//
//	// Set to debug level for development
//	logger.SetLevel(log.DEBUG)
//
//	// Set to error level for production
//	logger.SetLevel(log.ERROR)
func (l *Logger) SetLevel(level log.Lvl) {
	zlvl, elvl := MatchEchoLevel(level)

	l.setters = append(l.setters, WithLevel(elvl))
	l.level = elvl
	l.log = l.log.Level(zlvl)
}

// Prefix returns the current log prefix.
//
// This method returns the prefix string that is added to log messages.
// The prefix appears as a "prefix" field in the JSON output.
//
// Returns:
//   - string: current prefix, empty string if no prefix is set
//
// Example:
//
//	currentPrefix := logger.Prefix()
//	fmt.Printf("Current prefix: %s\n", currentPrefix)
func (l *Logger) Prefix() string {
	return l.prefix
}

// SetHeader sets the log header format.
//
// This method implements the Echo Logger interface but is not functional
// in this implementation since zerolog handles its own formatting.
// The method exists for interface compatibility only.
//
// Parameters:
//   - header: header format string (ignored)
//
// Note: This method does nothing and exists only for Echo interface compatibility.
func (l *Logger) SetHeader(_ string) {
	// not implemented - zerolog handles its own formatting
}

// SetPrefix changes the log prefix.
//
// This method sets a prefix that will be included as a "prefix" field
// in all subsequent log messages. The prefix is useful for identifying
// the source or context of log messages.
//
// Parameters:
//   - newPrefix: the new prefix string to use
//
// Example:
//
//	logger.SetPrefix("[API]")
//	logger.Info("Server started") // Will include "prefix": "[API]"
func (l *Logger) SetPrefix(newPrefix string) {
	l.setters = append(l.setters, WithPrefix(newPrefix))

	opts := newOptions(l.log, l.setters)

	l.prefix = newPrefix
	l.log = opts.zcontext.Logger()
}

// Unwrap returns the underlying zerolog.Logger instance.
//
// This method provides access to the wrapped zerolog.Logger for advanced
// usage scenarios where direct zerolog functionality is needed.
//
// Returns:
//   - zerolog.Logger: the underlying zerolog logger instance
//
// Example:
//
//	// Access zerolog-specific functionality
//	zlog := logger.Unwrap()
//	zlog.Info().Str("custom_field", "value").Msg("Advanced logging")
func (l *Logger) Unwrap() zerolog.Logger {
	return l.log
}

// logJSON is an internal helper method that logs a JSON object.
//
// This method takes a zerolog.Event and a map of key-value pairs,
// adds each pair as a field to the event, and logs it.
//
// Parameters:
//   - event: zerolog.Event to add fields to
//   - j: map containing key-value pairs to add as fields
func (l *Logger) logJSON(event *zerolog.Event, j log.JSON) {
	for k, v := range j {
		event = event.Interface(k, v)
	}

	event.Msg("")
}
