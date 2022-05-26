package easypost_test

func (c *ClientTests) TestApiKeysAll() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	apiKeys, err := client.GetAPIKeys()
	require.NoError(err)

	assert.NotNil(apiKeys)
}
