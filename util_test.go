package easypost

func (c *ClientTests) TestStringPointer() {
	assert := c.Assert()

	stringPointer := *StringPtr("test")
	assert.Equal(stringPointer, "test")
}

func (c *ClientTests) TestBoolPointer() {
	assert := c.Assert()

	boolPointer := *BoolPtr(true)
	assert.Equal(boolPointer, true)
}
