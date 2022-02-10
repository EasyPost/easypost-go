package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestCarrierTypes() {
	client := c.ProdClient()
	assert, require := c.Assert(), c.Require()
	carriers, err := client.GetCarrierTypes()
	require.NoError(err)

	var carrier *easypost.CarrierType
	for i := range carriers {
		if carriers[i].Type == "AustraliaPostAccount" {
			carrier = carriers[i]
			break
		}
	}
	require.NotNil(carrier)
	if assert.NotNil(carrier.Fields) {
		if assert.Contains(carrier.Fields.Credentials, "api_key") {
			assert.Equal(
				"Australia Post API Key",
				carrier.Fields.Credentials["api_key"].Label,
			)
		}
		if assert.Contains(carrier.Fields.Credentials, "api_secret") {
			assert.Equal(
				"Australia Post Secret Key",
				carrier.Fields.Credentials["api_secret"].Label,
			)
		}
		if assert.Contains(carrier.Fields.Credentials, "account_number") {
			assert.Equal(
				"Australia Post Account Number",
				carrier.Fields.Credentials["account_number"].Label,
			)
		}
		if assert.Contains(carrier.Fields.Credentials, "print_as_you_go") {
			assert.Equal(
				"Print as you go",
				carrier.Fields.Credentials["print_as_you_go"].Label,
			)
		}
	}
}
