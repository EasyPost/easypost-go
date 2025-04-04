package easypost

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"golang.org/x/text/unicode/norm"
)

// WebhookCustomHeader represents a custom header for an EasyPost webhook.
type WebhookCustomHeader struct {
	Name  string `json:"name" url:"name"`
	Value string `json:"value" url:"value"`
}

// A Webhook represents an EasyPost webhook callback URL.
type Webhook struct {
	ID            string                `json:"id,omitempty" url:"id,omitempty"`
	Object        string                `json:"object,omitempty" url:"object,omitempty"`
	Mode          string                `json:"mode,omitempty" url:"mode,omitempty"`
	URL           string                `json:"url,omitempty" url:"url,omitempty"`
	DisabledAt    *DateTime             `json:"disabled_at,omitempty" url:"disabled_at,omitempty"`
	WebhookSecret string                `json:"webhook_secret,omitempty" url:"webhook_secret,omitempty"`
	CustomHeaders []WebhookCustomHeader `json:"custom_headers,omitempty" url:"custom_headers,omitempty"`
}

// CreateUpdateWebhookOptions is used to specify parameters for creating and updating EasyPost webhooks.
type CreateUpdateWebhookOptions struct {
	URL           string                `json:"url,omitempty" url:"url,omitempty"`
	WebhookSecret string                `json:"webhook_secret,omitempty" url:"webhook_secret,omitempty"`
	CustomHeaders []WebhookCustomHeader `json:"custom_headers,omitempty" url:"custom_headers,omitempty"`
}

// createWebhookRequest is the request struct for creating webhooks (internal)
type createWebhookRequest struct {
	Webhook *Webhook `json:"webhook,omitempty" url:"webhook,omitempty"`
}

// updateWebhookRequest is the request struct for updating webhooks (internal)
type updateWebhookRequest struct {
	WebhookSecret string                `json:"webhook_secret,omitempty" url:"webhook_secret,omitempty"`
	CustomHeaders []WebhookCustomHeader `json:"custom_headers,omitempty" url:"custom_headers,omitempty"`
}

// listWebhooksResult is the response struct of listing webhooks (internal)
type listWebhooksResult struct {
	Webhooks *[]*Webhook `json:"webhooks,omitempty" url:"webhooks,omitempty"`
}

func (c *Client) composeCreateWebhookRequest(data *CreateUpdateWebhookOptions) *createWebhookRequest {
	return &createWebhookRequest{
		Webhook: &Webhook{
			URL:           data.URL,
			WebhookSecret: data.WebhookSecret,
			CustomHeaders: data.CustomHeaders,
		},
	}
}

func (c *Client) composeUpdateWebhookRequest(data *CreateUpdateWebhookOptions) *updateWebhookRequest {
	return &updateWebhookRequest{
		WebhookSecret: data.WebhookSecret,
		CustomHeaders: data.CustomHeaders,
	}
}

// CreateWebhook creates a new webhook with the provided details.
func (c *Client) CreateWebhook(data *CreateUpdateWebhookOptions) (out *Webhook, err error) {
	return c.CreateWebhookWithContext(context.Background(), data)
}

// CreateWebhookWithContext performs the same operation as CreateWebhook, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateWebhookWithContext(ctx context.Context,
	data *CreateUpdateWebhookOptions) (out *Webhook, err error) {
	req := c.composeCreateWebhookRequest(data)
	err = c.do(ctx, http.MethodPost, "webhooks", req, &out)
	return
}

// ListWebhooks returns all webhooks associated with the EasyPost account.
func (c *Client) ListWebhooks() (out []*Webhook, err error) {
	return c.ListWebhooksWithContext(context.Background())
}

// ListWebhooksWithContext performs the same operation as
// ListWebhooksWithContext, but allows specifying a context that can interrupt
// the request.
func (c *Client) ListWebhooksWithContext(ctx context.Context) (out []*Webhook, err error) {
	err = c.do(ctx, http.MethodGet, "webhooks", nil, &listWebhooksResult{Webhooks: &out})
	return
}

// GetWebhook retrieves a Webhook object with the given ID.
func (c *Client) GetWebhook(webhookID string) (out *Webhook, err error) {
	return c.GetWebhookWithContext(context.Background(), webhookID)
}

// GetWebhookWithContext performs the same operation as GetWebhook, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetWebhookWithContext(ctx context.Context, webhookID string) (out *Webhook, err error) {
	err = c.do(ctx, http.MethodGet, "webhooks/"+webhookID, nil, &out)
	return
}

// UpdateWebhook updates a webhook. Automatically re-enables webhook if it is disabled.
func (c *Client) UpdateWebhook(webhookID string, data *CreateUpdateWebhookOptions) (out *Webhook, err error) {
	return c.UpdateWebhookWithContext(context.Background(), webhookID, data)
}

// UpdateWebhookWithContext performs the same operation as UpdateWebhook, but
// allows specifying a context that can interrupt the request.
func (c *Client) UpdateWebhookWithContext(ctx context.Context,
	webhookID string, data *CreateUpdateWebhookOptions) (out *Webhook, err error) {
	req := c.composeUpdateWebhookRequest(data)
	err = c.do(ctx, http.MethodPatch, "webhooks/"+webhookID, req, &out)
	return
}

// DeleteWebhook removes a webhook.
func (c *Client) DeleteWebhook(webhookID string) error {
	return c.DeleteWebhookWithContext(context.Background(), webhookID)
}

// DeleteWebhookWithContext performs the same operation as DeleteWebhook, but
// allows specifying a context that can interrupt the request.
func (c *Client) DeleteWebhookWithContext(ctx context.Context, webhookID string) error {
	return c.do(ctx, http.MethodDelete, "webhooks/"+webhookID, nil, nil)
}

// ValidateWebhook validates a webhook by comparing the HMAC signature header sent from EasyPost to your shared secret.
// If the signatures do not match, an error will be raised signifying the webhook either did not originate
// from EasyPost or the secrets do not match. If the signatures do match, the `event_body` will be returned as JSON.
func (c *Client) ValidateWebhook(eventBody []byte, headers map[string]interface{}, webhookSecret string) (out *Event, err error) {
	return c.ValidateWebhookWithContext(context.Background(), eventBody, headers, webhookSecret)
}

// ValidateWebhookWithContext performs the same operation as ValidateWebhook, but
// allows specifying a context that can interrupt the request.
func (c *Client) ValidateWebhookWithContext(ctx context.Context, eventBody []byte, headers map[string]interface{}, webhookSecret string) (out *Event, err error) {
	easypostHmacSignature, signaturePresent := headers["X-Hmac-Signature"].(string)

	if signaturePresent {
		normalizedSecret := norm.NFKD.String(webhookSecret)
		encodedSecret := []byte(normalizedSecret)

		expectedSignature := hmac.New(sha256.New, encodedSecret)
		expectedSignature.Write(eventBody)

		digest := "hmac-sha256-hex=" + hex.EncodeToString(expectedSignature.Sum(nil))

		if hmac.Equal([]byte(digest), []byte(easypostHmacSignature)) {
			var webhookBody Event
			if err = json.Unmarshal(eventBody, &webhookBody); err != nil {
				return nil, err
			} else {
				return &webhookBody, nil
			}
		} else {
			return nil, newMismatchWebhookSignatureError()
		}
	} else {
		return nil, newMissingWebhookSignatureError()
	}
}
