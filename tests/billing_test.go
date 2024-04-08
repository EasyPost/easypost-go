package easypost_test

import "github.com/EasyPost/easypost-go/v4"

func GetBillingMockRequests() []easypost.MockRequest {
	return []easypost.MockRequest{
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "POST",
				UrlRegexPattern: "v2\\/bank_accounts\\/\\S*\\/charges$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{}`,
			},
		},
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "POST",
				UrlRegexPattern: "v2\\/credit_cards\\/\\S*\\/charges$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{}`,
			},
		},
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "DELETE",
				UrlRegexPattern: "v2\\/bank_accounts\\/\\S*$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{}`,
			},
		},
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "DELETE",
				UrlRegexPattern: "v2\\/credit_cards\\/\\S*$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{}`,
			},
		},
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "GET",
				UrlRegexPattern: "v2\\/payment_methods$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"id": "summary_123", "primary_payment_method": {"id": "pm_123", "object": "CreditCard", "last4": "1234"}, "secondary_payment_method": {"id": "pm_123", "object": "BankAccount", "bank_name": "Mock Bank"}}`,
			},
		},
	}
}

func (c *ClientTests) TestDeletePaymentMethod() {
	mockRequests := GetBillingMockRequests()

	client := c.MockClient(mockRequests)
	require := c.Require()

	err := client.DeletePaymentMethod(easypost.PrimaryPaymentMethodPriority)
	require.NoError(err)
}

func (c *ClientTests) TestFundWallet() {
	mockRequests := GetBillingMockRequests()

	client := c.MockClient(mockRequests)
	require := c.Require()

	err := client.FundWallet("2000", easypost.PrimaryPaymentMethodPriority)
	require.NoError(err)
}

func (c *ClientTests) TestRetrievePaymentMethods() {
	mockRequests := GetBillingMockRequests()

	client := c.MockClient(mockRequests)
	assert, require := c.Assert(), c.Require()

	paymentMethods, err := client.RetrievePaymentMethods()
	require.NoError(err)

	assert.True(paymentMethods.PrimaryPaymentMethod != nil)
	assert.True(paymentMethods.SecondaryPaymentMethod != nil)
}

func (c *ClientTests) TestGetPaymentMethodInfoByObjectType() {
	mockRequests := []easypost.MockRequest{
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "GET",
				UrlRegexPattern: "v2\\/payment_methods$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"id": "summary_123", "primary_payment_method": {"id": "pm_123", "object": "CreditCard"}, "secondary_payment_method": {"id": "pm_123", "object": "BankAccount"}}`,
			},
		},
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "DELETE",
				UrlRegexPattern: "v2\\/credit_cards\\/pm_123$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{}`,
			},
		},
	}

	client := c.MockClient(mockRequests)
	require := c.Require()

	// getPaymentMethodObjectType is a private method, so we can't test it directly, but we can test it via DeletePaymentMethod
	// The mocking setup here makes it so only /v2/credit_cards/pm_123 is allowed to be called
	// If the method works correctly without error, we can assume it's because it found the correct payment method type
	err := client.DeletePaymentMethod(easypost.PrimaryPaymentMethodPriority)
	require.NoError(err)
}

func (c *ClientTests) TestGetPaymentMethodInfoByLegacyIdPrefix() {
	mockRequests := []easypost.MockRequest{
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "GET",
				UrlRegexPattern: "v2\\/payment_methods$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{"id": "summary_123", "primary_payment_method": {"id": "card_123", "object": null}, "secondary_payment_method": {"id": "bank_123", "object": null}}`,
			},
		},
		{
			MatchRule: easypost.MockRequestMatchRule{
				Method:          "DELETE",
				UrlRegexPattern: "v2\\/bank_accounts\\/bank_123$",
			},
			ResponseInfo: easypost.MockRequestResponseInfo{
				StatusCode: 200,
				Body:       `{}`,
			},
		},
	}

	client := c.MockClient(mockRequests)
	require := c.Require()

	// getPaymentMethodObjectType is a private method, so we can't test it directly, but we can test it via DeletePaymentMethod
	// The mocking setup here makes it so only /v2/bank_accounts/bank_123 is allowed to be called
	// If the method works correctly without error, we can assume it's because it found the correct payment method type
	err := client.DeletePaymentMethod(easypost.SecondaryPaymentMethodPriority)
	require.NoError(err)
}
