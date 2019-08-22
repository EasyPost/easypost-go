package easypost_test

import (
	"testing"

	"github.com/EasyPost/easypost-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderCreateThenBuy(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	// We create an Order containing Shipment. Towards the end we assert on
	// Order and Parcel's values.
	to := &easypost.Address{
		Company: "Oakland Dmv Office",
		Name:    "Customer",
		Street1: "5300 Claremont Ave",
		City:    "Oakland",
		State:   "CA",
		Zip:     "94618",
		Country: "US",
		Phone:   "510-555-1212",
	}
	parcel := &easypost.Parcel{
		Weight: 21.2,
		Length: 12,
		Width:  12,
		Height: 3,
	}
	order, err := TestClient.CreateOrder(
		&easypost.Order{
			ToAddress: to,
			FromAddress: &easypost.Address{
				Company: "EasyPost",
				Street1: "One Montgomery St",
				Street2: "Ste 400",
				City:    "San Francisco",
				State:   "CA",
				Zip:     "94104",
				Phone:   "415-555-1212",
			},
			Shipments: []*easypost.Shipment{
				&easypost.Shipment{
					Parcel: parcel,
					Options: &easypost.ShipmentOptions{
						LabelFormat: "PDF",
					},
				},
				&easypost.Shipment{
					Parcel: &easypost.Parcel{
						Weight: 16,
						Length: 8,
						Width:  5,
						Height: 5,
					},
					Options: &easypost.ShipmentOptions{
						LabelFormat: "PDF",
					},
				},
			},
		},
	)
	require.NoError(err)
	order, err = TestClient.GetOrderRates(order.ID)
	require.NoError(err)

	assert.Len(order.Shipments, 2)

	var longBox, squareBox *easypost.Parcel
	for i := range order.Shipments {
		if parcel := order.Shipments[i].Parcel; assert.NotNil(parcel) {
			switch {
			case parcel.Height == 3:
				longBox = parcel
			case parcel.Width == 5:
				squareBox = parcel
			}
		}
	}
	if assert.NotNil(longBox) {
		assert.Equal(parcel.Height, longBox.Height)
		assert.Equal(parcel.Length, longBox.Length)
		assert.Equal(parcel.Width, longBox.Width)
		assert.Equal(parcel.Weight, longBox.Weight)
	}
	if assert.NotNil(squareBox) {
		assert.Equal(float64(5), squareBox.Height)
		assert.Equal(float64(8), squareBox.Length)
		assert.Equal(float64(5), squareBox.Width)
		assert.Equal(float64(16), squareBox.Weight)
	}

	order, err = TestClient.BuyOrder(order.ID, "USPS", "Priority")
	require.NoError(err)

	for i := range order.Shipments {
		shipment, err := TestClient.InsureShipment(
			order.Shipments[i].ID, "100.00",
		)
		if assert.NoError(err) {
			assert.NotEmpty(shipment.TrackingCode)
			assert.Equal("100.00", shipment.Insurance)
			if label := shipment.PostageLabel; assert.NotNil(label) {
				assert.Contains(
					label.LabelURL,
					"https://easypost-files.s3-us-west-2.amazonaws.com",
				)
			}
		}
	}
}
