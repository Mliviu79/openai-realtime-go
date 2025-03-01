package apierrs

import (
	"errors"
	"fmt"
	"testing"
)

func TestPermanentError(t *testing.T) {
	// Test creating a permanent error
	originalErr := errors.New("original error")
	permErr := Permanent(originalErr)

	// Check that it implements the error interface
	if permErr.Error() != originalErr.Error() {
		t.Errorf("Expected error message %q, got %q", originalErr.Error(), permErr.Error())
	}

	// Check that Unwrap returns the original error
	unwrappedErr := errors.Unwrap(permErr)
	if unwrappedErr != originalErr {
		t.Errorf("Expected Unwrap to return original error, got %v", unwrappedErr)
	}

	// Test passing nil
	nilPermErr := Permanent(nil)
	if nilPermErr != nil {
		t.Error("Expected Permanent(nil) to return nil")
	}

	// Test wrapping an already permanent error
	doublePermErr := Permanent(permErr)
	if doublePermErr != permErr {
		t.Error("Expected wrapping an already permanent error to return the same error")
	}

	// Test Is functionality
	if !errors.Is(permErr, &PermanentError{}) {
		t.Error("Expected Is to return true for PermanentError")
	}
}

func TestIsPermanent(t *testing.T) {
	// Create test errors
	regularErr := errors.New("regular error")
	permErr := Permanent(regularErr)
	wrappedPermErr := fmt.Errorf("wrapped: %w", permErr)

	// Test with nil
	if IsPermanent(nil) {
		t.Error("Expected IsPermanent(nil) to return false")
	}

	// Test with regular error
	if IsPermanent(regularErr) {
		t.Error("Expected IsPermanent(regularErr) to return false")
	}

	// Test with permanent error
	if !IsPermanent(permErr) {
		t.Error("Expected IsPermanent(permErr) to return true")
	}

	// Test with wrapped permanent error
	if !IsPermanent(wrappedPermErr) {
		t.Error("Expected IsPermanent(wrappedPermErr) to return true")
	}
}

func TestPermanentIf(t *testing.T) {
	// Create a regular error
	regularErr := errors.New("regular error")

	// Test with condition = true
	permErr := PermanentIf(regularErr, true)
	if !IsPermanent(permErr) {
		t.Error("Expected PermanentIf with true condition to return permanent error")
	}

	// Test with condition = false
	nonPermErr := PermanentIf(regularErr, false)
	if nonPermErr != regularErr {
		t.Error("Expected PermanentIf with false condition to return original error")
	}
	if IsPermanent(nonPermErr) {
		t.Error("Expected PermanentIf with false condition to not be permanent")
	}

	// Test with nil error
	nilErr := PermanentIf(nil, true)
	if nilErr != nil {
		t.Error("Expected PermanentIf(nil, true) to return nil")
	}

	// Test with nil error and false condition
	nilErrFalse := PermanentIf(nil, false)
	if nilErrFalse != nil {
		t.Error("Expected PermanentIf(nil, false) to return nil")
	}
}

func TestPermanentErrorChaining(t *testing.T) {
	// Test error chaining
	baseErr := errors.New("base error")
	permErr := Permanent(baseErr)
	wrappedErr := fmt.Errorf("wrapped: %w", permErr)
	furtherWrappedErr := fmt.Errorf("further wrapped: %w", wrappedErr)

	// Check that IsPermanent can see through the wrapping
	if !IsPermanent(furtherWrappedErr) {
		t.Error("Expected IsPermanent to detect permanent error through multiple wrappings")
	}

	// Check that we can get the original error message
	if !errors.Is(furtherWrappedErr, baseErr) {
		t.Error("Expected errors.Is to detect the base error through multiple wrappings")
	}

	// Check that we can get back the PermanentError
	var extractedPermErr *PermanentError
	if !errors.As(furtherWrappedErr, &extractedPermErr) {
		t.Error("Expected errors.As to extract the PermanentError")
	}

	// Check that the unwrapped error is correct
	unwrappedErr := extractedPermErr.Unwrap()
	if unwrappedErr != baseErr {
		t.Errorf("Expected unwrapped error to be %v, got %v", baseErr, unwrappedErr)
	}
}
