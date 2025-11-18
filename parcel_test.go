package easypost

import (
	"reflect"
	"strings"
)

func (c *ClientTests) TestParcelCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	parcel, err := client.CreateParcel(c.fixture.BasicParcel())
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&Parcel{}), reflect.TypeOf(parcel))
	assert.True(strings.HasPrefix(parcel.ID, "prcl_"))
	assert.Equal(15.4, parcel.Weight)
}

func (c *ClientTests) TestParcelRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	parcel, err := client.CreateParcel(c.fixture.BasicParcel())
	require.NoError(err)

	retrievedParcel, err := client.GetParcel(parcel.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&Parcel{}), reflect.TypeOf(retrievedParcel))
	assert.Equal(parcel, retrievedParcel)
}
