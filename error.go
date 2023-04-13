package easypost

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// APIError provides data on why an API request failed. It will be the type
// returned by Client methods when the HTTP API returns an HTTP error code. It
// will not be returned on network, parsing or context errors.
//
//	c := easypost.New(BadEasyPostAPIKey)
//	out, err := c.GetAPIKeys()
//	if err, ok := err.(*easypost.APIError); ok {
//		fmt.Println("EasyPost error code", err.Code)
//	}
type APIError struct {
	// Status is the HTTP text status of the response.
	Status string
	// StatusCode is the HTTP numerical status code of the response.
	StatusCode int
	// Code is a machine-readable error code returned by the API. It may be empty.
	Code string `json:"code,omitempty"`
	// Message is a human-readable error code returned by the API. It may be empty.
	Message interface{} `json:"message,omitempty"`
	// Field may be provided when the error relates to a specific field.
	Field string `json:"field,omitempty"`
	// Suggestion may be provided if the API can provide a suggestion to fix
	// the error.
	Suggestion string `json:"suggestion,omitempty"`
	// Errors may be provided if there are multiple errors, for example if
	// multiple fields have invalid values.
	Errors []*APIError `json:"errors,omitempty"`
}

func (e *APIError) UnmarshalJSON(data []byte) error {
	type alias APIError
	aux := &struct {
		Message interface{} `json:"message,omitempty"`
		*alias
	}{
		alias: (*alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// convert message to string
	messages := collectMessages(aux.Message, []string{})
	e.Message = strings.Join(messages, ", ")

	return nil
}

// Recursively traverses a JSON element to extract error messages and returns them as a comma-separated string.
func collectMessages(data interface{}, messages []string) []string {
	switch data := data.(type) {
	case []interface{}: // ["message", 123] or [{"key": "message"}, {"key": "message"}] or [ ["message", "message"], ["message", "message"] ]
		for _, value := range data {
			messages = collectMessages(value, messages)
		}
	case map[string]interface{}: // {"message": "value"} or {"message": ["value1", "value2"]} or {"message": [{"key": "value"}, {"key": "value"}]
		for _, value := range data {
			messages = collectMessages(value, messages)
		}
	default:
		messages = append(messages, fmt.Sprint(data))
	}

	return messages
}

type apiErrorResponse struct {
	Error *APIError `json:"error,omitempty"`
}

// Error provides a pretty printed string of an error object based on present data.
func (e *APIError) Error() string {
	if e.Message != "" {
		if e.Code != "" {
			return e.Code + " " + e.Message.(string)
		}
		return e.Message.(string)
	}
	if e.Code != "" {
		return e.Code
	}
	return fmt.Sprintf("%d %s", e.StatusCode, e.Status)
}

var EndOfPaginationError = errors.New("there are no more pages to retrieve")
