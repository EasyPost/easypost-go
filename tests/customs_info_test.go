package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v3"
)

func (c *ClientTests) TestCustomsInfoCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	customsInfo, err := client.CreateCustomsInfo(c.fixture.BasicCustomsInfo())
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.CustomsInfo{}), reflect.TypeOf(customsInfo))
	assert.True(strings.HasPrefix(customsInfo.ID, "cstinfo_"))
	assert.Equal("NOEEI 30.37(a)", customsInfo.EELPFC)
}

func (c *ClientTests) TestCustomsInfoRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	customsInfo, err := client.CreateCustomsInfo(c.fixture.BasicCustomsInfo())
	require.NoError(err)

	retrievedCustomsInfo, err := client.GetCustomsInfo(customsInfo.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.CustomsInfo{}), reflect.TypeOf(retrievedCustomsInfo))
	assert.Equal(customsInfo, retrievedCustomsInfo)
}
