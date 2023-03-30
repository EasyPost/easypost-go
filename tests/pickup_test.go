package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestPickupCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	pickupData := c.fixture.BasicPickup()
	pickupData.Shipment = shipment

	pickup, err := client.CreatePickup(pickupData)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Pickup{}), reflect.TypeOf(pickup))
	assert.True(strings.HasPrefix(pickup.ID, "pickup_"))
	assert.NotNil(pickup.PickupRates)
}

func (c *ClientTests) TestPickupRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	pickupData := c.fixture.BasicPickup()
	pickupData.Shipment = shipment

	pickup, err := client.CreatePickup(pickupData)
	require.NoError(err)

	retrievePickup, err := client.GetPickup(pickup.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Pickup{}), reflect.TypeOf(retrievePickup))
	assert.Equal(pickup.ID, retrievePickup.ID)
}

func (c *ClientTests) TestPickupBuy() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	pickupData := c.fixture.BasicPickup()
	pickupData.Shipment = shipment

	pickup, err := client.CreatePickup(pickupData)
	require.NoError(err)

	boughtPickup, err := client.BuyPickup(
		pickup.ID,
		&easypost.PickupRate{
			Carrier: c.fixture.USPS(),
			Service: c.fixture.PickupService(),
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Pickup{}), reflect.TypeOf(boughtPickup))
	assert.True(strings.HasPrefix(boughtPickup.ID, "pickup_"))
	assert.NotNil(pickup.Confirmation)
	assert.Equal("scheduled", boughtPickup.Status)
}

func (c *ClientTests) TestPickupCancel() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	pickupData := c.fixture.BasicPickup()
	pickupData.Shipment = shipment

	pickup, err := client.CreatePickup(pickupData)
	require.NoError(err)

	boughtPickup, err := client.BuyPickup(
		pickup.ID,
		&easypost.PickupRate{
			Carrier: c.fixture.USPS(),
			Service: c.fixture.PickupService(),
		},
	)
	require.NoError(err)

	cancelledPickup, err := client.CancelPickup(boughtPickup.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Pickup{}), reflect.TypeOf(cancelledPickup))
	assert.True(strings.HasPrefix(cancelledPickup.ID, "pickup_"))
	assert.Equal("canceled", cancelledPickup.Status)
}

func (c *ClientTests) TestPickupLowestRate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	pickupData := c.fixture.BasicPickup()
	pickupData.Shipment = shipment

	pickup, err := client.CreatePickup(pickupData)
	require.NoError(err)

	// Test lowest rate with no filters
	lowestRate, err := client.LowestPickupRate(pickup)
	require.NoError(err)
	assert.Equal("NextDay", lowestRate.Service)
	assert.Equal("0.00", lowestRate.Rate)
	assert.Equal("USPS", lowestRate.Carrier)

	// Test lowest rate with service filter (should error due to bad service)
	lowestRate, err = client.LowestPickupRateWithCarrierAndService(pickup, nil, []string{"BAD_SERVICE"})
	assert.Error(err)

	// Test lowest rate with carrier filter (should error due to bad carrier)
	lowestRate, err = client.LowestPickupRateWithCarrier(pickup, []string{"BAD_CARRIER"})
	assert.Error(err)
}

func (c *ClientTests) TestPickupAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	pickups, err := client.ListPickups(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	assert.LessOrEqual(len(pickups.Pickups), c.fixture.pageSize())
	assert.NotNil(pickups.HasMore)
	for _, pickup := range pickups.Pickups {
		assert.Equal(reflect.TypeOf(&easypost.Pickup{}), reflect.TypeOf(pickup))
	}
}

func (c *ClientTests) TestPickupsGetNextPage() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListPickups(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextPickupPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Pickups) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Pickups[len(firstPage.Pickups)-1].ID
			firstIdOfSecondPage := nextPage.Pickups[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIdOfSecondPage)
		}
	}()
	if err != nil {
		assert.Equal(err.Error(), easypost.EndOfPaginationError.Error())
		return
	}
}
