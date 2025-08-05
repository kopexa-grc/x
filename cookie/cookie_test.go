package cookie

import (
	"net/http"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		cookName string
		value    string
		config   *Config
		want     *http.Cookie
		wantNil  bool
	}{
		{
			name:     "basic cookie creation",
			cookName: "test",
			value:    "value123",
			config: &Config{
				Path:     "/",
				MaxAge:   3600,
				Secure:   true,
				HTTPOnly: true,
				SameSite: http.SameSiteLaxMode,
			},
			want: &http.Cookie{
				Name:     "test",
				Value:    "value123",
				Path:     "/",
				MaxAge:   3600,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			},
		},
		{
			name:     "cookie with config name fallback",
			cookName: "",
			value:    "value456",
			config: &Config{
				Name:     "session",
				Path:     "/api",
				Domain:   "example.com",
				MaxAge:   7200,
				Secure:   false,
				HTTPOnly: false,
				SameSite: http.SameSiteStrictMode,
			},
			want: &http.Cookie{
				Name:     "session",
				Value:    "value456",
				Path:     "/api",
				Domain:   "example.com",
				MaxAge:   7200,
				HttpOnly: false,
				Secure:   false,
				SameSite: http.SameSiteStrictMode,
			},
		},
		{
			name:     "cookie with negative MaxAge sets Expires",
			cookName: "delete_me",
			value:    "",
			config: &Config{
				MaxAge: -1,
			},
			want: &http.Cookie{
				Name:    "delete_me",
				Value:   "",
				MaxAge:  -1,
				Expires: time.Unix(1, 0),
			},
		},
		{
			name:     "no name provided and no config name",
			cookName: "",
			value:    "value",
			config:   &Config{},
			wantNil:  true,
		},
		{
			name:     "zero MaxAge does not set Expires",
			cookName: "temp",
			value:    "temp_value",
			config: &Config{
				MaxAge: 0,
			},
			want: &http.Cookie{
				Name:   "temp",
				Value:  "temp_value",
				MaxAge: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.cookName, tt.value, tt.config)

			if tt.wantNil {
				if got != nil {
					t.Errorf("New() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("New() returned nil, want non-nil cookie")
			}

			// Check basic fields
			if got.Name != tt.want.Name {
				t.Errorf("Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Value != tt.want.Value {
				t.Errorf("Value = %v, want %v", got.Value, tt.want.Value)
			}
			if got.Path != tt.want.Path {
				t.Errorf("Path = %v, want %v", got.Path, tt.want.Path)
			}
			if got.Domain != tt.want.Domain {
				t.Errorf("Domain = %v, want %v", got.Domain, tt.want.Domain)
			}
			if got.MaxAge != tt.want.MaxAge {
				t.Errorf("MaxAge = %v, want %v", got.MaxAge, tt.want.MaxAge)
			}
			if got.HttpOnly != tt.want.HttpOnly {
				t.Errorf("HttpOnly = %v, want %v", got.HttpOnly, tt.want.HttpOnly)
			}
			if got.Secure != tt.want.Secure {
				t.Errorf("Secure = %v, want %v", got.Secure, tt.want.Secure)
			}
			if got.SameSite != tt.want.SameSite {
				t.Errorf("SameSite = %v, want %v", got.SameSite, tt.want.SameSite)
			}

			// Check Expires field for specific cases
			if tt.config.MaxAge > 0 { // nolint:gocritic
				// For positive MaxAge, Expires should be set to a future time
				if got.Expires.IsZero() {
					t.Errorf("Expires should be set for positive MaxAge, but got zero time")
				}
			} else if tt.config.MaxAge < 0 {
				// For negative MaxAge, should be Unix epoch
				if !got.Expires.Equal(time.Unix(1, 0)) {
					t.Errorf("Expires = %v, want %v", got.Expires, time.Unix(1, 0))
				}
			} else {
				// For MaxAge = 0, Expires should be zero (not set)
				if !got.Expires.IsZero() {
					t.Errorf("Expires should be zero for MaxAge=0, got %v", got.Expires)
				}
			}
		})
	}
}

func TestNew_WithPositiveMaxAge(t *testing.T) {
	config := &Config{
		MaxAge: 3600, // 1 hour
	}

	before := time.Now()
	cookie := New("test", "value", config)
	after := time.Now()

	if cookie == nil {
		t.Fatal("New() returned nil")
	}

	if cookie.Expires.IsZero() {
		t.Error("Expires should be set for positive MaxAge")
	}

	// Check that expires is approximately 1 hour from now
	expected := before.Add(time.Hour)
	if cookie.Expires.Before(expected.Add(-time.Minute)) || cookie.Expires.After(after.Add(time.Hour+time.Minute)) {
		t.Errorf("Expires = %v, expected around %v", cookie.Expires, expected)
	}
}

func TestExpiresTime(t *testing.T) {
	tests := []struct {
		name        string
		maxAge      int
		wantExpires bool
		checkTime   bool
	}{
		{
			name:        "positive maxAge",
			maxAge:      3600,
			wantExpires: true,
			checkTime:   true,
		},
		{
			name:        "negative maxAge",
			maxAge:      -1,
			wantExpires: true,
			checkTime:   false,
		},
		{
			name:        "zero maxAge",
			maxAge:      0,
			wantExpires: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expires, ok := expiresTime(tt.maxAge)

			if ok != tt.wantExpires {
				t.Errorf("expiresTime() ok = %v, want %v", ok, tt.wantExpires)
			}

			if tt.checkTime && ok {
				// For positive maxAge, check that the time is in the future
				now := time.Now()
				expected := now.Add(time.Duration(tt.maxAge) * time.Second)

				// Allow for small timing differences
				if expires.Before(expected.Add(-time.Second)) || expires.After(expected.Add(time.Second)) {
					t.Errorf("expiresTime() = %v, expected around %v", expires, expected)
				}
			}

			if tt.maxAge < 0 && ok {
				// For negative maxAge, should return Unix epoch
				if !expires.Equal(time.Unix(1, 0)) {
					t.Errorf("expiresTime() = %v, want %v", expires, time.Unix(1, 0))
				}
			}
		})
	}
}

func TestConfig_SecurityDefaults(t *testing.T) {
	// Test that security-conscious defaults are being used appropriately
	config := &Config{
		Name:     "secure_session",
		Secure:   true,
		HTTPOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600,
	}

	cookie := New("", "secret_value", config)

	if cookie == nil {
		t.Fatal("New() returned nil")
	}

	if !cookie.Secure {
		t.Error("Expected Secure to be true for security test")
	}

	if !cookie.HttpOnly {
		t.Error("Expected HttpOnly to be true for security test")
	}

	if cookie.SameSite != http.SameSiteStrictMode {
		t.Error("Expected SameSite to be SameSiteStrictMode for security test")
	}
}

// Benchmark tests
func BenchmarkNew(b *testing.B) {
	config := &Config{
		Name:     "benchmark",
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		HTTPOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New("test", "value", config)
	}
}

func BenchmarkExpiresTime(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expiresTime(3600)
	}
}
