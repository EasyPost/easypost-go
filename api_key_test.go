package easypost

import (
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

func (c *ClientTests) TestApiKeysRetrieveChildUser() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	fakeChildUserId := "user_123456789"

	apiKeys, err := client.GetAPIKeysForUser(fakeChildUserId)
	require.Error(err)

	assert.Nil(apiKeys)
	assert.Equal(reflect.TypeOf(&FilteringError{}), reflect.TypeOf(err))
}

func (c *ClientTests) TestApiKeyLifecycle() {
	client := c.ReferralClient()
	assert, require := c.Assert(), c.Require()

	// Create an API key
	apiKey, err := client.CreateAPIKey("production")
	require.NoError(err)
	assert.NotNil(apiKey)
	assert.Contains(apiKey.ID, "ak_")
	assert.Equal("production", apiKey.Mode)

	// Disable the API key
	apiKey, err = client.DisableAPIKey(apiKey.ID)
	require.NoError(err)
	assert.False(apiKey.Active)

	// Enable the API key
	apiKey, err = client.EnableAPIKey(apiKey.ID)
	require.NoError(err)
	assert.True(apiKey.Active)

	// Delete the API key
	err = client.DeleteAPIKey(apiKey.ID)
	require.NoError(err)
}
