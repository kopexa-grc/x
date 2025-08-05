// Package cookie provides utilities for creating and configuring HTTP cookies
// with security best practices and flexible configuration options.
//
// # Overview
//
// The cookie package simplifies the creation of secure HTTP cookies by providing
// a configuration-driven approach to cookie creation. It supports all standard
// HTTP cookie attributes including security settings like Secure, HttpOnly,
// and SameSite.
//
// # Usage
//
// Basic cookie creation:
//
//	config := &cookie.Config{
//		Name:     "session",
//		Path:     "/",
//		MaxAge:   3600, // 1 hour
//		Secure:   true,
//		HTTPOnly: true,
//		SameSite: http.SameSiteLaxMode,
//	}
//
//	cookie := cookie.New("session", "abc123", config)
//
// # Security Considerations
//
// For production applications, it is recommended to:
// - Set Secure to true to ensure cookies are only sent over HTTPS
// - Set HTTPOnly to true to prevent JavaScript access
// - Use appropriate SameSite settings to prevent CSRF attacks
// - Set appropriate MaxAge values to limit cookie lifetime
//
// # Configuration
//
// The Config struct provides fine-grained control over cookie attributes:
// - Name: Default cookie name if not provided in New()
// - Domain: Cookie domain scope (defaults to responding server domain)
// - Path: Cookie path scope (defaults to responding URL path)
// - MaxAge: Cookie lifetime in seconds
// - Secure: HTTPS-only transmission
// - HTTPOnly: Prevents JavaScript access
// - SameSite: Cross-site request behavior
package cookie
