package easypost_test

import (
	"testing"

	"github.com/EasyPost/easypost-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParcelCreation(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	parcel, err := TestClient.CreateParcel(
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
