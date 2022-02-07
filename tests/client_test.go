package easypost_test

import (
	"github.com/EasyPost/easypost-go"
)

func (c *ClientTests) TestClientTimeout() {
	timeout1 := 1000
	timeout2 := 10
	tempClient := c.TestClient()
	// set the timeout as part of the client initialization
	client := &easypost.Client{
		APIKey:  tempClient.APIKey,
		Timeout: timeout1,
	}
	assert := c.Assert()
	// test that property has been set correctly in the EasyPost.Client instance
	assert.Equal(timeout1, client.Timeout)
	// override the Timeout property after initialization
	client.Timeout = timeout2
	// test that property has been changed correctly
	assert.Equal(timeout2, client.Timeout)
}
