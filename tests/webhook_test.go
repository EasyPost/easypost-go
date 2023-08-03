package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v3"
)

func (c *ClientTests) TestWebhookCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, err := client.CreateWebhookWithDetails(
		&easypost.CreateUpdateWebhookOptions{
			URL: c.fixture.WebhookUrl(),
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(webhook))
	assert.True(strings.HasPrefix(webhook.ID, "hook_"))
	assert.Equal(c.fixture.WebhookUrl(), webhook.URL)

	err = client.DeleteWebhook(webhook.ID) // we are deleting the webhook here, so we don't keep sending events to a dead webhook.
	require.NoError(err)

	require.NoError(err)
}

func (c *ClientTests) TestWebhookRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, err := client.CreateWebhookWithDetails(
		&easypost.CreateUpdateWebhookOptions{
			URL: c.fixture.WebhookUrl(),
		},
	)
	require.NoError(err)

	retrievedWebhook, err := client.GetWebhook(webhook.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(retrievedWebhook))
	assert.Equal(webhook, retrievedWebhook)

	err = client.DeleteWebhook(retrievedWebhook.ID) // we are deleting the webhook here, so we don't keep sending events to a dead webhook.
	require.NoError(err)

	require.NoError(err)
}

func (c *ClientTests) TestWebhookAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhooks, err := client.ListWebhooks()
	require.NoError(err)

	assert.LessOrEqual(len(webhooks), c.fixture.pageSize())
	for _, webhook := range webhooks {
		assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(webhook))
	}
}

func (c *ClientTests) TestWebhookDelete() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, err := client.CreateWebhookWithDetails(
		&easypost.CreateUpdateWebhookOptions{
			URL: c.fixture.WebhookUrl(),
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(webhook))
	assert.True(strings.HasPrefix(webhook.ID, "hook_"))
	assert.Equal(c.fixture.WebhookUrl(), webhook.URL)

	err = client.DeleteWebhook(webhook.ID)
	require.NoError(err)

	// This endpoint/method does not return anything, just make sure the request doesn't fail
	require.NoError(err)
}

func (c *ClientTests) TestWebhookUpdate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, err := client.CreateWebhookWithDetails(
		&easypost.CreateUpdateWebhookOptions{
			URL: c.fixture.WebhookUrl(),
		},
	)
	require.NoError(err)

	updatedWebhook, err := client.UpdateWebhook(webhook.ID, &easypost.CreateUpdateWebhookOptions{})
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(updatedWebhook))

	err = client.DeleteWebhook(updatedWebhook.ID)
	require.NoError(err)
}

func (c *ClientTests) TestWebhookCreateWithSecret() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, err := client.CreateWebhookWithDetails(
		&easypost.CreateUpdateWebhookOptions{
			URL:           c.fixture.WebhookUrl(),
			WebhookSecret: "123",
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(webhook))
	assert.True(strings.HasPrefix(webhook.ID, "hook_"))
	assert.Equal(c.fixture.WebhookUrl(), webhook.URL)

	err = client.DeleteWebhook(webhook.ID) // we are deleting the webhook here, so we don't keep sending events to a dead webhook.
	require.NoError(err)
}

func (c *ClientTests) TestWebhookUpdateWithSecret() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhook, err := client.CreateWebhookWithDetails(
		&easypost.CreateUpdateWebhookOptions{
			URL: c.fixture.WebhookUrl(),
		},
	)
	require.NoError(err)

	updatedWebhook, err := client.UpdateWebhook(webhook.ID,
		&easypost.CreateUpdateWebhookOptions{
			WebhookSecret: "123",
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Webhook{}), reflect.TypeOf(updatedWebhook))

	err = client.DeleteWebhook(updatedWebhook.ID)
	require.NoError(err)
}

func (c *ClientTests) TestValidateWebhook() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	webhookSecret := "sécret"
	expectedHmacSignature := "hmac-sha256-hex=e93977c8ccb20363d51a62b3fe1fc402b7829be1152da9e88cf9e8d07115a46b"
	headers := map[string]interface{}{
		"X-Hmac-Signature": expectedHmacSignature,
	}

	webhookBody, err := client.ValidateWebhook(c.fixture.EventBody(), headers, webhookSecret)
	require.NoError(err)

	assert.Equal(webhookBody.Description, "batch.created")
}

func (c *ClientTests) TestValidateWebhookInvalidSecret() {
	client := c.TestClient()
	assert := c.Assert()

	webhookSecret := "invalid_secret"
	headers := map[string]interface{}{
		"X-Hmac-Signature": "some-signature",
	}

	_, err := client.ValidateWebhook(c.fixture.EventBody(), headers, webhookSecret)
	assert.Error(err)
}

func (c *ClientTests) TestValidateWebhookMissingSecret() {
	client := c.TestClient()
	assert := c.Assert()

	webhookSecret := "123"
	headers := map[string]interface{}{
		"some-header": "some-value",
	}

	_, err := client.ValidateWebhook(c.fixture.EventBody(), headers, webhookSecret)
	assert.Error(err)
}
