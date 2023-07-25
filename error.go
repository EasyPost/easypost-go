package easypost

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Error represents an Error object returned by the EasyPost API.
//
// These are typically informational about why a request failed (server-side validation issues, missing data, etc.).
//
// This is different from the LibraryError class, which represents exceptions in the EasyPost library, such as bad HTTP status codes or local validation issues.
type Error struct {
	// Code is a machine-readable code of the problem encountered.
	Code string `json:"code,omitempty"`
	// Errors may be provided if there are multiple errors, for example if
	// multiple fields have invalid values.
	Errors []*Error `json:"errors,omitempty"`
	// Field may be provided when the error relates to a specific field.
	Field string `json:"field,omitempty"`
	// Message is a human-readable description of the problem encountered.
	Message interface{} `json:"message,omitempty"`
	// Suggestion may be provided if the API can provide a suggestion to fix
	// the error.
	Suggestion string `json:"suggestion,omitempty"`
}

func (e *Error) UnmarshalJSON(data []byte) error {
	type alias Error
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

// Error provides a pretty printed string of an Error object based on present data.
func (e *Error) Error() string {
	if e.Message != "" {
		if e.Code != "" {
			return e.Code + " " + e.Message.(string)
		}
		return e.Message.(string)
	}
	return e.Code
}

// LibraryError is the base type for all errors/exceptions in this EasyPost library.
type LibraryError struct {
	// Message is a human-readable error description.
	Message string
}

// Error provides a pretty printed string of a LibraryError object.
func (e *LibraryError) Error() string {
	return e.Message
}

// Local error types

type LocalError struct {
	LibraryError
}

type ValidationError struct {
	LocalError
}

var EndOfPaginationError = &LocalError{LibraryError{Message: NoPagesLeftToRetrieve}}

type FilteringError struct {
	LocalError
}

type InvalidObjectError struct {
	ValidationError
}

type InvalidParameterError struct {
	ValidationError
}

type JsonError struct {
	LocalError
}

type JsonDeserializationError struct {
	JsonError
}

type JsonSerializationError struct {
	JsonError
}

type JsonNoDataError struct {
	JsonError
}

type MissingParameterError struct {
	LocalError
}

type MissingPropertyError struct {
	LocalError
}

var MissingWebhookSignatureError = &LocalError{LibraryError{Message: MissingWebhookSignature}}
var MismatchWebhookSignatureError = &LocalError{LibraryError{Message: MismatchWebhookSignature}}

// API/HTTP error types

// APIError represents an error that occurred while communicating with the EasyPost API.
//
// This is typically due to a specific HTTP status code, such as 4xx or 5xx.
//
// This is different from the Error class, which represents information about what triggered the failed request.
//
// The information from the top-level Error class is used to generate this error, and any sub-errors are stored in the Errors field.
type APIError struct {
	LibraryError
	// Code is a machine-readable status of the problem encountered.
	Code string
	// StatusCode is the HTTP numerical status code of the response.
	StatusCode int
	// Errors may be provided if there are details about server-side issues that caused the API request to fail.
	Errors []*Error `json:"errors,omitempty"`
}

// Error provides a pretty printed string of an APIError object based on present data.
func (e *APIError) Error() string {
	if e.Message != "" {
		if e.Code != "" {
			return e.Code + " " + e.Message
		}
		return e.Message
	}
	if e.Code != "" {
		return e.Code
	}
	return fmt.Sprintf("%d %s", e.StatusCode, e.Code)
}

type BadRequestError struct {
	APIError
}

type ConnectionError struct {
	APIError
}

type GatewayTimeoutError struct {
	APIError
}

type InternalServerError struct {
	APIError
}

type InvalidRequestError struct {
	APIError
}

type MethodNotAllowedError struct {
	APIError
}

type NotFoundError struct {
	APIError
}

type PaymentError struct {
	APIError
}

type ProxyError struct {
	APIError
}

type RateLimitError struct {
	APIError
}

type RedirectError struct {
	APIError
}

type RetryError struct {
	APIError
}

type ServiceUnavailableError struct {
	APIError
}

type SSLError struct {
	APIError
}

type TimeoutError struct {
	APIError
}

type UnauthorizedError struct {
	APIError
}

type UnknownHttpError struct {
	APIError
}

// FromErrorResponse returns an APIError-based object based on the HTTP response.
// Do not pass a non-error response to this function.
func FromErrorResponse(response *http.Response) error {
	// build the base APIError object from the response
	apiError := &APIError{
		StatusCode: response.StatusCode,
	}

	// deserialize the response body into a temporary object
	buf, _ := ioutil.ReadAll(response.Body)
	aux := &struct {
		Error *Error `json:"error,omitempty"`
	}{}

	if json.Unmarshal(buf, &aux) == nil {
		// extract the details from the temporary object (top-level Error class) and store them in the APIError object
		apiError.Message = aux.Error.Message.(string)
		apiError.Code = aux.Error.Code
		apiError.Errors = aux.Error.Errors
	} else {
		// could not extract error details from the API response (or API did not return data, i.e. 1xx, 3xx or 5xx)
		if response.Status == "" {
			response.Status = ApiDidNotReturnErrorDetails
		}
		apiError.Message = response.Status
		apiError.Code = ApiErrorDetailsParsingError
		apiError.Errors = []*Error{}
	}

	// return the appropriate error type based on the status code
	switch response.StatusCode {
	case 0:
		return &ConnectionError{APIError: *apiError}
	case 100, 101, 102, 103:
		return &RetryError{APIError: *apiError}
	case 300, 301, 302, 303, 304, 305, 306, 307, 308:
		return &RedirectError{APIError: *apiError}
	case 400:
		return &BadRequestError{APIError: *apiError}
	case 401, 403:
		return &UnauthorizedError{APIError: *apiError}
	case 402:
		return &PaymentError{APIError: *apiError}
	case 404:
		return &NotFoundError{APIError: *apiError}
	case 405:
		return &MethodNotAllowedError{APIError: *apiError}
	case 407:
		return &ProxyError{APIError: *apiError}
	case 408:
		return &TimeoutError{APIError: *apiError}
	case 422:
		return &InvalidRequestError{APIError: *apiError}
	case 429:
		return &RateLimitError{APIError: *apiError}
	case 500:
		return &InternalServerError{APIError: *apiError}
	case 503:
		return &ServiceUnavailableError{APIError: *apiError}
	case 504:
		return &GatewayTimeoutError{APIError: *apiError}
	default:
		return &UnknownHttpError{APIError: *apiError}
	}
}
