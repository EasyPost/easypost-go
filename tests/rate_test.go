package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestRateRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.BasicShipment())
	require.NoError(err)

	rate, err := client.GetRate(shipment.Rates[0].ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Rate{}), reflect.TypeOf(rate))
	assert.True(strings.HasPrefix(rate.ID, "rate_"))
}

func (c *ClientTests) TestBetaStatelessRateRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	rates, err := client.BetaGetStatelessRates(c.fixture.BasicShipment())
	require.NoError(err)

	assert.Equal(reflect.TypeOf([]*easypost.StatelessRate{}), reflect.TypeOf(rates))
}

func (c *ClientTests) TestBetaStatelessRateGetLowest() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	rates, err := client.BetaGetStatelessRates(c.fixture.BasicShipment())
	require.NoError(err)

	lowestRate, err := client.LowestStatelessRate(rates)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(easypost.StatelessRate{}), reflect.TypeOf(lowestRate))
	assert.Equal("First", lowestRate.Service)
}
