package easypost_test

import (
	"testing"

	"github.com/EasyPost/easypost-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddressCreationVerification(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	// Create an address and then verify some fields to test whether it was
	// created just fine.
	address, err := TestClient.CreateAddress(
		&easypost.Address{
			Company: "EasyPost",
			Street1: "One Montgomery St",
			Street2: "Ste 400",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94104",
			Phone:   "415-555-1212",
		},
		nil,
	)
	require.NoError(err)
	address, err = TestClient.VerifyAddress(address.ID)
	require.NoError(err)
	require.NotNil(address)
	address, err = TestClient.GetAddress(address.ID)
	require.NoError(err)
	require.NotNil(address)
	assert.Equal("US", address.Country)
	assert.Empty(address.Email)
	assert.Empty(address.FederalTaxID)
	assert.Equal("CA", address.State)
	assert.Equal("94104-4533", address.Zip)
}

func TestAddressCreationWithVerify(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	// reate an address with a verify parameter to test that it verifies
	// accurately.
	address, err := TestClient.CreateAddress(
		&easypost.Address{
			Company: "EasyPost",
			Street1: "One Montgomery St",
			Street2: "Ste 400",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94104",
			Phone:   "415-555-1212",
		},
		&easypost.CreateAddressOptions{Verify: []string{"delivery"}},
	)
	require.NoError(err)
	assert.NotEmpty(address.ID)
	assert.Equal("1 MONTGOMERY ST STE 400", address.Street1)
	assert.Empty(address.Street2)
	assert.Equal("US", address.Country)
}

func TestAddressCreationWithVerifyFailure(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	// Create an address with a verify parameter to test that it fails
	// elegantly.
	address, err := TestClient.CreateAddress(
		&easypost.Address{
			Street1: "UNDELIEVRABLE ST",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94105",
			Country: "US",
			Company: "EasyPost",
			Phone:   "415-555-1212",
		},
		&easypost.CreateAddressOptions{Verify: []string{"delivery"}},
	)
	require.NoError(err)
	assert.NotEmpty(address.ID)
	require.NotNil(address.Verifications)
	require.NotNil(address.Verifications.Delivery)
	require.GreaterOrEqual(len(address.Verifications.Delivery.Errors), 2)
	assert.Equal(
		"Address not found", address.Verifications.Delivery.Errors[0].Message,
	)
	assert.Equal(
		"House number is missing",
		address.Verifications.Delivery.Errors[1].Message,
	)
}

func TestAddressCreationWithVerifyStrictFailure(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	// Create an address with a verify strict parameter to test that it fails
	// elegantly.
	address, err := TestClient.CreateAddress(
		&easypost.Address{
			Street1: "UNDELIEVRABLE ST",
			City:    "San Francisco",
			State:   "CA",
			Zip:     "94105",
			Country: "US",
			Company: "EasyPost",
			Phone:   "415-555-1212",
		},
		&easypost.CreateAddressOptions{VerifyStrict: []string{"delivery"}},
	)
	assert.Nil(address)
	require.Error(err)
	if assert.IsType((*easypost.APIError)(nil), err) {
		err := err.(*easypost.APIError)
		assert.Equal("ADDRESS.VERIFY.FAILURE", err.Code)
		assert.Equal("Unable to verify address.", err.Message)
		require.GreaterOrEqual(len(err.Errors), 2)
		assert.Equal("Address not found", err.Errors[0].Message)
		assert.Equal("House number is missing", err.Errors[1].Message)
	}
}

func TestAddressUnicode(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	// Create an address with unicode field and assert if it was created
	// correctly.
	address, err := TestClient.CreateAddress(
		&easypost.Address{State: "DELEGACIóN BENITO JUáREZ"}, nil,
	)
	require.NoError(err)
	assert.Equal("DELEGACIóN BENITO JUáREZ", address.State)
}
