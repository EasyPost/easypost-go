package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v4"
)

func (c *ClientTests) TestOrderCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	order, err := client.CreateOrder(c.fixture.BasicOrder())
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Order{}), reflect.TypeOf(order))
	assert.True(strings.HasPrefix(order.ID, "order_"))
	assert.NotNil(order.Rates)
}

func (c *ClientTests) TestOrderRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	order, err := client.CreateOrder(c.fixture.BasicOrder())
	require.NoError(err)

	retrievedOrder, err := client.GetOrder(order.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Order{}), reflect.TypeOf(retrievedOrder))
	assert.Equal(order.ID, retrievedOrder.ID)
}

func (c *ClientTests) TestOrderGetRates() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	order, err := client.CreateOrder(c.fixture.BasicOrder())
	require.NoError(err)

	rates, err := client.GetOrderRates(order.ID)
	require.NoError(err)

	ratesList := rates.Rates

	assert.Equal(reflect.TypeOf([]*easypost.Rate{}), reflect.TypeOf(ratesList))
	for _, rate := range ratesList {
		assert.Equal(reflect.TypeOf(&easypost.Rate{}), reflect.TypeOf(rate))
	}
}

func (c *ClientTests) TestOrderBuyRate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	order, err := client.CreateOrder(c.fixture.BasicOrder())
	require.NoError(err)

	boughtOrder, err := client.BuyOrder(order.ID, c.fixture.USPS(), c.fixture.USPSService())
	require.NoError(err)

	shipmentsList := boughtOrder.Shipments

	for _, shipment := range shipmentsList {
		assert.NotNil(shipment.PostageLabel)
	}
}

func (c *ClientTests) TestOrderLowestRate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	order, err := client.CreateOrder(c.fixture.BasicOrder())
	require.NoError(err)

	// Test lowest rate with no filters
	lowestRate, err := client.LowestOrderRate(order)
	require.NoError(err)
	assert.Equal("GroundAdvantage", lowestRate.Service)
	assert.Equal("11.33", lowestRate.Rate)
	assert.Equal("USPS", lowestRate.Carrier)

	// Test lowest rate with service filter (this rate is higher than the lowest but should filter)
	lowestRate, err = client.LowestOrderRateWithCarrierAndService(order, nil, []string{"Priority"})
	require.NoError(err)
	assert.Equal("Priority", lowestRate.Service)
	assert.Equal("13.79", lowestRate.Rate)
	assert.Equal("USPS", lowestRate.Carrier)

	// Test lowest rate with carrier filter (should error due to bad carrier)
	lowestRate, err = client.LowestOrderRateWithCarrier(order, []string{"BAD_CARRIER"})
	assert.Error(err)
}
