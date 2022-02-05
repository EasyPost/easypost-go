package easypost_test

import (
	"github.com/EasyPost/easypost-go"
)

func (c *ClientTests) TestLowestRate() {
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
	rate, err := client.LowestRate(shipment)

	require.NoError(err)
	assert.Equal("18.63", rate.Rate)
	assert.Equal("UPSDAP", rate.Carrier)
	assert.Equal("UPSStandard", rate.Service)

	carriers := []string{"DHLExpress", "USPS"}
	rateWithCarriersFilter, err := client.LowestRateWithCarrier(shipment, carriers)

	require.NoError(err)
	assert.Equal("20.38", rateWithCarriersFilter.Rate)
	assert.Equal("DHLExpress", rateWithCarriersFilter.Carrier)
	assert.Equal("ExpressWorldwide", rateWithCarriersFilter.Service)

	services := []string{"ExpressWorldwideNonDoc", "FirstClassPackageInternationalService"}
	rateWithCarriersAndServicesFilter, err := client.LowestRateWithCarrierAndService(shipment, carriers, services)
	require.NoError(err)
	assert.Equal("21.00", rateWithCarriersAndServicesFilter.Rate)
	assert.Equal("USPS", rateWithCarriersAndServicesFilter.Carrier)
	assert.Equal("FirstClassPackageInternationalService", rateWithCarriersAndServicesFilter.Service)
}
