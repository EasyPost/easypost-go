package easypost_test

import (
	"encoding/json"

	"github.com/EasyPost/easypost-go/v2"
	"github.com/dnaeon/go-vcr/cassette"
)

func cleanAPIKeys(data interface{}) {
	switch data := data.(type) {
	case map[string]interface{}:
		for k, v := range data {
			switch v := v.(type) {
			case string:
				if k == "key" {
					data[k] = "XXX"
				}
			case []interface{}, map[string]interface{}:
				cleanAPIKeys(v)
			}
		}
	case []interface{}:
		for i := range data {
			switch v := data[i].(type) {
			case []interface{}, map[string]interface{}:
				cleanAPIKeys(v)
			}
		}
	}
}

func (c *ClientTests) TestChildUserCreate() {
	c.recorder.AddFilter(func(i *cassette.Interaction) error {
		// filter out api keys.
		var data map[string]interface{}
		if json.Unmarshal([]byte(i.Response.Body), &data) != nil {
			return nil
		}
		cleanAPIKeys(data)
		buf, _ := json.Marshal(data)
		i.Response.Body = string(buf)
		return nil
	})

	client := c.ProdClient()
	assert := c.Assert()
	opts := &easypost.UserOptions{Name: easypost.StringPtr("Go Test User")}
	user, err := client.CreateUser(opts)
	if assert.NoError(err) {
		assert.NotEmpty(user.ID)
		if user, err = client.GetUser(user.ID); assert.NoError(err) {
			assert.Equal("Go Test User", user.Name)
		}
		*opts.Name = "Go Testing User"
		opts.ID = user.ID
		if user, err = client.UpdateUser(opts); assert.NoError(err) {
			if user, err = client.GetUser(user.ID); assert.NoError(err) {
				assert.Equal("Go Testing User", user.Name)
			}
		}
		assert.NoError(client.DeleteUser(user.ID))
	}

	opts.Password = easypost.StringPtr("super-secret-password")
	opts.PasswordConfirmation = opts.Password
	_, err = client.CreateUser(opts)
	if assert.IsType((*easypost.APIError)(nil), err) {
		err := err.(*easypost.APIError)
		assert.Equal(
			"Child users must be created via the API and must not contain an email or password.",
			err.Message,
		)
		assert.Equal("USER.INVALID", err.Code)
	}
}
