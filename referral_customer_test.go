package easypost

import (
	"errors"
	"reflect"
	"strings"
)

func (c *ClientTests) TestBetaReferralAddPaymentMethod() {
	client := c.ReferralClient()
	assert, require := c.Assert(), c.Require()

	_, err := client.BetaAddPaymentMethod("cus_123", "ba_123", PrimaryPaymentMethodPriority)
	require.Error(err)

	var eperr *APIError
	if errors.As(err, &eperr) {
		assert.Equal(422, eperr.StatusCode)
		assert.Equal("BILLING.INVALID_PAYMENT_GATEWAY_REFERENCE", eperr.Code)
		assert.Equal("Invalid Payment Gateway Reference.", eperr.Message)
	}
}

func (c *ClientTests) TestBetaReferralRefundByAmount() {
	client := c.ReferralClient()
	assert, require := c.Assert(), c.Require()

	_, err := client.BetaRefundByAmount(2000)
	require.Error(err)

	var eperr *APIError
	if errors.As(err, &eperr) {
		assert.Equal(422, eperr.StatusCode)
		assert.Equal("TRANSACTION.AMOUNT_INVALID", eperr.Code)
		assert.Equal("Refund amount is invalid. Please use a valid amount or escalate to finance.", eperr.Message)
	}
}

func (c *ClientTests) TestBetaReferralRefundByPaymentLogId() {
	client := c.ReferralClient()
	assert, require := c.Assert(), c.Require()

	_, err := client.BetaRefundByPaymentLog("paylog_...")
	require.Error(err)

	var eperr *APIError
	if errors.As(err, &eperr) {
		assert.Equal(422, eperr.StatusCode)
		assert.Equal("TRANSACTION.DOES_NOT_EXIST", eperr.Code)
		assert.Equal("We could not find a transaction with that id.", eperr.Message)
	}
}

func (c *ClientTests) TestReferralCustomerCreate() {
	client := c.PartnerClient()
	assert, require := c.Assert(), c.Require()

	user, err := client.CreateReferralCustomer(
		c.fixture.ReferralUser(),
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&ReferralCustomer{}), reflect.TypeOf(user))
	assert.True(strings.HasPrefix(user.ID, "user_"))
	assert.Equal("Test Referral", user.Name)
}

func (c *ClientTests) TestReferralCustomerAll() {
	client := c.PartnerClient()
	assert, require := c.Assert(), c.Require()

	referralCustomerCollection, err := client.ListReferralCustomers(
		&ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	assert.LessOrEqual(len(referralCustomerCollection.ReferralCustomers), c.fixture.pageSize())
	assert.NotNil(referralCustomerCollection.HasMore)
	for _, referralCustomer := range referralCustomerCollection.ReferralCustomers {
		assert.Equal(reflect.TypeOf(&ReferralCustomer{}), reflect.TypeOf(referralCustomer))
	}

	users, err := client.ListReferralCustomers(nil)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&ListReferralCustomersResult{}), reflect.TypeOf(users))
	assert.True(len(users.ReferralCustomers) > 0)
}

func (c *ClientTests) TestReferralCustomerUpdate() {
	client := c.PartnerClient()
	require := c.Require()

	user, err := client.CreateReferralCustomer(
		c.fixture.ReferralUser(),
	)
	require.NoError(err)

	email := "email@example.com"
	_, err = client.UpdateReferralCustomerEmail(user.ID, email)
	require.NoError(err)
}

func (c *ClientTests) TestReferralCustomerAddCreditCard() {
	client := c.PartnerClient()
	assert, require := c.Assert(), c.Require()

	referralAPIKey := c.ReferralAPIKey()

	creditCard, err := client.AddReferralCustomerCreditCard(referralAPIKey, c.fixture.TestCreditCard(), PrimaryPaymentMethodPriority)
	require.NoError(err)

	require.Equal(reflect.TypeOf(&PaymentMethodObject{}), reflect.TypeOf(creditCard))
	assert.True(strings.HasSuffix(c.fixture.TestCreditCard().Number, creditCard.Last4))
}

func (c *ClientTests) TestAddReferralCustomerCreditCardFromStripe() {
	client := c.ReferralClient()
	assert, require := c.Assert(), c.Require()

	referralAPIKey := c.ReferralAPIKey()

	_, err := client.AddReferralCustomerCreditCardFromStripe(
		referralAPIKey,
		c.fixture.BillingData().PaymentMethodID,
		PrimaryPaymentMethodPriority,
	)
	require.Error(err)

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		assert.Equal("Stripe::PaymentMethod does not exist for the specified reference_id", apiErr.Message)
	}
}

func (c *ClientTests) TestAddReferralCustomerBankAccountFromStripe() {
	client := c.ReferralClient()
	assert, require := c.Assert(), c.Require()

	referralAPIKey := c.ReferralAPIKey()

	_, err := client.AddReferralCustomerBankAccountFromStripe(
		referralAPIKey,
		c.fixture.BillingData().FinancialConnectionsID,
		c.fixture.BillingData().MandateData,
		PrimaryPaymentMethodPriority,
	)
	require.Error(err)

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		assert.Equal("account_holder_name must be present when creating a Financial Connections payment method", apiErr.Message)
	}
}

func (c *ClientTests) TestReferralCustomersGetNextPage() {
	client := c.PartnerClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListReferralCustomers(
		&ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextReferralCustomerPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.ReferralCustomers) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.ReferralCustomers[len(firstPage.ReferralCustomers)-1].ID
			firstIdOfSecondPage := nextPage.ReferralCustomers[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIdOfSecondPage)
		}
	}()
	if err != nil {
		var endOfPaginationErr *EndOfPaginationError
		if errors.As(err, &endOfPaginationErr) {
			assert.Equal(err.Error(), endOfPaginationErr.Error())
			return
		}
		require.NoError(err)
	}
}

func (c *ClientTests) TestBetaCreateCreditCardClientSecret() {
	client := c.ReferralClient()
	assert, require := c.Assert(), c.Require()

	response, err := client.BetaCreateCreditCardClientSecret()
	require.NoError(err)

	assert.True(strings.HasPrefix(response.ClientSecret, "seti_"))
}

func (c *ClientTests) TestBetaCreateBankAccountClientSecret() {
	client := c.ReferralClient()
	assert, require := c.Assert(), c.Require()

	response, err := client.BetaCreateBankAccountClientSecret()
	require.NoError(err)

	assert.True(strings.HasPrefix(response.ClientSecret, "fcsess_client_secret_"))
}
