package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
	"reflect"
	"strings"
)

func (c *ClientTests) TestRefundCreate() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, _ := client.CreateShipment(c.fixture.OneCallBuyShipment())

	retrievedShipment, _ := client.GetShipment(shipment.ID) // We need to retrieve the shipment so that the tracking_code has time to populate

	refund, _ := client.CreateRefund(
		map[string]interface{}{
			"carrier":        "USPS",
			"tracking_codes": retrievedShipment.TrackingCode,
		},
	)

	assert.True(strings.HasPrefix(refund[0].ID, "rfnd_"))
	assert.Equal("submitted", refund[0].Status)
}

func (c *ClientTests) TestRefundAll() {
	client := c.TestClient()
	assert := c.Assert()

	refunds, _ := client.ListRefunds(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)

	refundsList := refunds.Refunds

	assert.LessOrEqual(len(refundsList), c.fixture.pageSize())
	assert.NotNil(refunds.HasMore)
	for _, refund := range refundsList {
		assert.Equal(reflect.TypeOf(&easypost.Refund{}), reflect.TypeOf(refund))
	}
}

func (c *ClientTests) TestRefundRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	refunds, _ := client.ListRefunds(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)

	retrievedRefund, _ := client.GetRefund(refunds.Refunds[0].ID)

	assert.Equal(reflect.TypeOf(&easypost.Refund{}), reflect.TypeOf(retrievedRefund))
	assert.Equal(refunds.Refunds[0].ID, retrievedRefund.ID)
}
