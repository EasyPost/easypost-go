package easypost_test

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

// TestApiError tests that a bad API request returns an InvalidRequestError (a subclass of APIError), and that the
// error is parsed and pretty-printed correctly.
func (c *ClientTests) TestApiError() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	// Create a bad shipment so we can work with errors
	_, err := client.CreateShipment(&easypost.Shipment{})

	require.Error(err)
	if err, ok := err.(*easypost.InvalidRequestError); ok {
		assert.Equal(422, err.StatusCode)
		assert.Equal("PARAMETER.REQUIRED", err.Code)
		assert.Equal("Missing required parameter.", err.Message)
		assert.Equal(1, len(err.Errors))

		subError := err.Errors[0]
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
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}

	res.StatusCode = 0
	err := easypost.BuildErrorFromResponse(res)
	_, ok := err.(*easypost.ConnectionError)
	assert.True(ok)

	res.StatusCode = 100
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RetryError)
	assert.True(ok)

	res.StatusCode = 101
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RetryError)
	assert.True(ok)

	res.StatusCode = 102
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RetryError)
	assert.True(ok)

	res.StatusCode = 103
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RetryError)
	assert.True(ok)

	res.StatusCode = 300
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RedirectError)
	assert.True(ok)

	res.StatusCode = 301
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RedirectError)
	assert.True(ok)

	res.StatusCode = 302
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RedirectError)
	assert.True(ok)

	res.StatusCode = 303
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RedirectError)
	assert.True(ok)

	res.StatusCode = 304
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RedirectError)
	assert.True(ok)

	res.StatusCode = 305
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RedirectError)
	assert.True(ok)

	res.StatusCode = 306
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RedirectError)
	assert.True(ok)

	res.StatusCode = 307
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RedirectError)
	assert.True(ok)

	res.StatusCode = 308
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RedirectError)
	assert.True(ok)

	res.StatusCode = 400
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.BadRequestError)
	assert.True(ok)

	res.StatusCode = 401
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.UnauthorizedError)
	assert.True(ok)

	res.StatusCode = 402
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.PaymentError)
	assert.True(ok)

	res.StatusCode = 403
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.ForbiddenError)
	assert.True(ok)

	res.StatusCode = 404
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.NotFoundError)
	assert.True(ok)

	res.StatusCode = 405
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.MethodNotAllowedError)
	assert.True(ok)

	res.StatusCode = 407
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.ProxyError)
	assert.True(ok)

	res.StatusCode = 408
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.TimeoutError)
	assert.True(ok)

	res.StatusCode = 422
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.InvalidRequestError)
	assert.True(ok)

	res.StatusCode = 429
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.RateLimitError)
	assert.True(ok)

	res.StatusCode = 500
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.InternalServerError)
	assert.True(ok)

	res.StatusCode = 503
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.ServiceUnavailableError)
	assert.True(ok)

	res.StatusCode = 504
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.GatewayTimeoutError)
	assert.True(ok)

	res.StatusCode = 7000 // unaccounted for status code
	err = easypost.BuildErrorFromResponse(res)
	_, ok = err.(*easypost.UnknownHttpError)
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
		Body:       ioutil.NopCloser(strings.NewReader(fakeErrorResponse)),
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
		Body:       ioutil.NopCloser(strings.NewReader(fakeErrorResponse)),
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
		Body:       ioutil.NopCloser(strings.NewReader(fakeErrorResponse)),
	}

	err := easypost.BuildErrorFromResponse(res)

	errorMessage := err.Error()
	messages := []string{"message1", "message2", "message3", "message4", "message5", "message6"}
	for _, message := range messages {
		assert.Contains(errorMessage, message)
	}
}
