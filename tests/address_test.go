package easypost_test

import (
	"errors"
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v5"
)

func (c *ClientTests) TestAddressCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	address, err := client.CreateAddress(
		c.fixture.CaAddress1(),
		nil,
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(address))
	assert.True(strings.HasPrefix(address.ID, "adr_"))
	assert.Equal("388 Townsend St", address.Street1)
}

func (c *ClientTests) TestAddressCreateVerifyStrict() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	address, err := client.CreateAddress(
		c.fixture.CaAddress1(),
		&easypost.CreateAddressOptions{
			VerifyStrict: true,
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(address))
	assert.True(strings.HasPrefix(address.ID, "adr_"))
	assert.Equal("388 TOWNSEND ST APT 20", address.Street1)
}

func (c *ClientTests) TestAddressRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	address, err := client.CreateAddress(
		c.fixture.CaAddress1(),
		nil,
	)
	require.NoError(err)

	retrievedAddress, err := client.GetAddress(address.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(retrievedAddress))
	assert.Equal(address, retrievedAddress)
}

func (c *ClientTests) TestAddressAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	addresses, err := client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	assert.LessOrEqual(len(addresses.Addresses), c.fixture.pageSize())
	assert.NotNil(addresses.HasMore)
	for _, address := range addresses.Addresses {
		assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(address))
	}
}

func (c *ClientTests) TestAddressGetNextPage() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextAddressPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Addresses) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Addresses[len(firstPage.Addresses)-1].ID
			firstIdOfSecondPage := nextPage.Addresses[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIdOfSecondPage)
		}
	}()
	if err != nil {
		var endOfPaginationErr *easypost.EndOfPaginationError
		if errors.As(err, &endOfPaginationErr) {
			assert.Equal(err.Error(), endOfPaginationErr.Error())
			return
		}
		require.NoError(err)
	}
}

// We purposefully pass in slightly incorrect data to get the corrected address back once verified.
func (c *ClientTests) TestAddressCreateVerify() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	address, err := client.CreateAddress(
		c.fixture.IncorrectAddress(),
		&easypost.CreateAddressOptions{
			Verify: true,
		},
	)

	// Does return an address from CreateAddress even if requested verification fails
	require.NoError(err)
	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(address))

	// Delivery verification assertions
	deliveryVerification := address.Verifications.Delivery
	assert.False(deliveryVerification.Success)
	assert.Empty(deliveryVerification.Details)
	deliveryError := deliveryVerification.Errors[0]
	assert.Equal("E.ADDRESS.NOT_FOUND", deliveryError.Code)
	assert.Equal("address", deliveryError.Field)
	assert.Empty(deliveryError.Suggestion)
	assert.Equal("Address not found", deliveryError.Message)

	// ZIP4 verification assertions
	zip4Verification := address.Verifications.ZIP4
	assert.False(zip4Verification.Success)
	assert.Nil(zip4Verification.Details)
	zip4Error := zip4Verification.Errors[0]
	assert.Equal("E.ADDRESS.NOT_FOUND", zip4Error.Code)
	assert.Equal("address", zip4Error.Field)
	assert.Equal("", zip4Error.Suggestion)
	assert.Equal("Address not found", zip4Error.Message)
}

func (c *ClientTests) TestAddressCreateAndVerify() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	// Creating normally (without specifying "verify") will make the address, perform no verifications
	address, err := client.CreateAddress(
		c.fixture.IncorrectAddress(),
		&easypost.CreateAddressOptions{},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(address))
	assert.Nil(address.Verifications.Delivery)

	// Creating with verify would attempt to make the address and perform verifications
	address, err = client.CreateAndVerifyAddress(
		c.fixture.IncorrectAddress(),
		&easypost.CreateAddressOptions{
			Verify: true,
		},
	)

	// Does not return an address from CreateAndVerifyAddress if requested verification fails
	require.Error(err)
	assert.Nil(address)
}

func (c *ClientTests) TestAddressVerify() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	address, err := client.CreateAddress(
		c.fixture.CaAddress1(),
		nil,
	)
	require.NoError(err)

	verifiedAddress, err := client.VerifyAddress(address.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(verifiedAddress))
	assert.True(strings.HasPrefix(verifiedAddress.ID, "adr_"))
	assert.Equal("388 TOWNSEND ST APT 20", verifiedAddress.Street1)
}
