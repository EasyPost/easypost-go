package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
	"reflect"
	"strings"
)

func (c *ClientTests) TestParcelCreate() {
	client := c.TestClient()
	assert := c.Assert()

	parcel, _ := client.CreateParcel(c.fixture.BasicParcel())

	assert.Equal(reflect.TypeOf(&easypost.Parcel{}), reflect.TypeOf(parcel))
	assert.True(strings.HasPrefix(parcel.ID, "prcl_"))
	assert.Equal(15.4, parcel.Weight)
}

func (c *ClientTests) TestParcelRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	parcel, _ := client.CreateParcel(c.fixture.BasicParcel())

	retrievedParcel, _ := client.GetParcel(parcel.ID)

	assert.Equal(reflect.TypeOf(&easypost.Parcel{}), reflect.TypeOf(retrievedParcel))
	assert.Equal(parcel, retrievedParcel)
}
