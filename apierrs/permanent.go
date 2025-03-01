package apierrs

import (
	"errors"
)

// PermanentError signals that the operation should not be retried.
// Use this to wrap errors that are known to be permanent failures
// where retrying the operation would not help.
//
// Examples of permanent errors:
// - Authentication failures
// - Invalid configuration
// - Resource not found (when you're sure it doesn't exist)
// - Protocol errors
// - Data validation errors
//
// Usage:
//
//	if invalidConfig {
//	    return apierrs.Permanent(errors.New("invalid configuration"))
//	}
//
//	// In error handling:
//	if err != nil {
//	    if apierrs.IsPermanent(err) {
//	        // Don't retry
//	        return err
//	    }
//	    // Consider retrying for non-permanent errors
//	}
type PermanentError struct {
	Err error
}

// Error implements the error interface.
func (e *PermanentError) Error() string {
	return e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *PermanentError) Unwrap() error {
	return e.Err
}

// Is reports whether the target is a PermanentError.
func (e *PermanentError) Is(target error) bool {
	_, ok := target.(*PermanentError)
	return ok
}

// Permanent wraps the given err in a *PermanentError.
// If err is nil, nil is returned.
func Permanent(err error) error {
	if err == nil {
		return nil
	}

	// If it's already a permanent error, return it as is
	var perr *PermanentError
	if errors.As(err, &perr) {
		return err
	}

	return &PermanentError{
		Err: err,
	}
}

// IsPermanent checks if the given error is or contains a PermanentError.
// Returns false if err is nil.
func IsPermanent(err error) bool {
	var perr *PermanentError
	return err != nil && errors.As(err, &perr)
}

// PermanentIf wraps the error as permanent if the condition is true.
// This makes it easy to conditionally mark errors as permanent.
//
// Example:
//
//	return apierrs.PermanentIf(err, statusCode == 400)
func PermanentIf(err error, condition bool) error {
	if err == nil || !condition {
		return err
	}
	return Permanent(err)
}
