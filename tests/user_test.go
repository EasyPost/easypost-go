package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestUserCreate() {
	c.T().Skip("This endpoint returns the child user keys in plain text, do not run this test.")

	client := c.ProdClient()
	assert := c.Assert()

	userName := "Test User"

	user, _ := client.CreateUser(
		&easypost.UserOptions{
			Name: &userName,
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(user))
	assert.True(strings.HasPrefix(user.ID, "user_"))
	assert.Equal("Test User", user.Name)
}

func (c *ClientTests) TestUserRetrieve() {
	client := c.ProdClient()
	assert := c.Assert()

	retrievedUser, _ := client.GetUser(c.fixture.ChildUserID())

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(retrievedUser))
	assert.True(strings.HasPrefix(retrievedUser.ID, "user_"))
}

func (c *ClientTests) TestUserRetrieveMe() {
	client := c.ProdClient()
	assert := c.Assert()

	user, _ := client.RetrieveMe()

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(user))
	assert.True(strings.HasPrefix(user.ID, "user_"))
}

func (c *ClientTests) TestUserUpdate() {
	client := c.ProdClient()
	assert := c.Assert()

	user, _ := client.RetrieveMe()

	test_phone := "5555555555"

	updatedUser, _ := client.UpdateUser(
		&easypost.UserOptions{
			ID:          user.ID,
			PhoneNumber: &test_phone,
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(updatedUser))
	assert.True(strings.HasPrefix(updatedUser.ID, "user_"))
	assert.Equal(test_phone, updatedUser.PhoneNumber)
}

func (c *ClientTests) TestUserUpdateBrand() {
	client := c.ProdClient()
	assert := c.Assert()

	color := "#123456"

	user, _ := client.RetrieveMe()

	brand, _ := client.UpdateBrand(
		map[string]interface{}{
			"color": color,
		},
		user.ID,
	)

	assert.Equal(reflect.TypeOf(&easypost.Brand{}), reflect.TypeOf(brand))
	assert.True(strings.HasPrefix(brand.ID, "brd_"))
	assert.Equal(color, brand.Color)
}

func (c *ClientTests) TestUserDelete() {
	c.T().Skip("Due to our inability to create child users securely, we must also skip deleting them as we cannot replace the deleted ones easily.")
	client := c.ProdClient()
	require := c.Require()

	userName := "Test User"

	user, _ := client.CreateUser(
		&easypost.UserOptions{
			Name: &userName,
		},
	)

	err := client.DeleteUser(user.ID)

	require.NoError(err)
}

func (c *ClientTests) TestUserApiKeysAll() {
	c.T().Skip("API keys are returned as plaintext, do not run this test.")
	client := c.ProdClient()
	assert := c.Assert()

	apiKeys, _ := client.GetAPIKeys()

	assert.NotNil(apiKeys)
}
