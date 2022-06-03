package easypost_test

func (c *ClientTests) TestPaymentMethodsAll() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	paymentMethods, err := client.ListPaymentMethods()
	require.NoError(err)

	assert.True(paymentMethods.PrimaryPaymentMethod != nil)
	assert.True(paymentMethods.SecondaryPaymentMethod != nil)
}
