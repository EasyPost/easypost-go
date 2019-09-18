package easypost_test

import (
	"testing"

	"github.com/EasyPost/easypost-go"
	"github.com/stretchr/testify/assert"
)

func TestChildUserCreate(t *testing.T) {
	if ProdClient.APIKey == "" {
		t.Skip("no production API key")
	}
	assert := assert.New(t)
	opts := &easypost.UserOptions{Name: easypost.StringPtr("Go Test User")}
	user, err := ProdClient.CreateUser(opts)
	if assert.NoError(err) {
		assert.NotEmpty(user.ID)
		if user, err = ProdClient.GetUser(user.ID); assert.NoError(err) {
			assert.Equal("Go Test User", user.Name)
		}
		*opts.Name = "Go Testing User"
		if user, err = ProdClient.UpdateUser(opts); assert.NoError(err) {
			if user, err = ProdClient.GetUser(user.ID); assert.NoError(err) {
				assert.Equal("Go Testing User", user.Name)
			}
		}
	}

	opts.Password = easypost.StringPtr("super-secret-password")
	opts.PasswordConfirmation = opts.Password
	_, err = ProdClient.CreateUser(opts)
	if assert.IsType((*easypost.APIError)(nil), err) {
		err := err.(*easypost.APIError)
		assert.Equal(
			"Child users must be created via the API and must not contain an email or password.",
			err.Message,
		)
		assert.Equal("USER.INVALID", err.Code)
	}
}
