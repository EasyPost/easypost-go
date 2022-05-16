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

	user, err := client.CreateUser(
		&easypost.UserOptions{
			Name: &userName,
		},
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(user))
	assert.True(strings.HasPrefix(user.ID, "user_"))
	assert.Equal("Test User", user.Name)
}

func (c *ClientTests) TestUserRetrieve() {
	client := c.ProdClient()
	assert := c.Assert()

	authenticatedUser, err := client.RetrieveMe()
	if err != nil {
		c.T().Error(err)
		return
	}

	retrievedUser, err := client.GetUser(authenticatedUser.Children[0].ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(retrievedUser))
	assert.True(strings.HasPrefix(retrievedUser.ID, "user_"))
}

func (c *ClientTests) TestUserRetrieveMe() {
	client := c.ProdClient()
	assert := c.Assert()

	user, err := client.RetrieveMe()
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(user))
	assert.True(strings.HasPrefix(user.ID, "user_"))
}

func (c *ClientTests) TestUserUpdate() {
	client := c.ProdClient()
	assert := c.Assert()

	user, err := client.RetrieveMe()
	if err != nil {
		c.T().Error(err)
		return
	}

	test_phone := "5555555555"

	updatedUser, err := client.UpdateUser(
		&easypost.UserOptions{
			ID:          user.ID,
			PhoneNumber: &test_phone,
		},
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(updatedUser))
	assert.True(strings.HasPrefix(updatedUser.ID, "user_"))
	assert.Equal(test_phone, updatedUser.PhoneNumber)
}

func (c *ClientTests) TestUserUpdateBrand() {
	client := c.ProdClient()
	assert := c.Assert()

	color := "#123456"

	user, err := client.RetrieveMe()
	if err != nil {
		c.T().Error(err)
		return
	}

	brand, err := client.UpdateBrand(
		map[string]interface{}{
			"color": color,
		},
		user.ID,
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Brand{}), reflect.TypeOf(brand))
	assert.True(strings.HasPrefix(brand.ID, "brd_"))
	assert.Equal(color, brand.Color)
}

func (c *ClientTests) TestUserDelete() {
	c.T().Skip("Due to our inability to create child users securely, we must also skip deleting them as we cannot replace the deleted ones easily.")
	client := c.ProdClient()
	require := c.Require()

	userName := "Test User"

	user, err := client.CreateUser(
		&easypost.UserOptions{
			Name: &userName,
		},
	)
	if err != nil {
		c.T().Error(err)
		return
	}

	err = client.DeleteUser(user.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	require.NoError(err)
}

func (c *ClientTests) TestUserApiKeysAll() {
	c.T().Skip("API keys are returned as plaintext, do not run this test.")
	client := c.ProdClient()
	assert := c.Assert()

	apiKeys, err := client.GetAPIKeys()
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.NotNil(apiKeys)
}
