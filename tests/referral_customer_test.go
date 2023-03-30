package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestBetaReferralAddPaymentMethod() {
	client := c.ReferralClient()
	assert, require := c.Assert(), c.Require()

	_, err := client.BetaAddPaymentMethod("cus_123", "ba_123", easypost.PrimaryPaymentMethodPriority)
	require.Error(err)
	if err, ok := err.(*easypost.APIError); ok {
		assert.Equal(422, err.StatusCode)
		assert.Equal("BILLING.INVALID_PAYMENT_GATEWAY_REFERENCE", err.Code)
		assert.Equal("Invalid Payment Gateway Reference.", err.Message)
	}
}

func (c *ClientTests) TestBetaReferralRefundByAmount() {
	client := c.ReferralClient()
	assert, require := c.Assert(), c.Require()

	_, err := client.BetaRefundByAmount(2000)
	require.Error(err)
	if err, ok := err.(*easypost.APIError); ok {
		assert.Equal(422, err.StatusCode)
		assert.Equal("TRANSACTION.AMOUNT_INVALID", err.Code)
		assert.Equal("Refund amount is invalid. Please use a valid amount or escalate to finance.", err.Message)
	}
}

func (c *ClientTests) TestBetaReferralRefundByPaymentLogId() {
	client := c.ReferralClient()
	assert, require := c.Assert(), c.Require()

	_, err := client.BetaRefundByPaymentLog("paylog_...")
	require.Error(err)
	if err, ok := err.(*easypost.APIError); ok {
		assert.Equal(422, err.StatusCode)
		assert.Equal("TRANSACTION.DOES_NOT_EXIST", err.Code)
		assert.Equal("We could not find a transaction with that id.", err.Message)
	}
}

func (c *ClientTests) TestReferralCustomerCreate() {
	client := c.PartnerClient()
	assert, require := c.Assert(), c.Require()

	user, err := client.CreateReferralCustomer(
		c.fixture.ReferralUser(),
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.ReferralCustomer{}), reflect.TypeOf(user))
	assert.True(strings.HasPrefix(user.ID, "user_"))
	assert.Equal("Test Referral", user.Name)
}

func (c *ClientTests) TestReferralCustomerAll() {
	client := c.PartnerClient()
	assert, require := c.Assert(), c.Require()

	referralCustomerCollection, err := client.ListReferralCustomers(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	assert.LessOrEqual(len(referralCustomerCollection.ReferralCustomers), c.fixture.pageSize())
	assert.NotNil(referralCustomerCollection.HasMore)
	for _, referralCustomer := range referralCustomerCollection.ReferralCustomers {
		assert.Equal(reflect.TypeOf(&easypost.ReferralCustomer{}), reflect.TypeOf(referralCustomer))
	}

	users, err := client.ListReferralCustomers(nil)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.ListReferralCustomersResult{}), reflect.TypeOf(users))
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

	creditCard, err := client.AddReferralCustomerCreditCard(referralAPIKey, c.fixture.TestCreditCard(), easypost.PrimaryPaymentMethodPriority)
	require.NoError(err)

	require.Equal(reflect.TypeOf(&easypost.PaymentMethodObject{}), reflect.TypeOf(creditCard))
	assert.True(strings.HasSuffix(c.fixture.TestCreditCard().Number, creditCard.Last4))
}

func (c *ClientTests) TestReferralCustomersGetNextPage() {
	client := c.PartnerClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListReferralCustomers(
		&easypost.ListOptions{
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
		assert.Equal(err.Error(), easypost.EndOfPaginationError.Error())
		return
	}
}
