# Logger Package

A comprehensive logging toolkit built on [Zerolog](https://github.com/rs/zerolog) that provides structured, high-performance logging for Go applications within [Kopexa](https://kopexa.com)'s ecosystem. This package offers CLI-optimized output, buffered writers, context-aware logging, and debugging utilities.

## Features

- **Multiple Output Modes** - CLI-optimized console output, JSON logging, and GCP-compatible JSON formatting
- **CLI-Optimized Display** - Custom console writer with colors, compact formatting, and cross-platform support
- **Buffered Output** - Thread-safe buffered writer for TUI applications that need to pause/resume log output
- **Context Integration** - Request-scoped logging with automatic request ID generation and custom tags
- **Performance Tracing** - Function execution time tracking with minimal overhead
- **Debug Utilities** - Pretty-printed JSON/YAML output and file dumping for development
- **Environment Detection** - Automatic log level configuration from DEBUG/TRACE environment variables
- **Cross-Platform** - Windows-compatible color output and formatting

## Installation

```bash
go get github.com/kopexa-grc/x/logger
```

## Quick Start

### Basic CLI Logging

```go
package main

import (
    "github.com/kopexa-grc/x/logger"
    "github.com/rs/zerolog/log"
)

func main() {
    // Set up CLI-optimized logging
    logger.Set("info")
    logger.CliLogger()
    
    log.Info().Msg("Application started")
    log.Debug().Str("version", "1.0.0").Msg("Debug info")
}
```

### JSON Logging for Production

```go
package main

import (
    "os"
    
    "github.com/kopexa-grc/x/logger"
    "github.com/rs/zerolog/log"
)

func main() {
    // Set up structured JSON logging
    logger.UseJSONLogging(os.Stdout)
    logger.Set("info")
    
    log.Info().
        Str("service", "my-api").
        Str("version", "1.0.0").
        Msg("Service started")
}
```

### GCP-Compatible Logging

```go
package main

import (
    "os"
    
    "github.com/kopexa-grc/x/logger"
    "github.com/rs/zerolog/log"
)

func main() {
    // Set up GCP-compatible JSON logging
    logger.UseGCPJSONLogging(os.Stdout)
    logger.Set("info")
    
    log.Info().
        Str("service", "my-service").
        Msg("Service started")
}
```

## Configuration

### Log Levels

The package supports standard log levels with easy configuration:

```go
// Set log level from string
logger.Set("debug")    // Sets global level to debug
logger.Set("info")     // Sets global level to info (default)
logger.Set("warn")     // Sets global level to warn
logger.Set("error")    // Sets global level to error
logger.Set("trace")    // Sets global level to trace

// Get current log level
level := logger.GetLevel()
fmt.Printf("Current log level: %s", level)
```

### Environment-Based Configuration

```go
// Automatically detect log level from environment variables
if level, ok := logger.GetEnvLogLevel(); ok {
    logger.Set(level)
}

// This checks for:
// DEBUG=true or DEBUG=1 -> sets debug level
// TRACE=true or TRACE=1 -> sets trace level
```

### Test Environment Setup

```go
func TestMyFunction(t *testing.T) {
    // Set up logging for tests
    logger.InitTestEnv()
    
    // Your test code here
}
```

## Output Modes

### CLI Output (Console Writer)

Perfect for command-line applications with colorized, human-readable output:

```go
// Full CLI output with colors and symbols
logger.CliLogger()

// Compact CLI output
logger.CliCompactLogger(os.Stderr)

// Standard zerolog console writer
logger.StandardZerologLogger()
```

**CLI Output Features:**
- Colorized log levels with symbols (→, !, ×)
- Relative file paths from current working directory
- Cross-platform color support (including Windows)
- Compact mode for space-constrained environments

### JSON Output

Structured logging for production environments:

```go
// Standard JSON logging
logger.UseJSONLogging(os.Stdout)

// GCP-compatible JSON (with proper field names)
logger.UseGCPJSONLogging(os.Stdout)
```

**GCP JSON Output includes:**
- `severity` instead of `level`
- `timestamp` with RFC3339Nano format
- Compatible with Google Cloud Logging

## Context-Aware Logging

### Request-Scoped Logging

Automatically add request IDs to all log entries within a request context:

```go
package main

import (
    "context"
    
    "github.com/kopexa-grc/x/logger"
)

func handleRequest(reqID string) {
    // Create request-scoped context
    ctx := logger.RequestScopedContext(context.Background(), reqID)
    
    // Get logger from context
    log := logger.FromContext(ctx)
    
    // All log entries will include the request ID
    log.Info().Msg("Processing request")
    log.Debug().Str("step", "validation").Msg("Validating input")
}
```

### Custom Tags

Add custom tags to log entries within a context:

```go
func processUser(ctx context.Context, userID string) {
    // Add custom tag to context
    logger.AddTag(ctx, "user_id", userID)
    logger.AddTag(ctx, "operation", "user_update")
    
    // Get logger with tags
    log := logger.FromContext(ctx)
    
    // Log entries will include custom tags under "ctags"
    log.Info().Msg("Processing user update")
}
```

## Buffered Output

For TUI applications that need to control when logs are displayed:

```go
package main

import (
    "os"
    
    "github.com/kopexa-grc/x/logger"
    "github.com/rs/zerolog/log"
)

func main() {
    // Use the default buffered writer
    logger.SetWriter(logger.LogOutputWriter)
    
    // In your TUI application:
    // Pause logging output
    if bw, ok := logger.LogOutputWriter.(*logger.BufferedWriter); ok {
        bw.Pause()
        
        // Show TUI interface...
        
        // Resume logging output
        bw.Resume()
    }
}
```

**BufferedWriter Methods:**
- `Pause()` - Buffer log output instead of writing directly
- `Resume()` - Flush buffered content and resume direct writing
- Thread-safe with read-write mutex protection

## Performance Tracing

Track function execution times with minimal overhead:

```go
package main

import (
    "time"
    
    "github.com/kopexa-grc/x/logger"
)

func expensiveOperation() {
    defer logger.FuncDur(time.Now(), "mypackage.expensiveOperation")
    
    // Your function logic here
    time.Sleep(100 * time.Millisecond)
}

// Output (at trace level):
// TRC func=mypackage.expensiveOperation took=100.234ms logger.FuncDur>
```

## Debug Utilities

### Pretty JSON Output

Display formatted JSON in CLI for debugging:

```go
package main

import "github.com/kopexa-grc/x/logger"

func debugAPI() {
    data := map[string]interface{}{
        "user":    "john_doe",
        "status":  "active",
        "profile": map[string]string{
            "email": "john@example.com",
            "role":  "admin",
        },
    }
    
    // Only outputs when debug level is enabled
    logger.DebugJSON(data)
    
    // Only outputs when trace level is enabled
    logger.TraceJSON(data)
    
    // Get pretty JSON string
    prettyStr := logger.PrettyJSON(data)
    fmt.Println(prettyStr)
}
```

### Debug File Dumps

Automatically dump objects to files when debugging:

```go
package main

import "github.com/kopexa-grc/x/logger"

func debugWithDumps() {
    data := map[string]interface{}{
        "request":  "complex_data",
        "response": "detailed_response",
    }
    
    // Dumps to ./mondoo-debug-request.json when DEBUG=1 or TRACE=1
    logger.DebugDumpJSON("request", data)
    
    // Dumps to ./mondoo-debug-response.yaml when DEBUG=1 or TRACE=1  
    logger.DebugDumpYAML("response", data)
}
```

**Debug Dump Features:**
- Only creates files when DEBUG=1/true or TRACE=1/true
- Automatic file naming with configurable prefix
- Support for both JSON and YAML formats
- Properly indented output for readability

## Advanced Usage

### Custom Writers

Use custom output writers for specialized logging needs:

```go
package main

import (
    "bytes"
    "os"
    
    "github.com/kopexa-grc/x/logger"
)

func customWriter() {
    // Log to memory buffer
    var buf bytes.Buffer
    logger.SetWriter(&buf)
    
    // Log to multiple destinations
    multiWriter := io.MultiWriter(os.Stdout, &buf)
    logger.SetWriter(multiWriter)
    
    // Use custom buffered writer
    customBuffered := logger.NewBufferedWriter(os.Stderr)
    logger.SetWriter(customBuffered)
}
```

### Integration with HTTP Middleware

Example integration with HTTP request logging:

```go
package main

import (
    "net/http"
    
    "github.com/kopexa-grc/x/logger"
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get or generate request ID
        reqID := r.Header.Get("X-Request-ID")
        
        // Create request-scoped context
        ctx := logger.RequestScopedContext(r.Context(), reqID)
        
        // Add request tags
        logger.AddTag(ctx, "method", r.Method)
        logger.AddTag(ctx, "path", r.URL.Path)
        
        // Get logger and log request
        log := logger.FromContext(ctx)
        log.Info().Msg("Request started")
        
        // Continue with scoped context
        next.ServeHTTP(w, r.WithContext(ctx))
        
        log.Info().Msg("Request completed")
    })
}
```

### Environment Configuration

```go
package main

import (
    "os"
    
    "github.com/kopexa-grc/x/logger"
)

func setupFromEnv() {
    // Check environment variables
    if level, ok := logger.GetEnvLogLevel(); ok {
        logger.Set(level)
    } else {
        // Default to info level
        logger.Set("info")
    }
    
    // Configure output based on environment
    if os.Getenv("LOG_FORMAT") == "json" {
        if os.Getenv("CLOUD_PROVIDER") == "gcp" {
            logger.UseGCPJSONLogging(os.Stdout)
        } else {
            logger.UseJSONLogging(os.Stdout)
        }
    } else {
        logger.CliLogger()
    }
}
```

## API Reference

### Global Configuration

| Function | Description |
|----------|-------------|
| `Set(level string)` | Set global log level (error, warn, info, debug, trace) |
| `GetLevel() string` | Get current global log level |
| `GetEnvLogLevel() (string, bool)` | Get log level from DEBUG/TRACE env vars |
| `InitTestEnv()` | Configure logging for test environment |

### Output Configuration

| Function | Description |
|----------|-------------|
| `SetWriter(w io.Writer)` | Set custom writer for global logger |
| `UseJSONLogging(out io.Writer)` | Enable JSON logging |
| `UseGCPJSONLogging(out io.Writer)` | Enable GCP-compatible JSON logging |
| `CliLogger()` | Enable CLI-optimized console output |
| `CliCompactLogger(out io.Writer)` | Enable compact CLI output |
| `StandardZerologLogger()` | Use standard zerolog console writer |

### Context Operations

| Function | Description |
|----------|-------------|
| `RequestScopedContext(ctx context.Context, reqID string) context.Context` | Create request-scoped logging context |
| `FromContext(ctx context.Context) *zerolog.Logger` | Get logger from context |
| `WithTagsContext(ctx context.Context) context.Context` | Add tags support to context |
| `AddTag(ctx context.Context, name, value string)` | Add custom tag to context |
| `GetTags(ctx context.Context) map[string]string` | Get all tags from context |

### Buffered Writer

| Function | Description |
|----------|-------------|
| `NewBufferedWriter(out io.Writer) io.Writer` | Create new buffered writer |
| `(*BufferedWriter).Pause()` | Pause output and start buffering |
| `(*BufferedWriter).Resume()` | Resume output and flush buffer |

### Performance & Debug

| Function | Description |
|----------|-------------|
| `FuncDur(start time.Time, name string)` | Log function execution time (use with defer) |
| `DebugJSON(obj interface{})` | Pretty print JSON at debug level |
| `TraceJSON(obj interface{})` | Pretty print JSON at trace level |
| `PrettyJSON(obj interface{}) string` | Get pretty JSON string |
| `DebugDumpJSON(name string, obj interface{})` | Dump JSON to file in debug mode |
| `DebugDumpYAML(name string, obj interface{})` | Dump YAML to file in debug mode |

### Global Variables

| Variable | Type | Description |
|----------|------|-------------|
| `LogOutputWriter` | `io.Writer` | Default buffered writer (stderr) |
| `Debug` | `bool` | Debug mode flag |
| `DumpLocal` | `string` | File dump prefix for debug dumps |

## Dependencies

- [rs/zerolog](https://github.com/rs/zerolog): High-performance structured logging
- [muesli/termenv](https://github.com/muesli/termenv): Terminal environment detection and styling
- [google/uuid](https://github.com/google/uuid): UUID generation for request IDs
- [hokaccha/go-prettyjson](https://github.com/hokaccha/go-prettyjson): Pretty JSON formatting
- [sigs.k8s.io/yaml](https://sigs.k8s.io/yaml): YAML marshaling

## About Kopexa

This package is part of [Kopexa](https://kopexa.com), a comprehensive GRC (Governance, Risk, and Compliance) software platform that helps organizations manage their compliance requirements, assess risks, and maintain governance standards.

## License

Licensed under the Business Source License 1.1 (BUSL-1.1).

Copyright (c) Kopexa GmbH