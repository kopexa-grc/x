package cookie

import (
	"net/http"
	"time"
)

// Config configures http.Cookie creation.
//
// Config provides comprehensive control over HTTP cookie attributes including
// security settings, scope, and lifetime. It supports all standard cookie
// attributes defined in RFC 6265.
//
// Example usage:
//
//	config := &Config{
//		Name:     "session",
//		Path:     "/",
//		MaxAge:   3600,
//		Secure:   true,
//		HTTPOnly: true,
//		SameSite: http.SameSiteLaxMode,
//	}
type Config struct {
	// Name specifies the default cookie name.
	// If empty, the name parameter in New() must be provided.
	Name string

	// Domain specifies the cookie domain scope.
	// Leave empty for requested resource scope.
	// Defaults to the domain name of the responding server when unset.
	Domain string

	// Path specifies the cookie path scope.
	// Leave empty for requested resource scope.
	// Defaults to the path of the responding URL when unset.
	Path string

	// MaxAge controls the cookie lifetime in seconds.
	// MaxAge=0 means no 'Max-Age' attribute specified (session cookie).
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge int

	// Secure indicates that the cookie may only be transferred over HTTPS.
	// Recommended to be true for production applications.
	Secure bool

	// HTTPOnly indicates that the browser should prohibit non-HTTP
	// (i.e. JavaScript) cookie access. Recommended to be true to prevent XSS attacks.
	HTTPOnly bool

	// SameSite controls whether the cookie is sent with cross-site requests.
	// Use SameSiteLaxMode or SameSiteStrictMode to prevent CSRF attacks.
	SameSite http.SameSite
}

// New returns a new http.Cookie with the given name, value, and properties from config.
//
// The cookie name is determined by the name parameter. If name is empty,
// config.Name is used as a fallback. If both are empty, nil is returned.
//
// The cookie's Expires field is automatically set based on the MaxAge value:
//   - MaxAge > 0: Expires is set to current time + MaxAge seconds
//   - MaxAge < 0: Expires is set to Unix epoch (deletes the cookie)
//   - MaxAge = 0: Expires is not set (session cookie)
//
// Example:
//
//	config := &Config{
//		Path:     "/",
//		MaxAge:   3600,
//		Secure:   true,
//		HTTPOnly: true,
//		SameSite: http.SameSiteLaxMode,
//	}
//	cookie := New("session", "abc123", config)
//
// Returns nil if no cookie name can be determined.
func New(name, value string, config *Config) *http.Cookie {
	cookieName := name
	if cookieName == "" {
		cookieName = config.Name
	}

	if cookieName == "" {
		return nil
	}

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    value,
		Path:     config.Path,
		Domain:   config.Domain,
		MaxAge:   config.MaxAge,
		HttpOnly: config.HTTPOnly,
		Secure:   config.Secure,
		SameSite: config.SameSite,
	}

	if expires, ok := expiresTime(config.MaxAge); ok {
		cookie.Expires = expires
	}

	return cookie
}

// expiresTime converts a maxAge time in seconds to a time.Time in the future.
//
// This function implements the cookie expiration logic as defined in RFC 6265.
// It returns the calculated expiration time and a boolean indicating whether
// the expiration time should be set.
//
// Behavior:
//   - maxAge > 0: Returns current time + maxAge seconds, true
//   - maxAge < 0: Returns Unix epoch time (Jan 1, 1970), true (deletes cookie)
//   - maxAge = 0: Returns zero time, false (session cookie, no expiration)
//
// Reference: http://golang.org/src/net/http/cookie.go?s=618:801#L23
func expiresTime(maxAge int) (time.Time, bool) {
	if maxAge > 0 {
		d := time.Duration(maxAge) * time.Second
		return time.Now().Add(d), true
	}
	if maxAge < 0 {
		return time.Unix(1, 0), true
	}

	return time.Time{}, false
}
