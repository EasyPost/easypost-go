package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestRefundCreate() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	retrievedShipment, err := client.GetShipment(shipment.ID) // We need to retrieve the shipment so that the tracking_code has time to populate
	if err != nil {
		c.T().Error(err)
		return
	}

	refund, err := client.CreateRefund(
		map[string]interface{}{
			"carrier":        "USPS",
			"tracking_codes": retrievedShipment.TrackingCode,
		},
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.True(strings.HasPrefix(refund[0].ID, "rfnd_"))
	assert.Equal("submitted", refund[0].Status)
}

func (c *ClientTests) TestRefundAll() {
	client := c.TestClient()
	assert := c.Assert()

	refunds, err := client.ListRefunds(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	if err != nil {
		c.T().Error(err)
		return
	}

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

	refunds, err := client.ListRefunds(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	retrievedRefund, err := client.GetRefund(refunds.Refunds[0].ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Refund{}), reflect.TypeOf(retrievedRefund))
	assert.Equal(refunds.Refunds[0].ID, retrievedRefund.ID)
}
