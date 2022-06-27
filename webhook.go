package easypost

import (
	"context"
	"time"
)

// A Webhook represents an EasyPost webhook callback URL.
type Webhook struct {
	ID            string     `json:"id,omitempty"`
	Object        string     `json:"object,omitempty"`
	Mode          string     `json:"mode,omitempty"`
	URL           string     `json:"url,omitempty"`
	DisabledAt    *time.Time `json:"disabled_at,omitempty"`
	WebhookSecret string     `json:"webhook_secret,omitempty"`
}

// UpdateWebhookOptions is used to specify parameters for updating EasyPost webhooks.
type UpdateWebhookOptions struct {
	URL           string `json:"url,omitempty"`
	WebhookSecret string `json:"webhook_secret,omitempty"`
}

type createUpdateWebhookRequest struct {
	Webhook *Webhook `json:"webhook,omitempty"`
}

type listWebhooksResult struct {
	Webhooks *[]*Webhook `json:"webhooks,omitempty"`
}

// CreateWebhook creates a new webhook with the given URL.
func (c *Client) CreateWebhook(u string) (out *Webhook, err error) {
	return c.CreateWebhookWithContext(context.Background(), u)
}

// CreateWebhookWithSecret creates a new webhook with the given URL and secret.
func (c *Client) CreateWebhookWithSecret(u string, s string) (out *Webhook, err error) {
	return c.CreateWebhookWithSecretWithContext(context.Background(), u, s)
}

// CreateWebhookWithContext performs the same operation as CreateWebhook, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateWebhookWithContext(ctx context.Context, u string) (out *Webhook, err error) {
	req := &createUpdateWebhookRequest{Webhook: &Webhook{URL: u}}
	err = c.post(ctx, "webhooks", req, &out)
	return
}

// CreateWebhookWithSecretWithContext performs the same operation as CreateWebhookWithSecret, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateWebhookWithSecretWithContext(ctx context.Context, u string, s string) (out *Webhook, err error) {
	req := &createUpdateWebhookRequest{Webhook: &Webhook{URL: u, WebhookSecret: s}}
	err = c.post(ctx, "webhooks", req, &out)
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
	err = c.get(ctx, "webhooks", &listWebhooksResult{Webhooks: &out})
	return
}

// GetWebhook retrieves a Webhook object with the given ID.
func (c *Client) GetWebhook(webhookID string) (out *Webhook, err error) {
	return c.GetWebhookWithContext(context.Background(), webhookID)
}

// GetWebhookWithContext performs the same operation as GetWebhook, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetWebhookWithContext(ctx context.Context, webhookID string) (out *Webhook, err error) {
	err = c.get(ctx, "webhooks/"+webhookID, &out)
	return
}

// Deprecated: Use UpdateWebhook instead.
// EnableWebhook re-enables a disabled webhook.
func (c *Client) EnableWebhook(webhookID string) (out *Webhook, err error) {
	return c.UpdateWebhook(webhookID, &UpdateWebhookOptions{})
}

// Deprecated: Use UpdateWebhookWithContext instead.
// EnableWebhookWithContext performs the same operation as EnableWebhook, but
// allows specifying a context that can interrupt the request.
func (c *Client) EnableWebhookWithContext(ctx context.Context, webhookID string) (out *Webhook, err error) {
	return c.UpdateWebhookWithContext(ctx, webhookID, &UpdateWebhookOptions{})
}

// UpdateWebhook updates a webhook. Automatically re-enables webhook if it is disabled.
func (c *Client) UpdateWebhook(webhookID string, data *UpdateWebhookOptions) (out *Webhook, err error) {
	return c.UpdateWebhookWithContext(context.Background(), webhookID, data)
}

// UpdateWebhookWithContext performs the same operation as UpdateWebhook, but
// allows specifying a context that can interrupt the request.
func (c *Client) UpdateWebhookWithContext(ctx context.Context, webhookID string, data *UpdateWebhookOptions) (out *Webhook, err error) {
	req := &createUpdateWebhookRequest{Webhook: &Webhook{URL: data.URL, WebhookSecret: data.WebhookSecret}}
	err = c.patch(ctx, "webhooks/"+webhookID, req, &out)
	return
}

// DeleteWebhook removes a webhook.
func (c *Client) DeleteWebhook(webhookID string) error {
	return c.DeleteWebhookWithContext(context.Background(), webhookID)
}

// DeleteWebhookWithContext performs the same operation as DeleteWebhook, but
// allows specifying a context that can interrupt the request.
func (c *Client) DeleteWebhookWithContext(ctx context.Context, webhookID string) error {
	return c.del(ctx, "webhooks/"+webhookID)
}
