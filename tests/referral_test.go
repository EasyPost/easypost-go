package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

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

	creditCard, err := client.AddReferralCustomerCreditCard(ReferralAPIKey, c.fixture.TestCreditCard(), easypost.PrimaryPaymentMethodPriority)
	require.NoError(err)

	require.Equal(reflect.TypeOf(&easypost.PaymentMethodObject{}), reflect.TypeOf(creditCard))
	assert.True(strings.HasSuffix(c.fixture.TestCreditCard().Number, creditCard.Last4))
}
