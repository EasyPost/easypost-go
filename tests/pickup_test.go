package easypost_test

import (
	"time"

	"github.com/EasyPost/easypost-go"
)

func noonOnNextMonday() time.Time {
	now := time.Now()
	date := time.Date(
		now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location(),
	)
	wd := date.Weekday()
	if wd < time.Monday {
		return date.AddDate(0, 0, int(time.Monday-wd))
	}
	return date.AddDate(0, 0, 7-int(wd-time.Monday))
}

func (c *ClientTests) TestPickupBatch() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()
	// Create a Batch containing multiple Shipments. Then we try to buy a
	// Pickup and assert if it was bought.
	from, err := client.CreateAddress(
		&easypost.Address{
			Name:    "Homer Simpson",
			Company: "EasyPost",
			Street1: "1 MONTGOMERY ST STE 400",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94104",
			Phone:   "415-555-1212",
		},
		&easypost.CreateAddressOptions{Verify: []string{"delivery"}},
	)
	require.NoError(err)

	batch, err := client.CreateAndBuyBatch(
		&easypost.Shipment{
			ToAddress: &easypost.Address{
				Name:    "Bugs Bunny",
				Street1: "4000 Warner Blvd",
				City:    "Burbank",
				State:   "CA",
				Zip:     "91522",
				Phone:   "818-555-1212",
			},
			FromAddress: from,
			Parcel:      &easypost.Parcel{Weight: 10.2},
			Carrier:     "USPS",
			Service:     "Priority",
		},
		&easypost.Shipment{
			ToAddress: &easypost.Address{
				Name:    "Bugs Bunny",
				Street1: "4000 Warner Blvd",
				City:    "Burbank",
				State:   "CA",
				Zip:     "91522",
				Phone:   "818-555-1212",
			},
			FromAddress: &easypost.Address{
				Name:            "Yosemite Sam",
				Company:         "EasyPost",
				Street1:         "1 Market St",
				CarrierFacility: "San Francisco",
				State:           "CA",
				Zip:             "94105",
				Phone:           "415-555-1212",
			},
			Parcel:  &easypost.Parcel{Weight: 10.2},
			Carrier: "USPS",
			Service: "Priority",
		},
	)
	require.NoError(err)

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

	// Insure the shipments after purchase
	for i := range batch.Shipments {
		_, err := client.InsureShipment(batch.Shipments[i].ID, "100.00")
		assert.NoError(err)
	}

	minDatetime := noonOnNextMonday()
	maxDatetime := minDatetime.AddDate(0, 0, 1)

	pickup, err := client.CreatePickup(
		&easypost.Pickup{
			Address:          from,
			Batch:            batch,
			Reference:        "internal_id_1234",
			MinDatetime:      &minDatetime,
			MaxDatetime:      &maxDatetime,
			IsAccountAddress: true,
			Instructions:     "Special pickup instructions",
		},
	)
	require.NoError(err)
	require.NotEmpty(pickup.PickupRates)

	_, err = client.BuyPickup(pickup.ID, pickup.PickupRates[0])
	require.NoError(err)
}

func (c *ClientTests) TestSinglePickup() {
	client := c.TestClient()
	require := c.Require()

	from, err := client.CreateAddress(
		&easypost.Address{
			Name:    "Homer Simpson",
			Company: "EasyPost",
			Street1: "1 MONTGOMERY ST STE 400",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94104",
			Phone:   "415-456-7890",
		},
		&easypost.CreateAddressOptions{Verify: []string{"delivery"}},
	)
	require.NoError(err)

	shipment, err := client.CreateShipment(
		&easypost.Shipment{
			ToAddress: &easypost.Address{
				Name:    "Bugs Bunny",
				Street1: "4000 Warner Blvd",
				City:    "Burbank",
				State:   "CA",
				Zip:     "91522",
				Phone:   "818-555-1212",
			},
			FromAddress: from,
			Parcel:      &easypost.Parcel{Weight: 21.2},
		},
	)
	require.NoError(err)
	require.NotEmpty(shipment.Rates)

	shipment, err = client.BuyShipment(
		shipment.ID, shipment.Rates[0], "100.00",
	)
	require.NoError(err)

	minDatetime := noonOnNextMonday()
	maxDatetime := minDatetime.AddDate(0, 0, 1)

	pickup, err := client.CreatePickup(
		&easypost.Pickup{
			Address:          from,
			Shipment:         shipment,
			Reference:        "internal_id_1234",
			MinDatetime:      &minDatetime,
			MaxDatetime:      &maxDatetime,
			IsAccountAddress: true,
			Instructions:     "Special pickup instructions",
		},
	)
	require.NoError(err)
	require.NotEmpty(pickup.PickupRates)
	// This is probably a bug in the API. It always returns a rate, but isn't
	// always a valid one.
	_, err = client.BuyPickup(pickup.ID, pickup.PickupRates[0])
	if err != nil {
		require.Contains(
			err.Error(), "schedule and change requests must contain",
		)
	}
}
