// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package echolog

import (
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// Config defines configuration options for the echolog logging middleware.
//
// This struct contains all configuration options for customizing the behavior
// of the logging middleware, including request logging, error handling,
// and performance monitoring features.
type Config struct {
	// Logger is the echolog instance to use for logging.
	// If nil, a default logger writing to os.Stdout with timestamps will be created.
	Logger *Logger

	// Skipper defines a function to skip the entire middleware for certain requests.
	// If nil, middleware.DefaultSkipper is used (never skips).
	//
	// Example:
	//   Skipper: func(c echo.Context) bool {
	//       return c.Path() == "/health" // Skip logging for health checks
	//   }
	Skipper middleware.Skipper

	// AfterNextSkipper defines a function to skip logging after the request handler executes.
	// This allows the middleware to set up context but skip the final logging step.
	// If nil, middleware.DefaultSkipper is used (never skips).
	AfterNextSkipper middleware.Skipper

	// BeforeNext is executed before calling the next handler in the chain.
	// This can be used to perform setup operations with access to the enriched logger.
	BeforeNext middleware.BeforeFunc

	// Enricher allows adding custom fields to the logger for each request.
	// The function receives the Echo context and logger context, returning an enriched logger.
	//
	// Example:
	//   Enricher: func(c echo.Context, logger zerolog.Context) zerolog.Context {
	//       if userID := getUserID(c); userID != "" {
	//           return logger.Str("user_id", userID)
	//       }
	//       return logger
	//   }
	Enricher Enricher

	// RequestIDHeader specifies the HTTP header name to read/write request IDs.
	// Defaults to "X-Request-ID" if empty.
	RequestIDHeader string

	// RequestIDKey specifies the JSON field name for request IDs in log output.
	// Defaults to "request_id" if empty.
	RequestIDKey string

	// NestKey specifies a JSON field name to nest all request details under.
	// If empty, request details are logged at the root level.
	//
	// Example with NestKey="request":
	//   {"level":"info","request":{"method":"GET","uri":"/api/users"},"message":"request details"}
	NestKey string

	// HandleError determines whether errors should be propagated to Echo's error handler.
	// If true, errors from handlers will be passed to echo.Context.Error().
	// If false, errors are returned without additional processing.
	HandleError bool

	// RequestLatencyLimit sets a duration threshold for slow request detection.
	// Requests exceeding this duration will be logged at RequestLatencyLevel.
	// If zero, latency-based level escalation is disabled.
	RequestLatencyLimit time.Duration

	// RequestLatencyLevel specifies the log level for slow requests.
	// Only used when RequestLatencyLimit is non-zero and exceeded.
	// Common values: zerolog.WarnLevel, zerolog.ErrorLevel
	RequestLatencyLevel zerolog.Level
}

// Enricher is a function type for adding custom fields to request loggers.
//
// Enricher functions are called for each request and can add context-specific
// information to the logger. They receive the Echo context and a zerolog context,
// and should return an enriched zerolog context.
//
// The enricher is called before the request handler executes, allowing you to
// add information that will be included in all log messages for that request.
//
// Example implementations:
//
//	// Add user information to logs
//	func UserEnricher(c echo.Context, logger zerolog.Context) zerolog.Context {
//		if user := getCurrentUser(c); user != nil {
//			return logger.Str("user_id", user.ID).Str("user_role", user.Role)
//		}
//		return logger
//	}
//
//	// Add request metadata
//	func MetadataEnricher(c echo.Context, logger zerolog.Context) zerolog.Context {
//		return logger.
//			Str("client_ip", c.RealIP()).
//			Str("user_agent", c.Request().UserAgent()).
//			Int64("content_length", c.Request().ContentLength)
//	}
type Enricher func(c echo.Context, logger zerolog.Context) zerolog.Context

// Context wraps echo.Context with an enhanced logger instance.
//
// This wrapper provides access to a logger that has been enriched with
// request-specific information like request IDs and custom fields added
// through the Enricher function.
type Context struct {
	echo.Context         // embedded Echo context
	logger       *Logger // request-specific logger instance
}

// NewContext creates a new Context wrapper with the provided Echo context and logger.
//
// This function wraps an Echo context with an enhanced logger instance,
// making the logger available through the Context.Logger() method.
//
// Parameters:
//   - ctx: the Echo context to wrap
//   - logger: the logger instance to associate with this context
//
// Returns:
//   - *Context: wrapped context with enhanced logging capabilities
//
// Example:
//
//	// Typically used internally by the middleware
//	enhancedCtx := echolog.NewContext(c, enrichedLogger)
func NewContext(ctx echo.Context, logger *Logger) *Context {
	return &Context{ctx, logger}
}

// Logger returns the logger instance associated with this context.
//
// This method provides access to the request-specific logger that has been
// enriched with request IDs and custom fields. It implements the Echo
// Logger interface, making it compatible with existing Echo logging patterns.
//
// Returns:
//   - echo.Logger: the enriched logger instance
//
// Example:
//
//	func myHandler(c echo.Context) error {
//		c.Logger().Info("Processing request")
//		return c.JSON(200, map[string]string{"status": "ok"})
//	}
func (c *Context) Logger() echo.Logger {
	return c.logger
}

// LoggingMiddleware returns an Echo middleware function that provides comprehensive request logging.
//
// This middleware captures detailed information about each HTTP request including timing,
// status codes, request/response sizes, and other metadata. It enriches the logger with
// request-specific information and makes it available to downstream handlers.
//
// Key features:
//   - Automatic request ID tracking
//   - Latency monitoring with configurable thresholds
//   - Request/response size logging
//   - Error handling and propagation
//   - Custom field enrichment
//   - Context-aware logging
//
// Parameters:
//   - config: Configuration struct defining middleware behavior
//
// Returns:
//   - echo.MiddlewareFunc: middleware function for use with Echo
//
// Example:
//
//	e := echo.New()
//	logger := echolog.New(os.Stdout, echolog.WithTimestamp())
//
//	// Basic configuration
//	e.Use(echolog.LoggingMiddleware(echolog.Config{
//		Logger: logger,
//	}))
//
//	// Advanced configuration with slow request detection
//	e.Use(echolog.LoggingMiddleware(echolog.Config{
//		Logger:              logger,
//		RequestLatencyLimit: 500 * time.Millisecond,
//		RequestLatencyLevel: zerolog.WarnLevel,
//		HandleError:         true,
//		Enricher: func(c echo.Context, logger zerolog.Context) zerolog.Context {
//			return logger.Str("endpoint", c.Path())
//		},
//	}))
func LoggingMiddleware(config Config) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	if config.AfterNextSkipper == nil {
		config.AfterNextSkipper = middleware.DefaultSkipper
	}

	if config.Logger == nil {
		config.Logger = New(os.Stdout, WithTimestamp())
	}

	if config.RequestIDKey == "" {
		config.RequestIDKey = "request_id"
	}

	if config.RequestIDHeader == "" {
		config.RequestIDHeader = "X-Request-ID"
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			var err error

			req := c.Request()
			start := time.Now()

			logger := config.Logger

			logger = enrichLogger(c, logger, config)

			id := getRequestID(c, config)
			if id != "" {
				logger = From(logger.log, WithField(config.RequestIDKey, id))
			}

			// The request context is retrieved and set to the logger's context
			// the context is then set to the request, and a new context is created with the logger
			c.SetRequest(req.WithContext(logger.WithContext(req.Context())))
			c = NewContext(c, logger)

			if config.BeforeNext != nil {
				config.BeforeNext(c)
			}

			if err = next(c); err != nil {
				if config.HandleError {
					c.Error(err)
				}
			}

			if config.AfterNextSkipper(c) {
				return nil
			}

			logEvent(c, logger, config, start, err)

			return err
		}
	}
}

// getRequestID extracts the request ID from HTTP headers.
//
// This function looks for a request ID in both request and response headers,
// checking the request headers first and falling back to response headers.
// This supports scenarios where request IDs are set by upstream proxies
// or generated by other middleware.
//
// Parameters:
//   - c: Echo context containing the request and response
//   - config: middleware configuration containing header name
//
// Returns:
//   - string: request ID if found, empty string otherwise
func getRequestID(c echo.Context, config Config) string {
	req := c.Request()
	res := c.Response()

	id := req.Header.Get(config.RequestIDHeader)
	if id == "" {
		id = res.Header().Get(config.RequestIDHeader)
	}

	return id
}

// enrichLogger applies custom enrichment to the logger using the configured Enricher function.
//
// This function creates a new logger instance with additional fields added by the
// Enricher function. The enricher can add context-specific information like user IDs,
// session information, or other request metadata.
//
// Parameters:
//   - c: Echo context for the current request
//   - logger: base logger instance to enrich
//   - config: middleware configuration containing the Enricher function
//
// Returns:
//   - *Logger: enriched logger instance, or original logger if no enricher configured
func enrichLogger(c echo.Context, logger *Logger, config Config) *Logger {
	if config.Enricher != nil {
		logger = From(logger.log)
		logger.log = config.Enricher(c, logger.log.With()).Logger()
	}

	return logger
}

// logEvent generates the final log entry for a completed HTTP request.
//
// This function creates a comprehensive log entry containing request metadata,
// timing information, response details, and error information. It intelligently
// selects the appropriate log level based on error conditions and latency thresholds.
//
// The function logs the following information:
//   - Request method, URI, and protocol
//   - Client IP and user agent
//   - Response status code and size
//   - Request processing latency
//   - Error details (if any)
//   - Content length for request and response
//
// Log level selection:
//   - Error level: when an error occurred during processing
//   - Configurable level: when request latency exceeds the configured threshold
//   - Default level: for normal request processing
//
// Parameters:
//   - c: Echo context containing request and response information
//   - logger: logger instance to use for output
//   - config: middleware configuration including latency settings
//   - start: timestamp when request processing began
//   - err: error that occurred during processing (nil if no error)
//
// Example output (with NestKey=""):
//
//	{
//	  "level": "info",
//	  "request_id": "abc123",
//	  "remote_ip": "192.168.1.1",
//	  "method": "GET",
//	  "uri": "/api/users",
//	  "status": 200,
//	  "latency_human": "15.2ms",
//	  "bytes_in": "0",
//	  "bytes_out": "1024",
//	  "message": "request details"
//	}
//
// Example output (with NestKey="request"):
//
//	{
//	  "level": "info",
//	  "request": {
//	    "remote_ip": "192.168.1.1",
//	    "method": "GET",
//	    "uri": "/api/users",
//	    "status": 200,
//	    "latency_human": "15.2ms"
//	  },
//	  "message": "request details"
//	}
func logEvent(c echo.Context, logger *Logger, config Config, start time.Time, err error) {
	req := c.Request()
	res := c.Response()
	stop := time.Now()
	latency := stop.Sub(start)

	var mainEvt *zerolog.Event
	// this is the error that's passed in as input from the middleware func

	switch {
	case err != nil:
		mainEvt = logger.log.WithLevel(zerolog.ErrorLevel).Str("error", err.Error())
	case config.RequestLatencyLimit != 0 && latency > config.RequestLatencyLimit:
		mainEvt = logger.log.WithLevel(config.RequestLatencyLevel)
	default:
		mainEvt = logger.log.WithLevel(logger.log.GetLevel())
	}

	var evt *zerolog.Event

	if config.NestKey != "" {
		evt = zerolog.Dict()
	} else {
		evt = mainEvt
	}

	evt.Str("remote_ip", c.RealIP())
	evt.Str("host", req.Host)
	evt.Str("method", req.Method)
	evt.Str("uri", req.RequestURI)
	evt.Str("user_agent", req.UserAgent())
	evt.Int("status", res.Status)
	evt.Str("referer", req.Referer())
	evt.Str("latency_human", latency.String())
	evt.Str("client_ip", c.RealIP())
	evt.Str("request_protocol", req.Proto)

	cl := req.Header.Get(echo.HeaderContentLength)
	if cl == "" {
		cl = "0"
	}

	evt.Str("bytes_in", cl)
	evt.Str("bytes_out", strconv.FormatInt(res.Size, 10))

	if config.NestKey != "" {
		mainEvt.Dict(config.NestKey, evt).Msg("request details")
	} else {
		mainEvt.Msg("request details")
	}
}
