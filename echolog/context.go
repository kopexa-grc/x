// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package echolog

import (
	"context"

	"github.com/rs/zerolog"
)

// WithContext returns a new context.Context with the logger attached.
//
// This method embeds the logger into the provided context, making it available
// to downstream functions via the Ctx function. The logger can be retrieved
// anywhere in the call chain where the context is available.
//
// Parameters:
//   - ctx: the context to attach the logger to
//
// Returns:
//   - context.Context: new context with the logger embedded
//
// Example:
//
//	logger := echolog.New(os.Stdout, echolog.WithTimestamp())
//	ctx := logger.WithContext(context.Background())
//
//	// Pass context to other functions
//	processRequest(ctx)
//
//	func processRequest(ctx context.Context) {
//		// Retrieve logger from context
//		log := echolog.Ctx(ctx)
//		log.Info().Msg("processing request")
//	}
func (l Logger) WithContext(ctx context.Context) context.Context {
	zerologger := l.Unwrap()

	// Check if zerologger is uninitialized by comparing its level
	if zerologger.GetLevel() == zerolog.NoLevel {
		return ctx // Return the original context if zerologger is uninitialized
	}

	return zerologger.WithContext(ctx)
}

// Ctx retrieves a zerolog.Logger from the provided context.
//
// This function extracts a logger that was previously attached to the context
// using WithContext. If no logger is found in the context or if the context
// is nil, it returns a default logger instance.
//
// This is particularly useful in middleware and request handlers where you
// want to access a logger that has been enriched with request-specific
// information like request IDs.
//
// Parameters:
//   - ctx: the context to extract the logger from
//
// Returns:
//   - *zerolog.Logger: logger instance from context or a default logger
//
// Example:
//
//	func handleRequest(c echo.Context) error {
//		// Get logger from request context (set by middleware)
//		log := echolog.Ctx(c.Request().Context())
//
//		log.Info().
//			Str("endpoint", c.Path()).
//			Str("method", c.Request().Method).
//			Msg("processing request")
//
//		return c.JSON(200, map[string]string{"status": "ok"})
//	}
//
// Note: This function always returns a valid logger, never nil.
func Ctx(ctx context.Context) *zerolog.Logger {
	if ctx == nil {
		defaultLogger := zerolog.New(nil) // Create a default logger
		return &defaultLogger
	}

	return zerolog.Ctx(ctx)
}
