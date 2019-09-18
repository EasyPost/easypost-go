package easypost_test

import (
	"net/http"
	"testing"

	"github.com/EasyPost/easypost-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhooks(t *testing.T) {
	assert, require := assert.New(t), require.New(t)
	url := "example.com/" + unique()
	webhook, err := TestClient.CreateWebhook(url)
	require.NoError(err)
	assert.NotEmpty(webhook.ID)
	assert.Equal("test", webhook.Mode)
	assert.Equal("http://"+url, webhook.URL)
	assert.Empty(webhook.DisabledAt)

	webhook2, err := TestClient.GetWebhook(webhook.ID)
	require.NoError(err)
	assert.Equal(webhook.ID, webhook2.ID)

	webhook3, err := TestClient.EnableWebhook(webhook.ID)
	require.NoError(err)
	require.Equal(webhook.ID, webhook3.ID)

	webhooks, err := TestClient.ListWebhooks()
	require.NoError(err)
	ids := make([]string, len(webhooks))
	for i := range webhooks {
		ids[i] = webhooks[i].ID
	}
	assert.Contains(ids, webhook.ID)

	require.NoError(TestClient.DeleteWebhook(webhook.ID))

	_, err = TestClient.GetWebhook(webhook.ID)
	if assert.IsType((*easypost.APIError)(nil), err) {
		err := err.(*easypost.APIError)
		assert.Equal(err.StatusCode, http.StatusNotFound)
	}
}
