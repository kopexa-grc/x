# multierr

A powerful Go package for handling multiple errors efficiently and ergonomically, developed as part of [Kopexa](https://kopexa.com)'s enterprise-grade Go ecosystem.

## Features

- **Efficient Error Collection** - Collect multiple errors from batch operations, parallel processing, or complex workflows
- **Automatic Nil Filtering** - Automatically filters out nil errors during collection
- **Smart Deduplication** - Remove duplicate errors based on their string representation
- **Flexible Filtering** - Filter errors using custom predicates for advanced error management
- **Comprehensive Formatting** - Human-readable error formatting with clear structure
- **Go 1.13+ Compatibility** - Full support for error unwrapping and modern error handling patterns
- **Zero Dependencies** - Lightweight implementation with no external dependencies
- **Thread-Safe Design** - Safe for use in concurrent environments with proper synchronization

## Installation

```bash
go get github.com/kopexa-grc/x/multierr
```

## Quick Start

### Basic Error Collection

```go
package main

import (
    "errors"
    "fmt"
    "log"
    
    "github.com/kopexa-grc/x/multierr"
)

func main() {
    var errs multierr.Errors
    
    // Collect errors from multiple operations
    errs.Add(validateUser("john@example.com"))
    errs.Add(validateAccount("account123"))
    errs.Add(validatePermissions("admin"))
    
    // Check if any errors occurred
    if !errs.IsEmpty() {
        log.Printf("Validation failed: %v", errs.Deduplicate())
    } else {
        fmt.Println("All validations passed!")
    }
}

func validateUser(email string) error {
    if email == "" {
        return errors.New("email is required")
    }
    return nil
}

func validateAccount(id string) error {
    if id == "" {
        return errors.New("account ID is required")
    }
    return nil
}

func validatePermissions(role string) error {
    if role != "admin" && role != "user" {
        return errors.New("invalid role")
    }
    return nil
}
```

### Error Wrapping

```go
package main

import (
    "errors"
    "fmt"
    "os"
    
    "github.com/kopexa-grc/x/multierr"
)

func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return multierr.Wrap(err, "failed to open configuration file")
    }
    defer file.Close()
    
    // Process file...
    if err := processData(file); err != nil {
        return multierr.Wrap(err, "failed to process file data")
    }
    
    return nil
}

func processData(file *os.File) error {
    // Simulate processing error
    return errors.New("invalid file format")
}

func main() {
    if err := processFile("config.json"); err != nil {
        fmt.Printf("Error: %v\n", err)
        // Output: Error: failed to process file data: invalid file format
        
        // Works with Go 1.13+ error handling
        if errors.Is(err, os.ErrNotExist) {
            fmt.Println("Config file not found")
        }
    }
}
```

## Advanced Usage

### Error Filtering

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "strings"
    
    "github.com/kopexa-grc/x/multierr"
)

func main() {
    var errs multierr.Errors
    
    // Collect various types of errors
    errs.Add(errors.New("validation error: invalid email"))
    errs.Add(context.DeadlineExceeded)
    errs.Add(errors.New("network error: connection failed"))
    errs.Add(errors.New("validation error: invalid age"))
    
    // Filter out timeout errors
    nonTimeoutErrors := errs.Filter(func(e error) bool {
        return errors.Is(e, context.DeadlineExceeded)
    })
    
    // Filter out validation errors
    nonValidationErrors := errs.Filter(func(e error) bool {
        return strings.Contains(e.Error(), "validation")
    })
    
    fmt.Printf("Non-timeout errors: %d\n", len(nonTimeoutErrors.Errors))
    fmt.Printf("Non-validation errors: %d\n", len(nonValidationErrors.Errors))
}
```

### Batch Processing Example

```go
package main

import (
    "errors"
    "fmt"
    "log"
    
    "github.com/kopexa-grc/x/multierr"
)

type User struct {
    ID    string
    Email string
    Age   int
}

func processBatch(users []User) error {
    var errs multierr.Errors
    
    for _, user := range users {
        if err := processUser(user); err != nil {
            // Wrap error with context
            wrappedErr := multierr.Wrap(err, fmt.Sprintf("failed to process user %s", user.ID))
            errs.Add(wrappedErr)
        }
    }
    
    // Return deduplicated errors or nil
    return errs.Deduplicate()
}

func processUser(user User) error {
    var errs multierr.Errors
    
    // Validate user data
    errs.Add(validateEmail(user.Email))
    errs.Add(validateAge(user.Age))
    errs.Add(validateID(user.ID))
    
    if !errs.IsEmpty() {
        return errs.Deduplicate()
    }
    
    // Process user...
    return nil
}

func validateEmail(email string) error {
    if email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(email, "@") {
        return errors.New("invalid email format")
    }
    return nil
}

func validateAge(age int) error {
    if age < 0 {
        return errors.New("age cannot be negative")
    }
    if age > 150 {
        return errors.New("age is unrealistic")
    }
    return nil
}

func validateID(id string) error {
    if id == "" {
        return errors.New("ID is required")
    }
    return nil
}

func main() {
    users := []User{
        {ID: "1", Email: "john@example.com", Age: 25},
        {ID: "2", Email: "invalid-email", Age: 30},
        {ID: "", Email: "jane@example.com", Age: -5},
    }
    
    if err := processBatch(users); err != nil {
        log.Printf("Batch processing failed:\n%s", err.Error())
    } else {
        fmt.Println("All users processed successfully!")
    }
}
```

### Parallel Processing with Error Collection

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"
    
    "github.com/kopexa-grc/x/multierr"
)

func processItemsConcurrently(items []string) error {
    var (
        errs multierr.Errors
        mu   sync.Mutex
        wg   sync.WaitGroup
    )
    
    for _, item := range items {
        wg.Add(1)
        go func(item string) {
            defer wg.Done()
            
            if err := processItem(item); err != nil {
                mu.Lock()
                errs.Add(multierr.Wrap(err, fmt.Sprintf("item %s", item)))
                mu.Unlock()
            }
        }(item)
    }
    
    wg.Wait()
    return errs.Deduplicate()
}

func processItem(item string) error {
    // Simulate processing time
    time.Sleep(100 * time.Millisecond)
    
    // Simulate random errors
    if len(item) == 0 {
        return errors.New("empty item")
    }
    if item == "error" {
        return errors.New("processing failed")
    }
    
    return nil
}

func main() {
    items := []string{"item1", "item2", "", "error", "item5"}
    
    if err := processItemsConcurrently(items); err != nil {
        fmt.Printf("Some items failed to process:\n%s", err.Error())
    } else {
        fmt.Println("All items processed successfully!")
    }
}
```

## API Reference

### Types

#### `Errors`

The main type for collecting and managing multiple errors.

```go
type Errors struct {
    Errors []error
}
```

#### `withMessage`

Internal type for wrapping errors with additional context (created by `Wrap` function).

### Functions

#### `Wrap(err error, message string) error`

Wraps an error with additional context message, compatible with Go 1.13+ error handling.

#### `(*Errors) Add(err ...error)`

Adds one or more errors to the collection, automatically filtering nil errors.

#### `(*Errors) Filter(f func(error) bool) *Errors`

Returns a new collection containing errors that do NOT match the predicate.

#### `(*Errors) Error() string`

Formats all errors into a human-readable string with bullet points.

#### `(Errors) Deduplicate() error`

Removes duplicate errors and returns nil, single error, or Errors collection.

#### `(*Errors) IsEmpty() bool`

Returns true if the collection contains no errors (safe with nil receivers).

## Error Output Format

The package produces well-formatted error messages:

```
3 errors occurred:
    * failed to process user 2: invalid email format
    * failed to process user 3: ID is required
    * failed to process user 3: age cannot be negative
```

## Performance

- **Efficient Collection**: Minimal allocations during error accumulation
- **O(n) Deduplication**: Uses map-based deduplication for optimal performance
- **Optimized Formatting**: Uses `strings.Builder` for efficient string concatenation
- **Memory Conscious**: Only stores unique error instances during deduplication

## Best Practices

1. **Use for Batch Operations**: Ideal for validating multiple items or processing batches
2. **Combine with Filtering**: Use filtering to separate different types of errors
3. **Wrap with Context**: Add meaningful context when wrapping errors
4. **Check IsEmpty()**: Always check if errors exist before processing
5. **Deduplicate Before Returning**: Call `Deduplicate()` to clean up duplicate errors
6. **Thread Safety**: Use proper synchronization when collecting errors concurrently

## Compatibility

- **Go Version**: Requires Go 1.13+ for full error unwrapping support
- **Error Handling**: Compatible with `errors.Is()`, `errors.As()`, and `errors.Unwrap()`
- **Legacy Support**: Provides `Cause()` method for compatibility with older error handling patterns

## Contributing

This package is part of Kopexa's enterprise Go ecosystem. For contribution guidelines and support, visit [kopexa.com](https://kopexa.com).

## License

Licensed under the Business Source License 1.1 (BUSL-1.1).

Copyright (c) Kopexa GmbH

---

Built with ❤️ by the Kopexa engineering team. Visit [kopexa.com](https://kopexa.com) to learn more about our enterprise solutions.