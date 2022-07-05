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
		assert.Equal("SHIPMENT.INVALID_PARAMS", err.Code)
		assert.Equal("Unable to create shipment, one or more parameters were invalid.", err.Message)
		// assert.Equal("TODO", err.Errors[0]) // TODO: You cannot currently get a single object from the Errors key due to bad deserialization
		// assert.Equal("TODO", err.Errors[1]) // TODO: You cannot currently get a single object from the Errors key due to bad deserialization
	}

	// Assert that the pretty printed error is the same
	errorString := err.Error()

	assert.Equal("SHIPMENT.INVALID_PARAMS Unable to create shipment, one or more parameters were invalid.", errorString)
}
