package easypost_test

func (c *ClientTests) TestError() {
	client := c.TestClient()
	require := c.Require()

	// Create a bad shipment so we can work with errors
	_, err := client.CreateShipment(nil)

	require.Error(err)
}
