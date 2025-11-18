package easypost

func (c *ClientTests) TestCreateCustomerPortalAccountLink() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	childUsers, err := client.ListChildUsers(
		&ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	accountLink, err := client.CreateCustomerPortalAccountLink(
		&CustomerPortalAccountLinkParameters{
			SessionType: "account_onboarding",
			UserId:      childUsers.Children[0].ID,
			RefreshUrl:  "https://example.com/refresh",
			ReturnUrl:   "https://example.com/return",
		},
	)
	require.NoError(err)

	assert.Equal("CustomerPortalAccountLink", accountLink.Object)
}
