// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package echolog_test

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kopexa-grc/x/echolog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestLoggerIntegration(t *testing.T) {
	// Create Echo instance
	e := echo.New()
	e.HideBanner = true

	// Create a buffer to capture log output
	buf := new(bytes.Buffer)

	// Configure logger options
	setters := []echolog.ConfigSetter{
		echolog.WithLevel(log.DEBUG),
		echolog.WithTimestamp(),
		echolog.WithCaller(),
	}

	// Create logger with the buffer as output
	logger := echolog.New(io.MultiWriter(buf, os.Stdout), setters...)

	// Set the logger on the Echo instance
	e.Logger = logger

	// Set up middleware
	e.Use(echolog.LoggingMiddleware(echolog.Config{
		Logger:          logger,
		RequestIDHeader: "X-Request-ID",
		RequestIDKey:    "request_id",
		HandleError:     true,
	}))

	// Define test route
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	// Start server on a free port
	l, err := net.Listen("tcp", ":0") // nolint:gosec,noctx
	assert.NoError(t, err)

	// Get the actual port assigned by the OS
	addr := l.Addr().String()

	// Configure server to use this listener
	go func() {
		e.Listener = l
		e.Start("") // nolint:errcheck
	}()
	defer e.Close()

	// Allow server time to start
	time.Sleep(100 * time.Millisecond)

	// Create a real HTTP client
	client := &http.Client{}

	// Create a request to the server
	req, err := http.NewRequest(http.MethodGet, "http://"+addr, nil) // nolint:noctx
	assert.NoError(t, err)
	req.Header.Set("X-Request-ID", "test-request-id")

	// Send the request
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)

	// Assert response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// Assert response body
	assert.Equal(t, "Hello, World!", bodyString)

	// Check log output contains expected information
	logOutput := buf.String()
	assert.Contains(t, logOutput, "\"request_id\":\"test-request-id\"") // Request ID should be in the log
	assert.Contains(t, logOutput, "\"method\":\"GET\"")                 // Method should be in the log
	assert.Contains(t, logOutput, "\"uri\":\"/\"")                      // Path should be in the log
	assert.Contains(t, logOutput, "\"status\":200")                     // Status code should be in the log
}
