// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package logger provides a comprehensive logging toolkit built on zerolog
// for high-performance, structured logging in Go applications.
//
// The package offers multiple output modes including CLI-optimized console output,
// JSON logging, and GCP-compatible formatting. It features buffered writers for
// TUI applications, context-aware logging with request IDs, performance tracing,
// and extensive debugging utilities.
//
// # Quick Start
//
// For CLI applications:
//
//	logger.Set("info")
//	logger.CliLogger()
//	log.Info().Msg("Application started")
//
// For production JSON logging:
//
//	logger.UseJSONLogging(os.Stdout)
//	logger.Set("info")
//	log.Info().Str("service", "my-api").Msg("Service started")
//
// # Context-Aware Logging
//
// Create request-scoped logging contexts:
//
//	ctx := logger.RequestScopedContext(context.Background(), "req-123")
//	log := logger.FromContext(ctx)
//	log.Info().Msg("Processing request") // Includes request ID
//
// # Performance Tracing
//
// Track function execution times:
//
//	func myFunction() {
//		defer logger.FuncDur(time.Now(), "mypackage.myFunction")
//		// Function logic here
//	}
//
// # Debug Utilities
//
// Pretty-print JSON for debugging:
//
//	logger.DebugJSON(complexObject) // Only outputs at debug level
//	logger.DebugDumpJSON("debug", obj) // Dumps to file when DEBUG=1
//
// The package is part of the Kopexa GRC platform ecosystem and provides
// enterprise-grade logging capabilities with minimal performance overhead.
package logger
