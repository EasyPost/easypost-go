package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
	"reflect"
	"strings"
)

func (c *ClientTests) TestCreateRefund() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, _ := client.CreateShipment(OneCallBuyShipment())

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

func (c *ClientTests) TestAllRefund() {
	client := c.TestClient()
	assert := c.Assert()

	refunds, _ := client.ListRefunds(
		&easypost.ListOptions{
			PageSize: 5,
		},
	)

	refundsList := refunds.Refunds

	assert.LessOrEqual(len(refundsList), 5)
	assert.NotNil(refunds.HasMore)
	for _, refund := range refundsList {
		assert.Equal(reflect.TypeOf(&easypost.Refund{}), reflect.TypeOf(refund))
	}
}

func (c *ClientTests) TestRetrieveRefund() {
	client := c.TestClient()
	assert := c.Assert()

	refunds, _ := client.ListRefunds(
		&easypost.ListOptions{
			PageSize: 5,
		},
	)

	retrievedRefund, _ := client.GetRefund(refunds.Refunds[0].ID)

	assert.Equal(reflect.TypeOf(&easypost.Refund{}), reflect.TypeOf(retrievedRefund))
	assert.Equal(refunds.Refunds[0].ID, retrievedRefund.ID)
}

func BasicAddress() *easypost.Address {
	return &easypost.Address{
		Name:    "Jack Sparrow",
		Company: "EasyPost",
		Street1: "388 Townsend St",
		Street2: "Apt 20",
		City:    "San Francisco",
		State:   "CA",
		Zip:     "94107",
		Phone:   "5555555555",
	}
}

func BasicParcel() *easypost.Parcel {
	return &easypost.Parcel{
		Length: 10,
		Width:  8,
		Height: 4,
		Weight: 15.4,
	}
}

func OneCallBuyShipment() *easypost.Shipment {
	return &easypost.Shipment{
		ToAddress:         BasicAddress(),
		FromAddress:       BasicAddress(),
		Parcel:            BasicParcel(),
		Service:           "First",
		CarrierAccountIDs: []string{"ca_e606176ddb314afe896733636dba2f3b"},
		Carrier:           "USPS",
	}
}
