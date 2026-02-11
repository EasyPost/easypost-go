package easypost

func GetFedExRegistrationMockRequests() []MockRequest {
	return []MockRequest{
		{
			MatchRule: MockRequestMatchRule{
				Method:          "POST",
				UrlRegexPattern: "v2\\/fedex_registrations\\/\\S*\\/address$",
			},
			ResponseInfo: MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"email_address":null,"options":["SMS","CALL","INVOICE"],"phone_number":"***-***-9721"}`,
			},
		},
		{
			MatchRule: MockRequestMatchRule{
				Method:          "POST",
				UrlRegexPattern: "v2\\/fedex_registrations\\/\\S*\\/pin$",
			},
			ResponseInfo: MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"message":"sent secured Pin"}`,
			},
		},
		{
			MatchRule: MockRequestMatchRule{
				Method:          "POST",
				UrlRegexPattern: "v2\\/fedex_registrations\\/\\S*\\/pin\\/validate$",
			},
			ResponseInfo: MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"id":"ca_123","object":"CarrierAccount","type":"FedexAccount","credentials":{"account_number":"123456789","mfa_key":"123456789-XXXXX"}}`,
			},
		},
		{
			MatchRule: MockRequestMatchRule{
				Method:          "POST",
				UrlRegexPattern: "v2\\/fedex_registrations\\/\\S*\\/invoice$",
			},
			ResponseInfo: MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"id":"ca_123","object":"CarrierAccount","type":"FedexAccount","credentials":{"account_number":"123456789","mfa_key":"123456789-XXXXX"}}`,
			},
		},
	}
}

func (c *ClientTests) TestRegisterFedExAddress() {
	mockRequests := GetFedExRegistrationMockRequests()
	client := c.MockClient(mockRequests)
	assert, require := c.Assert(), c.Require()

	fedexAccountNumber := "123456789"
	addressValidation := map[string]interface{}{
		"name":         "BILLING NAME",
		"street1":      "1234 BILLING STREET",
		"city":         "BILLINGCITY",
		"state":        "ST",
		"postal_code":  "12345",
		"country_code": "US",
	}

	easypostDetails := map[string]interface{}{
		"carrier_account_id": "ca_123",
	}

	params := map[string]interface{}{
		"address_validation": addressValidation,
		"easypost_details":   easypostDetails,
	}

	response, err := client.RegisterFedExAddress(fedexAccountNumber, params)
	require.NoError(err)

	assert.Equal("", response.EmailAddress)
	assert.Contains(response.Options, "SMS")
	assert.Contains(response.Options, "CALL")
	assert.Contains(response.Options, "INVOICE")
	assert.Equal("***-***-9721", response.PhoneNumber)
}

func (c *ClientTests) TestRequestFedExPin() {
	mockRequests := GetFedExRegistrationMockRequests()
	client := c.MockClient(mockRequests)
	assert, require := c.Assert(), c.Require()

	fedexAccountNumber := "123456789"

	response, err := client.RequestFedExPin(fedexAccountNumber, "SMS")
	require.NoError(err)

	assert.Equal("sent secured Pin", response.Message)
}

func (c *ClientTests) TestValidateFedExPin() {
	mockRequests := GetFedExRegistrationMockRequests()
	client := c.MockClient(mockRequests)
	assert, require := c.Assert(), c.Require()

	fedexAccountNumber := "123456789"
	pinValidation := map[string]interface{}{
		"pin_code": "123456",
		"name":     "BILLING NAME",
	}

	easypostDetails := map[string]interface{}{
		"carrier_account_id": "ca_123",
	}

	params := map[string]interface{}{
		"pin_validation":   pinValidation,
		"easypost_details": easypostDetails,
	}

	response, err := client.ValidateFedExPin(fedexAccountNumber, params)
	require.NoError(err)

	assert.Equal("ca_123", response.ID)
	assert.Equal("CarrierAccount", response.Object)
	assert.Equal("FedexAccount", response.Type)
	assert.Equal("123456789", response.Credentials["account_number"])
	assert.Equal("123456789-XXXXX", response.Credentials["mfa_key"])
}

func (c *ClientTests) TestSubmitFedExInvoice() {
	mockRequests := GetFedExRegistrationMockRequests()
	client := c.MockClient(mockRequests)
	assert, require := c.Assert(), c.Require()

	fedexAccountNumber := "123456789"
	invoiceValidation := map[string]interface{}{
		"name":             "BILLING NAME",
		"invoice_number":   "INV-12345",
		"invoice_date":     "2025-12-08",
		"invoice_amount":   "100.00",
		"invoice_currency": "USD",
	}

	easypostDetails := map[string]interface{}{
		"carrier_account_id": "ca_123",
	}

	params := map[string]interface{}{
		"invoice_validation": invoiceValidation,
		"easypost_details":   easypostDetails,
	}

	response, err := client.SubmitFedExInvoice(fedexAccountNumber, params)
	require.NoError(err)

	assert.Equal("ca_123", response.ID)
	assert.Equal("CarrierAccount", response.Object)
	assert.Equal("FedexAccount", response.Type)
	assert.Equal("123456789", response.Credentials["account_number"])
	assert.Equal("123456789-XXXXX", response.Credentials["mfa_key"])
}
