package easypost

import (
	"errors"
	"io"
	"net/http"
	"strings"
)

func (c *ClientTests) TestErrorInheritance() {
	assert, require := c.Assert(), c.Require()

	client := c.TestClient()
	client.APIKey = "not-a-real-key" // Invalid API key will result in 403 Forbidden error

	_, err := client.CreateAddress(
		&Address{
			Street1: "",
		},
		&CreateAddressOptions{},
	)

	require.Error(err)

	// LibraryError is the base of all errors thrown in this library, error should be able to be cast to this
	var libraryError *LibraryError
	if !errors.As(err, &libraryError) {
		assert.Fail("Error should be a LibraryError")
	}

	// APIError is the base of all errors thrown due to the API, error should be able to be cast to this
	var apiError *APIError
	if !errors.As(err, &apiError) {
		assert.Fail("Error should be an APIError")
	}
	// Should be able to access APIError-specific fields now that error has been cast
	assert.Equal(403, apiError.StatusCode)

	// LocalError is the base of all errors thrown due to the library itself, and a subtype of LibraryError, but separate from APIError, error should NOT be able to be cast to this
	var localError *LocalError
	if errors.As(err, &localError) {
		assert.Fail("Error should not be a LocalError")
	}

	// ForbiddenError is an error due to a 403 status code, and a subtype of APIError, error should be able to be cast to this
	var forbiddenError *ForbiddenError
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
	_, err := client.CreateShipment(&Shipment{})

	require.Error(err)

	var invalidRequestError *InvalidRequestError
	if errors.As(err, &invalidRequestError) {
		assert.Equal(422, invalidRequestError.StatusCode)
		assert.Equal("PARAMETER.REQUIRED", invalidRequestError.Code)
		assert.Equal("Missing required parameter.", invalidRequestError.Message)
		if errorsList, ok := invalidRequestError.Errors.([]interface{}); ok {
			assert.Equal(1, len(errorsList))
			if fieldError, ok := errorsList[0].(*FieldError); ok {
				assert.Equal("shipment", fieldError.Field)
				assert.Equal("cannot be blank", fieldError.Message)
			}
		}
	}

	// Assert that the pretty printed error is the same
	errorString := err.Error()
	assert.Equal("PARAMETER.REQUIRED Missing required parameter.", errorString)
}

// TestApiErrorAlternativeFormat tests that we assign properties of an error correctly when returned via the alternative format.
// NOTE: Claims (among other things) uses the alternative errors format.
func (c *ClientTests) TestApiErrorAlternativeFormat() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	createClaimParams := c.fixture.BasicClaim()
	createClaimParams.TrackingCode = "123" // Intentionally pass a bad tracking code
	_, err := client.CreateClaim(createClaimParams)

	require.Error(err)

	var notFoundError *NotFoundError
	if errors.As(err, &notFoundError) {
		assert.Equal(404, notFoundError.StatusCode)
		assert.Equal("NOT_FOUND", notFoundError.Code)
		assert.Equal("The requested resource could not be found.", notFoundError.Message)
		if errorsList, ok := notFoundError.Errors.([]interface{}); ok {
			assert.Equal("No eligible insurance found with provided tracking code.", errorsList[0])
		}
	}

	// Assert that the pretty printed error is the same
	errorString := err.Error()
	assert.Equal("NOT_FOUND The requested resource could not be found.", errorString)
}

// TestApiErrorStatusCodes tests that the correct API error type is determined for each HTTP status code.
func (c *ClientTests) TestApiErrorStatusCodes() {
	assert := c.Assert()

	res := &http.Response{
		StatusCode: 0,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	res.StatusCode = 0
	err := BuildErrorFromResponse(res)
	ok := errors.As(err, new(*ConnectionError))
	assert.True(ok)

	res.StatusCode = 100
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RetryError))
	assert.True(ok)

	res.StatusCode = 101
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RetryError))
	assert.True(ok)

	res.StatusCode = 102
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RetryError))
	assert.True(ok)

	res.StatusCode = 103
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RetryError))
	assert.True(ok)

	res.StatusCode = 300
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RedirectError))
	assert.True(ok)

	res.StatusCode = 301
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RedirectError))
	assert.True(ok)

	res.StatusCode = 302
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RedirectError))
	assert.True(ok)

	res.StatusCode = 303
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RedirectError))
	assert.True(ok)

	res.StatusCode = 304
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RedirectError))
	assert.True(ok)

	res.StatusCode = 305
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RedirectError))
	assert.True(ok)

	res.StatusCode = 306
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RedirectError))
	assert.True(ok)

	res.StatusCode = 307
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RedirectError))
	assert.True(ok)

	res.StatusCode = 308
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RedirectError))
	assert.True(ok)

	res.StatusCode = 400
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*BadRequestError))
	assert.True(ok)

	res.StatusCode = 401
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*UnauthorizedError))
	assert.True(ok)

	res.StatusCode = 402
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*PaymentError))
	assert.True(ok)

	res.StatusCode = 403
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*ForbiddenError))
	assert.True(ok)

	res.StatusCode = 404
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*NotFoundError))
	assert.True(ok)

	res.StatusCode = 405
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*MethodNotAllowedError))
	assert.True(ok)

	res.StatusCode = 407
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*ProxyError))
	assert.True(ok)

	res.StatusCode = 408
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*TimeoutError))
	assert.True(ok)

	res.StatusCode = 422
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*InvalidRequestError))
	assert.True(ok)

	res.StatusCode = 429
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*RateLimitError))
	assert.True(ok)

	res.StatusCode = 500
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*InternalServerError))
	assert.True(ok)

	res.StatusCode = 503
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*ServiceUnavailableError))
	assert.True(ok)

	res.StatusCode = 504
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*GatewayTimeoutError))
	assert.True(ok)

	res.StatusCode = 7000 // unaccounted for status code
	err = BuildErrorFromResponse(res)
	ok = errors.As(err, new(*UnknownHttpError))
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

	err := BuildErrorFromResponse(res)

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

	err := BuildErrorFromResponse(res)

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

	err := BuildErrorFromResponse(res)

	errorMessage := err.Error()
	messages := []string{"message1", "message2", "message3", "message4", "message5", "message6"}
	for _, message := range messages {
		assert.Contains(errorMessage, message)
	}
}
