package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v3"
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
	assert.Equal("DhlEcsAccount", carrierAccount.Type)

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountCreateWithCustomWorkflow() {
	assert, require := c.Assert(), c.Require()

	client := c.ProdClient()

	carrierAccount := c.fixture.BasicCarrierAccount()
	carrierAccount.Type = "FedexAccount"
	// we need to include data in this interface, otherwise it will be omitted during the API call
	carrierAccount.RegistrationData = map[string]interface{}{
		"some": "data",
	}

	_, err := client.CreateCarrierAccount(carrierAccount)

	// We're sending bad data to the API, so we expect an error
	require.Error(err)
	assert.Equal(reflect.TypeOf(&easypost.InvalidRequestError{}), reflect.TypeOf(err))
	assert.Equal(422, err.(*easypost.InvalidRequestError).StatusCode)
	assert.NotEmpty(err.(*easypost.InvalidRequestError).Errors)
	// We expect one of the sub-errors to be regarding a missing field
	errorFound := false
	for _, err := range err.(*easypost.InvalidRequestError).Errors {
		if err.Field == "account_number" && err.Message == "must be present and a string" {
			errorFound = true
			break
		}
	}

	assert.True(errorFound)
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
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	testDescription := "My custom description"

	carrierAccount, err := c.GenerateCarrierAccount()
	require.NoError(err)

	carrierAccount.Description = testDescription

	updatedCarrierAccount, err := client.UpdateCarrierAccount(
		&easypost.CarrierAccount{
			ID:          carrierAccount.ID,
			Description: testDescription,
		},
	)
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
