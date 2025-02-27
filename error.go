package easypost

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// LibraryError is the base type for all errors/exceptions in this EasyPost library.
type LibraryError struct {
	// Message is a human-readable error description.
	Message interface{}
}

// Error provides a pretty printed string of a LibraryError object.
func (e *LibraryError) Error() string {
	return e.Message.(string)
}

func (e *LibraryError) UnmarshalJSON(data []byte) error {
	type alias LibraryError
	tmpError := &struct {
		Message interface{} `json:"message,omitempty" url:"message,omitempty"`
		*alias
	}{
		alias: (*alias)(e),
	}
	if err := json.Unmarshal(data, &tmpError); err != nil {
		return err
	}

	// convert message to string
	messages := collectMessages(tmpError.Message, []string{})
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

// Local error types

// LocalError represents an error caused by the EasyPost library itself, such as validation or JSON serialization issues.
type LocalError struct {
	LibraryError // subtype of LibraryError
}

// Unwrap returns the underlying LibraryError error.
func (e *LocalError) Unwrap() error {
	return &e.LibraryError
}

// EndOfPaginationErrorType is raised when there are no more pages to retrieve.
// TODO: This type will be renamed to EndOfPaginationError in a future release to match the other error types once the EndOfPaginationError helper is removed.
type EndOfPaginationErrorType struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError error.
func (e *EndOfPaginationErrorType) Unwrap() error {
	return &e.LocalError
}

// newEndOfPaginationError returns a new EndOfPaginationErrorType object.
func newEndOfPaginationError() *EndOfPaginationErrorType {
	return &EndOfPaginationErrorType{LocalError{LibraryError{Message: NoPagesLeftToRetrieve}}}
}

// EndOfPaginationError is a singleton instance of EndOfPaginationErrorType.
// Deprecated: This helper will be removed in a future release. For access to the underlying message, use easypost.NoPagesLeftToRetrieve instead.
var EndOfPaginationError = newEndOfPaginationError()

// FilteringError is raised when there is an issue while running a filtering operation.
type FilteringError struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError error.
func (e *FilteringError) Unwrap() error {
	return &e.LocalError
}

// newFilteringError returns a new FilteringError object with the given message.
func newFilteringError(message string) *FilteringError {
	return &FilteringError{LocalError{LibraryError{Message: message}}}
}

// InvalidObjectError is raised when an object is invalid.
type InvalidObjectError struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError error.
func (e *InvalidObjectError) Unwrap() error {
	return &e.LocalError
}

// newInvalidObjectError returns a new InvalidObjectError object with the given message.
func newInvalidObjectError(message string) *InvalidObjectError {
	return &InvalidObjectError{LocalError{LibraryError{Message: message}}}
}

// MissingPropertyError is raised when a required property is missing.
type MissingPropertyError struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError error.
func (e *MissingPropertyError) Unwrap() error {
	return &e.LocalError
}

// newMissingPropertyError returns a new MissingPropertyError object with the given property.
func newMissingPropertyError(property string) *MissingPropertyError {
	message := MissingProperty + property
	return &MissingPropertyError{LocalError{LibraryError{Message: message}}}
}

// MissingWebhookSignatureErrorType is raised when a webhook does not contain a valid HMAC signature.
// TODO: This type will be renamed to MissingWebhookSignatureError in a future release to match the other error types once the MissingWebhookSignatureError helper is removed.
type MissingWebhookSignatureErrorType struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError error.
func (e *MissingWebhookSignatureErrorType) Unwrap() error {
	return &e.LocalError
}

// newMissingWebhookSignatureError returns a new MissingWebhookSignatureErrorType object.
func newMissingWebhookSignatureError() *MissingWebhookSignatureErrorType {
	return &MissingWebhookSignatureErrorType{LocalError{LibraryError{Message: MissingWebhookSignature}}}
}

// MissingWebhookSignatureError is raised when a webhook does not contain a valid HMAC signature.
// Deprecated: This helper will be removed in a future release. For access to the underlying message, use easypost.MissingWebhookSignature instead.
var MissingWebhookSignatureError = newMissingWebhookSignatureError()

// MismatchWebhookSignatureErrorType is raised when a webhook received did not originate from EasyPost or had a webhook secret mismatch.
// TODO: This type will be renamed to MismatchWebhookSignatureError in a future release to match the other error types once the MismatchWebhookSignatureError helper is removed.
type MismatchWebhookSignatureErrorType struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError error.
func (e *MismatchWebhookSignatureErrorType) Unwrap() error {
	return &e.LocalError
}

// newMismatchWebhookSignatureError returns a new MismatchWebhookSignatureErrorType object.
func newMismatchWebhookSignatureError() *MismatchWebhookSignatureErrorType {
	return &MismatchWebhookSignatureErrorType{LocalError{LibraryError{Message: MismatchWebhookSignature}}}
}

// MismatchWebhookSignatureError is raised when a webhook received did not originate from EasyPost or had a webhook secret mismatch.
// Deprecated: This helper will be removed in a future release. For access to the underlying message, use easypost.MismatchWebhookSignature instead.
var MismatchWebhookSignatureError = newMismatchWebhookSignatureError()

// ExternalApiError represents an error caused by an external API, such as a 3rd party HTTP API (not EasyPost).
type ExternalApiError struct {
	LibraryError // subtype of LibraryError
}

// Unwrap returns the underlying LibraryError object.
func (e *ExternalApiError) Unwrap() error {
	return &e.LibraryError
}

// newExternalApiError returns a new ExternalApiError object with the given message.
func newExternalApiError(message string) *ExternalApiError {
	return &ExternalApiError{LibraryError{Message: message}}
}

// InvalidFunctionError is raised when a function call is invalid or not allowed.
type InvalidFunctionError struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError error.
func (e *InvalidFunctionError) Unwrap() error {
	return &e.LocalError
}

// newInvalidFunctionError returns a new InvalidFunctionError object with the given message.
func newInvalidFunctionError(message string) *InvalidFunctionError {
	return &InvalidFunctionError{LocalError{LibraryError{Message: message}}}
}

// API/HTTP error types

// APIError represents an error that occurred while communicating with the EasyPost API.
//
// This is typically due to a specific HTTP status code, such as 4xx or 5xx.
//
// The information from the top-level `error` key is used to generate this error, and any sub-errors are stored in the Errors field.
type APIError struct {
	LibraryError // subtype of LibraryError
	// Code is a machine-readable status of the problem encountered.
	Code string `json:"code,omitempty" url:"code,omitempty"`
	// StatusCode is the HTTP numerical status code of the response.
	StatusCode int
	// Errors may be provided if there are details about server-side issues that caused the API request to fail.
	Errors interface{} `json:"errors,omitempty" url:"errors,omitempty"`
}

// Error provides a pretty printed string of an APIError object based on present data.
func (e *APIError) Error() string {
	var message string
	switch msg := e.Message.(type) {
	case string:
		message = msg
	case []interface{}:
		var messages []string
		for _, m := range msg {
			messages = append(messages, fmt.Sprint(m))
		}
		message = strings.Join(messages, ", ")
	case map[string]interface{}:
		message = extractMessagesFromMap(msg)
	default:
		message = fmt.Sprint(msg)
	}

	if message != "" {
		if e.Code != "" {
			return e.Code + " " + message
		}
		return message
	}
	if e.Code != "" {
		return e.Code
	}
	return fmt.Sprintf("%d %s", e.StatusCode, e.Code)
}

func extractMessagesFromMap(m map[string]interface{}) string {
	var messages []string
	for _, v := range m {
		switch val := v.(type) {
		case string:
			messages = append(messages, val)
		case []interface{}:
			for _, item := range val {
				messages = append(messages, fmt.Sprint(item))
			}
		case map[string]interface{}:
			messages = append(messages, extractMessagesFromMap(val))
		default:
			messages = append(messages, fmt.Sprint(val))
		}
	}
	return strings.Join(messages, ", ")
}

func (e *APIError) Unwrap() error {
	return &e.LibraryError
}

// FieldError represents a FieldError object returned by the EasyPost API.
//
// These are typically informational about why a request failed (server-side validation issues, missing data, etc.).
type FieldError struct {
	// Field may be provided when the error relates to a specific field.
	Field string `json:"field,omitempty" url:"field,omitempty"`
	// Message is a human-readable description of the problem encountered.
	Message interface{} `json:"message,omitempty" url:"message,omitempty"`
	// Suggestion is an occasional insight on how to correct the error
	Suggestion string `json:"suggestion,omitempty" url:"suggestion,omitempty"`
}

// BadRequestError is raised when the API returns a 400 status code.
type BadRequestError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *BadRequestError) Unwrap() error {
	return &e.APIError
}

// ConnectionError is raised when the API returns a 0 status code.
type ConnectionError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *ConnectionError) Unwrap() error {
	return &e.APIError
}

// GatewayTimeoutError is raised when the API returns a 504 status code.
type GatewayTimeoutError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *GatewayTimeoutError) Unwrap() error {
	return &e.APIError
}

// InternalServerError is raised when the API returns a 500 status code.
type InternalServerError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *InternalServerError) Unwrap() error {
	return &e.APIError
}

// InvalidRequestError is raised when the API returns a 422 status code.
type InvalidRequestError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *InvalidRequestError) Unwrap() error {
	return &e.APIError
}

// MethodNotAllowedError is raised when the API returns a 405 status code.
type MethodNotAllowedError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *MethodNotAllowedError) Unwrap() error {
	return &e.APIError
}

// NotFoundError is raised when the API returns a 404 status code.
type NotFoundError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *NotFoundError) Unwrap() error {
	return &e.APIError
}

// PaymentError is raised when the API returns a 402 status code.
type PaymentError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *PaymentError) Unwrap() error {
	return &e.APIError
}

// ProxyError is raised when the API returns a 407 status code.
type ProxyError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *ProxyError) Unwrap() error {
	return &e.APIError
}

// RateLimitError is raised when the API returns a 429 status code.
type RateLimitError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *RateLimitError) Unwrap() error {
	return &e.APIError
}

// RedirectError is raised when the API returns a 3xx status code.
type RedirectError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *RedirectError) Unwrap() error {
	return &e.APIError
}

// RetryError is raised when the API returns a 1xx status code.
type RetryError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *RetryError) Unwrap() error {
	return &e.APIError
}

// ServiceUnavailableError is raised when the API returns a 503 status code.
type ServiceUnavailableError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *ServiceUnavailableError) Unwrap() error {
	return &e.APIError
}

// SSLError is raised when there is an issue with the SSL certificate.
type SSLError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *SSLError) Unwrap() error {
	return &e.APIError
}

// TimeoutError is raised when the API returns a 408 status code.
type TimeoutError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *TimeoutError) Unwrap() error {
	return &e.APIError
}

// UnauthorizedError is raised when the API returns a 401 status code.
type UnauthorizedError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *UnauthorizedError) Unwrap() error {
	return &e.APIError
}

// ForbiddenError is raised when the API returns a 403 status code.
type ForbiddenError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *ForbiddenError) Unwrap() error {
	return &e.APIError
}

// UnknownHttpError is raised when the API returns an unrecognized status code.
type UnknownHttpError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError error.
func (e *UnknownHttpError) Unwrap() error {
	return &e.APIError
}

// BuildErrorFromResponse returns an APIError-based object based on the HTTP response.
// Do not pass a non-error response to this function.
func BuildErrorFromResponse(response *http.Response) error {
	// build the base APIError object from the response
	apiError := &APIError{
		StatusCode: response.StatusCode,
	}

	// deserialize the response body into a temporary object
	buf, _ := io.ReadAll(response.Body)
	tmpError := &struct {
		Error *struct {
			Code    string      `json:"code,omitempty"`
			Message interface{} `json:"message,omitempty"`
			Errors  interface{} `json:"errors,omitempty"`
		} `json:"error,omitempty"`
	}{}

	if json.Unmarshal(buf, &tmpError) == nil {
		// extract the details from the temporary object (top-level Error class) and store them in the APIError object
		if tmpError.Error != nil {
			apiError.Code = tmpError.Error.Code
			apiError.Message = tmpError.Error.Message
			if tmpError.Error.Errors != nil {
				switch errors := tmpError.Error.Errors.(type) {
				case []interface{}:
					var errorList []interface{}
					for _, err := range errors {
						switch err := err.(type) {
						case map[string]interface{}:
							fieldError := &FieldError{}
							if field, ok := err["field"].(string); ok {
								fieldError.Field = field
							}
							if message, ok := err["message"].(string); ok {
								fieldError.Message = message
							}
							if suggestion, ok := err["suggestion"].(string); ok {
								fieldError.Suggestion = suggestion
							}
							errorList = append(errorList, fieldError)
						default:
							errorList = append(errorList, fmt.Sprint(err))
						}
					}
					apiError.Errors = errorList
				default:
					apiError.Errors = errors
				}
			}
		}
	} else {
		// could not extract error details from the API response (or API did not return data, i.e. 1xx, 3xx or 5xx)
		if response.Status == "" {
			response.Status = ApiDidNotReturnErrorDetails
		}
		apiError.Message = response.Status
		apiError.Code = ApiErrorDetailsParsingError
		apiError.Errors = []*FieldError{}
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
	case 401:
		return &UnauthorizedError{APIError: *apiError}
	case 402:
		return &PaymentError{APIError: *apiError}
	case 403:
		return &ForbiddenError{APIError: *apiError}
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
	case 502, 504:
		return &GatewayTimeoutError{APIError: *apiError}
	default:
		return &UnknownHttpError{APIError: *apiError}
	}
}
