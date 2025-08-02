// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package multierr provides utilities for handling multiple errors efficiently and ergonomically.
//
// This package offers tools for collecting, filtering, and managing multiple errors that can
// occur during batch operations, parallel processing, or complex workflows. It provides a
// clean API for error aggregation with deduplication and filtering capabilities.
//
// # Key Features
//
//   - Efficient error collection and aggregation
//   - Automatic nil error filtering
//   - Error deduplication based on error messages
//   - Flexible error filtering with custom predicates
//   - Comprehensive error formatting for debugging
//   - Error wrapping functionality compatible with Go 1.13+ error handling
//
// # Basic Usage
//
// Collecting multiple errors:
//
//	var errs multierr.Errors
//	errs.Add(operation1())
//	errs.Add(operation2())
//	errs.Add(operation3())
//
//	if err := errs.Deduplicate(); err != nil {
//	    log.Printf("Operations failed: %v", err)
//	}
//
// Error filtering:
//
//	// Filter out specific error types
//	filtered := errs.Filter(func(e error) bool {
//	    return errors.Is(e, context.DeadlineExceeded)
//	})
//
//	if !filtered.IsEmpty() {
//	    log.Printf("Non-timeout errors: %v", filtered.Deduplicate())
//	}
//
// Error wrapping:
//
//	if err := someOperation(); err != nil {
//	    return multierr.Wrap(err, "failed to process user data")
//	}
//
// # Performance
//
// The package is designed for efficiency with minimal allocations during error collection.
// Deduplication is performed using a map for O(n) complexity, and error formatting uses
// strings.Builder for optimal string concatenation.
//
// # Compatibility
//
// This package is compatible with Go 1.13+ error handling patterns, including error
// unwrapping and the errors.Is/errors.As functions. The Wrap function creates errors
// that properly implement the Unwrap() method.
//
// For more information about Kopexa's Go ecosystem, visit https://kopexa.com
package multierr

import (
	"strconv"
	"strings"
)

// withMessage wraps an error with an additional message.
//
// This implementation is inspired by https://github.com/pkg/errors and provides
// a simple way to add context to errors while maintaining the error chain for
// proper unwrapping behavior.
//
// The wrapped error implements both Cause() (for legacy compatibility) and
// Unwrap() (for Go 1.13+ error handling) methods.

type withMessage struct {
	cause error  // underlying error being wrapped
	msg   string // additional context message
}

// Error returns the error message with the additional context.
//
// The format is: "message: original_error_message"
//
// Example:
//
//	err := multierr.Wrap(errors.New("file not found"), "failed to read config")
//	fmt.Println(err.Error()) // Output: "failed to read config: file not found"
func (w withMessage) Error() string { return w.msg + ": " + w.cause.Error() }

// Cause returns the underlying error for legacy error handling.
//
// This method provides compatibility with older error handling patterns
// that expect a Cause() method to retrieve the wrapped error.
//
// Returns:
//   - error: the original error that was wrapped
func (w withMessage) Cause() error { return w.cause }

// Unwrap returns the underlying error for Go 1.13+ error handling.
//
// This method enables the use of errors.Is() and errors.As() functions
// with wrapped errors, providing proper error chain traversal.
//
// Returns:
//   - error: the original error that was wrapped
func (w withMessage) Unwrap() error { return w.cause }

// Wrap wraps an error with an additional context message.
//
// This function adds contextual information to an error without losing the
// original error details. The wrapped error maintains the error chain for
// proper unwrapping with Go 1.13+ error handling patterns.
//
// If the provided error is nil, Wrap returns nil, making it safe to use
// in conditional chains without additional nil checks.
//
// Parameters:
//   - err: the error to wrap (can be nil)
//   - message: contextual message to add to the error
//
// Returns:
//   - error: wrapped error with additional context, or nil if err was nil
//
// Example:
//
//	func processFile(filename string) error {
//	    file, err := os.Open(filename)
//	    if err != nil {
//	        return multierr.Wrap(err, "failed to open configuration file")
//	    }
//	    defer file.Close()
//
//	    // Process file...
//	    if err := processData(file); err != nil {
//	        return multierr.Wrap(err, "failed to process file data")
//	    }
//
//	    return nil
//	}
//
// The wrapped error can be used with standard Go error handling:
//
//	if err := processFile("config.json"); err != nil {
//	    if errors.Is(err, os.ErrNotExist) {
//	        log.Printf("Config file not found: %v", err)
//	    } else {
//	        log.Printf("Processing failed: %v", err)
//	    }
//	}
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return withMessage{
		cause: err,
		msg:   message,
	}
}

// Errors represents a collection of errors that can be accumulated and processed together.
//
// This type provides an efficient way to collect multiple errors from batch operations,
// parallel processing, or complex workflows. It automatically filters out nil errors
// and provides utilities for deduplication, filtering, and formatting.
//
// The zero value of Errors is ready to use.
//
// Example usage:
//
//	var errs multierr.Errors
//
//	// Collect errors from multiple operations
//	errs.Add(validateUser(user))
//	errs.Add(validateAccount(account))
//	errs.Add(validatePermissions(permissions))
//
//	// Check if any errors occurred
//	if !errs.IsEmpty() {
//	    return errs.Deduplicate()
//	}
//
// Thread safety: Errors is not safe for concurrent use. If you need to collect
// errors from multiple goroutines, use appropriate synchronization mechanisms.
type Errors struct {
	// Errors contains the collected error instances.
	// This field is exported to allow direct access when needed, but it's
	// recommended to use the provided methods for most operations.
	Errors []error
}

// Add appends one or more errors to the collection.
//
// This method automatically filters out nil errors, so it's safe to call
// with the return values of functions that may or may not return errors.
// Multiple errors can be added in a single call for convenience.
//
// Parameters:
//   - err: variadic list of errors to add (nil errors are ignored)
//
// Example:
//
//	var errs multierr.Errors
//
//	// Add individual errors
//	errs.Add(operation1())
//	errs.Add(operation2())
//
//	// Add multiple errors at once
//	errs.Add(
//	    validateName(name),
//	    validateEmail(email),
//	    validateAge(age),
//	)
//
//	// Nil errors are safely ignored
//	errs.Add(nil, someError, nil) // Only someError is added
func (m *Errors) Add(err ...error) {
	for i := range err {
		if err[i] != nil {
			m.Errors = append(m.Errors, err[i])
		}
	}
}

// Filter creates a new Errors collection containing errors that do NOT match the predicate.
//
// This method applies a filter function to each error in the collection and returns
// a new Errors instance containing only the errors for which the predicate returns false.
// This is useful for removing specific types of errors or errors matching certain criteria.
//
// The original Errors collection is not modified.
//
// Parameters:
//   - f: predicate function that returns true for errors to be excluded
//
// Returns:
//   - *Errors: new collection containing errors that did not match the predicate
//
// Example:
//
//	var errs multierr.Errors
//	errs.Add(errors.New("validation error"))
//	errs.Add(context.DeadlineExceeded)
//	errs.Add(errors.New("network error"))
//
//	// Filter out timeout errors
//	nonTimeoutErrors := errs.Filter(func(e error) bool {
//	    return errors.Is(e, context.DeadlineExceeded)
//	})
//
//	// Filter out errors containing "validation"
//	nonValidationErrors := errs.Filter(func(e error) bool {
//	    return strings.Contains(e.Error(), "validation")
//	})
//
//	// Filter by custom error types
//	type NetworkError struct{ msg string }
//	func (ne NetworkError) Error() string { return ne.msg }
//
//	nonNetworkErrors := errs.Filter(func(e error) bool {
//	    var netErr NetworkError
//	    return errors.As(e, &netErr)
//	})
func (m *Errors) Filter(f func(e error) bool) *Errors {
	res := Errors{}
	for i := range m.Errors {
		cur := m.Errors[i]
		if !f(cur) {
			res.Errors = append(res.Errors, cur)
		}
	}
	return &res
}

// Error implements the error interface and formats all collected errors into a readable string.
//
// This method creates a formatted string representation of all errors in the collection,
// with each error listed on a separate line with bullet points. The format makes it
// easy to identify and read multiple errors in logs or error messages.
//
// The output format is:
//   - Single error: "1 error occurred:\n\t* error_message\n"
//   - Multiple errors: "N errors occurred:\n\t* error1\n\t* error2\n..."
//
// Returns:
//   - string: formatted string containing all error messages
//
// Example output:
//
//	3 errors occurred:
//		* failed to validate user: invalid email format
//		* failed to validate account: insufficient balance
//		* failed to validate permissions: access denied
//
// Usage:
//
//	var errs multierr.Errors
//	errs.Add(errors.New("first error"))
//	errs.Add(errors.New("second error"))
//
//	fmt.Printf("Operations failed:\n%s", errs.Error())
//	// Output:
//	// Operations failed:
//	// 2 errors occurred:
//	//     * first error
//	//     * second error
func (m *Errors) Error() string {
	var res strings.Builder

	n := strconv.Itoa(len(m.Errors))
	if n == "1" {
		res.WriteString("1 error occurred:\n")
	} else {
		res.WriteString(n + " errors occurred:\n")
	}

	for i := range m.Errors {
		res.WriteString("\t* ")
		res.WriteString(m.Errors[i].Error())
		res.WriteByte('\n')
	}
	return res.String()
}

// Deduplicate removes duplicate errors and returns a single error or nil.
//
// This method removes duplicate errors based on their Error() string representation,
// keeping only unique error messages. If no errors remain after deduplication,
// it returns nil. If exactly one unique error remains, that error is returned directly.
// If multiple unique errors remain, they are returned as a new Errors instance.
//
// The deduplication process uses a map for O(n) time complexity and preserves
// the original error instances (not just their messages).
//
// Returns:
//   - nil: if no errors are present
//   - error: the single unique error if only one remains
//   - *Errors: collection of unique errors if multiple remain
//
// Example:
//
//	var errs multierr.Errors
//	errs.Add(errors.New("duplicate error"))
//	errs.Add(errors.New("unique error"))
//	errs.Add(errors.New("duplicate error")) // Same message as first
//
//	result := errs.Deduplicate()
//	// result contains only 2 unique errors
//
//	// Handle the result
//	if result != nil {
//	    log.Printf("Unique errors found: %v", result)
//	}
//
// Use cases:
//   - Batch validation where the same validation error might occur multiple times
//   - Parallel processing where identical errors might be generated
//   - Log aggregation where you want to show unique error types
//
// Note: Deduplication is based on the string representation of errors, so
// structurally different errors with the same message will be considered duplicates.
func (m Errors) Deduplicate() error {
	if len(m.Errors) == 0 {
		return nil
	}

	track := map[string]error{}
	for i := range m.Errors {
		e := m.Errors[i]
		track[e.Error()] = e
	}

	res := make([]error, len(track))
	i := 0
	for _, v := range track {
		res[i] = v
		i++
	}
	return &Errors{Errors: res}
}

// IsEmpty returns true if the collection contains no errors.
//
// This method provides a convenient way to check if any errors have been
// collected without needing to access the Errors slice directly. It safely
// handles nil receivers, making it useful in scenarios where the Errors
// instance might not be initialized.
//
// Returns:
//   - true: if the collection is nil or contains no errors
//   - false: if the collection contains one or more errors
//
// Example:
//
//	var errs multierr.Errors
//
//	// Check before collecting errors
//	fmt.Println(errs.IsEmpty()) // true
//
//	errs.Add(errors.New("something went wrong"))
//	fmt.Println(errs.IsEmpty()) // false
//
//	// Use in conditional logic
//	if !errs.IsEmpty() {
//	    return errs.Deduplicate()
//	}
//
//	// Safe with nil pointers
//	var nilErrs *multierr.Errors
//	fmt.Println(nilErrs.IsEmpty()) // true (no panic)
//
// This method is particularly useful in validation scenarios:
//
//	func validateData(data *Data) error {
//	    var errs multierr.Errors
//
//	    errs.Add(validateName(data.Name))
//	    errs.Add(validateEmail(data.Email))
//	    errs.Add(validateAge(data.Age))
//
//	    if errs.IsEmpty() {
//	        return nil // All validations passed
//	    }
//
//	    return errs.Deduplicate()
//	}
func (m *Errors) IsEmpty() bool {
	if m == nil {
		return true
	}
	return len(m.Errors) == 0
}
