package session

import (
	"encoding/json"
	"fmt"
)

//-----------------------------------------------------------------------------
// Tool Types
//-----------------------------------------------------------------------------

// FunctionDefinition defines a function that the model can call
// Deprecated: This struct is no longer used as the API expects a flattened structure.
// Use Tool directly instead.
type FunctionDefinition struct {
	// Name is the name of the function
	Name string `json:"name"`

	// Description explains what the function does
	Description string `json:"description"`

	// Parameters defines the inputs to the function
	Parameters json.RawMessage `json:"parameters"`
}

// FunctionChoice specifies a particular function to use
type FunctionChoice struct {
	// Name is the name of the function to use
	Name string `json:"name"`
}

// Tool represents a function that the model can call
type Tool struct {
	// Type is always "function" for now
	Type string `json:"type"`

	// Name is the name of the function
	Name string `json:"name"`

	// Description explains what the function does
	Description string `json:"description"`

	// Parameters defines the inputs to the function
	Parameters json.RawMessage `json:"parameters"`
}

// ToolChoice represents how the model should choose tools
type ToolChoice string

const (
	// ToolChoiceAuto lets the model decide when to use tools
	ToolChoiceAuto ToolChoice = "auto"

	// ToolChoiceNone prevents the model from using tools
	ToolChoiceNone ToolChoice = "none"

	// ToolChoiceRequired forces the model to use a tool
	ToolChoiceRequired ToolChoice = "required"

	// ToolChoiceFunction specifies a particular function to use
	ToolChoiceFunction ToolChoice = "function"
)

// ToolChoiceObj represents tool selection configuration
type ToolChoiceObj struct {
	// Type specifies how the model should choose tools
	Type ToolChoice `json:"type"`

	// Function specifies a particular function to use (when Type is "function")
	Function *FunctionChoice `json:"function,omitempty"`
}

// MarshalJSON implements custom JSON marshaling for ToolChoiceObj
func (tc ToolChoiceObj) MarshalJSON() ([]byte, error) {
	if tc.Type == "auto" || tc.Type == "none" || tc.Type == "required" {
		return json.Marshal(tc.Type)
	}

	if tc.Type == "function" && tc.Function != nil {
		return json.Marshal(map[string]interface{}{
			"type":     tc.Type,
			"function": tc.Function,
		})
	}

	return nil, fmt.Errorf("invalid tool choice: type=%s", tc.Type)
}

// UnmarshalJSON implements custom JSON unmarshaling for ToolChoiceObj
func (tc *ToolChoiceObj) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a string first
	var typeStr string
	if err := json.Unmarshal(data, &typeStr); err == nil {
		if typeStr == "auto" || typeStr == "none" || typeStr == "required" {
			tc.Type = ToolChoice(typeStr)
			tc.Function = nil
			return nil
		}
		return fmt.Errorf("invalid tool choice type: %s", typeStr)
	}

	// If not a string, try to unmarshal as an object
	var obj struct {
		Type     string          `json:"type"`
		Function *FunctionChoice `json:"function"`
	}

	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	if obj.Type != "function" {
		return fmt.Errorf("invalid tool choice type: %s", obj.Type)
	}

	tc.Type = ToolChoice(obj.Type)
	tc.Function = obj.Function
	return nil
}
