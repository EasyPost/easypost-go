package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
	"reflect"
)

func (c *ClientTests) TestShipmentCreateWithCarbonOffset() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipmentWithCarbonOffset(c.fixture.BasicCarbonOffsetShipment())
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Shipment{}), reflect.TypeOf(shipment))
	assert.NotNil(shipment.Rates)

	rate := shipment.Rates[0]
	assert.NotNil(rate.CarbonOffset)
}

func (c *ClientTests) TestShipmentBuyWithCarbonOffset() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.FullCarbonOffsetShipment())
	require.NoError(err)

	lowestRate, err := client.LowestShipmentRate(shipment)
	require.NoError(err)

	boughtShipment, err := client.BuyShipmentWithCarbonOffset(shipment.ID, &lowestRate, "")
	require.NoError(err)

	assert.NotNil(boughtShipment.PostageLabel)

	assert.NotNil(boughtShipment.Fees)
	carbonOffsetExists := false
	for _, fee := range boughtShipment.Fees {
		if fee.Type == "CarbonOffsetFee" {
			carbonOffsetExists = true
		}
	}
	assert.True(carbonOffsetExists)
}

func (c *ClientTests) TestShipmentOneCallBuyWithCarbonOffset() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipmentWithCarbonOffset(c.fixture.OneCallBuyCarbonOffsetShipment())
	require.NoError(err)

	assert.NotNil(shipment.PostageLabel)

	assert.NotNil(shipment.Fees)
	carbonOffsetExists := false
	for _, fee := range shipment.Fees {
		if fee.Type == "CarbonOffsetFee" {
			carbonOffsetExists = true
		}
	}
	assert.True(carbonOffsetExists)
}

func (c *ClientTests) TestShipmentReRateWithCarbonOffset() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyCarbonOffsetShipment())
	require.NoError(err)

	baseRates := shipment.Rates

	newRatesWithCarbon, err := client.RerateShipmentWithCarbonOffset(shipment.ID)
	require.NoError(err)

	newRatesWithoutCarbon, err := client.RerateShipment(shipment.ID)
	require.NoError(err)

	baseRate := baseRates[0]
	newRateWithCarbon := newRatesWithCarbon[0]
	newRateWithoutCarbon := newRatesWithoutCarbon[0]

	assert.Nil(baseRate.CarbonOffset)
	assert.NotNil(newRateWithCarbon.CarbonOffset)
	assert.Nil(newRateWithoutCarbon.CarbonOffset)
}
