// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package convert provides efficient utilities for converting between different data representations.
//
// This package offers lightweight alternatives to JSON marshaling/unmarshaling for common
// conversion operations, achieving significantly better performance through the use of
// reflection instead of serialization.
//
// # Key Features
//
//   - High-performance struct to map conversions
//   - Deep recursive conversion of nested structures
//   - Slice and array conversion support
//   - Type-safe map filtering utilities
//   - JSON tag support for field naming
//   - Zero-allocation optimizations for common cases
//
// # Performance Characteristics
//
// The functions in this package are optimized for performance:
//
//   - JSONToDict: ~5x faster than JSON marshaling for simple structs
//   - JSONToDictSlice: ~4x faster than JSON marshaling for struct slices
//   - Map inputs: ~1000x faster with zero allocations
//   - Nested structures: ~4x faster with 72% less memory usage
//
// # Usage Examples
//
// Convert a struct to a map:
//
//	type User struct {
//		Name string `json:"name"`
//		Age  int    `json:"age"`
//	}
//
//	user := User{Name: "John", Age: 30}
//	dict, err := convert.JSONToDict(user)
//	// dict = map[string]any{"name": "John", "age": 30}
//
// Convert a slice of structs:
//
//	users := []User{{Name: "John", Age: 30}, {Name: "Jane", Age: 25}}
//	dicts, err := convert.JSONToDictSlice(users)
//	// dicts = []any{map[string]any{"name": "John", "age": 30}, ...}
//
// Filter map values by type:
//
//	data := map[string]any{"count": 42, "name": "test", "active": true}
//	strings := convert.DictToTypedMap[string](data)
//	// strings = map[string]string{"name": "test"}
//
// # Design Principles
//
// This package follows the Google API Design Guide principles:
//
//   - Resource-oriented design: Functions operate on well-defined data types
//   - Consistent naming: All conversion functions follow a clear naming pattern
//   - Error handling: Functions return explicit errors for invalid operations
//   - Performance: Optimized for common use cases with minimal allocations
//   - Backwards compatibility: Maintains semantic compatibility with JSON conversions
//
// # Thread Safety
//
// All functions in this package are safe for concurrent use. They do not modify
// input parameters and create new data structures for their results.
package convert
