# Text Message Example

This example demonstrates how to correctly send text messages in a conversation session using the OpenAI Realtime API Go client.

## Purpose

This example is designed to verify that the messaging client correctly uses the `input_text` content type when sending text messages to the OpenAI Realtime API. It creates a simple conversation session, sends a text message, and then prints the model's response.

## How It Works

1. The example initializes a client with your OpenAI API key
2. Creates a new conversation session with the appropriate model and modality
3. Connects to the session using WebSockets
4. Sends a simple "Hello, world!" text message
5. **Requests the model to generate a response using `SendResponseCreate`**
6. Processes and prints the model's response as it streams back

> **Important Note**: When using the OpenAI Realtime API, you must explicitly request the model to generate a response after sending your message(s). This is done via the `SendResponseCreate` method, which triggers the model to process all previously sent messages and generate a response.

## Running the Example

To run this example:

1. Set your OpenAI API key as an environment variable:
   ```
   export OPENAI_API_KEY=your_api_key_here
   ```

2. Run the example:
   ```
   go run examples/text_message/text_message_example.go
   ```

## Expected Output

When the example runs successfully, you should see:

```
Creating session...
Created session with ID: sess_abc123...
Connecting to session...
Sending text message...
Requesting model response...
Waiting for response...
[Model's streaming response will appear here]
Response complete
```

If the message is sent with the correct content type, the model will respond normally. If there's an issue with the content type, you'll see an error message indicating an invalid content type.

## Two-Step Process for Getting Responses

The OpenAI Realtime API requires a two-step process to get a response:

1. Send your message(s) using methods like `SendText`, `SendAudio`, etc.
2. Request the model to generate a response using `SendResponseCreate` with appropriate configuration

This design allows you to:
- Send multiple messages before requesting a response
- Configure response parameters like modalities, temperature, etc.
- Control exactly when the model should start generating a response 