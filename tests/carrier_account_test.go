package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) GenerateCarrierAccount() (*easypost.CarrierAccount, error) {
	client := c.ProdClient()

	carrierAccount, err := client.CreateCarrierAccount(
		&easypost.CarrierAccount{
			Type: "UpsAccount",
			Credentials: map[string]string{
				"account_number":        "A1A1A1",
				"user_id":               "USERID",
				"password":              "PASSWORD",
				"access_license_number": "ALN",
			},
		},
	)

	return carrierAccount, err
}

func (c *ClientTests) TestCarrierAccountCreate() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	carrierAccount, _ := c.GenerateCarrierAccount()

	assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(carrierAccount))
	assert.True(strings.HasPrefix(carrierAccount.ID, "ca_"))
	assert.Equal("UpsAccount", carrierAccount.Type)

	err := client.DeleteCarrierAccount(carrierAccount.ID)

	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountRetrieve() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	carrierAccount, _ := c.GenerateCarrierAccount()

	retrievedCarrierAccount, _ := client.GetCarrierAccount(carrierAccount.ID)
	assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(retrievedCarrierAccount))
	assert.Equal(carrierAccount, retrievedCarrierAccount)

	err := client.DeleteCarrierAccount(carrierAccount.ID)

	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountAll() {
	client := c.ProdClient()
	assert := c.Assert()

	carrierAccounts, _ := client.ListCarrierAccounts()

	for _, carrierAccount := range carrierAccounts {
		assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(carrierAccount))
	}
}

func (c *ClientTests) TestCarrierAccountUpdate() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	testDescription := "My custom description"

	carrierAccount, _ := c.GenerateCarrierAccount()

	carrierAccount.Description = testDescription

	updatedCarrierAccount, _ := client.UpdateCarrierAccount(carrierAccount)

	assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(updatedCarrierAccount))
	assert.True(strings.HasPrefix(updatedCarrierAccount.ID, "ca_"))
	assert.Equal(testDescription, updatedCarrierAccount.Description)

	err := client.DeleteCarrierAccount(carrierAccount.ID)

	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountDelete() {
	client := c.ProdClient()
	require := c.Require()

	carrierAccount, _ := c.GenerateCarrierAccount()

	err := client.DeleteCarrierAccount(carrierAccount.ID)

	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountTypes() {
	client := c.ProdClient()
	assert := c.Assert()

	carriersAccountTypes, _ := client.GetCarrierTypes()

	assert.Equal(reflect.TypeOf([]*easypost.CarrierType{}), reflect.TypeOf(carriersAccountTypes))
}
