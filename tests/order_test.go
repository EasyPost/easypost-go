package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestOrderCreate() {
	client := c.TestClient()
	assert := c.Assert()

	order, _ := client.CreateOrder(c.fixture.BasicOrder())

	assert.Equal(reflect.TypeOf(&easypost.Order{}), reflect.TypeOf(order))
	assert.True(strings.HasPrefix(order.ID, "order_"))
	assert.NotNil(order.Rates)
}

func (c *ClientTests) TestOrderRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	order, _ := client.CreateOrder(c.fixture.BasicOrder())

	retrievedOrder, _ := client.GetOrder(order.ID)

	assert.Equal(reflect.TypeOf(&easypost.Order{}), reflect.TypeOf(retrievedOrder))
	assert.Equal(order.ID, retrievedOrder.ID)
}

func (c *ClientTests) TestOrderGetRates() {
	client := c.TestClient()
	assert := c.Assert()

	order, _ := client.CreateOrder(c.fixture.BasicOrder())

	rates, _ := client.GetOrderRates(order.ID)

	ratesList := rates.Rates

	assert.Equal(reflect.TypeOf([]*easypost.Rate{}), reflect.TypeOf(ratesList))
	for _, rate := range ratesList {
		assert.Equal(reflect.TypeOf(&easypost.Rate{}), reflect.TypeOf(rate))
	}
}

func (c *ClientTests) TestOrderBuyRate() {
	client := c.TestClient()
	assert := c.Assert()

	order, _ := client.CreateOrder(c.fixture.BasicOrder())

	boughtOrder, _ := client.BuyOrder(order.ID, c.fixture.USPS(), c.fixture.USPSService())

	shipmentsList := boughtOrder.Shipments

	for _, shipment := range shipmentsList {
		assert.NotNil(shipment.PostageLabel)
	}
}
