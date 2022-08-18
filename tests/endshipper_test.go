package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestEndShipperCreate() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	endshipper, err := client.CreateEndShipper(c.fixture.CaAddress1())
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(endshipper))
	assert.True(strings.HasPrefix(endshipper.ID, "es"))
	assert.Equal("388 TOWNSEND ST APT 20", endshipper.Street1)
}

func (c *ClientTests) TestEndShipperRetrieve() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	endshipper, err := client.CreateEndShipper(c.fixture.CaAddress1())
	require.NoError(err)

	retrievedEndShipper, err := client.GetEndShipper(endshipper.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(retrievedEndShipper))
	assert.True(strings.HasPrefix(retrievedEndShipper.ID, "es"))
	assert.Equal(endshipper.Street1, retrievedEndShipper.Street1)
}

func (c *ClientTests) TestEndShipperAll() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	endshippers, err := client.ListEndShippers(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	assert.LessOrEqual(len(endshippers.EndShippers), c.fixture.pageSize())
	assert.NotNil(endshippers.HasMore)
	for _, endshipper := range endshippers.EndShippers {
		assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(endshipper))
	}
}

func (c *ClientTests) TestEndShipperUpdate() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	endshipper, err := client.CreateEndShipper(c.fixture.CaAddress1())
	require.NoError(err)

	endshipper.Name = "Captain Sparrow"
	endshipper.Company = "EasyPost"
	endshipper.Street1 = "388 Townsend St"
	endshipper.Street2 = "Apt 20"
	endshipper.City = "San Francisco"
	endshipper.State = "CA"
	endshipper.Zip = "94107"
	endshipper.Country = "US"
	endshipper.Phone = "9999999999"
	endshipper.Email = "test@example.com"

	updatedEndShipper, err := client.UpdateEndShippers(endshipper)

	require.NoError(err)

	assert.Equal("CAPTAIN SPARROW", updatedEndShipper.Name)
}
