package easypost_test

func (c *ClientTests) TestCreditCardFund() {
	c.T().Skip("Skipping due to the lack of an available real credit card in tests")

	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	creditCard, err := client.FundCreditCard("20", "primary")
	require.NoError(err)

	assert.True(creditCard != nil)
}

func (c *ClientTests) TestDeleteCreditCard() {
	c.T().Skip("Skipping due to the lack of an available real credit card in tests")

	client := c.ProdClient()
	require := c.Require()

	err := client.DeleteCreditCard("card_1234")
	require.NoError(err)
}
