package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestScanFormCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	scanform, err := client.CreateScanForm(shipment.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.ScanForm{}), reflect.TypeOf(scanform))
	assert.True(strings.HasPrefix(scanform.ID, "sf_"))
}

func (c *ClientTests) TestScanFormRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	scanform, err := client.CreateScanForm(shipment.ID)
	require.NoError(err)

	retrievedScanform, err := client.GetScanForm(scanform.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.ScanForm{}), reflect.TypeOf(retrievedScanform))
	assert.Equal(scanform, retrievedScanform)
}

func (c *ClientTests) TestScanFormAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	scanforms, err := client.ListScanForms(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	scanformsList := scanforms.ScanForms

	assert.LessOrEqual(len(scanformsList), c.fixture.pageSize())
	assert.NotNil(scanforms.HasMore)
	for _, scanform := range scanformsList {
		assert.Equal(reflect.TypeOf(&easypost.ScanForm{}), reflect.TypeOf(scanform))
	}
}
