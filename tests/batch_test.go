package easypost_test

import (
	"time"

	"github.com/EasyPost/easypost-go"
)

func (c *ClientTests) TestBatchCreateAndBuy() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()
	// We create Address and Parcel objects. We then try to create a Batch
	// containing a shipment. Finally, we assert on saved and returned data.
	from, err := client.CreateAddress(
		&easypost.Address{
			Company: "EasyPost",
			Street1: "One Montgomery St",
			Street2: "Ste 400",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94104",
			Phone:   "415-555-1212",
		},
		nil,
	)
	require.NoError(err)

	parcel, err := client.CreateParcel(
		&easypost.Parcel{PredefinedPackage: "RegionalRateBoxA", Weight: 64},
	)
	require.NoError(err)

	shipments := []*easypost.Shipment{
		{
			ToAddress: &easypost.Address{
				Name:    "Bugs Bunny",
				Street1: "4000 Warner Blvd",
				City:    "Burbank",
				State:   "CA",
				Zip:     "91522",
				Phone:   "818-555-1212",
			},
			FromAddress: from,
			Parcel:      parcel,
			Reference:   "1234567890",
			Carrier:     "USPS",
			Service:     "Priority",
		},
	}

	// create batch of shipments.
	batch, err := client.CreateAndBuyBatch(shipments...)
	require.NoError(err)
	assert.Equal(1, batch.NumShipments)

	// Poll while waiting for the batch to purchase the shipments.
	for done := false; !done; {
		time.Sleep(time.Second)
		switch batch.State {
		case "creating", "queued_for_purchase", "purchasing":
			batch, err = client.GetBatch(batch.ID)
			require.NoError(err)
		default:
			done = true
		}
	}

	// Insure the shipments after purchase.
	for i := range batch.Shipments {
		_, err := client.InsureShipment(batch.Shipments[i].ID, "100.00")
		assert.NoError(err)
	}

	require.Len(batch.Shipments, 1)
	assert.Equal("postage_purchased", batch.Shipments[0].BatchStatus)

	shipment, err := client.GetShipment(batch.Shipments[0].ID)
	require.NoError(err)

	if addr := shipment.BuyerAddress; assert.NotNil(addr) {
		assert.Equal("Burbank", addr.City)
		assert.Equal("US", addr.Country)
		assert.Equal("Bugs Bunny", addr.Name)
		assert.Equal("8185551212", addr.Phone)
		assert.Equal("4000 Warner Blvd", addr.Street1)
	}

	// Assert on fees.
	if fees := shipment.Fees; assert.GreaterOrEqual(len(fees), 3) {
		assert.Equal("0.00000", fees[0].Amount)
		assert.Equal("7.92000", fees[1].Amount)
		assert.Equal("1.00000", fees[2].Amount)
	}

	// Assert on parcel.
	if parcel := shipment.Parcel; assert.NotNil(parcel) {
		assert.Equal("RegionalRateBoxA", parcel.PredefinedPackage)
		assert.Equal(float64(64), parcel.Weight)
	}

	// Assert on tracker.
	if tracker := shipment.Tracker; assert.NotNil(tracker) {
		assert.NotEmpty(tracker.TrackingCode)
		assert.NotEmpty(tracker.ShipmentID)
	}
}
