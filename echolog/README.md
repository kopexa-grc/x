# echolog

A powerful [Zerolog](https://github.com/rs/zerolog) wrapper for the [Echo](https://echo.labstack.com/) web framework that provides structured JSON logging with minimal performance overhead within [kopexa](https://kopexa.com)'s ecosystem.

Credits: https://github.com/ziflex/lecho

## Features

- **Full Echo Logger Interface Compatibility** - Drop-in replacement for Echo's default logger  
- **High Performance** - Built on zerolog's zero-allocation JSON logger  
- **Flexible Configuration** - Extensive options for customization  
- **Rich Middleware** - Request logging with request ID tracking and context enrichment  
- **Smart Error Handling** - Configurable error propagation and latency-based log levels  
- **Context Integration** - Seamless integration with Go's context package  

## Installation

```bash
go get github.com/kopexa-grc/x/echolog
```

## Quick Start

### Basic Usage

```go
package main

import (
    "os"
    
    "github.com/kopexa-grc/x/echolog"
    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    e.Logger = echolog.New(os.Stdout)
    
    e.GET("/", func(c echo.Context) error {
        c.Logger().Info("Hello from Echo!")
        return c.String(200, "Hello, World!")
    })
    
    e.Logger.Fatal(e.Start(":8080"))
}
```

### Using Existing Zerolog Instance

```go
package main

import (
    "os"
    
    "github.com/kopexa-grc/x/echolog"
    "github.com/labstack/echo/v4"
    "github.com/rs/zerolog"
)

func main() {
    // Create a zerolog instance with custom configuration
    log := zerolog.New(os.Stdout).With().
        Timestamp().
        Str("service", "my-api").
        Logger()
    
    e := echo.New()
    e.Logger = echolog.From(log)
    
    e.Logger.Fatal(e.Start(":8080"))
}
```

## Configuration Options

The logger supports extensive configuration through functional options:

```go
package main

import (
    "os"
    
    "github.com/kopexa-grc/x/echolog"
    "github.com/labstack/echo/v4"
    "github.com/labstack/gommon/log"
)

func main() {
    e := echo.New()
    e.Logger = echolog.New(
        os.Stdout,
        echolog.WithLevel(log.DEBUG),
        echolog.WithFields(map[string]interface{}{
            "service": "my-api",
            "version": "1.0.0",
        }),
        echolog.WithTimestamp(),
        echolog.WithCaller(),
        echolog.WithPrefix("=� API"),
    )
    
    e.Logger.Fatal(e.Start(":8080"))
}
```

### Available Options

- `WithLevel(level)` - Set log level (DEBUG, INFO, WARN, ERROR)
- `WithField(name, value)` - Add a single field to all log entries
- `WithFields(map[string]interface{})` - Add multiple fields to all log entries
- `WithTimestamp()` - Add timestamp to log entries
- `WithCaller()` - Add caller information (file:line)
- `WithCallerWithSkipFrameCount(skip)` - Add caller with custom skip count
- `WithPrefix(prefix)` - Add a prefix to log entries
- `WithHook(hook)` - Add a zerolog hook
- `WithHookFunc(hookFunc)` - Add a zerolog hook function

## Middleware

### Basic Request Logging

```go
package main

import (
    "os"
    
    "github.com/kopexa-grc/x/echolog"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/labstack/gommon/log"
)

func main() {
    e := echo.New()
    
    logger := echolog.New(
        os.Stdout,
        echolog.WithLevel(log.DEBUG),
        echolog.WithTimestamp(),
        echolog.WithCaller(),
    )
    e.Logger = logger
    
    // Add request ID middleware first
    e.Use(middleware.RequestID())
    
    // Add logging middleware
    e.Use(echolog.LoggingMiddleware(echolog.Config{
        Logger: logger,
    }))
    
    e.GET("/", func(c echo.Context) error {
        // Both Echo and zerolog interfaces work
        c.Logger().Info("Using Echo interface")
        echolog.Ctx(c.Request().Context()).Info().Msg("Using zerolog interface")
        
        return c.String(200, "Hello, World!")
    })
    
    e.Logger.Fatal(e.Start(":8080"))
}
```

**Output:**
```json
{"level":"info","request_id":"abc123","remote_ip":"127.0.0.1","host":"localhost:8080","method":"GET","uri":"/","user_agent":"curl/7.68.0","status":200,"referer":"","latency_human":"156.291�s","client_ip":"127.0.0.1","request_protocol":"HTTP/1.1","bytes_in":"0","bytes_out":"13","message":"request details"}
```

### Latency-Based Log Levels

Automatically escalate log level for slow requests:

```go
e.Use(echolog.LoggingMiddleware(echolog.Config{
    Logger: logger,
    RequestLatencyLevel: zerolog.WarnLevel,
    RequestLatencyLimit: 500 * time.Millisecond,
}))
```

### Nested Logging

Group request details under a sub-dictionary:

```go
e.Use(echolog.LoggingMiddleware(echolog.Config{
    Logger: logger,
    NestKey: "request",
}))
```

**Output:**
```json
{"level":"info","request":{"remote_ip":"127.0.0.1","method":"GET","uri":"/","status":200},"message":"request details"}
```

### Request Enrichment

Add custom fields to log entries using the Enricher function:

```go
e.Use(echolog.LoggingMiddleware(echolog.Config{
    Logger: logger,
    Enricher: func(c echo.Context, logger zerolog.Context) zerolog.Context {
        // Add user ID if available
        if userID := c.Get("user_id"); userID != nil {
            return logger.Str("user_id", userID.(string))
        }
        
        // Add request size
        return logger.Int64("request_size", c.Request().ContentLength)
    },
}))
```

### Error Handling

Control whether errors are propagated to Echo's error handler:

```go
e.Use(echolog.LoggingMiddleware(echolog.Config{
    Logger: logger,
    HandleError: true, // Propagate errors to Echo's error handler
}))
```

## Advanced Usage

### Context Integration

Access the logger from any context:

```go
func myHandler(c echo.Context) error {
    // Get logger from request context
    logger := echolog.Ctx(c.Request().Context())
    
    logger.Info().
        Str("action", "processing_request").
        Int("user_id", 123).
        Msg("Processing user request")
    
    return c.JSON(200, map[string]string{"status": "ok"})
}
```

### Custom Request ID Headers

Configure custom request ID headers:

```go
e.Use(echolog.LoggingMiddleware(echolog.Config{
    Logger: logger,
    RequestIDHeader: "X-Correlation-ID",
    RequestIDKey: "correlation_id",
}))
```

### Multiple Output Writers

Log to multiple destinations:

```go
import "github.com/rs/zerolog"

consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
multiWriter := zerolog.MultiLevelWriter(consoleWriter, os.Stderr)

logger := echolog.New(multiWriter)
```

### File Logging with Rotation

Using [Lumberjack](https://github.com/natefinch/lumberjack) for log rotation:

```go
import "gopkg.in/natefinch/lumberjack.v2"

logger := echolog.New(&lumberjack.Logger{
    Filename:   "/var/log/myapp/app.log",
    MaxSize:    500, // megabytes
    MaxBackups: 3,
    MaxAge:     28, // days
    Compress:   true,
})
```

## Level Conversion Utilities

Convert between Echo and Zerolog log levels:

```go
import (
    "fmt"
    
    "github.com/kopexa-grc/x/echolog"
    "github.com/labstack/gommon/log"
    "github.com/rs/zerolog"
)

func main() {
    // Convert Echo level to Zerolog level
    zeroLevel, echoLevel := echolog.MatchEchoLevel(log.WARN)
    fmt.Println("Zerolog:", zeroLevel, "Echo:", echoLevel)
    
    // Convert Zerolog level to Echo level
    echoLevel, zeroLevel = echolog.MatchZeroLevel(zerolog.InfoLevel)
    fmt.Println("Echo:", echoLevel, "Zerolog:", zeroLevel)
}
```

## Complete Example

Here's a comprehensive example showcasing most features:

```go
package main

import (
    "context"
    "os"
    "time"
    
    "github.com/kopexa-grc/x/echolog"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/labstack/gommon/log"
    "github.com/rs/zerolog"
)

func main() {
    e := echo.New()
    
    // Create logger with comprehensive configuration
    logger := echolog.New(
        os.Stdout,
        echolog.WithLevel(log.DEBUG),
        echolog.WithTimestamp(),
        echolog.WithCaller(),
        echolog.WithFields(map[string]interface{}{
            "service": "example-api",
            "version": "1.0.0",
        }),
    )
    e.Logger = logger
    
    // Middleware stack
    e.Use(middleware.RequestID())
    e.Use(echolog.LoggingMiddleware(echolog.Config{
        Logger: logger,
        RequestLatencyLevel: zerolog.WarnLevel,
        RequestLatencyLimit: 100 * time.Millisecond,
        HandleError: true,
        Enricher: func(c echo.Context, logger zerolog.Context) zerolog.Context {
            return logger.Str("endpoint", c.Path())
        },
    }))
    
    // Routes
    e.GET("/", homeHandler)
    e.GET("/users/:id", getUserHandler)
    e.POST("/users", createUserHandler)
    
    e.Logger.Fatal(e.Start(":8080"))
}

func homeHandler(c echo.Context) error {
    c.Logger().Info("Home endpoint accessed")
    return c.JSON(200, map[string]string{"message": "Welcome!"})
}

func getUserHandler(c echo.Context) error {
    userID := c.Param("id")
    logger := echolog.Ctx(c.Request().Context())
    
    logger.Info().Str("user_id", userID).Msg("Fetching user")
    
    // Simulate user lookup
    user := map[string]interface{}{
        "id":   userID,
        "name": "John Doe",
    }
    
    return c.JSON(200, user)
}

func createUserHandler(c echo.Context) error {
    logger := echolog.Ctx(c.Request().Context())
    
    // Log the creation attempt
    logger.Info().Msg("Creating new user")
    
    // Simulate user creation
    newUser := map[string]interface{}{
        "id":   "123",
        "name": "Jane Doe",
    }
    
    logger.Info().
        Str("user_id", "123").
        Msg("User created successfully")
    
    return c.JSON(201, newUser)
}
```

## Performance

echolog is built on top of zerolog, which is designed for high performance:

- Zero allocation JSON logger
- Lazy evaluation of expensive operations
- Structured logging without reflection
- Minimal overhead over raw zerolog

## License

Licensed under the Business Source License 1.1 (BUSL-1.1).

Copyright (c) Kopexa GmbH