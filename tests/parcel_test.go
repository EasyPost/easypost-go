package easypost_test

import (
	"github.com/EasyPost/easypost-go"
)

func (c *ClientTests) TestParcelCreation() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()
	parcel, err := client.CreateParcel(
		&easypost.Parcel{
			PredefinedPackage: "RegionalRateBoxA",
			Length:            10.2,
			Width:             7.8,
			Height:            4.3,
			Weight:            21.2,
		},
	)
	require.NoError(err)
	assert.Equal("RegionalRateBoxA", parcel.PredefinedPackage)
	assert.Equal(float64(10.2), parcel.Length)
	assert.Equal(float64(7.8), parcel.Width)
	assert.Equal(float64(4.3), parcel.Height)
	assert.Equal(float64(21.2), parcel.Weight)
}
