package easypost_test

import (
	"errors"
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v5"
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

func (c *ClientTests) TestCarrierAccountCreateWithCustomWorkflow() { // FedExAccount
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := c.fixture.BasicCarrierAccount()
	createParameters.Type = "FedexAccount"
	// we need to include data in this interface, otherwise it will be omitted during the API call
	createParameters.RegistrationData = map[string]interface{}{
		"some": "data",
	}

	_, err := client.CreateCarrierAccount(createParameters)

	// We're sending bad data to the API, so we expect an error
	require.Error(err)

	var invalidRequestError *easypost.InvalidRequestError
	if errors.As(err, &invalidRequestError) {
		assert.Equal(reflect.TypeOf(&easypost.InvalidRequestError{}), reflect.TypeOf(err))
		assert.Equal(422, invalidRequestError.StatusCode)
		// We expect one of the sub-errors to be regarding a missing field
		if errorsList, ok := invalidRequestError.Errors.([]interface{}); ok {
			if fieldError, ok := errorsList[0].(*easypost.FieldError); ok {
				assert.Equal("shipping_streets", fieldError.Field)
				assert.Equal("must be present and a string", fieldError.Message)
			}
		}
	}
}

func (c *ClientTests) TestCarrierAccountPreventUsersUsingUpsAccountForGenericCreation() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := c.fixture.BasicCarrierAccount()
	createParameters.Type = "UpsAccount"

	_, err := client.CreateCarrierAccount(createParameters)
	require.Error(err)

	assert.Equal(reflect.TypeOf(&easypost.InvalidFunctionError{}), reflect.TypeOf(err))
}

func (c *ClientTests) TestCarrierAccountCreateUps() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := &easypost.UpsCarrierAccountCreationParameters{
		AccountNumber: "123456789",
		Type:          "UpsAccount",
	}

	carrierAccount, err := client.CreateUpsCarrierAccount(createParameters)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(carrierAccount))
	assert.True(strings.HasPrefix(carrierAccount.ID, "ca_"))
	assert.Equal("UpsAccount", carrierAccount.Type)

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountPreventUsersUsingNotUpsAccountForUpsCreation() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := &easypost.UpsCarrierAccountCreationParameters{
		AccountNumber: "123456789",
		Type:          "NotUpsAccount",
	}

	_, err := client.CreateUpsCarrierAccount(createParameters)
	require.Error(err)

	assert.Equal(reflect.TypeOf(&easypost.InvalidFunctionError{}), reflect.TypeOf(err))
}

func (c *ClientTests) TestCarrierAccountCreateOauth() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := &easypost.CarrierAccount{
		Type: "AmazonShippingAccount",
	}

	carrierAccount, err := client.CreateOauthCarrierAccount(createParameters)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(carrierAccount))
	assert.True(strings.HasPrefix(carrierAccount.ID, "ca_"))
	assert.Equal("AmazonShippingAccount", carrierAccount.Type)

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountPreventUsersUsingNotOauthAccountForOauthCreation() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := &easypost.CarrierAccount{
		Type: "NotOauthAccount",
	}

	_, err := client.CreateOauthCarrierAccount(createParameters)
	require.Error(err)

	assert.Equal(reflect.TypeOf(&easypost.InvalidFunctionError{}), reflect.TypeOf(err))
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

	updateParameters := &easypost.CarrierAccount{
		ID:          carrierAccount.ID,
		Description: testDescription,
	}

	updatedCarrierAccount, err := client.UpdateCarrierAccount(updateParameters)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(updatedCarrierAccount))
	assert.True(strings.HasPrefix(updatedCarrierAccount.ID, "ca_"))
	assert.Equal(testDescription, updatedCarrierAccount.Description)

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountPreventUsersUsingUpsAccountForGenericUpdate() {
	id := "ca_123"

	mockRequests := []easypost.MockRequest{
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "GET",
				UrlRegexPattern: "v2\\/carrier_accounts/" + id + "$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"id": "` + id + `", "type": "UpsAccount"}`,
			},
		},
	}

	client := c.MockClient(mockRequests)
	assert, require := c.Assert(), c.Require()

	updateParameters := &easypost.CarrierAccount{
		ID:          id,
		Description: "My custom description",
	}

	_, err := client.UpdateCarrierAccount(updateParameters)
	require.Error(err)

	assert.Equal(reflect.TypeOf(&easypost.InvalidFunctionError{}), reflect.TypeOf(err))
}

func (c *ClientTests) TestCarrierAccountUpdateUps() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := &easypost.UpsCarrierAccountCreationParameters{
		AccountNumber: "123456789",
		Type:          "UpsAccount",
	}

	carrierAccount, err := client.CreateUpsCarrierAccount(createParameters)
	require.NoError(err)

	updateParameters := &easypost.UpsCarrierAccountUpdateParameters{
		AccountNumber: "987654321",
	}

	updatedCarrierAccount, err := client.UpdateUpsCarrierAccount(carrierAccount.ID, updateParameters)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.CarrierAccount{}), reflect.TypeOf(updatedCarrierAccount))
	assert.True(strings.HasPrefix(updatedCarrierAccount.ID, "ca_"))

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountPreventUsersUsingNotUpsAccountForUpsUpdate() {
	id := "ca_123"

	mockRequests := []easypost.MockRequest{
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "GET",
				UrlRegexPattern: "v2\\/carrier_accounts/" + id + "$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"id": "` + id + `", "type": "NotUpsAccount"}`,
			},
		},
	}

	client := c.MockClient(mockRequests)
	assert, require := c.Assert(), c.Require()

	updateParameters := &easypost.UpsCarrierAccountUpdateParameters{
		AccountNumber: "987654321",
	}

	_, err := client.UpdateUpsCarrierAccount(id, updateParameters)
	require.Error(err)

	assert.Equal(reflect.TypeOf(&easypost.InvalidFunctionError{}), reflect.TypeOf(err))
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
