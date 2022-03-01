package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
	"reflect"
	"strings"
)

func (c *ClientTests) TestRateRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, _ := client.CreateShipment(c.fixture.BasicShipment())

	rate, _ := client.GetRate(shipment.Rates[0].ID)

	assert.Equal(reflect.TypeOf(&easypost.Rate{}), reflect.TypeOf(rate))
	assert.True(strings.HasPrefix(rate.ID, "rate_"))
}
