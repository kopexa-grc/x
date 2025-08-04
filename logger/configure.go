// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package logger provides flexible, environment-aware logging configuration for applications.
// It supports automatic detection of Kubernetes and Docker environments and configures
// the logger accordingly for best practices in each environment.
//
// Features:
//   - Automatic log format selection (console, JSON, GCP JSON)
//   - Environment variable and file-based detection for Docker/Kubernetes
//   - Version tagging in logs
//   - Simple API for integration in main.go
package logger

import (
	"os"

	"github.com/rs/zerolog/log"
)

// Configure sets up the global logger based on the current environment and version.
// It detects Kubernetes and Docker environments and chooses the appropriate log format.
// The log level can be set via environment variables (DEBUG, TRACE).
// The version string is added to all log entries.
// Returns the effective log level and any error encountered.
func Configure(version string) (logLevel string, err error) {
	logLevel = "info"

	// Env-Level has precedence
	if levelFromEnv, ok := GetEnvLogLevel(); ok {
		logLevel = levelFromEnv
	}

	// Logging-Format based on environment
	switch {
	case isKubernetes():
		UseGCPJSONLogging(os.Stdout) // structured JSON-Logger for K8s
	case isDocker():
		UseJSONLogging(os.Stdout) // plain JSON for Docker
	default:
		CliLogger() // nicer console logger locally
	}

	Set(logLevel)

	log.Logger = log.Logger.With().Str("version", version).Logger()

	log.Debug().Str("level", logLevel).Msg("logger configured")

	return logLevel, nil
}

// isKubernetes returns true if the application is running in a Kubernetes environment.
// It checks for the presence of a typical service account file or the KUBERNETES_SERVICE_HOST env var.
func isKubernetes() bool {
	// K8s typical Serviceaccount file
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		return true
	}

	// Alternative check via ENV-Var
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}

	return false
}

// isDocker returns true if the application is running in a Docker container.
// It checks for the RUNNING_IN_DOCKER env var or the presence of /.dockerenv.
func isDocker() bool {
	// Special case: ENV-Var set in Dockerfile
	if os.Getenv("RUNNING_IN_DOCKER") == "true" {
		return true
	}

	// docker created a .dockerenv file at the root of the container
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}
