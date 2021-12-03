package easypost_test

import (
	"github.com/EasyPost/easypost-go"
)

func (c *ClientTests) TestInsuranceCreation() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()
	to, err := client.CreateAddress(
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

	from, err := client.CreateAddress(
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

	insurance, _ := client.CreateInsurance(
		&easypost.Insurance{
			ToAddress:    to,
			FromAddress:  from,
			TrackingCode: "EZ2000000002",
			Carrier:      "USPS",
			Amount:       "101.00",
		},
	)
	assert.NotNil(insurance.ToAddress)
	assert.NotNil(insurance.FromAddress)
	assert.NotNil(insurance.Tracker)
	assert.Equal("EZ2000000002", insurance.TrackingCode)
	assert.Equal("101.00000", insurance.Amount)

	insurance2, err := client.GetInsurance(insurance.ID)
	require.NoError(err)
	assert.Equal(insurance.ID, insurance2.ID)
	assert.NotNil(insurance2.ToAddress)
	assert.NotNil(insurance2.FromAddress)
	assert.NotNil(insurance2.TrackingCode)
	assert.Equal(insurance.Amount, insurance2.Amount)
	assert.NotNil(insurance2.Tracker)

	res, err := client.ListInsurances(&easypost.ListOptions{PageSize: 5})
	require.NoError(err)
	assert.Len(res.Insurances, 5)
	assert.True(res.HasMore)
}
