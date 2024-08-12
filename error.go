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
	Code string `json:"code,omitempty" url:"code,omitempty"`
	// Errors may be provided if there are multiple errors, for example if
	// multiple fields have invalid values.
	Errors []*Error `json:"errors,omitempty" url:"errors,omitempty"`
	// Field may be provided when the error relates to a specific field.
	Field string `json:"field,omitempty" url:"field,omitempty"`
	// Message is a human-readable description of the problem encountered.
	Message interface{} `json:"message,omitempty" url:"message,omitempty"`
	// Suggestion may be provided if the API can provide a suggestion to fix
	// the error.
	Suggestion string `json:"suggestion,omitempty" url:"suggestion,omitempty"`
}

func (e *Error) UnmarshalJSON(data []byte) error {
	type alias Error
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

// LocalError represents an error caused by the EasyPost library itself, such as validation or JSON serialization issues.
type LocalError struct {
	LibraryError // subtype of LibraryError
}

// Unwrap returns the underlying LibraryError object.
func (e *LocalError) Unwrap() error {
	return &e.LibraryError
}

// EndOfPaginationError is raised when there are no more pages to retrieve.
type EndOfPaginationError struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError object.
func (e *EndOfPaginationError) Unwrap() error {
	return &e.LibraryError
}

// newEndOfPaginationError returns a new EndOfPaginationError object.
func newEndOfPaginationError() *EndOfPaginationError {
	return &EndOfPaginationError{LocalError{LibraryError{Message: NoPagesLeftToRetrieve}}}
}

// FilteringError is raised when there is an issue while running a filtering operation.
type FilteringError struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError object.
func (e *FilteringError) Unwrap() error {
	return &e.LibraryError
}

// newFilteringError returns a new FilteringError object with the given message.
func newFilteringError(message string) *FilteringError {
	return &FilteringError{LocalError{LibraryError{Message: message}}}
}

// InvalidObjectError is raised when an object is invalid.
type InvalidObjectError struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError object.
func (e *InvalidObjectError) Unwrap() error {
	return &e.LibraryError
}

// newInvalidObjectError returns a new InvalidObjectError object with the given message.
func newInvalidObjectError(message string) *InvalidObjectError {
	return &InvalidObjectError{LocalError{LibraryError{Message: message}}}
}

// MissingPropertyError is raised when a required property is missing.
type MissingPropertyError struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError object.
func (e *MissingPropertyError) Unwrap() error {
	return &e.LibraryError
}

// newMissingPropertyError returns a new MissingPropertyError object with the given property.
func newMissingPropertyError(property string) *MissingPropertyError {
	message := MissingProperty + property
	return &MissingPropertyError{LocalError{LibraryError{Message: message}}}
}

// MissingWebhookSignatureError is raised when a webhook does not contain a valid HMAC signature.
type MissingWebhookSignatureError struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError object.
func (e *MissingWebhookSignatureError) Unwrap() error {
	return &e.LocalError
}

// newMissingWebhookSignatureError returns a new MissingWebhookSignatureError object.
func newMissingWebhookSignatureError() *MissingWebhookSignatureError {
	return &MissingWebhookSignatureError{LocalError{LibraryError{Message: MissingWebhookSignature}}}
}

// MismatchWebhookSignatureError is raised when a webhook received did not originate from EasyPost or had a webhook secret mismatch.
type MismatchWebhookSignatureError struct {
	LocalError // subtype of LocalError
}

// Unwrap returns the underlying LocalError object.
func (e *MismatchWebhookSignatureError) Unwrap() error {
	return &e.LocalError
}

// newMismatchWebhookSignatureError returns a new MismatchWebhookSignatureError object.
func newMismatchWebhookSignatureError() *MismatchWebhookSignatureError {
	return &MismatchWebhookSignatureError{LocalError{LibraryError{Message: MismatchWebhookSignature}}}
}

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

// Unwrap returns the underlying LocalError object.
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
// This is different from the Error class, which represents information about what triggered the failed request.
//
// The information from the top-level Error class is used to generate this error, and any sub-errors are stored in the Errors field.
type APIError struct {
	LibraryError // subtype of LibraryError
	// Code is a machine-readable status of the problem encountered.
	Code string
	// StatusCode is the HTTP numerical status code of the response.
	StatusCode int
	// Errors may be provided if there are details about server-side issues that caused the API request to fail.
	Errors []*Error `json:"errors,omitempty" url:"errors,omitempty"`
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

func (e *APIError) Unwrap() error {
	return &e.LibraryError
}

// BadRequestError is raised when the API returns a 400 status code.
type BadRequestError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *BadRequestError) Unwrap() error {
	return &e.APIError
}

// ConnectionError is raised when the API returns a 0 status code.
type ConnectionError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *ConnectionError) Unwrap() error {
	return &e.APIError
}

// GatewayTimeoutError is raised when the API returns a 504 status code.
type GatewayTimeoutError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *GatewayTimeoutError) Unwrap() error {
	return &e.APIError
}

// InternalServerError is raised when the API returns a 500 status code.
type InternalServerError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *InternalServerError) Unwrap() error {
	return &e.APIError
}

// InvalidRequestError is raised when the API returns a 422 status code.
type InvalidRequestError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *InvalidRequestError) Unwrap() error {
	return &e.APIError
}

// MethodNotAllowedError is raised when the API returns a 405 status code.
type MethodNotAllowedError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *MethodNotAllowedError) Unwrap() error {
	return &e.APIError
}

// NotFoundError is raised when the API returns a 404 status code.
type NotFoundError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *NotFoundError) Unwrap() error {
	return &e.APIError
}

// PaymentError is raised when the API returns a 402 status code.
type PaymentError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *PaymentError) Unwrap() error {
	return &e.APIError
}

// ProxyError is raised when the API returns a 407 status code.
type ProxyError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *ProxyError) Unwrap() error {
	return &e.APIError
}

// RateLimitError is raised when the API returns a 429 status code.
type RateLimitError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *RateLimitError) Unwrap() error {
	return &e.APIError
}

// RedirectError is raised when the API returns a 3xx status code.
type RedirectError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *RedirectError) Unwrap() error {
	return &e.APIError
}

// RetryError is raised when the API returns a 1xx status code.
type RetryError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *RetryError) Unwrap() error {
	return &e.APIError
}

// ServiceUnavailableError is raised when the API returns a 503 status code.
type ServiceUnavailableError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *ServiceUnavailableError) Unwrap() error {
	return &e.APIError
}

// SSLError is raised when there is an issue with the SSL certificate.
type SSLError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *SSLError) Unwrap() error {
	return &e.APIError
}

// TimeoutError is raised when the API returns a 408 status code.
type TimeoutError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *TimeoutError) Unwrap() error {
	return &e.APIError
}

// UnauthorizedError is raised when the API returns a 401 status code.
type UnauthorizedError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *UnauthorizedError) Unwrap() error {
	return &e.APIError
}

// ForbiddenError is raised when the API returns a 403 status code.
type ForbiddenError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
func (e *ForbiddenError) Unwrap() error {
	return &e.APIError
}

// UnknownHttpError is raised when the API returns an unrecognized status code.
type UnknownHttpError struct {
	APIError // subtype of APIError
}

// Unwrap returns the underlying APIError object.
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
	buf, _ := ioutil.ReadAll(response.Body)
	tmpError := &struct {
		Error *Error `json:"error,omitempty" url:"error,omitempty"`
	}{}

	if json.Unmarshal(buf, &tmpError) == nil {
		// extract the details from the temporary object (top-level Error class) and store them in the APIError object
		apiError.Message = tmpError.Error.Message.(string)
		apiError.Code = tmpError.Error.Code
		apiError.Errors = tmpError.Error.Errors
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
