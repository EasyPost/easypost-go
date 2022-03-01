package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestWebhookCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, _ := client.CreateWebhook("http://example.com")

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(webhook))
	assert.True(strings.HasPrefix(webhook.ID, "hook_"))
	assert.Equal("http://example.com", webhook.URL)

	err := client.DeleteWebhook(webhook.ID) // we are deleting the webhook here so we don't keep sending events to a dead webhook.

	require.NoError(err)
}

func (c *ClientTests) TestWebhookRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, _ := client.CreateWebhook("http://example.com")

	retrievedWebhook, _ := client.GetWebhook(webhook.ID)

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(retrievedWebhook))
	assert.Equal(webhook, retrievedWebhook)

	err := client.DeleteWebhook(retrievedWebhook.ID) // we are deleting the webhook here, so we don't keep sending events to a dead webhook.

	require.NoError(err)
}

func (c *ClientTests) TestWebhookAll() {
	client := c.TestClient()
	assert := c.Assert()

	webhooks, _ := client.ListWebhooks()

	assert.LessOrEqual(len(webhooks), c.fixture.pageSize())
	for _, webhook := range webhooks {
		assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(webhook))
	}
}

func (c *ClientTests) TestWebhookDelete() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, _ := client.CreateWebhook("http://example.com")

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(webhook))
	assert.True(strings.HasPrefix(webhook.ID, "hook_"))
	assert.Equal("http://example.com", webhook.URL)

	err := client.DeleteWebhook(webhook.ID)

	// This endpoint/method does not return anything, just make sure the request doesn't fail
	require.NoError(err)
}

func (c *ClientTests) TestWebhookUpdate() {
	c.T().Skip("Cannot be easily tested - requires a disabled webhook.")
}
