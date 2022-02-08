package easypost_test

import (
	"strconv"

	"github.com/EasyPost/easypost-go"
)

func (c *ClientTests) TestScanFormCreateAndRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()
	shipment, err := client.CreateShipment(
		&easypost.Shipment{
			ToAddress: &easypost.Address{
				Name:    "Elmer Fudd",
				Street1: "55 Sparks St.",
				City:    "Ottawa",
				State:   "ON",
				Zip:     "k1p5a5",
				Country: "CA",
				Phone:   "613-555-1212",
			},
			FromAddress: &easypost.Address{
				Company: "EasyPost",
				Street1: "One Montgomery St",
				Street2: "Ste 400",
				City:    "San Francisco",
				State:   "CA",
				Zip:     "94104",
				Phone:   "415-456-7890",
			},
			Parcel: &easypost.Parcel{
				Length: 10.2,
				Width:  7.8,
				Height: 4.3,
				Weight: 21.2,
			},
			CustomsInfo: &easypost.CustomsInfo{
				CustomsCertify:    true,
				CustomsSigner:     "Wile E. Coyote",
				ContentsType:      "gift",
				EELPFC:            "NOEEI 30.37(a)",
				NonDeliveryOption: "return",
				RestrictionType:   "none",
				CustomsItems: []*easypost.CustomsItem{
					{
						Description:    "EasyPost t-shirts",
						HSTariffNumber: "123456",
						OriginCountry:  "US",
						Quantity:       2,
						Value:          96.27,
						Weight:         21.1,
					},
				},
			},
		},
	)
	require.NoError(err)
	require.NotEmpty(shipment.Rates)
	var rate *easypost.Rate
	for i := range shipment.Rates {
		if rate == nil {
			rate = shipment.Rates[i]
		} else {
			x, _ := strconv.ParseFloat(shipment.Rates[i].Rate, 64)
			y, _ := strconv.ParseFloat(rate.Rate, 64)
			if x < y {
				rate = shipment.Rates[i]
			}
		}
	}
	shipment, err = client.BuyShipment(shipment.ID, rate, "100.00")
	require.NoError(err)

	scanForm, err := client.CreateScanForm(shipment.ID)
	require.NoError(err)
	assert.NotEmpty(scanForm.ID)
	if assert.NotEmpty(scanForm.TrackingCodes) {
		assert.Equal(shipment.TrackingCode, scanForm.TrackingCodes[0])
	}

	scanForm2, err := client.GetScanForm(scanForm.ID)
	require.NoError(err)
	assert.Equal(scanForm.ID, scanForm2.ID)

	scanForms, err := client.ListScanForms(
		&easypost.ListOptions{PageSize: 2},
	)
	require.NoError(err)
	assert.NotEmpty(scanForms.ScanForms)
}
