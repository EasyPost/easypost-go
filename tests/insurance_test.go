package easypost_test

import (
	"testing"

	"github.com/EasyPost/easypost-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsuranceCreation(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	to, err := TestClient.CreateAddress(
		&easypost.Address{
			Name:    "Bugs Bunny",
			Street1: "4000 Warner Blvd",
			City:    "Burbank",
			State:   "CA",
			Zip:     "91522",
			Phone:   "818-555-1212",
		},
		nil,
	)
	require.NoError(err)

	from, err := TestClient.CreateAddress(
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

	insurance, err := TestClient.CreateInsurance(
		&easypost.Insurance{
			ToAddress:    to,
			FromAddress:  from,
			TrackingCode: "EZ2000000002",
			Carrier:      "USPS",
			Amount:       "101.00",
		},
	)
	require.NoError(err)

	assert.NotNil(insurance.ToAddress)
	assert.NotNil(insurance.FromAddress)
	assert.NotNil(insurance.Tracker)
	assert.Equal("EZ2000000002", insurance.TrackingCode)
	assert.Equal("101.00000", insurance.Amount)

	insurance2, err := TestClient.GetInsurance(insurance.ID)
	require.NoError(err)
	assert.Equal(insurance.ID, insurance2.ID)
	assert.NotNil(insurance2.ToAddress)
	assert.NotNil(insurance2.FromAddress)
	assert.NotNil(insurance2.TrackingCode)
	assert.Equal(insurance.Amount, insurance2.Amount)
	assert.NotNil(insurance2.Tracker)

	res, err := TestClient.ListInsurances(&easypost.ListOptions{PageSize: 5})
	require.NoError(err)
	assert.Len(res.Insurances, 5)
	assert.True(res.HasMore)
}
