package easypost

import (
	"reflect"
	"strings"
)

func (c *ClientTests) TestRateRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.BasicShipment())
	require.NoError(err)

	rate, err := client.GetRate(shipment.Rates[0].ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&Rate{}), reflect.TypeOf(rate))
	assert.True(strings.HasPrefix(rate.ID, "rate_"))
}

func (c *ClientTests) TestBetaStatelessRateRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	response, err := client.BetaGetStatelessRates(c.fixture.BasicShipment())
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&GetStatelessRatesResult{}), reflect.TypeOf(response))
	assert.NotEmpty(response.Rates)
	for _, rate := range response.Rates {
		assert.Equal(reflect.TypeOf(&StatelessRate{}), reflect.TypeOf(rate))
	}
}

func (c *ClientTests) TestBetaStatelessRateGetLowest() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	response, err := client.BetaGetStatelessRates(c.fixture.BasicShipment())
	require.NoError(err)

	lowestRate, err := client.LowestStatelessRate(response.Rates)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(StatelessRate{}), reflect.TypeOf(lowestRate))
	assert.Equal("GroundAdvantage", lowestRate.Service)
}

func (c *ClientTests) TestBetaStatelessRateGetLowestError() {
	client := c.TestClient()
	require := c.Require()

	response, err := client.BetaGetStatelessRates(c.fixture.BasicShipment())
	require.NoError(err)

	// Bad carrier
	_, err = client.LowestStatelessRateWithCarrier(response.Rates, []string{"BadCarrier"})
	require.Error(err)

	// Bad service
	_, err = client.LowestStatelessRateWithCarrierAndService(response.Rates, []string{"USPS"}, []string{"BadService"})
	require.Error(err)
}
