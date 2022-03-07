package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestAddressCreate() {
	client := c.TestClient()
	assert := c.Assert()

	address, _ := client.CreateAddress(
		c.fixture.BasicAddress(),
		nil,
	)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(address))
	assert.True(strings.HasPrefix(address.ID, "adr_"))
	assert.Equal("388 Townsend St", address.Street1)
}

func (c *ClientTests) TestAddressCreateVerifyStrict() {
	client := c.TestClient()
	assert := c.Assert()

	address, _ := client.CreateAddress(
		c.fixture.BasicAddress(),
		&easypost.CreateAddressOptions{
			VerifyStrict: []string{"true"},
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(address))
	assert.True(strings.HasPrefix(address.ID, "adr_"))
	assert.Equal("388 TOWNSEND ST APT 20", address.Street1)
}

func (c *ClientTests) TestAddressRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	address, _ := client.CreateAddress(
		c.fixture.BasicAddress(),
		nil,
	)

	retrievedAddress, _ := client.GetAddress(address.ID)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(retrievedAddress))
	assert.Equal(address, retrievedAddress)
}

func (c *ClientTests) TestAddressAll() {
	client := c.TestClient()
	assert := c.Assert()

	addresses, _ := client.ListAddresses(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)

	assert.LessOrEqual(len(addresses.Addresses), c.fixture.pageSize())
	assert.NotNil(addresses.HasMore)
	for _, address := range addresses.Addresses {
		assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(address))
	}
}

// We purposefully pass in slightly incorrect data to get the corrected address back once verified.
func (c *ClientTests) TestAddressCreateVerify() {
	client := c.TestClient()
	assert := c.Assert()

	address, _ := client.CreateAddress(
		c.fixture.IncorrectAddressToVerify(),
		&easypost.CreateAddressOptions{
			Verify: []string{"delivery"},
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(address))
	assert.True(strings.HasPrefix(address.ID, "adr_"))
	assert.Equal("417 MONTGOMERY ST STE 500", address.Street1)
}

func (c *ClientTests) TestAddressVerify() {
	client := c.TestClient()
	assert := c.Assert()

	address, _ := client.CreateAddress(
		c.fixture.BasicAddress(),
		nil,
	)

	verifiedAddress, _ := client.VerifyAddress(address.ID)

	assert.Equal(reflect.TypeOf(&easypost.Address{}), reflect.TypeOf(verifiedAddress))
	assert.True(strings.HasPrefix(verifiedAddress.ID, "adr_"))
	assert.Equal("388 TOWNSEND ST APT 20", verifiedAddress.Street1)
}
