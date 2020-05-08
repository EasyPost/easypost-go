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
