package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestCreateClient() {
	require := c.Require()

	client := easypost.New("")

	_, err := client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.Error(err) // Throw an error when no API key present
}
