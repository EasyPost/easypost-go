package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
)

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
