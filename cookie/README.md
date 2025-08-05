# Cookie Package

A Go package for creating secure HTTP cookies with flexible configuration options.

## Overview

The `cookie` package provides a simple and secure way to create HTTP cookies with proper configuration for modern web applications. It emphasizes security best practices while maintaining flexibility for various use cases.

## Features

- **Security-first design**: Built-in support for Secure, HttpOnly, and SameSite attributes
- **Flexible configuration**: Configure cookies through a simple Config struct
- **Automatic expiration handling**: Converts MaxAge to Expires time automatically
- **Production-ready**: Follows security best practices for web applications
- **Zero dependencies**: Uses only Go standard library

## Installation

```bash
go get github.com/kopexa-grc/x/cookie
```

## Quick Start

```go
package main

import (
    "net/http"
    "github.com/kopexa-grc/x/cookie"
)

func main() {
    // Create a secure session cookie
    config := &cookie.Config{
        Name:     "session",
        Path:     "/",
        MaxAge:   3600, // 1 hour
        Secure:   true,
        HTTPOnly: true,
        SameSite: http.SameSiteLaxMode,
    }

    // Create the cookie
    sessionCookie := cookie.New("", "user_session_token", config)
    
    // Use in HTTP handler
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        http.SetCookie(w, sessionCookie)
        w.Write([]byte("Cookie set!"))
    })
}
```

## Configuration Options

The `Config` struct provides comprehensive control over cookie attributes:

```go
type Config struct {
    Name     string           // Default cookie name
    Domain   string           // Cookie domain scope
    Path     string           // Cookie path scope
    MaxAge   int              // Cookie lifetime in seconds
    Secure   bool             // HTTPS-only transmission
    HTTPOnly bool             // Prevent JavaScript access
    SameSite http.SameSite    // Cross-site request behavior
}
```

### MaxAge Behavior

- `MaxAge > 0`: Cookie expires in the specified number of seconds
- `MaxAge = 0`: No Max-Age attribute (session cookie)
- `MaxAge < 0`: Cookie is deleted immediately (expires at Unix epoch)

### SameSite Options

- `http.SameSiteDefaultMode`: Browser default behavior
- `http.SameSiteLaxMode`: Sent with same-site requests and top-level navigation
- `http.SameSiteStrictMode`: Only sent with same-site requests
- `http.SameSiteNoneMode`: Sent with all requests (requires Secure=true)

## Usage Examples

### Basic Session Cookie

```go
config := &cookie.Config{
    Name:     "session",
    Path:     "/",
    MaxAge:   86400, // 24 hours
    Secure:   true,
    HTTPOnly: true,
    SameSite: http.SameSiteLaxMode,
}

sessionCookie := cookie.New("", "abc123def456", config)
```

### Remember Me Cookie

```go
config := &cookie.Config{
    Name:     "remember_token",
    Path:     "/",
    MaxAge:   2592000, // 30 days
    Secure:   true,
    HTTPOnly: true,
    SameSite: http.SameSiteStrictMode,
}

rememberCookie := cookie.New("", "long_lived_token", config)
```

### Delete Cookie

```go
config := &cookie.Config{
    Name:   "session",
    Path:   "/",
    MaxAge: -1, // Delete immediately
}

deleteCookie := cookie.New("", "", config)
```

### API Token Cookie

```go
config := &cookie.Config{
    Name:     "api_token",
    Path:     "/api",
    Domain:   ".example.com",
    MaxAge:   3600,
    Secure:   true,
    HTTPOnly: true,
    SameSite: http.SameSiteStrictMode,
}

apiCookie := cookie.New("", "api_token_value", config)
```

## Security Best Practices

### Production Recommendations

1. **Always use Secure=true** in production to ensure cookies are only sent over HTTPS
2. **Set HTTPOnly=true** to prevent JavaScript access and XSS attacks
3. **Use appropriate SameSite settings** to prevent CSRF attacks
4. **Set reasonable MaxAge values** to limit exposure if compromised
5. **Use specific Path and Domain** to limit cookie scope

### Example Secure Configuration

```go
// Recommended production configuration
config := &cookie.Config{
    Name:     "secure_session",
    Path:     "/",
    Domain:   ".yourdomain.com",
    MaxAge:   3600,
    Secure:   true,
    HTTPOnly: true,
    SameSite: http.SameSiteStrictMode,
}
```

## API Reference

### Functions

#### `New(name, value string, config *Config) *http.Cookie`

Creates a new HTTP cookie with the specified name, value, and configuration.

**Parameters:**
- `name`: Cookie name (if empty, uses `config.Name`)
- `value`: Cookie value
- `config`: Cookie configuration

**Returns:**
- `*http.Cookie`: Configured cookie, or `nil` if no name is provided

**Example:**
```go
cookie := cookie.New("session", "value123", &cookie.Config{
    MaxAge:   3600,
    Secure:   true,
    HTTPOnly: true,
})
```

## Testing

Run the test suite:

```bash
go test ./cookie
```

Run tests with coverage:

```bash
go test -cover ./cookie
```

Run benchmarks:

```bash
go test -bench=. ./cookie
```

## Contributing

Contributions are welcome! Please ensure that:

1. All tests pass
2. Code coverage remains high
3. Security best practices are maintained
4. Documentation is updated for new features

## License

This project is licensed under the Business Source License 1.1 (BUSL-1.1) - see the LICENSE file for details.

**Summary**: This is a source-available license that allows non-commercial use, viewing, and analysis of the code. Commercial use requires separate licensing. For commercial licenses or individual agreements, please contact: hello@kopexa.com

## Changelog

See [CHANGELOG.md](../CHANGELOG.md) for version history and changes.
