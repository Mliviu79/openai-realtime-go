# Changes from Original Repository

This document outlines the major changes and improvements made in this fork of the [WqyJh/go-openai-realtime](https://github.com/WqyJh/go-openai-realtime) library.

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