package easypost_test

import (
	"github.com/EasyPost/easypost-go/v5"
)

func (c *ClientTests) TestGetLumaPromise() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipmentData := c.fixture.BasicShipment()
	lumaRequest := easypost.LumaRequest{
		Shipment:        *shipmentData,
		RulesetName:     c.fixture.LumaRulesetName(),
		PlannedShipDate: c.fixture.LumaPlannedShipDate(),
	}

	response, err := client.GetLumaPromise(&lumaRequest)
	require.NoError(err)

	assert.NotNil(response.LumaSelectedRate)
}
