package easypost_test

import (
	"net/http"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestWebhooks() {
	assert, require := c.Assert(), c.Require()
	client := c.TestClient()
	// May have to update this if re-recording the casset for this test:
	url := "example.com/2020050722012281"
	webhook, err := client.CreateWebhook(url)
	require.NoError(err)
	assert.NotEmpty(webhook.ID)
	assert.Equal("test", webhook.Mode)
	assert.Equal("http://"+url, webhook.URL)
	assert.Empty(webhook.DisabledAt)

	webhook2, err := client.GetWebhook(webhook.ID)
	require.NoError(err)
	assert.Equal(webhook.ID, webhook2.ID)

	webhook3, err := client.EnableWebhook(webhook.ID)
	require.NoError(err)
	require.Equal(webhook.ID, webhook3.ID)

	webhooks, err := client.ListWebhooks()
	require.NoError(err)
	ids := make([]string, len(webhooks))
	for i := range webhooks {
		ids[i] = webhooks[i].ID
	}
	assert.Contains(ids, webhook.ID)

	require.NoError(client.DeleteWebhook(webhook.ID))

	_, err = client.GetWebhook(webhook.ID)
	if assert.IsType((*easypost.APIError)(nil), err) {
		err := err.(*easypost.APIError)
		assert.Equal(err.StatusCode, http.StatusNotFound)
	}
}
