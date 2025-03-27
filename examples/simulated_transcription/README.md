# Simulated Transcription Example

This example demonstrates how to send pretend "transcribed audio" directly to a conversation session, **without** actually using a transcription service.

## Purpose

Unlike the dual_session example (which shows a complete audio-to-response pipeline), this example focuses only on the conversation part. It simulates what you might do if you:

1. Already have transcribed text from another source
2. Want to send that text to a conversation model as if it were transcribed audio
3. Need to format the messages appropriately

This approach is useful when:
- You're using a different transcription service
- You're processing pre-recorded transcriptions
- You're building a prototype and want to simulate the transcription part

## How It Works

1. The example creates a single conversation session
2. It sends an initial setup message to the model
3. It loops through a list of simulated "transcriptions"
4. For each transcription, it:
   - Formats the transcription as if it came from a transcription service
   - Sends it to the model
   - Requests a response
   - Processes and displays the response

## Key Features

This example illustrates:
- Proper formatting of "transcription" messages
- The two-step message sending process (send text, then request response)
- Processing of streamed responses
- Handling different message types

## Running the Example

1. Set your OpenAI API key:
   ```bash
   export OPENAI_API_KEY=your_api_key_here
   ```

2. Run the example:
   ```bash
   go run examples/simulated_transcription/simulated_transcription_example.go
   ```

## Expected Output

The example will send three simulated transcriptions, and you'll see:

```
Creating conversation session...
Created conversation session with ID: sess_xyz123...
Connecting to session...
Sent: Hello, I'm going to send you some audio transcriptions. Please respond to them as they come in.
Requesting initial response...
Waiting for AI response...
AI is generating a response...
Hello! I'll be happy to respond to your audio transcriptions. Feel free to send them whenever you're ready.
Response complete
Full response: Hello! I'll be happy to respond to your audio transcriptions. Feel free to send them whenever you're ready.

Sending simulated transcription 1: 'This is a test of the transcription system.'
Requesting response...
Waiting for AI response...
AI is generating a response...
I received your transcription: "This is a test of the transcription system." The transcription system appears to be working correctly! The text has been accurately captured.
Response complete
Full response: I received your transcription: "This is a test of the transcription system." The transcription system appears to be working correctly! The text has been accurately captured.

[additional transcriptions and responses...]
```

## Differences from Dual Session Example

The dual_session example demonstrates a complete pipeline:
- Real audio → Transcription Session → Text → Conversation Session → Response

This simulated example skips the first half:
- [No Audio] → [No Transcription Session] → Simulated Text → Conversation Session → Response

Use this example if you just want to see how to format and send transcriptions to the conversation model without dealing with actual audio and transcription. 