package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestRateRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, err := client.CreateShipment(c.fixture.BasicShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	rate, err := client.GetRate(shipment.Rates[0].ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Rate{}), reflect.TypeOf(rate))
	assert.True(strings.HasPrefix(rate.ID, "rate_"))
}
