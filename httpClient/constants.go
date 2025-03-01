package httpClient

// APIType represents the type of API (Azure or OpenAI)
type APIType string

const (
	// APITypeOpenAI represents the standard OpenAI API
	APITypeOpenAI APIType = "OPEN_AI"

	// APITypeAzure represents the Azure OpenAI API
	APITypeAzure APIType = "AZURE"

	// OpenaiRealtimeAPIURLv1 is the base URL for the OpenAI Realtime API.
	OpenaiRealtimeAPIURLv1 = "wss://api.openai.com/v1/realtime"

	// OpenaiAPIURLv1 is the base URL for the OpenAI API.
	OpenaiAPIURLv1 = "https://api.openai.com/v1"

	// azureAPIVersion20241001Preview is the API version for Azure.
	azureAPIVersion20241001Preview = "2024-10-01-preview"
)
