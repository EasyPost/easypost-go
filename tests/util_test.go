package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestStringPointer() {
	assert := c.Assert()

	stringPointer := *easypost.StringPtr("test")
	assert.Equal(stringPointer, "test")
}

func (c *ClientTests) TestBoolPointer() {
	assert := c.Assert()

	boolPointer := *easypost.BoolPtr(true)
	assert.Equal(boolPointer, true)
}
