package easypost_test

import (
	"strconv"
	"strings"

	"github.com/EasyPost/easypost-go"
)

func (c *ClientTests) TestShipmentCreation() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	to, err := client.CreateAddress(
		&easypost.Address{
			Name:    "Elmer Fudd",
			Street1: "55 Sparks St.",
			City:    "Ottawa",
			State:   "ON",
			Zip:     "k1p5a5",
			Country: "CA",
			Phone:   "613-555-1212",
		},
		nil,
	)
	require.NoError(err)

	from, err := client.CreateAddress(
		&easypost.Address{
			Company: "EasyPost",
			Street1: "One Montgomery St",
			Street2: "Ste 400",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94104",
			Phone:   "415-456-7890",
		},
		nil,
	)
	require.NoError(err)

	parcel, err := client.CreateParcel(
		&easypost.Parcel{
			Length: 10.2,
			Width:  7.8,
			Height: 4.3,
			Weight: 21.2,
		},
	)
	require.NoError(err)

	customsItem, err := client.CreateCustomsItem(
		&easypost.CustomsItem{
			Description:    "EasyPost t-shirts",
			HSTariffNumber: "123456",
			OriginCountry:  "US",
			Quantity:       2,
			Value:          96.27,
			Weight:         21.1,
		},
	)
	require.NoError(err)

	customsInfo, err := client.CreateCustomsInfo(
		&easypost.CustomsInfo{
			CustomsCertify:    true,
			CustomsSigner:     "Wile E. Coyote",
			ContentsType:      "gift",
			EELPFC:            "NOEEI 30.37(a)",
			NonDeliveryOption: "return",
			RestrictionType:   "none",
			CustomsItems:      []*easypost.CustomsItem{customsItem},
		},
	)
	require.NoError(err)

	shipment, err := client.CreateShipment(
		&easypost.Shipment{
			ToAddress:   to,
			FromAddress: from,
			Parcel:      parcel,
			CustomsInfo: customsInfo,
		},
	)
	require.NoError(err)

	if assert.NotEmpty(shipment.BuyerAddress) {
		assert.Equal(to.Country, shipment.BuyerAddress.Country)
		assert.Equal(to.Phone, shipment.BuyerAddress.Phone)
		assert.Equal(to.Street1, shipment.BuyerAddress.Street1)
		assert.Equal(to.Zip, shipment.BuyerAddress.Zip)
	}
	if assert.NotEmpty(shipment.Parcel) {
		assert.Equal(parcel.Height, shipment.Parcel.Height)
		assert.Equal(parcel.Weight, shipment.Parcel.Weight)
		assert.Equal(parcel.Width, shipment.Parcel.Width)
	}
	if assert.NotEmpty(shipment.CustomsInfo) {
		assert.Equal(
			customsInfo.ContentsExplanation,
			shipment.CustomsInfo.ContentsExplanation,
		)
		if assert.NotEmpty(shipment.CustomsInfo.CustomsItems) {
			assert.Equal(
				customsItem.Description,
				shipment.CustomsInfo.CustomsItems[0].Description,
			)
			assert.Equal(
				customsItem.HSTariffNumber,
				shipment.CustomsInfo.CustomsItems[0].HSTariffNumber,
			)
			assert.Equal(
				customsItem.Value, shipment.CustomsInfo.CustomsItems[0].Value,
			)
		}
	}

	require.NotEmpty(shipment.Rates)
	var rate *easypost.Rate
	for _, r := range shipment.Rates {
		if !strings.EqualFold(r.Carrier, "USPS") &&
			!strings.EqualFold(r.Carrier, "ups") {
			continue
		}
		if !strings.EqualFold(r.Service, "priorityMAILInternational") {
			continue
		}
		if rate == nil {
			rate = r
		} else {
			x, _ := strconv.ParseFloat(r.Rate, 64)
			y, _ := strconv.ParseFloat(rate.Rate, 64)
			if x < y {
				rate = r
			}
		}
	}
	require.NotNil(rate)
	shipment, err = client.BuyShipment(shipment.ID, rate, "100.00")
	require.NoError(err)

	assert.NotEmpty(shipment.TrackingCode)
	assert.Equal("100.00", shipment.Insurance)
	if assert.NotNil(shipment.PostageLabel) {
		assert.Contains(
			shipment.PostageLabel.LabelURL,
			"https://easypost-files.s3-us-west-2.amazonaws.com",
		)
	}
}

func (c *ClientTests) TestShipmentList() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	out, err := client.ListShipments(&easypost.ListShipmentsOptions{
		Purchased:       easypost.BoolPtr(false),
		IncludeChildren: easypost.BoolPtr(false),
	})
	require.NoError(err)
	require.NotEmpty(out)
	assert.True(len(out.Shipments) > 0)
}

func (c *ClientTests) TestShipmentSmartrates() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	to, err := client.CreateAddress(
		&easypost.Address{
			Name:    "Elmer Fudd",
			Street1: "179 N Harbor Dr",
			City:    "Redondo Beach",
			State:   "CA",
			Zip:     "90277",
			Country: "US",
			Phone:   "613-555-1212",
		},
		nil,
	)
	require.NoError(err)

	from, err := client.CreateAddress(
		&easypost.Address{
			Company: "EasyPost",
			Street1: "One Montgomery St",
			Street2: "Ste 400",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94104",
			Phone:   "415-456-7890",
		},
		nil,
	)
	require.NoError(err)

	parcel, err := client.CreateParcel(
		&easypost.Parcel{
			Length: 10.2,
			Width:  7.8,
			Height: 4.3,
			Weight: 21.2,
		},
	)
	require.NoError(err)

	shipment, err := client.CreateShipment(
		&easypost.Shipment{
			ToAddress:   to,
			FromAddress: from,
			Parcel:      parcel,
		},
	)
	require.NoError(err)
	require.NotEmpty(shipment.Rates)

	smartrates, err := client.GetShipmentSmartrates(shipment.ID)
	require.NoError(err)
	assert.Equal(shipment.Rates[0].ID, smartrates[0].ID)
	assert.Equal(smartrates[0].TimeInTransit.Percentile50, 1)
	assert.Equal(smartrates[0].TimeInTransit.Percentile75, 2)
	assert.Equal(smartrates[0].TimeInTransit.Percentile85, 2)
	assert.Equal(smartrates[0].TimeInTransit.Percentile90, 3)
	assert.Equal(smartrates[0].TimeInTransit.Percentile95, 3)
	assert.Equal(smartrates[0].TimeInTransit.Percentile97, 4)
	assert.Equal(smartrates[0].TimeInTransit.Percentile99, 5)
}

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
