package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestCustomsInfoCreate() {
	client := c.TestClient()
	assert := c.Assert()

	customsInfo, err := client.CreateCustomsInfo(c.fixture.BasicCustomsInfo())
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.CustomsInfo{}), reflect.TypeOf(customsInfo))
	assert.True(strings.HasPrefix(customsInfo.ID, "cstinfo_"))
	assert.Equal("NOEEI 30.37(a)", customsInfo.EELPFC)
}

func (c *ClientTests) TestCustomsInfoRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	customsInfo, err := client.CreateCustomsInfo(c.fixture.BasicCustomsInfo())
	if err != nil {
		c.T().Error(err)
		return
	}

	retrievedCustomsInfo, err := client.GetCustomsInfo(customsInfo.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.CustomsInfo{}), reflect.TypeOf(retrievedCustomsInfo))
	assert.Equal(customsInfo, retrievedCustomsInfo)
}
