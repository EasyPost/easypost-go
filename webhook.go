package easypost

import (
	"context"
	"time"
)

// A Webhook represents an EasyPost webhook callback URL.
type Webhook struct {
	ID         string     `json:"id,omitempty"`
	Object     string     `json:"object,omitempty"`
	Mode       string     `json:"mode,omitempty"`
	URL        string     `json:"url,omitempty"`
	DisabledAt *time.Time `json:"disabled_at,omitempty"`
}

type createWebhookRequest struct {
	Webhook *Webhook `json:"webhook,omitempty"`
}

// CreateWebhook creates a new webhook with the given URL.
func (c *Client) CreateWebhook(u string) (out *Webhook, err error) {
	req := &createWebhookRequest{Webhook: &Webhook{URL: u}}
	err = c.post(context.Background(), "webhooks", req, &out)
	return
}

// CreateWebhookWithContext performs the same operation as CreateWebhook, but
// allows specifying a context that can interrupt the request.
func (c *Client) CreateWebhookWithContext(ctx context.Context, u string) (out *Webhook, err error) {
	req := &createWebhookRequest{Webhook: &Webhook{URL: u}}
	err = c.post(ctx, "webhooks", req, &out)
	return
}

type listWebhooksResult struct {
	Webhooks *[]*Webhook `json:"webhooks,omitempty"`
}

// ListWebhooks returns all webhooks associated with the EasyPost account.
func (c *Client) ListWebhooks() (out []*Webhook, err error) {
	err = c.get(context.Background(), "webhooks", &listWebhooksResult{Webhooks: &out})
	return
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
	err = c.get(context.Background(), "webhooks/"+webhookID, &out)
	return
}

// GetWebhookWithContext performs the same operation as GetWebhook, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetWebhookWithContext(ctx context.Context, webhookID string) (out *Webhook, err error) {
	err = c.get(ctx, "webhooks/"+webhookID, &out)
	return
}

// EnableWebhook re-enables a disabled webhook.
func (c *Client) EnableWebhook(webhookID string) (out *Webhook, err error) {
	err = c.put(context.Background(), "webhooks/"+webhookID, nil, &out)
	return
}

// EnableWebhookWithContext performs the same operation as EnableWebhook, but
// allows specifying a context that can interrupt the request.
func (c *Client) EnableWebhookWithContext(ctx context.Context, webhookID string) (out *Webhook, err error) {
	err = c.put(ctx, "webhooks/"+webhookID, nil, &out)
	return
}

// DeleteWebhook removes a webhook.
func (c *Client) DeleteWebhook(webhookID string) error {
	return c.del(context.Background(), "webhooks/"+webhookID)
}

// DeleteWebhookWithContext performs the same operation as DeleteWebhook, but
// allows specifying a context that can interrupt the request.
func (c *Client) DeleteWebhookWithContext(ctx context.Context, webhookID string) error {
	return c.del(ctx, "webhooks/"+webhookID)
}
