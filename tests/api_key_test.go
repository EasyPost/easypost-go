package easypost_test

import (
	"github.com/EasyPost/easypost-go/v4"
	"reflect"
)

func (c *ClientTests) TestApiKeysAll() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	apiKeys, err := client.GetAPIKeys()
	require.NoError(err)

	assert.NotNil(apiKeys)
}

func (c *ClientTests) TestApiKeysForParentUser() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	authenticatedUser, err := client.RetrieveMe()
	require.NoError(err)

	apiKeys, err := client.GetAPIKeysForUser(authenticatedUser.ID)
	require.NoError(err)

	assert.NotNil(apiKeys)
}

func (c *ClientTests) TestApiKeysForChildUser() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	fakeChildUserId := "user_123456789"

	apiKeys, err := client.GetAPIKeysForUser(fakeChildUserId)
	require.Error(err)

	assert.Nil(apiKeys)
	assert.Equal(reflect.TypeOf(&easypost.FilteringError{}), reflect.TypeOf(err))
}
