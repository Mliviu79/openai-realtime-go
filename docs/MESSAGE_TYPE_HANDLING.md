# Message Type Handling in openai-realtime-go

This document outlines best practices for handling message types in the openai-realtime-go package.

## Problem

Go's type system treats string-like custom types (like `outgoing.OutMsgType` and `incoming.RcvdMsgType`) 
as distinct from regular strings. When converting these types to strings using a direct type conversion 
(`string(someType)`), Go treats numeric values as Unicode code points rather than converting them to 
their string representation. This can lead to unexpected behavior in logging and error messages.

## Solution

The package provides several methods for safely handling message types:

### 1. For Outgoing Message Types

Use the `OutMsgType()` method provided by the `OutMsg` interface:

```go
// Get the type as a string
msgType := outMsg.OutMsgType()

// Log the message type
log.Info().Str("type", msgType).Msg("Sent message")
```

### 2. For Incoming Message Types

Use the `RcvdMsgType()` method and convert the result to a string:

```go
// Get the message type as a string
messageTypeStr := string(msg.RcvdMsgType())

// Log the message type
log.Info().Str("message_type", messageTypeStr).Msg("Received message")
```

### 3. For WebSocket Message Types

Use explicit conversion to int for logging:

```go
// Log WebSocket message type
log.Info().Int("message_type", int(msgType)).Msg("Received WebSocket message")
```

### 4. For Binary Data

Use standard string conversion for byte arrays:

```go
// Convert binary data to string for logging
safeDataString := string(data)

// Log the data
log.Info().Str("data", safeDataString).Msg("Processed data")
```

### 5. Utility Functions

The package provides utility functions in `messages/types/converters.go` for safe conversion:

```go
// For byte arrays
safeString := types.SafeStringFromBytes(data)

// For integers
safeString := types.SafeStringFromInt(msgType)

// For any message type
safeString := types.SafeStringFromMessageType(anyTypeValue)
```

## Best Practices

1. **Never** use direct type casting (`string(numericValue)`) for numeric message types
2. **Always** use the appropriate interface method (`OutMsgType()` or `RcvdMsgType()`)
3. When comparing message types, convert constants to strings first:
   ```go
   audioBufferAppendType := string(outgoing.OutMsgTypeAudioBufferAppend)
   if msgType != audioBufferAppendType { ... }
   ```
4. Use the utility functions in `messages/types/converters.go` for complex cases

## Examples

### Good Example

```go
// Get message type as string
msgType := outMsg.OutMsgType()
audioBufferAppendType := string(outgoing.OutMsgTypeAudioBufferAppend)

// Compare as strings
if msgType != audioBufferAppendType {
    log.Info().Str("type", msgType).Msg("Processed message")
}
```

### Bad Example

```go
// WRONG: Direct casting of numeric types
log.Info().Str("type", string(msgType)).Msg("Processed message")

// WRONG: Comparing different types
if outMsg.OutMsgType() != outgoing.OutMsgTypeAudioBufferAppend { ... }
```

By following these practices, we ensure consistent and correct handling of message types throughout the codebase. 
