// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package echolog

import (
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
)

// Options holds configuration settings for a Logger instance.
//
// This struct contains the internal configuration state used during logger
// construction and reconfiguration. It manages the zerolog context, logging
// level, and prefix settings.
type Options struct {
	zcontext zerolog.Context // zerolog context for building the logger
	level    log.Lvl         // current log level in Echo format
	prefix   string          // log prefix string
}

// ConfigSetter is a function type for configuring Logger options.
//
// ConfigSetter functions are used to apply configuration settings to a Logger
// instance during construction or reconfiguration. They follow the functional
// options pattern, allowing for flexible and composable configuration.
//
// Example usage:
//
//	logger := echolog.New(os.Stdout,
//		echolog.WithLevel(log.DEBUG),
//		echolog.WithTimestamp(),
//		echolog.WithFields(map[string]interface{}{
//			"service": "api",
//		}),
//	)
type ConfigSetter func(opts *Options)

// newOptions creates a new Options instance with the given zerolog.Logger and ConfigSetters.
//
// This internal function initializes the Options struct with default values derived
// from the provided logger and applies all configuration setters in order.
//
// Parameters:
//   - log: base zerolog.Logger to derive initial settings from
//   - setters: slice of ConfigSetter functions to apply
//
// Returns:
//   - *Options: configured Options instance
func newOptions(log zerolog.Logger, setters []ConfigSetter) *Options {
	elvl, _ := MatchZeroLevel(log.GetLevel())

	opts := &Options{
		zcontext: log.With(),
		level:    elvl,
	}

	for _, set := range setters {
		set(opts)
	}

	return opts
}

// WithLevel returns a ConfigSetter that sets the logging level.
//
// This function configures the minimum log level that will be output.
// Messages below this level will be discarded for performance.
//
// Parameters:
//   - level: the logging level to set (log.DEBUG, log.INFO, log.WARN, log.ERROR, log.OFF)
//
// Returns:
//   - ConfigSetter: function that applies the level configuration
//
// Example:
//
//	// Set debug level for development
//	logger := echolog.New(os.Stdout, echolog.WithLevel(log.DEBUG))
//	
//	// Set error level for production
//	logger := echolog.New(os.Stdout, echolog.WithLevel(log.ERROR))
func WithLevel(level log.Lvl) ConfigSetter {
	return func(opts *Options) {
		zlvl, elvl := MatchEchoLevel(level)

		opts.zcontext = opts.zcontext.Logger().Level(zlvl).With()
		opts.level = elvl
	}
}

// WithField returns a ConfigSetter that adds a single field to all log messages.
//
// This function adds a key-value pair that will be included in every log message
// produced by the logger. This is useful for adding context like service name,
// version, or other identifying information.
//
// Parameters:
//   - name: the field name/key
//   - value: the field value (can be any type)
//
// Returns:
//   - ConfigSetter: function that applies the field configuration
//
// Example:
//
//	logger := echolog.New(os.Stdout,
//		echolog.WithField("service", "api"),
//		echolog.WithField("version", "1.0.0"),
//		echolog.WithField("instance_id", instanceID),
//	)
func WithField(name string, value any) ConfigSetter {
	return func(opts *Options) {
		opts.zcontext = opts.zcontext.Interface(name, value)
	}
}

// WithFields returns a ConfigSetter that adds multiple fields to all log messages.
//
// This function adds a map of key-value pairs that will be included in every
// log message produced by the logger. This is more efficient than multiple
// WithField calls when adding several fields.
//
// Parameters:
//   - fields: map containing key-value pairs to add to all log messages
//
// Returns:
//   - ConfigSetter: function that applies the fields configuration
//
// Example:
//
//	logger := echolog.New(os.Stdout,
//		echolog.WithFields(map[string]interface{}{
//			"service":     "user-api",
//			"version":     "2.1.0",
//			"environment": "production",
//			"region":      "us-west-2",
//		}),
//	)
func WithFields(fields map[string]any) ConfigSetter {
	return func(opts *Options) {
		opts.zcontext = opts.zcontext.Fields(fields)
	}
}

// WithTimestamp returns a ConfigSetter that adds timestamps to log messages.
//
// This function enables automatic timestamp generation for each log message.
// The timestamp is added as a "time" field in RFC3339 format.
//
// Returns:
//   - ConfigSetter: function that applies timestamp configuration
//
// Example:
//
//	// Logger with timestamps enabled
//	logger := echolog.New(os.Stdout, echolog.WithTimestamp())
//	logger.Info("message") // Output includes: "time":"2023-01-01T12:00:00Z"
func WithTimestamp() ConfigSetter {
	return func(opts *Options) {
		opts.zcontext = opts.zcontext.Timestamp()
	}
}

// WithCaller returns a ConfigSetter that adds caller information to log messages.
//
// This function enables automatic inclusion of the file name and line number
// where the log message was generated. The caller information is added as
// a "caller" field.
//
// Returns:
//   - ConfigSetter: function that applies caller information configuration
//
// Example:
//
//	// Logger with caller information enabled
//	logger := echolog.New(os.Stdout, echolog.WithCaller())
//	logger.Info("message") // Output includes: "caller":"main.go:25"
//
// Note: Adding caller information has a performance cost and should be used
// primarily during development or debugging.
func WithCaller() ConfigSetter {
	return func(opts *Options) {
		opts.zcontext = opts.zcontext.Caller()
	}
}

// WithCallerWithSkipFrameCount returns a ConfigSetter that adds caller information with custom frame skipping.
//
// This function is similar to WithCaller but allows you to specify how many
// stack frames to skip when determining the caller. This is useful when the
// logger is wrapped by other functions and you want to show the actual caller
// rather than the wrapper function.
//
// Parameters:
//   - skipFrameCount: number of stack frames to skip (0 = current function)
//
// Returns:
//   - ConfigSetter: function that applies caller configuration with frame skipping
//
// Example:
//
//	// Skip 2 frames to show the actual caller instead of wrapper functions
//	logger := echolog.New(os.Stdout, echolog.WithCallerWithSkipFrameCount(2))
//	
//	// When called from a wrapper function, this will show the original caller
//	func logWrapper(msg string) {
//		logger.Info(msg) // Shows caller of logWrapper, not logWrapper itself
//	}
func WithCallerWithSkipFrameCount(skipFrameCount int) ConfigSetter {
	return func(opts *Options) {
		opts.zcontext = opts.zcontext.CallerWithSkipFrameCount(skipFrameCount)
	}
}

// WithPrefix returns a ConfigSetter that adds a prefix to all log messages.
//
// This function sets a prefix string that will be included as a "prefix" field
// in all log messages. Prefixes are useful for identifying the source, component,
// or context of log messages.
//
// Parameters:
//   - prefix: the prefix string to add to all log messages
//
// Returns:
//   - ConfigSetter: function that applies the prefix configuration
//
// Example:
//
//	// Add a component prefix
//	logger := echolog.New(os.Stdout, echolog.WithPrefix("[AUTH]"))
//	logger.Info("user login") // Output includes: "prefix":"[AUTH]"
//	
//	// Add multiple context prefixes
//	apiLogger := echolog.New(os.Stdout, echolog.WithPrefix("[API]"))
//	dbLogger := echolog.New(os.Stdout, echolog.WithPrefix("[DB]"))
func WithPrefix(prefix string) ConfigSetter {
	return func(opts *Options) {
		opts.zcontext = opts.zcontext.Str("prefix", prefix)
	}
}

// WithHook returns a ConfigSetter that adds a zerolog.Hook to the logger.
//
// Hooks allow you to customize log processing by implementing the zerolog.Hook
// interface. Hooks can modify log events before they are written, filter events,
// or perform side effects like sending alerts.
//
// Parameters:
//   - hook: a zerolog.Hook implementation
//
// Returns:
//   - ConfigSetter: function that applies the hook configuration
//
// Example:
//
//	// Custom hook that adds hostname to all log entries
//	type HostnameHook struct{}
//	
//	func (h HostnameHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
//		hostname, _ := os.Hostname()
//		e.Str("hostname", hostname)
//	}
//	
//	logger := echolog.New(os.Stdout, echolog.WithHook(HostnameHook{}))
func WithHook(hook zerolog.Hook) ConfigSetter {
	return func(opts *Options) {
		opts.zcontext = opts.zcontext.Logger().Hook(hook).With()
	}
}

// WithHookFunc returns a ConfigSetter that adds a hook function to the logger.
//
// This is a convenience function for adding a hook using a function instead of
// implementing the zerolog.Hook interface. The function will be called for each
// log event and can modify the event before it's written.
//
// Parameters:
//   - hook: a function that matches the zerolog.HookFunc signature
//
// Returns:
//   - ConfigSetter: function that applies the hook function configuration
//
// Example:
//
//	// Add a hook function that adds request ID to all logs
//	hookFunc := func(e *zerolog.Event, level zerolog.Level, msg string) {
//		if requestID := getRequestID(); requestID != "" {
//			e.Str("request_id", requestID)
//		}
//	}
//	
//	logger := echolog.New(os.Stdout, echolog.WithHookFunc(hookFunc))
func WithHookFunc(hook zerolog.HookFunc) ConfigSetter {
	return func(opts *Options) {
		opts.zcontext = opts.zcontext.Logger().Hook(hook).With()
	}
}
