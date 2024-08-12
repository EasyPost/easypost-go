package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v4"
)

func (c *ClientTests) TestRefundCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	retrievedShipment, err := client.GetShipment(shipment.ID) // We need to retrieve the shipment so that the tracking_code has time to populate
	require.NoError(err)

	refund, err := client.CreateRefund(
		map[string]interface{}{
			"carrier":        "USPS",
			"tracking_codes": [1]string{retrievedShipment.TrackingCode},
		},
	)
	require.NoError(err)

	assert.True(strings.HasPrefix(refund[0].ID, "rfnd_"))
	assert.Equal("submitted", refund[0].Status)
}

func (c *ClientTests) TestRefundAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	refunds, err := client.ListRefunds(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	refundsList := refunds.Refunds

	assert.LessOrEqual(len(refundsList), c.fixture.pageSize())
	assert.NotNil(refunds.HasMore)
	for _, refund := range refundsList {
		assert.Equal(reflect.TypeOf(&easypost.Refund{}), reflect.TypeOf(refund))
	}
}

func (c *ClientTests) TestRefundRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	refunds, err := client.ListRefunds(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	retrievedRefund, err := client.GetRefund(refunds.Refunds[0].ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Refund{}), reflect.TypeOf(retrievedRefund))
	assert.Equal(refunds.Refunds[0].ID, retrievedRefund.ID)
}

func (c *ClientTests) TestRefundsGetNextPage() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListRefunds(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextRefundPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Refunds) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Refunds[len(firstPage.Refunds)-1].ID
			firstIdOfSecondPage := nextPage.Refunds[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIdOfSecondPage)
		}
	}()
	if err != nil {
		assert.Equal(err.Error(), easypost.NoPagesLeftToRetrieve)
		return
	}
}
