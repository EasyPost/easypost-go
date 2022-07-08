package easypost_test

import "github.com/EasyPost/easypost-go/v2"

func (c *ClientTests) TestDeletePaymentMethod() {
	c.T().Skip("Skipping due to the lack of an available real payment method in tests")

	client := c.ProdClient()
	require := c.Require()

	err := client.DeletePaymentMethod(easypost.PrimaryPaymentMethodPriority)
	require.NoError(err)
}

func (c *ClientTests) TestFundWallet() {
	c.T().Skip("Skipping due to the lack of an available real payment method in tests")

	client := c.ProdClient()
	require := c.Require()

	err := client.FundWallet("2000", easypost.PrimaryPaymentMethodPriority)
	require.NoError(err)
}

func (c *ClientTests) TestRetrievePaymentMethods() {
	c.T().Skip("Skipping due to having to manually add and remove a payment method from the account")

	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	paymentMethods, err := client.RetrievePaymentMethods()
	require.NoError(err)

	assert.True(paymentMethods.PrimaryPaymentMethod != nil)
	assert.True(paymentMethods.SecondaryPaymentMethod != nil)
}
