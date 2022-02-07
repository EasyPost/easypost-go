package easypost_test

import (
	"github.com/EasyPost/easypost-go"
)

func GenerateTestShipment() *easypost.Shipment {

	lowestUSPS := &easypost.Rate{Rate: "1.0", Carrier: "USPS", Service: "ParcelSelect"}
	highestUSPS := &easypost.Rate{Rate: "10.0", Carrier: "USPS", Service: "Priority"}
	lowestUPS := &easypost.Rate{Rate: "2.0", Carrier: "UPS", Service: "ParcelSelect"}
	highestUPS := &easypost.Rate{Rate: "20.0", Carrier: "UPS", Service: "Priority"}
	highestFedex := &easypost.Rate{Rate: "50.0", Carrier: "FedEx", Service: "Overnight"}

	var rates []*easypost.Rate
	rates = append(rates, lowestUSPS, highestUSPS, lowestUPS, highestUPS, highestFedex)
	shipment := &easypost.Shipment{
		Rates: rates,
	}
	return shipment
}

func (c *ClientTests) TestLowestRateWithoutPreference() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()
	shipment := *GenerateTestShipment()

	rate, err := client.LowestRate(&shipment)

	require.NoError(err)
	assert.Equal("1.0", rate.Rate)
	assert.Equal("USPS", rate.Carrier)
	assert.Equal("ParcelSelect", rate.Service)

}

func (c *ClientTests) TestLowestRateWithCarrierPreference() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()
	shipment := *GenerateTestShipment()
	carrier := []string{"FedEx", "UPS"}
	rate, err := client.LowestRateWithCarrier(&shipment, carrier)

	require.NoError(err)
	assert.Equal("2.0", rate.Rate)
	assert.Equal("UPS", rate.Carrier)
	assert.Equal("ParcelSelect", rate.Service)

}

func (c *ClientTests) TestLowestRateWithCarrierAndServicePreference() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()
	shipment := *GenerateTestShipment()
	carrier := []string{"FedEx", "UPS"}
	service := []string{"Overnight"}
	rate, err := client.LowestRateWithCarrierAndService(&shipment, carrier, service)

	require.NoError(err)
	assert.Equal("50.0", rate.Rate)
	assert.Equal("FedEx", rate.Carrier)
	assert.Equal("Overnight", rate.Service)

}

func (c *ClientTests) TestNoRateAvailable() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()
	shipment := *GenerateTestShipment()
	carrier := []string{"Cainiao"}
	rate, err := client.LowestRateWithCarrier(&shipment, carrier)

	require.Error(err)
	assert.Equal(easypost.Rate{}, rate)

}
