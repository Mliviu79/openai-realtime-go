# Changes from Original Repository

This document outlines the major changes and improvements made in this fork of the [WqyJh/go-openai-realtime](https://github.com/WqyJh/go-openai-realtime) library.

## Latest Updates (April 2025)

### v0.1.3

- **Fixed SendText Method**: Updated the `SendText` method in the messaging client to use `factory.InputTextContent` instead of `factory.TextContent` to properly match OpenAI's API expectations
- **Test Updates**: Updated the test suite to reflect the new content type used in the `SendText` method
- **Improved Examples**:
  - Added new `dual_session` example demonstrating how to run both transcription and conversation sessions simultaneously
  - Added new `simulated_transcription` example showing how to process simulated transcriptions
  - Updated examples to use the two-step process (send message + request response) required by the API
  - Enhanced documentation in all examples
- **Bug Fixes**: Addressed issues with failing tests and linter errors
- **Documentation**: Added detailed READMEs to all examples explaining their purpose and usage

- **OpenAI Realtime API Update**: Added support for the latest OpenAI Realtime API features:
  - Added dedicated support for transcription sessions via the new `/realtime/transcription_sessions` endpoint
  - Added ability to update transcription sessions via the `transcription_session.update` event
  - Added ability to request log probabilities in transcription results
  - Added new transcription models: `gpt-4o-transcribe` and `gpt-4o-mini-transcribe`
  - Added new voice options: `fable`, `onyx`, and `nova`
  - Added semantic VAD support with eagerness levels for more natural turn detection
  - Added input audio noise reduction for near-field and far-field microphones

- **Message Handling Fixes**: Fixed issues in `messages/types/response.go` and `messages/incoming/session_test.go`.
- **Session Management Enhancements**: Improved session type definitions in `session/session_types.go`.
- **Tool Functionality Updates**: Enhanced tool functionality in `session/tools.go`.

## Architecture Changes

### Package Structure

The original library used a completely flat structure with all functionality in the root directory (over 25 files including `client.go`, `server_event.go`, `client_event.go`, etc.). This fork has been completely reorganized with a modular package layout:

- `openaiClient`: Main client package for API connection (formerly `client.go`)
- `session`: Session management and configuration (formerly scattered across multiple files)
- `messages`: Message types and handling, with subpackages:
  - `incoming`: Incoming message types (formerly part of `server_event.go`)
  - `outgoing`: Outgoing message types (formerly part of `client_event.go`)
  - `types`: Shared message type definitions (formerly in `types.go`)
  - `factory`: Message factory functions
- `messaging`: High-level messaging interface
- `ws`: WebSocket connection management (formerly in `conn.go` and `ws.go`)
- `httpClient`: HTTP client for REST API endpoints (formerly in `api.go`)
- `logger`: Logging utilities (formerly in `log.go`)
- `apierrs`: Error handling (formerly in `permanent_error.go`)

### API Design

- **Type Safety**: Strong typing for all API calls and responses
- **Interface-Based Design**: Clear interfaces for each component
- **Options Pattern**: Functional options pattern for configurable API calls
- **Context Support**: Full support for Go contexts throughout the library

## Feature Enhancements

### Session Management

- More comprehensive session management with full parameter support
- Support for both REST API and WebSocket session creation/updates

### Message Handling

- Complete implementation of all 28 incoming message types
- Complete implementation of all 9 outgoing message types
- Strong typing for all message components
- Consistent error handling

### Audio Support

- Better audio format handling
- Support for all voice types
- Audio transcription configuration

### Tool/Function Calling

- Enhanced support for function definitions
- Improved type safety for tool parameters

### Turn Detection

- Full support for Voice Activity Detection (VAD)
- Configurable thresholds and behavior

## Documentation Improvements

- **Package Documentation**: Comprehensive godoc for all packages
- **Type Documentation**: Detailed documentation for all types and interfaces
- **Function Documentation**: Parameter and return value documentation
- **Examples**: Code examples for common use cases
- **README**: More comprehensive README with usage examples
- **Comprehensive Example**: Full example demonstrating all features

## Model Support

Added support for additional models:
- `gpt-4o-realtime-preview-2024-12-17`
- `gpt-4o-mini-realtime-preview`
- `gpt-4o-mini-realtime-preview-2024-12-17`

## Package Relationship

- **API Compatibility**: This fork is not backwards compatible with the original library due to the extensive refactoring
- **Module Path**: Uses a new module path to avoid conflicts with the original library

## Compatibility Changes

- Requires Go 1.23+ (up from 1.19 in the original)
- Not backwards compatible with the original API due to the extensive refactoring
- Improved error handling and return types

## Code Quality

- Comprehensive testing
- Improved error handling
- Better logging with configurable levels
- Adherence to Go best practices and idioms 
