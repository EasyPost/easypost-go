package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
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
			VerifyStrict: []string{"true"},
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

// We purposefully pass in slightly incorrect data to get the corrected address back once verified.
func (c *ClientTests) TestAddressCreateVerify() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	address, err := client.CreateAddress(
		c.fixture.IncorrectAddress(),
		&easypost.CreateAddressOptions{
			Verify: []string{"delivery"},
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(address))
	assert.True(strings.HasPrefix(address.ID, "adr_"))
	assert.Equal("417 MONTGOMERY ST FL 5", address.Street1)
}

func (c *ClientTests) TestAddressCreateAndVerify() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	address, err := client.CreateAndVerifyAddress(
		c.fixture.IncorrectAddress(),
		&easypost.CreateAddressOptions{
			Verify: []string{"delivery"},
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(address))
	assert.True(strings.HasPrefix(address.ID, "adr_"))
	assert.Equal("417 MONTGOMERY ST FL 5", address.Street1)
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
