package easypost

func (c *ClientTests) TestCreateEmbeddablesSession() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	childUsers, err := client.ListChildUsers(
		&ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	session, err := client.CreateEmbeddablesSession(
		&EmbeddablesSessionParameters{
			OriginHost: "https://example.com",
			UserId:     childUsers.Children[0].ID,
		},
	)
	require.NoError(err)

	assert.Equal("EmbeddablesSession", session.Object)
}
