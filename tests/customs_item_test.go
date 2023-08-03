package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v3"
)

func (c *ClientTests) TestCustomsItemCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	customsItem, err := client.CreateCustomsItem(c.fixture.BasicCustomsItem())
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.CustomsItem{}), reflect.TypeOf(customsItem))
	assert.True(strings.HasPrefix(customsItem.ID, "cstitem_"))
	assert.Equal(23.25, customsItem.Value)
}

func (c *ClientTests) TestCustomsItemRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	customsItem, err := client.CreateCustomsItem(c.fixture.BasicCustomsItem())
	require.NoError(err)

	retrievedCustomsItem, err := client.GetCustomsItem(customsItem.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.CustomsItem{}), reflect.TypeOf(retrievedCustomsItem))
	assert.Equal(customsItem, retrievedCustomsItem)
}
