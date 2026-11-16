package easypost

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func (c *ClientTests) GenerateCarrierAccount() (*CarrierAccount, error) {
	client := c.ProdClient()

	carrierAccount, err := client.CreateCarrierAccount(c.fixture.BasicCarrierAccount())

	return carrierAccount, err
}

func (c *ClientTests) TestCarrierAccountCreate() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	carrierAccount, err := c.GenerateCarrierAccount()
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&CarrierAccount{}), reflect.TypeOf(carrierAccount))
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

	var invalidRequestError *InvalidRequestError
	if errors.As(err, &invalidRequestError) {
		assert.Equal(reflect.TypeOf(&InvalidRequestError{}), reflect.TypeOf(err))
		assert.Equal(422, invalidRequestError.StatusCode)
		// We expect one of the sub-errors to be regarding a missing field
		if errorsList, ok := invalidRequestError.Errors.([]interface{}); ok {
			if fieldError, ok := errorsList[0].(*FieldError); ok {
				assert.Equal("shipping_streets", fieldError.Field)
				assert.Equal("must be present and a string", fieldError.Message)
			}
		}
	}
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingFedexDefaultAccount() {
	c.testCarrierAccountCreateURLRouting("FedexDefaultAccount", "carrier_accounts/open")
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingGofoWalletAccount() {
	c.testCarrierAccountCreateURLRouting("GofoWalletAccount", "carrier_accounts/open")
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingUpsDapAccount() {
	c.testCarrierAccountCreateURLRouting("UpsDapAccount", "carrier_accounts/open")
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingUniUniWalletAccount() {
	c.testCarrierAccountCreateURLRouting("UniUniWalletAccount", "carrier_accounts/open")
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingVehoWalletAccount() {
	c.testCarrierAccountCreateURLRouting("VehoWalletAccount", "carrier_accounts/open")
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingFedexAccount() {
	c.testCarrierAccountCreateURLRouting("FedexAccount", "carrier_accounts/register")
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingFedexSmartpostAccount() {
	c.testCarrierAccountCreateURLRouting("FedexSmartpostAccount", "carrier_accounts/register")
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingAmazonShippingAccount() {
	c.testCarrierAccountCreateURLRouting("AmazonShippingAccount", "carrier_accounts/register_oauth")
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingUpsAccount() {
	c.testCarrierAccountCreateURLRouting("UpsAccount", "carrier_accounts/register_oauth")
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingUpsMailInnovationsAccount() {
	c.testCarrierAccountCreateURLRouting("UpsMailInnovationsAccount", "carrier_accounts/register_oauth")
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingUpsSurepostAccount() {
	c.testCarrierAccountCreateURLRouting("UpsSurepostAccount", "carrier_accounts/register_oauth")
}

func (c *ClientTests) TestCarrierAccountCreateUrlRoutingUspsShipAccount() {
	c.testCarrierAccountCreateURLRouting("UspsShipAccount", "carrier_accounts/register_oauth")
}

func (c *ClientTests) testCarrierAccountCreateURLRouting(carrierAccountType, expectedEndpoint string) {
	assert, require := c.Assert(), c.Require()

	client := c.ProdClient()
	capturedURL := ""

	hookID := fmt.Sprintf("route-%s", carrierAccountType)
	requestSubscriber := RequestHookEventSubscriber{
		Callback: func(ctx context.Context, event RequestHookEvent) error {
			if capturedURL == "" && strings.HasPrefix(event.URL.Path, "/v2/carrier_accounts") {
				capturedURL = event.URL.Path
			}
			return nil
		},
		HookEventSubscriber: HookEventSubscriber{
			ID: hookID,
		},
	}
	client.Hooks.AddRequestEventSubscriber(requestSubscriber)

	var err error
	var carrierAccount *CarrierAccount
	switch expectedEndpoint {
	case "carrier_accounts/register_oauth":
		if isUpsCarrierAccountType(carrierAccountType) {
			carrierAccount, err = client.CreateUpsCarrierAccount(&UpsCarrierAccountCreationParameters{
				Type:          carrierAccountType,
				AccountNumber: "123456789",
			})
		} else {
			carrierAccount, err = client.CreateOauthCarrierAccount(&CarrierAccount{
				Type: carrierAccountType,
			})
		}
	default:
		carrierAccount, err = client.CreateCarrierAccount(&CarrierAccount{
			Type: carrierAccountType,
			RegistrationData: map[string]interface{}{
				"some": "data",
			},
		})
	}

	var apiError *APIError
	if errors.As(err, &apiError) {
		// Intentionally no status assertion: this test validates URL routing only.
	} else if err != nil {
		require.NoError(err)
	} else {
		deleteErr := client.DeleteCarrierAccount(carrierAccount.ID)
		require.NoError(deleteErr)
	}

	assert.True(strings.HasSuffix(capturedURL, expectedEndpoint), capturedURL)

	client.Hooks.RemoveRequestEventSubscriber(requestSubscriber)
}

func isUpsCarrierAccountType(carrierAccountType string) bool {
	for _, upsCarrierAccountType := range getUpsCarrierAccountTypes() {
		if carrierAccountType == upsCarrierAccountType {
			return true
		}
	}

	return false
}

func (c *ClientTests) TestCarrierAccountPreventUsersUsingUpsAccountForGenericCreation() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := c.fixture.BasicCarrierAccount()
	createParameters.Type = "UpsAccount"

	_, err := client.CreateCarrierAccount(createParameters)
	require.Error(err)

	assert.Equal(reflect.TypeOf(&InvalidFunctionError{}), reflect.TypeOf(err))
}

func (c *ClientTests) TestCarrierAccountCreateUps() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := &UpsCarrierAccountCreationParameters{
		AccountNumber: "123456789",
		Type:          "UpsAccount",
	}

	carrierAccount, err := client.CreateUpsCarrierAccount(createParameters)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&CarrierAccount{}), reflect.TypeOf(carrierAccount))
	assert.True(strings.HasPrefix(carrierAccount.ID, "ca_"))
	assert.Equal("UpsAccount", carrierAccount.Type)

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountPreventUsersUsingNotUpsAccountForUpsCreation() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := &UpsCarrierAccountCreationParameters{
		AccountNumber: "123456789",
		Type:          "NotUpsAccount",
	}

	_, err := client.CreateUpsCarrierAccount(createParameters)
	require.Error(err)

	assert.Equal(reflect.TypeOf(&InvalidFunctionError{}), reflect.TypeOf(err))
}

func (c *ClientTests) TestCarrierAccountCreateOauth() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := &CarrierAccount{
		Type: "AmazonShippingAccount",
	}

	carrierAccount, err := client.CreateOauthCarrierAccount(createParameters)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&CarrierAccount{}), reflect.TypeOf(carrierAccount))
	assert.True(strings.HasPrefix(carrierAccount.ID, "ca_"))
	assert.Equal("AmazonShippingAccount", carrierAccount.Type)

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountPreventUsersUsingNotOauthAccountForOauthCreation() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := &CarrierAccount{
		Type: "NotOauthAccount",
	}

	_, err := client.CreateOauthCarrierAccount(createParameters)
	require.Error(err)

	assert.Equal(reflect.TypeOf(&InvalidFunctionError{}), reflect.TypeOf(err))
}

func (c *ClientTests) TestCarrierAccountRetrieve() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	carrierAccount, err := c.GenerateCarrierAccount()
	require.NoError(err)

	retrievedCarrierAccount, err := client.GetCarrierAccount(carrierAccount.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&CarrierAccount{}), reflect.TypeOf(retrievedCarrierAccount))
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
		assert.Equal(reflect.TypeOf(&CarrierAccount{}), reflect.TypeOf(carrierAccount))
	}
}

func (c *ClientTests) TestCarrierAccountUpdate() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	testDescription := "My custom description"

	carrierAccount, err := c.GenerateCarrierAccount()
	require.NoError(err)

	updateParameters := &CarrierAccount{
		ID:          carrierAccount.ID,
		Description: testDescription,
	}

	updatedCarrierAccount, err := client.UpdateCarrierAccount(updateParameters)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&CarrierAccount{}), reflect.TypeOf(updatedCarrierAccount))
	assert.True(strings.HasPrefix(updatedCarrierAccount.ID, "ca_"))
	assert.Equal(testDescription, updatedCarrierAccount.Description)

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountPreventUsersUsingUpsAccountForGenericUpdate() {
	id := "ca_123"

	mockRequests := []MockRequest{
		{
			MatchRule: MockRequestMatchRule{
				Method:          "GET",
				URLRegexPattern: "v2\\/carrier_accounts/" + id + "$",
			},
			ResponseInfo: MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"id": "` + id + `", "type": "UpsAccount"}`,
			},
		},
	}

	client := c.MockClient(mockRequests)
	assert, require := c.Assert(), c.Require()

	updateParameters := &CarrierAccount{
		ID:          id,
		Description: "My custom description",
	}

	_, err := client.UpdateCarrierAccount(updateParameters)
	require.Error(err)

	assert.Equal(reflect.TypeOf(&InvalidFunctionError{}), reflect.TypeOf(err))
}

func (c *ClientTests) TestCarrierAccountUpdateUps() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	createParameters := &UpsCarrierAccountCreationParameters{
		AccountNumber: "123456789",
		Type:          "UpsAccount",
	}

	carrierAccount, err := client.CreateUpsCarrierAccount(createParameters)
	require.NoError(err)

	updateParameters := &UpsCarrierAccountUpdateParameters{
		AccountNumber: "987654321",
	}

	updatedCarrierAccount, err := client.UpdateUpsCarrierAccount(carrierAccount.ID, updateParameters)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&CarrierAccount{}), reflect.TypeOf(updatedCarrierAccount))
	assert.True(strings.HasPrefix(updatedCarrierAccount.ID, "ca_"))

	err = client.DeleteCarrierAccount(carrierAccount.ID)
	require.NoError(err)
}

func (c *ClientTests) TestCarrierAccountPreventUsersUsingNotUpsAccountForUpsUpdate() {
	id := "ca_123"

	mockRequests := []MockRequest{
		{
			MatchRule: MockRequestMatchRule{
				Method:          "GET",
				URLRegexPattern: "v2\\/carrier_accounts/" + id + "$",
			},
			ResponseInfo: MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"id": "` + id + `", "type": "NotUpsAccount"}`,
			},
		},
	}

	client := c.MockClient(mockRequests)
	assert, require := c.Assert(), c.Require()

	updateParameters := &UpsCarrierAccountUpdateParameters{
		AccountNumber: "987654321",
	}

	_, err := client.UpdateUpsCarrierAccount(id, updateParameters)
	require.Error(err)

	assert.Equal(reflect.TypeOf(&InvalidFunctionError{}), reflect.TypeOf(err))
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

	assert.Equal(reflect.TypeOf([]*CarrierType{}), reflect.TypeOf(carriersAccountTypes))
}
