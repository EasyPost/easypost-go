package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestWebhookCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, err := client.CreateWebhook(c.fixture.WebhookUrl())
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(webhook))
	assert.True(strings.HasPrefix(webhook.ID, "hook_"))
	assert.Equal(c.fixture.WebhookUrl(), webhook.URL)

	err = client.DeleteWebhook(webhook.ID) // we are deleting the webhook here so we don't keep sending events to a dead webhook.
	if err != nil {
		c.T().Error(err)
		return
	}

	require.NoError(err)
}

func (c *ClientTests) TestWebhookRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, err := client.CreateWebhook(c.fixture.WebhookUrl())
	if err != nil {
		c.T().Error(err)
		return
	}

	retrievedWebhook, err := client.GetWebhook(webhook.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(retrievedWebhook))
	assert.Equal(webhook, retrievedWebhook)

	err = client.DeleteWebhook(retrievedWebhook.ID) // we are deleting the webhook here, so we don't keep sending events to a dead webhook.
	if err != nil {
		c.T().Error(err)
		return
	}

	require.NoError(err)
}

func (c *ClientTests) TestWebhookAll() {
	client := c.TestClient()
	assert := c.Assert()

	webhooks, err := client.ListWebhooks()
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.LessOrEqual(len(webhooks), c.fixture.pageSize())
	for _, webhook := range webhooks {
		assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(webhook))
	}
}

func (c *ClientTests) TestWebhookDelete() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, err := client.CreateWebhook(c.fixture.WebhookUrl())
	if err != nil {
		c.T().Error(err)
		return
	}

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(webhook))
	assert.True(strings.HasPrefix(webhook.ID, "hook_"))
	assert.Equal(c.fixture.WebhookUrl(), webhook.URL)

	err = client.DeleteWebhook(webhook.ID)
	if err != nil {
		c.T().Error(err)
		return
	}

	// This endpoint/method does not return anything, just make sure the request doesn't fail
	require.NoError(err)
}

func (c *ClientTests) TestWebhookUpdate() {
	c.T().Skip("Cannot be easily tested - requires a disabled webhook.")
}
