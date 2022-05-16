package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestScanFormCreate() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	scanform, err := client.CreateScanForm(shipment.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.ScanForm{}), reflect.TypeOf(scanform))
	assert.True(strings.HasPrefix(scanform.ID, "sf_"))
}

func (c *ClientTests) TestScanFormRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	if err != nil {
		c.T().Error(err)
		return
	}

	scanform, err := client.CreateScanForm(shipment.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	retrievedScanform, err := client.GetScanForm(scanform.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.ScanForm{}), reflect.TypeOf(retrievedScanform))
	assert.Equal(scanform, retrievedScanform)
}

func (c *ClientTests) TestScanFormAll() {
	client := c.TestClient()
	assert := c.Assert()

	scanforms, err := client.ListScanForms(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	scanformsList := scanforms.ScanForms

	assert.LessOrEqual(len(scanformsList), c.fixture.pageSize())
	assert.NotNil(scanforms.HasMore)
	for _, scanform := range scanformsList {
		assert.Equal(reflect.TypeOf(&easypost.ScanForm{}), reflect.TypeOf(scanform))
	}
}
