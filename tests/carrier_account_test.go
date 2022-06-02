package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) GenerateCarrierAccount() (*easypost.CarrierAccount, error) {
	client := c.ProdClient()

	carrierAccount, err := client.CreateCarrierAccount(c.fixture.BasicCarrierAccount())

	return carrierAccount, err
}

func (c *ClientTests) TestCarrierAccountCreate() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	carrierAccount, err := c.GenerateCarrierAccount()
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(carrierAccount))
	assert.True(strings.HasPrefix(carrierAccount.ID, "ca_"))
	assert.Equal("UpsAccount", carrierAccount.Type)

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountRetrieve() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	carrierAccount, err := c.GenerateCarrierAccount()
	require.NoError(err)

	retrievedCarrierAccount, err := client.GetCarrierAccount(carrierAccount.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(retrievedCarrierAccount))
	assert.Equal(carrierAccount, retrievedCarrierAccount)

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountAll() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	carrierAccounts, err := client.ListCarrierAccounts()
	require.NoError(err)

	for _, carrierAccount := range carrierAccounts {
		assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(carrierAccount))
	}
}

func (c *ClientTests) TestCarrierAccountUpdate() {
	c.T().Skip("VCR can't match this cassette properly")
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	testDescription := "My custom description"

	carrierAccount, err := c.GenerateCarrierAccount()
	require.NoError(err)

	carrierAccount.Description = testDescription

	updatedCarrierAccount, err := client.UpdateCarrierAccount(carrierAccount)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(updatedCarrierAccount))
	assert.True(strings.HasPrefix(updatedCarrierAccount.ID, "ca_"))
	assert.Equal(testDescription, updatedCarrierAccount.Description)

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountDelete() {
	client := c.ProdClient()
	require := c.Require()

	carrierAccount, err := c.GenerateCarrierAccount()
	require.NoError(err)

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountTypes() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	carriersAccountTypes, err := client.GetCarrierTypes()
	require.NoError(err)

	assert.Equal(reflect.TypeOf([]*easypost.CarrierType{}), reflect.TypeOf(carriersAccountTypes))
}
