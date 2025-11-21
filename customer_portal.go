package easypost

import (
	"context"
	"net/http"
)

// A CustomerPortalAccountLink object represents an object containing an AccountLink.
type CustomerPortalAccountLink struct {
	Object    string    `json:"object,omitempty" url:"object,omitempty"`
	Link      string    `json:"link,omitempty" url:"link,omitempty"`
	CreatedAt *DateTime `json:"created_at,omitempty" url:"created_at,omitempty"`
	ExpiresAt *DateTime `json:"expires_at,omitempty" url:"expires_at,omitempty"`
}

// CustomerPortalAccountLinkParameters is used to specify parameters for creating an AccountLink.
type CustomerPortalAccountLinkParameters struct {
	SessionType string                 `json:"session_type,omitempty" url:"session_type,omitempty"`
	UserId      string                 `json:"user_id,omitempty" url:"user_id,omitempty"`
	RefreshUrl  string                 `json:"refresh_url,omitempty" url:"refresh_url,omitempty"`
	ReturnUrl   string                 `json:"return_url,omitempty" url:"return_url,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" url:"metadata,omitempty"`
}

// CreateCustomerPortalAccountLink generates a one-time Customer Portal URL.
func (c *Client) CreateCustomerPortalAccountLink(in *CustomerPortalAccountLinkParameters) (out *CustomerPortalAccountLink, err error) {
	return c.CreateCustomerPortalAccountLinkWithContext(context.Background(), in)
}

// CreateCustomerPortalAccountLinkWithContext performs the same operation as CreateCustomerPortalAccountLink, but allows specifying a context that can interrupt the request.
func (c *Client) CreateCustomerPortalAccountLinkWithContext(ctx context.Context, in *CustomerPortalAccountLinkParameters) (out *CustomerPortalAccountLink, err error) {
	err = c.do(ctx, http.MethodPost, "customer_portal/account_link", in, &out)
	return
}
