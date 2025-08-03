# Cache Package

A Redis client wrapper that provides configuration management and health checking capabilities for caching operations in Go applications. This package is part of [Kopexa](https://kopexa.com), a comprehensive GRC (Governance, Risk, and Compliance) software platform.

## Features

- **Configuration-driven setup**: Simple configuration struct for Redis client initialization
- **Connection pooling**: Built-in support for connection pool management
- **Health checking**: Ready-to-use health check functionality
- **Authentication support**: Username/password authentication
- **Timeout controls**: Configurable read, write, and dial timeouts
- **Retry logic**: Configurable retry attempts with backoff

## Installation

```bash
go get github.com/redis/go-redis/v9
```

## Quick Start

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/kopexa-grc/x/cache"
)

func main() {
    // Create configuration for the cache package
    config := cache.Config{
        Address:     "localhost:6379",
        DB:          0,
        DialTimeout: 5 * time.Second,
        MaxRetries:  3,
    }

    // Create Redis client using our cache package
    client := cache.New(config)
    defer client.Close()

    // Set up health checking
    healthFn := cache.Healthcheck(client)
    ctx := context.Background()
    
    // Check if Redis is healthy
    if err := healthFn(ctx); err != nil {
        log.Fatal("Redis connection failed:", err)
    }

    // Now use the Redis client (this returns a *redis.Client)
    err := client.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        log.Fatal(err)
    }

    val, err := client.Get(ctx, "key").Result()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Value: %s", val)
}
```

## Configuration

The `Config` struct provides comprehensive configuration options for Redis connections:

### Basic Configuration

```go
config := cache.Config{
    Enabled:  true,                    // Enable/disable Redis client
    Address:  "localhost:6379",        // Redis server address
    DB:       0,                       // Redis database number
}
```

### Authentication

```go
config := cache.Config{
    Address:  "redis.example.com:6379",
    Username: "myuser",                // Redis username (optional)
    Password: "mypassword",            // Redis password (optional)
    Name:     "myapp-client",          // Client name for monitoring
}
```

### Timeouts and Retries

```go
config := cache.Config{
    DialTimeout:  5 * time.Second,     // Connection timeout
    ReadTimeout:  3 * time.Second,     // Read operation timeout
    WriteTimeout: 3 * time.Second,     // Write operation timeout
    MaxRetries:   3,                   // Maximum retry attempts
}
```

### Connection Pooling

```go
config := cache.Config{
    MinIdleConns:   5,                 // Minimum idle connections
    MaxIdleConns:   10,                // Maximum idle connections
    MaxActiveConns: 100,               // Maximum active connections
}
```

## Configuration Reference

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Enabled` | `bool` | `true` | Enable/disable Redis client |
| `Address` | `string` | `"localhost:6379"` | Redis server host:port |
| `Name` | `string` | `""` | Client name for monitoring |
| `Username` | `string` | `""` | Redis username |
| `Password` | `string` | `""` | Redis password |
| `DB` | `int` | `0` | Redis database number |
| `DialTimeout` | `time.Duration` | `5s` | Connection establishment timeout |
| `ReadTimeout` | `time.Duration` | `0` | Socket read timeout (0=default 3s, -1=no timeout, -2=disable) |
| `WriteTimeout` | `time.Duration` | `0` | Socket write timeout (0=default 3s, -1=no timeout, -2=disable) |
| `MaxRetries` | `int` | `3` | Maximum retry attempts (-1=disable retries) |
| `MinIdleConns` | `int` | `0` | Minimum idle connections in pool |
| `MaxIdleConns` | `int` | `0` | Maximum idle connections in pool |
| `MaxActiveConns` | `int` | `0` | Maximum active connections (0=unlimited) |

## Health Checking

The package provides a built-in health check function that can be integrated into your application's health monitoring:

```go
import "github.com/kopexa-grc/x/cache"

// Initialize with cache package
config := cache.Config{Address: "localhost:6379"}
client := cache.New(config)
healthFn := cache.Healthcheck(client)

// Use in HTTP health endpoint
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    if err := healthFn(ctx); err != nil {
        http.Error(w, "Redis unhealthy: "+err.Error(), http.StatusServiceUnavailable)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
})
```

## Advanced Usage

### With Configuration from Environment Variables

If you're using a configuration library like [koanf](https://github.com/knadh/koanf), the struct tags are already set up:

```go
import (
    "github.com/knadh/koanf/v2"
    "github.com/knadh/koanf/providers/env"
)

var config cache.Config
k := koanf.New(".")
k.Load(env.Provider("REDIS_", ".", func(s string) string {
    return strings.Replace(strings.ToLower(s), "redis_", "", 1)
}), nil)
k.Unmarshal("", &config)

client := cache.New(config)
```

### With JSON Configuration

```go
import "encoding/json"

configJSON := `{
    "address": "redis.example.com:6379",
    "db": 1,
    "username": "myuser",
    "password": "mypassword",
    "dialTimeout": "10s",
    "maxRetries": 5
}`

var config cache.Config
json.Unmarshal([]byte(configJSON), &config)
client := cache.New(config)
```

## Error Handling

The package uses the standard go-redis error handling. Common patterns:

```go
// Check if key exists
val, err := client.Get(ctx, "key").Result()
if err == redis.Nil {
    // Key does not exist
    log.Println("Key does not exist")
} else if err != nil {
    // Other error
    log.Fatal(err)
} else {
    // Key exists
    log.Printf("Value: %s", val)
}
```

## Dependencies

- [go-redis/redis/v9](https://github.com/redis/go-redis): Redis client for Go

## About Kopexa

This package is part of [Kopexa](https://kopexa.com), a comprehensive GRC (Governance, Risk, and Compliance) software platform that helps organizations manage their compliance requirements, assess risks, and maintain governance standards.

## License

This package is part of the Kopexa GRC platform and follows the same license terms.