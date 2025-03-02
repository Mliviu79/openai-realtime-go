// Package types provides shared type definitions and utilities for the OpenAI Realtime API
package types

import (
	"fmt"
)

// SafeStringFromBytes converts byte data to a string safely for logging purposes
// This avoids the issue of interpreting binary data as Unicode characters
func SafeStringFromBytes(data []byte) string {
	return string(data)
}

// SafeStringFromInt converts an integer to a string safely for logging purposes
// This avoids the issue of interpreting an integer as a Unicode code point
func SafeStringFromInt(i int) string {
	return fmt.Sprintf("%d", i)
}

// SafeStringFromMessageType converts any message type constant to its string representation
// This function is a catch-all for any type of message constant that needs to be
// converted to a string for logging or display purposes
func SafeStringFromMessageType(msgType any) string {
	switch v := msgType.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	case int:
		return fmt.Sprintf("%d", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// BinaryData is a wrapper for byte slices with safe string conversion
type BinaryData []byte

// String returns a safe string representation of the binary data
func (d BinaryData) String() string {
	return string(d)
}

// NewBinaryData creates a new BinaryData from a byte slice
func NewBinaryData(data []byte) BinaryData {
	return BinaryData(data)
}
