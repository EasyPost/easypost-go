package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestCustomsItemCreate() {
	client := c.TestClient()
	assert := c.Assert()

	customsItem, _ := client.CreateCustomsItem(c.fixture.BasicCustomsItem())

	assert.Equal(reflect.TypeOf(&easypost.CustomsItem{}), reflect.TypeOf(customsItem))
	assert.True(strings.HasPrefix(customsItem.ID, "cstitem_"))
	assert.Equal(23.0, customsItem.Value)
}

func (c *ClientTests) TestCustomsItemRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	customsItem, _ := client.CreateCustomsItem(c.fixture.BasicCustomsItem())

	retrievedCustomsItem, _ := client.GetCustomsItem(customsItem.ID)

	assert.Equal(reflect.TypeOf(&easypost.CustomsItem{}), reflect.TypeOf(retrievedCustomsItem))
	assert.Equal(customsItem, retrievedCustomsItem)
}
