package easypost_test

import (
	"encoding/json"
	"testing"

	"github.com/EasyPost/easypost-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestAPIError(t *testing.T) {
	assert := assert.New(t)
	buf := []byte(`
{
	"error": {
		"code":    "ADDRESS.VERIFY.FAILURE",
		"message": "Unable to verify address.",
		"errors":  [
			{
				"code":       "E.STREET.INVALID",
				"field":      "street1",
				"message":    "Street is invalid",
				"suggestion": "Larkspur Cres"
			},
			{
				"code":       "E.HOUSE_NUMBER.INVALID",
				"field":      "street1",
				"message":    "House number is invalid",
				"suggestion": null
			}
		]
	}
}`,
	)

	apiErr := &easypost.APIError{
		Status: "Unprocessable Entity", StatusCode: 422,
	}
	type apiErrorResponse struct {
		Error *easypost.APIError `json:"error,omitempty"`
	}
	assert.NoError(json.Unmarshal(buf, &apiErrorResponse{Error: apiErr}))
	assert.Equal("Unprocessable Entity", apiErr.Status)
	assert.Equal(422, apiErr.StatusCode)
	assert.Equal("ADDRESS.VERIFY.FAILURE", apiErr.Code)
	assert.Equal("Unable to verify address.", apiErr.Message)
	assert.Empty(apiErr.Field)
	assert.Empty(apiErr.Suggestion)
	if assert.Len(apiErr.Errors, 2) {
		assert.Empty(apiErr.Errors[0].Status)
		assert.Empty(apiErr.Errors[0].StatusCode)
		assert.Equal("E.STREET.INVALID", apiErr.Errors[0].Code)
		assert.Equal("Street is invalid", apiErr.Errors[0].Message)
		assert.Equal("street1", apiErr.Errors[0].Field)
		assert.Equal("Larkspur Cres", apiErr.Errors[0].Suggestion)
		assert.Empty(apiErr.Errors[0].Errors)
		assert.Empty(apiErr.Errors[1].Status)
		assert.Empty(apiErr.Errors[1].StatusCode)
		assert.Equal("E.HOUSE_NUMBER.INVALID", apiErr.Errors[1].Code)
		assert.Equal("House number is invalid", apiErr.Errors[1].Message)
		assert.Equal("street1", apiErr.Errors[1].Field)
		assert.Empty(apiErr.Errors[1].Suggestion)
		assert.Empty(apiErr.Errors[1].Errors)
	}
}
