// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package echolog

import (
	"github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
)

// Mapping tables for converting between Echo and zerolog log levels.
//
// These maps provide bidirectional conversion between Echo's log levels
// and zerolog's log levels, ensuring consistent behavior across both
// logging interfaces.
var (
	// echoLevels maps Echo log levels to corresponding zerolog levels.
	// This mapping is used when converting from Echo's logging interface
	// to zerolog's internal representation.
	echoLevels = map[log.Lvl]zerolog.Level{
		log.DEBUG: zerolog.DebugLevel, // Debug level for detailed diagnostic information
		log.INFO:  zerolog.InfoLevel,  // Info level for general application information
		log.WARN:  zerolog.WarnLevel,  // Warning level for potentially harmful situations
		log.ERROR: zerolog.ErrorLevel, // Error level for error conditions
		log.OFF:   zerolog.NoLevel,    // Off level disables logging
	}

	// zeroLevels maps zerolog log levels to corresponding Echo levels.
	// This mapping is used when converting from zerolog's internal representation
	// back to Echo's logging interface. Note that zerolog.TraceLevel maps to
	// log.DEBUG since Echo doesn't have a trace level.
	zeroLevels = map[zerolog.Level]log.Lvl{
		zerolog.TraceLevel: log.DEBUG, // Trace maps to Debug (Echo has no trace level)
		zerolog.DebugLevel: log.DEBUG, // Debug level for detailed diagnostic information
		zerolog.InfoLevel:  log.INFO,  // Info level for general application information
		zerolog.WarnLevel:  log.WARN,  // Warning level for potentially harmful situations
		zerolog.ErrorLevel: log.ERROR, // Error level for error conditions
		zerolog.NoLevel:    log.OFF,   // No level corresponds to logging disabled
	}
)

// MatchEchoLevel converts an Echo log level to the corresponding zerolog level.
//
// This function takes an Echo log level and returns both the equivalent zerolog
// level and the original Echo level. If the provided Echo level is not recognized,
// it defaults to zerolog.NoLevel and log.OFF (logging disabled).
//
// Parameters:
//   - level: Echo log level to convert (log.DEBUG, log.INFO, log.WARN, log.ERROR, log.OFF)
//
// Returns:
//   - zerolog.Level: corresponding zerolog level
//   - log.Lvl: the original Echo level (for consistency)
//
// Example:
//
//	zeroLevel, echoLevel := echolog.MatchEchoLevel(log.WARN)
//	fmt.Printf("Echo %v maps to zerolog %v\n", echoLevel, zeroLevel)
//	// Output: Echo WARN maps to zerolog warn
func MatchEchoLevel(level log.Lvl) (zerolog.Level, log.Lvl) {
	zlvl, found := echoLevels[level]

	if found {
		return zlvl, level
	}

	return zerolog.NoLevel, log.OFF
}

// MatchZeroLevel converts a zerolog level to the corresponding Echo log level.
//
// This function takes a zerolog level and returns both the equivalent Echo
// level and the original zerolog level. If the provided zerolog level is not
// recognized, it defaults to log.OFF and zerolog.NoLevel (logging disabled).
//
// Note that zerolog.TraceLevel maps to log.DEBUG since Echo doesn't have a
// dedicated trace level.
//
// Parameters:
//   - level: zerolog level to convert
//
// Returns:
//   - log.Lvl: corresponding Echo log level
//   - zerolog.Level: the original zerolog level (for consistency)
//
// Example:
//
//	echoLevel, zeroLevel := echolog.MatchZeroLevel(zerolog.InfoLevel)
//	fmt.Printf("Zerolog %v maps to Echo %v\n", zeroLevel, echoLevel)
//	// Output: Zerolog info maps to Echo INFO
//
//	// Trace level maps to DEBUG
//	echoLevel, zeroLevel = echolog.MatchZeroLevel(zerolog.TraceLevel)
//	fmt.Printf("Zerolog %v maps to Echo %v\n", zeroLevel, echoLevel)
//	// Output: Zerolog trace maps to Echo DEBUG
func MatchZeroLevel(level zerolog.Level) (log.Lvl, zerolog.Level) {
	elvl, found := zeroLevels[level]

	if found {
		return elvl, level
	}

	return log.OFF, zerolog.NoLevel
}
