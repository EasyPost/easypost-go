package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestInsuranceCreate() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, _ := client.CreateShipment(c.fixture.OneCallBuyShipment())

	insurance, _ := client.CreateInsurance(
		&easypost.Insurance{
			ToAddress:    c.fixture.BasicAddress(),
			FromAddress:  c.fixture.BasicAddress(),
			TrackingCode: shipment.TrackingCode,
			Carrier:      c.fixture.USPS(),
			Amount:       "100",
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.Insurance{}), reflect.TypeOf(insurance))
	assert.True(strings.HasPrefix(insurance.ID, "ins_"))
	assert.Equal("100.00000", insurance.Amount)
}

func (c *ClientTests) TestInsuranceRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, _ := client.CreateShipment(c.fixture.OneCallBuyShipment())

	insurance, _ := client.CreateInsurance(
		&easypost.Insurance{
			ToAddress:    c.fixture.BasicAddress(),
			FromAddress:  c.fixture.BasicAddress(),
			TrackingCode: shipment.TrackingCode,
			Carrier:      c.fixture.USPS(),
			Amount:       "100",
		},
	)
	retrievedInsurance, _ := client.GetInsurance(insurance.ID)

	assert.Equal(reflect.TypeOf(&easypost.Insurance{}), reflect.TypeOf(retrievedInsurance))
	assert.Equal(insurance, retrievedInsurance)
}

func (c *ClientTests) TestInsuranceAll() {
	client := c.TestClient()
	assert := c.Assert()

	insurances, _ := client.ListInsurances(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)

	insurancesList := insurances.Insurances

	assert.LessOrEqual(len(insurancesList), c.fixture.pageSize())
	assert.NotNil(insurances.HasMore)
	for _, insurance := range insurancesList {
		assert.Equal(reflect.TypeOf(&easypost.Insurance{}), reflect.TypeOf(insurance))
	}
}
