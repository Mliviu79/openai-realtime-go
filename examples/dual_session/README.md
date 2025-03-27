# Dual Session Example

This example demonstrates how to run two different types of OpenAI Realtime API sessions simultaneously:
- A **transcription session** that processes audio and converts it to text
- A **conversation session** that responds to user messages using the GPT-4o model

This flow mimics a real voice assistant application where a user's speech is transcribed and then processed by an AI model.

## What This Example Shows

1. **Real Voice Assistant Flow**: Demonstrates a complete audio-to-response pipeline, similar to voice assistants like Siri or Alexa
2. **Concurrent Session Management**: How to create and connect to both a transcription session and a conversation session at the same time
3. **Inter-Session Communication**: How to pass transcribed text from the transcription session to the conversation session
4. **Parallel Processing**: Using goroutines to handle messages from both sessions concurrently
5. **Graceful Shutdown**: Properly cleaning up resources when the program exits
6. **Two-Step Response Process**: Demonstrating the required flow of sending a message and then explicitly requesting a response

## How It Works

1. The example creates a transcription session with the OpenAI API
2. It also creates a conversation session with the OpenAI API
3. Instead of sending real audio, it **simulates** the transcription process
4. The example directly injects simulated transcriptions into the conversation session
5. The conversation session sends these simulated transcriptions to the AI model and requests responses
6. The AI responds to the transcriptions as if they were real user messages

## Direct Simulation Approach

This example uses a direct simulation approach for clarity and reliability:

```go
// In a real app, we would send real audio:
// 1. audioChunk := captureRealAudioFromMicrophone()
// 2. msgClient.SendAudio(ctx, encodeToBase64(audioChunk), "")
// 3. Wait for the API to send back a transcription

// Instead, we directly simulate receiving transcriptions:
fmt.Printf("[SIMULATION] Pretending user said: '%s'\n", phrase)
fmt.Println("[SIMULATION] Sending audio chunk (simulated)")
fmt.Println("[SIMULATION] Audio processing started (simulated)")
// Skip sending actual audio and waiting for transcription
fmt.Printf("[SIMULATION] Received transcription: %s\n", phrase)

// Then send the simulated transcription to the conversation session
transcriptionCh <- phrase
```

This approach ensures the example works reliably without requiring real audio input or depending on OpenAI's transcription processing of empty audio packets.

### Why We Use Direct Simulation

1. **Reliability**: Sending empty audio packets doesn't trigger voice activity detection
2. **Clarity**: The simulation explicitly shows what's happening at each step
3. **Reproducibility**: The example produces the same results each time
4. **Simplicity**: No need to include actual audio files or real microphone input

In a real application, you would:
1. Create both session types just like in this example
2. Capture real audio from a microphone
3. Send that audio to the transcription session
4. Process the actual transcriptions from the API
5. Send those transcriptions to the conversation session

### Important: Requesting Responses

The OpenAI Realtime API requires a two-step process to get a response:
1. Send your message(s) using methods like `SendText`
2. Request the model to generate a response using `SendResponseCreate` with appropriate configuration

Without the second step, the model will not generate any response to your messages.

## Running the Example

1. Set your OpenAI API key:
   ```bash
   export OPENAI_API_KEY=your_api_key_here
   ```

2. Run the example:
   ```bash
   go run examples/dual_session/dual_session_example.go
   ```

## Expected Output

The example simulates sending three transcriptions and getting AI responses for each. You'll see output similar to:

```
Starting transcription session...
Created transcription session with ID: sess_BFeQ...

[IMPORTANT] This example uses direct simulation instead of real audio processing
[SIMULATION] In a real app, you would send actual audio and process real transcriptions
[SIMULATION] For demonstration purposes, we're simulating the transcription process

Starting conversation session...
Created conversation session with ID: sess_ABCd...
[CONVERSATION] Sent: Hello, I'm going to send you some audio transcriptions...
[CONVERSATION] AI is generating a response...
[CONVERSATION] AI: Hello! I'm ready to receive your audio transcriptions and respond to them. Please go ahead and send the audio when you're ready.
[CONVERSATION] Response complete
[CONVERSATION] Full response: Hello! I'm ready to receive your audio transcriptions and respond to them. Please go ahead and send the audio when you're ready.

====== SIMULATED AUDIO CHUNK 1 ======
[SIMULATION] Pretending user said: 'This is a test of the transcription system.'
[SIMULATION] Sending audio chunk (simulated)
[SIMULATION] Audio processing started (simulated)
[SIMULATION] Received transcription: This is a test of the transcription system.
[SIMULATION] Sent transcription to conversation session

====== PROCESSING TRANSCRIPTION ======
[CONVERSATION] Sent transcription: This is a test of the transcription system.
[CONVERSATION] Waiting for AI response...
[CONVERSATION] AI is generating a response...
[CONVERSATION] AI: I've received your transcription: "This is a test of the transcription system." The transcription system appears to be working correctly! The audio was successfully converted to text. If you'd like to test it further, you can try sending another audio sample with different content.
[CONVERSATION] Response complete
[CONVERSATION] Full response: I've received your transcription: "This is a test of the transcription system." The transcription system appears to be working correctly! The audio was successfully converted to text. If you'd like to test it further, you can try sending another audio sample with different content.
```

The example runs for about 60 seconds and then exits automatically. During this time, you'll see real AI-generated responses to your simulated transcriptions, demonstrating the full capability of both session types working together.

## Key Implementation Details

- **Direct Simulation**: Instead of relying on real audio processing, we simulate transcriptions for reliability
- **Separate WebSocket Connections**: Each session type uses its own WebSocket connection
- **Response Requests**: After sending messages to the conversation session, we explicitly request a response with `SendResponseCreate`
- **Concurrent Goroutines**: Multiple goroutines handle message sending, receiving, and processing
- **Channel-Based Communication**: Transcription results are passed between sessions using Go channels

This example provides a reliable demonstration of how to structure your application to use both transcription and conversation sessions in a complete voice assistant flow. 