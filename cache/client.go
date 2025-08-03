// Package cache provides a Redis client wrapper with configuration management
// and health checking capabilities for caching operations.
package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config represents the configuration settings for a Redis client connection.
// It includes all necessary parameters for establishing and managing a Redis
// connection with proper timeouts, authentication, and connection pooling.
type Config struct {
	// Enabled specifies whether the Redis client should be active.
	// Default: true
	Enabled bool `json:"enabled" koanf:"enabled" default:"true"`

	// Address is the Redis server host and port in the format "host:port".
	// Default: "localhost:6379"
	Address string `json:"address" koanf:"address" default:"localhost:6379"`

	// Name is an optional client name identifier for the Redis connection.
	// This can be useful for debugging and monitoring purposes.
	Name string `json:"name" koanf:"name" default:""`

	// Username for Redis authentication. Leave empty if authentication is not required.
	Username string `json:"username" koanf:"username"`

	// Password for Redis authentication. Must match the password configured on the Redis server.
	Password string `json:"password" koanf:"password"`

	// DB is the Redis database number to select after connecting.
	// Default: 0 (uses the default database)
	DB int `json:"db" koanf:"db" default:"0"`

	// DialTimeout is the maximum time to wait when establishing new connections.
	// Default: 5s
	DialTimeout time.Duration `json:"dialTimeout" koanf:"dialTimeout" default:"5s"`

	// ReadTimeout is the timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Supported values:
	//   - 0: default timeout (3 seconds)
	//   - -1: no timeout (block indefinitely)
	//   - -2: disables SetReadDeadline calls completely
	// Default: 0
	ReadTimeout time.Duration `json:"readTimeout" koanf:"readTimeout" default:"0"`

	// WriteTimeout is the timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking. Supported values:
	//   - 0: default timeout (3 seconds)
	//   - -1: no timeout (block indefinitely)
	//   - -2: disables SetWriteDeadline calls completely
	// Default: 0
	WriteTimeout time.Duration `json:"writeTimeout" koanf:"writeTimeout" default:"0"`

	// MaxRetries is the maximum number of retries before giving up on a command.
	// Default: 3 retries; -1 disables retries
	MaxRetries int `json:"maxRetries" koanf:"maxRetries" default:"3"`

	// MinIdleConns is the minimum number of idle connections to maintain.
	// Useful when establishing new connections is slow.
	// Default: 0 (idle connections are not closed by default)
	MinIdleConns int `json:"minIdleConns" koanf:"minIdleConns" default:"0"`

	// MaxIdleConns is the maximum number of idle connections to maintain.
	// Default: 0 (idle connections are not closed by default)
	MaxIdleConns int `json:"maxIdleConns" koanf:"maxIdleConns" default:"0"`

	// MaxActiveConns is the maximum number of connections allocated by the pool.
	// When zero, there is no limit on the number of connections in the pool.
	// Default: 0 (unlimited)
	MaxActiveConns int `json:"maxActiveConns" koanf:"maxActiveConns" default:"0"`
}

// New creates and returns a new Redis client using the provided configuration.
// The client is configured with connection pooling, timeouts, and authentication
// settings as specified in the Config struct.
//
// Example:
//
//	config := cache.Config{
//		Address: "localhost:6379",
//		DB: 0,
//		DialTimeout: 5 * time.Second,
//	}
//	client := cache.New(config)
func New(c Config) *redis.Client {
	opts := &redis.Options{
		Addr:            c.Address,
		DB:              c.DB,
		DialTimeout:     c.DialTimeout,
		ReadTimeout:     c.ReadTimeout,
		WriteTimeout:    c.WriteTimeout,
		MaxRetries:      c.MaxRetries,
		MinIdleConns:    c.MinIdleConns,
		MaxIdleConns:    c.MaxIdleConns,
		MaxActiveConns:  c.MaxActiveConns,
		DisableIdentity: true,
	}

	// optional fields
	if c.Name != "" {
		opts.ClientName = c.Name
	}

	if c.Username != "" {
		opts.Username = c.Username
	}

	if c.Password != "" {
		opts.Password = c.Password
	}

	return redis.NewClient(opts)
}

// Healthcheck returns a health check function that pings the Redis client
// to verify the connection is working properly. This function can be used
// in health check endpoints or monitoring systems.
//
// The returned function accepts a context and returns an error if the
// Redis connection is not healthy.
//
// Example:
//
//	client := cache.New(config)
//	healthFn := cache.Healthcheck(client)
//	err := healthFn(ctx)
//	if err != nil {
//		// Handle unhealthy Redis connection
//	}
func Healthcheck(c *redis.Client) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		// check if its alive
		if err := c.Ping(ctx).Err(); err != nil {
			return err
		}

		return nil
	}
}
