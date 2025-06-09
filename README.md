# Vicidial Go Client Library

This Go library provides an API client for the Vicidial contact center system, focusing on call control, transfers, and agent state management.

## Project Structure

The codebase is organized into several files based on functionality:

- **session.go**: Contains the `VicidialSession` struct definition, `NewSession`, `Logout`, and basic session management functions
- **call_control.go**: Functions for call control like `Hangup`, `LiveHangup`, `SetDisposition`
- **transfer.go**: Functions related to transfers, including `MainXferSendRedirect`, `BlindTransfer`, and three-way calling
- **keepalive.go**: Functions for keeping the session alive and processing conference status
- **events.go**: Agent event handling for external integrations
- **api.go**: Common API request functions and helpers
- **constants.go**: API endpoint constants and other application constants

## Features

- Session management (login, logout)
- Call control (hangup, live hangup, blind transfer, consultative transfer)
- Agent status management
- Three-way calling support
- Call disposition handling
- Agent event notifications
- DTMF sending
- Call logging
- Automatic keepalive functionality

## Core Functions

### Session Management
- `NewSession`: Initialize a new agent session
- `Logout`: End an agent session

### Call Control 
- `Hangup`: Standard agent hangup
- `LiveHangup`: Hangup specific channels
- `SetDisposition`: Set call disposition
- `Ready`: Set agent to ready state
- `Pause`: Set agent to paused state
- `SetAgentStatus`: Update agent status

### Transfers
- `MainXferSendRedirect`: Handle transfers based on type
- `BlindTransfer`: Blind transfer to a number
- `LocalCloserTransfer`: Transfer to in-group
- `ParkCustomerDial`: Park customer and dial for 3-way
- `Leave3WayCall`: Leave a three-way call
- `SendManualDial`: Initiate a manual dial

### Keepalives
- `KeepAliveLoop`: Maintains agent session with periodic checks
- `ProcessConfExtenActions`: Process server commands

## Getting Started

### Installation

```
go get github.com/yourusername/vicidial-go-client
```

### Basic Usage

```go
package main

import (
	"fmt"
	"time"
	
	"github.com/yourusername/vicidial-go-client/agent/vicidial"
)

func main() {
	// Initialize configuration
	config := &vicidial.VicidialSettings{
		ServerURL:  "https://your-vicidial-server",
		User:       "agent-username",
		Pass:       "agent-password",
		CampaignID: "campaign-id",
	}
	
	// Create a new session
	session, err := vicidial.NewSession(config)
	if err != nil {
		fmt.Printf("Failed to create session: %v\n", err)
		return
	}
	
	// Start keepalive in background
	stopChan := make(chan struct{})
	defer close(stopChan)
	go session.KeepAliveLoop(stopChan, 5*time.Second)
	
	// Set agent status to ready
	_, err = session.Ready()
	if err != nil {
		fmt.Printf("Failed to set status to ready: %v\n", err)
	}
	
	// Perform a blind transfer
	resp, err := session.BlindTransfer("1234567890", "BLIND_XFER", "")
	if err != nil {
		fmt.Printf("Failed to transfer: %v\n", err)
	}
	
	// Set disposition after call
	_, err = session.SetDisposition("CALLBK", "", "123456")
	if err != nil {
		fmt.Printf("Failed to set disposition: %v\n", err)
	}
	
	// Logout when done
	session.Logout()
}
```

## Dependencies

- github.com/rs/zerolog - For logging
- golang.org/x/net/html - For HTML parsing during login

## License

[Your License Here] 
