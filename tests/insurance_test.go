package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestInsuranceCreate() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	insuranceData := c.fixture.BasicInsurance()
	insuranceData.TrackingCode = shipment.TrackingCode

	insurance, err := client.CreateInsurance(insuranceData)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Insurance{}), reflect.TypeOf(insurance))
	assert.True(strings.HasPrefix(insurance.ID, "ins_"))
	assert.Equal("100.00000", insurance.Amount)
}

func (c *ClientTests) TestInsuranceRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	insuranceData := c.fixture.BasicInsurance()
	insuranceData.TrackingCode = shipment.TrackingCode

	insurance, err := client.CreateInsurance(insuranceData)
	if err != nil {
		c.T().Error(err)
		return
	}

	retrievedInsurance, err := client.GetInsurance(insurance.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Insurance{}), reflect.TypeOf(retrievedInsurance))
	assert.Equal(insurance, retrievedInsurance)
}

func (c *ClientTests) TestInsuranceAll() {
	client := c.TestClient()
	assert := c.Assert()

	insurances, err := client.ListInsurances(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	insurancesList := insurances.Insurances

	assert.LessOrEqual(len(insurancesList), c.fixture.pageSize())
	assert.NotNil(insurances.HasMore)
	for _, insurance := range insurancesList {
		assert.Equal(reflect.TypeOf(&easypost.Insurance{}), reflect.TypeOf(insurance))
	}
}
