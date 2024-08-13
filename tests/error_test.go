package easypost_test

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/EasyPost/easypost-go/v4"
)

func (c *ClientTests) TestErrorInheritance() {
	assert, require := c.Assert(), c.Require()

	client := c.TestClient()
	client.APIKey = "not-a-real-key" // Invalid API key will result in 403 Forbidden error

	_, err := client.CreateAddress(
		&easypost.Address{
			Street1: "",
		},
		&easypost.CreateAddressOptions{},
	)

	require.Error(err)

	// LibraryError is the base of all errors thrown in this library, error should be able to be cast to this
	var libraryError *easypost.LibraryError
	if !errors.As(err, &libraryError) {
		assert.Fail("Error should be a LibraryError")
	}

	// APIError is the base of all errors thrown due to the API, error should be able to be cast to this
	var apiError *easypost.APIError
	if !errors.As(err, &apiError) {
		assert.Fail("Error should be an APIError")
	}
	// Should be able to access APIError-specific fields now that error has been cast
	assert.Equal(403, apiError.StatusCode)

	// LocalError is the base of all errors thrown due to the library itself, and a subtype of LibraryError, but separate from APIError, error should NOT be able to be cast to this
	var localError *easypost.LocalError
	if errors.As(err, &localError) {
		assert.Fail("Error should not be a LocalError")
	}

	// ForbiddenError is an error due to a 403 status code, and a subtype of APIError, error should be able to be cast to this
	var forbiddenError *easypost.ForbiddenError
	if !errors.As(err, &forbiddenError) {
		assert.Fail("Error should be a ForbiddenError")
	}

	// As a generic error, we can't access specific fields
	// ex. Can't access err.StatusCode
	// Once casted to a specific error type, we can access specific fields of that error type
	// ex. Can access forbiddenError.StatusCode or apiError.StatusCode
	assert.Equal(403, forbiddenError.StatusCode)
	assert.Equal(403, apiError.StatusCode)
}

// TestApiError tests that a bad API request returns an InvalidRequestError (a subclass of APIError), and that the
// error is parsed and pretty-printed correctly.
func (c *ClientTests) TestApiError() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	// Create a bad shipment so we can work with errors
	_, err := client.CreateShipment(&easypost.Shipment{})

	require.Error(err)

	var invalidRequestError *easypost.InvalidRequestError
	if errors.As(err, &invalidRequestError) {
		assert.Equal(422, invalidRequestError.StatusCode)
		assert.Equal("PARAMETER.REQUIRED", invalidRequestError.Code)
		assert.Equal("Missing required parameter.", invalidRequestError.Message)
		assert.Equal(1, len(invalidRequestError.Errors))

		subError := invalidRequestError.Errors[0]
		assert.Equal("shipment", subError.Field)
		assert.Equal("cannot be blank", subError.Message)
	}

	// Assert that the pretty printed error is the same
	errorString := err.Error()
	assert.Equal("PARAMETER.REQUIRED Missing required parameter.", errorString)
}

// TestApiErrorStatusCodes tests that the correct API error type is determined for each HTTP status code.
func (c *ClientTests) TestApiErrorStatusCodes() {
	assert := c.Assert()

	res := &http.Response{
		StatusCode: 0,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	res.StatusCode = 0
	err := easypost.BuildErrorFromResponse(res)
	ok := errors.As(err, new(*easypost.ConnectionError))
	assert.True(ok)

	res.StatusCode = 100
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RetryError))
	assert.True(ok)

	res.StatusCode = 101
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RetryError))
	assert.True(ok)

	res.StatusCode = 102
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RetryError))
	assert.True(ok)

	res.StatusCode = 103
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RetryError))
	assert.True(ok)

	res.StatusCode = 300
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RedirectError))
	assert.True(ok)

	res.StatusCode = 301
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RedirectError))
	assert.True(ok)

	res.StatusCode = 302
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RedirectError))
	assert.True(ok)

	res.StatusCode = 303
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RedirectError))
	assert.True(ok)

	res.StatusCode = 304
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RedirectError))
	assert.True(ok)

	res.StatusCode = 305
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RedirectError))
	assert.True(ok)

	res.StatusCode = 306
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RedirectError))
	assert.True(ok)

	res.StatusCode = 307
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RedirectError))
	assert.True(ok)

	res.StatusCode = 308
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RedirectError))
	assert.True(ok)

	res.StatusCode = 400
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.BadRequestError))
	assert.True(ok)

	res.StatusCode = 401
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.UnauthorizedError))
	assert.True(ok)

	res.StatusCode = 402
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.PaymentError))
	assert.True(ok)

	res.StatusCode = 403
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.ForbiddenError))
	assert.True(ok)

	res.StatusCode = 404
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.NotFoundError))
	assert.True(ok)

	res.StatusCode = 405
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.MethodNotAllowedError))
	assert.True(ok)

	res.StatusCode = 407
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.ProxyError))
	assert.True(ok)

	res.StatusCode = 408
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.TimeoutError))
	assert.True(ok)

	res.StatusCode = 422
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.InvalidRequestError))
	assert.True(ok)

	res.StatusCode = 429
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.RateLimitError))
	assert.True(ok)

	res.StatusCode = 500
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.InternalServerError))
	assert.True(ok)

	res.StatusCode = 503
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.ServiceUnavailableError))
	assert.True(ok)

	res.StatusCode = 504
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.GatewayTimeoutError))
	assert.True(ok)

	res.StatusCode = 7000 // unaccounted for status code
	err = easypost.BuildErrorFromResponse(res)
	ok = errors.As(err, new(*easypost.UnknownHttpError))
	assert.True(ok)
}

// TestApiErrorMessageParseArray tests that the internal API error parsing works correctly when the error message is an array.
func (c *ClientTests) TestApiErrorMessageParseArray() {
	assert := c.Assert()

	fakeErrorResponse := `{
		"error": {
			"code": "UNPROCESSABLE_ENTITY",
			"message": ["Bad format", "Bad format 2"],
			"errors": []
		}
	}`

	res := &http.Response{
		StatusCode: 422,
		Body:       io.NopCloser(strings.NewReader(fakeErrorResponse)),
	}

	err := easypost.BuildErrorFromResponse(res)

	errorMessage := err.Error()
	assert.Equal("UNPROCESSABLE_ENTITY Bad format, Bad format 2", errorMessage)
}

// TestErrorMessageParseMap tests that the internal API error parsing works correctly when the error message is a map.
func (c *ClientTests) TestErrorMessageParseMap() {
	assert := c.Assert()

	fakeErrorResponse := `{
		"error": {
			"code": "UNPROCESSABLE_ENTITY",
			"message": {
				"errors": [
					"Bad format", "Bad format 2"
				]
			},
			"errors": []
		}
	}`

	res := &http.Response{
		StatusCode: 422,
		Body:       io.NopCloser(strings.NewReader(fakeErrorResponse)),
	}

	err := easypost.BuildErrorFromResponse(res)

	errorMessage := err.Error()
	assert.Equal("UNPROCESSABLE_ENTITY Bad format, Bad format 2", errorMessage)
}

// TestErrorMessageParseExtreme tests that the internal API error parsing works correctly when the error message is a combination of maps and arrays.
func (c *ClientTests) TestErrorMessageParseExtreme() {
	assert := c.Assert()

	fakeErrorResponse := `{
		"error": {
			"code": "UNPROCESSABLE_ENTITY",
			"message": {
				"errors": [
					{
						"message1": "message1",
						"errors": ["message2", "message3"],
						"errors2": {
							"key": {
								"key2": "message4"
							}
						}
					},
					"message5",
					{
						"message6": "message6"
					}
				]
			},
			"errors": []
		}
	}`

	res := &http.Response{
		StatusCode: 422,
		Body:       io.NopCloser(strings.NewReader(fakeErrorResponse)),
	}

	err := easypost.BuildErrorFromResponse(res)

	errorMessage := err.Error()
	messages := []string{"message1", "message2", "message3", "message4", "message5", "message6"}
	for _, message := range messages {
		assert.Contains(errorMessage, message)
	}
}
