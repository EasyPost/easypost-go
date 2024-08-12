package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v4"
)

func (c *ClientTests) TestUserCreate() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	userName := "Test User"

	user, err := client.CreateUser(
		&easypost.UserOptions{
			Name: &userName,
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(user))
	assert.True(strings.HasPrefix(user.ID, "user_"))
	assert.Equal("Test User", user.Name)
}

func (c *ClientTests) TestUserRetrieve() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	authenticatedUser, err := client.RetrieveMe()
	require.NoError(err)

	retrievedUser, err := client.GetUser(authenticatedUser.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(retrievedUser))
	assert.True(strings.HasPrefix(retrievedUser.ID, "user_"))
}

func (c *ClientTests) TestUserRetrieveMe() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	user, err := client.RetrieveMe()
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(user))
	assert.True(strings.HasPrefix(user.ID, "user_"))
}

func (c *ClientTests) TestUserAllChildUsers() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	childUsers, err := client.ListChildUsers(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	assert.LessOrEqual(len(childUsers.Children), c.fixture.pageSize())
	assert.NotNil(childUsers.HasMore)
	for _, user := range childUsers.Children {
		assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(user))
	}
}

func (c *ClientTests) TestUserGetNextChildUserPage() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListChildUsers(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextChildUserPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Children) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Children[len(firstPage.Children)-1].ID
			firstIdOfSecondPage := nextPage.Children[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIdOfSecondPage)
		}
	}()
	if err != nil {
		assert.Equal(err.Error(), easypost.NoPagesLeftToRetrieve)
		return
	}
}

func (c *ClientTests) TestUserUpdate() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	user, err := client.RetrieveMe()
	require.NoError(err)

	test_name := "Test User"

	updatedUser, err := client.UpdateUser(
		&easypost.UserOptions{
			ID:   user.ID,
			Name: &test_name,
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.User{}), reflect.TypeOf(updatedUser))
	assert.True(strings.HasPrefix(updatedUser.ID, "user_"))
	assert.Equal(test_name, updatedUser.Name)
}

func (c *ClientTests) TestUserUpdateBrand() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()

	color := "#123456"

	user, err := client.RetrieveMe()
	require.NoError(err)

	brand, err := client.UpdateBrand(
		map[string]interface{}{
			"color": color,
		},
		user.ID,
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Brand{}), reflect.TypeOf(brand))
	assert.True(strings.HasPrefix(brand.ID, "brd_"))
	assert.Equal(color, brand.Color)
}

func (c *ClientTests) TestUserDelete() {
	client := c.ProdClient()
	require := c.Require()

	userName := "Test User"

	user, err := client.CreateUser(
		&easypost.UserOptions{
			Name: &userName,
		},
	)
	require.NoError(err)

	err = client.DeleteUser(user.ID)
	require.NoError(err)

	require.NoError(err)
}
