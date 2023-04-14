package easypost_test

import (
	"encoding/json"
	"github.com/EasyPost/easypost-go/v2"
)

type apiErrorResponse struct {
	Error *easypost.APIError `json:"error,omitempty"`
}

func (c *ClientTests) TestError() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	// Create a bad shipment so we can work with errors
	_, err := client.CreateShipment(&easypost.Shipment{})

	require.Error(err)
	if err, ok := err.(*easypost.APIError); ok {
		assert.Equal(422, err.StatusCode)
		assert.Equal("PARAMETER.REQUIRED", err.Code)
		assert.Equal("Missing required parameter.", err.Message)
		// assert.Equal("TODO", err.Errors[0]) // TODO: You cannot currently get a single object from the Errors key due to bad deserialization
		// assert.Equal("TODO", err.Errors[1]) // TODO: You cannot currently get a single object from the Errors key due to bad deserialization
	}

	// Assert that the pretty printed error is the same
	errorString := err.Error()

	assert.Equal("PARAMETER.REQUIRED Missing required parameter.", errorString)
}

func (c *ClientTests) TestErrorParseArray() {
	assert := c.Assert()

	fakeErrorResponse := `{
		"error": {
			"code": "UNPROCESSABLE_ENTITY",
			"message": ["Bad format", "Bad format 2"],
			"errors": []
		}
	}`

	buf := []byte(fakeErrorResponse)
	apiErr := &easypost.APIError{Status: "Status", StatusCode: 404}
	err := json.Unmarshal(buf, &apiErrorResponse{Error: apiErr})
	if err != nil {
		panic(err)
	}

	errorMessage := apiErr.Message

	assert.Equal("Bad format, Bad format 2", errorMessage)
}

func (c *ClientTests) TestErrorParseMap() {
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

	buf := []byte(fakeErrorResponse)
	apiErr := &easypost.APIError{Status: "Status", StatusCode: 404}
	err := json.Unmarshal(buf, &apiErrorResponse{Error: apiErr})
	if err != nil {
		panic(err)
	}

	errorMessage := apiErr.Message

	assert.Equal("Bad format, Bad format 2", errorMessage)
}

func (c *ClientTests) TestErrorParseExtreme() {
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

	buf := []byte(fakeErrorResponse)
	apiErr := &easypost.APIError{Status: "Status", StatusCode: 404}
	err := json.Unmarshal(buf, &apiErrorResponse{Error: apiErr})
	if err != nil {
		panic(err)
	}

	errorMessage := apiErr.Message

	assert.Equal("message1, message2, message3, message4, message5, message6", errorMessage)
}
