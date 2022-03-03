package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
	"reflect"
	"strings"
)

func (c *ClientTests) TestScanFormCreate() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, _ := client.CreateShipment(c.fixture.OneCallBuyShipment())

	scanform, _ := client.CreateScanForm(shipment.ID)

	assert.Equal(reflect.TypeOf(&easypost.ScanForm{}), reflect.TypeOf(scanform))
	assert.True(strings.HasPrefix(scanform.ID, "sf_"))
}

func (c *ClientTests) TestScanFormRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	shipment, _ := client.CreateShipment(c.fixture.OneCallBuyShipment())

	scanform, _ := client.CreateScanForm(shipment.ID)

	retrievedScanform, _ := client.GetScanForm(scanform.ID)

	assert.Equal(reflect.TypeOf(&easypost.ScanForm{}), reflect.TypeOf(retrievedScanform))
	assert.Equal(scanform, retrievedScanform)
}

func (c *ClientTests) TestScanFormAll() {
	client := c.TestClient()
	assert := c.Assert()

	scanforms, _ := client.ListScanForms(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)

	scanformsList := scanforms.ScanForms

	assert.LessOrEqual(len(scanformsList), c.fixture.pageSize())
	assert.NotNil(scanforms.HasMore)
	for _, scanform := range scanformsList {
		assert.Equal(reflect.TypeOf(&easypost.ScanForm{}), reflect.TypeOf(scanform))
	}
}
