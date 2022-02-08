package easypost

import (
	"context"
	"time"
)

// APIKey represents a single API key.
type APIKey struct {
	Object    string     `json:"object,omitempty"`
	Mode      string     `json:"mode,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Key       string     `json:"key,omitempty"`
}

// APIKeys contains information about a list of API keys for the given user and
// any child users.
type APIKeys struct {
	ID       string     `json:"id,omitempty"`
	Children []*APIKeys `json:"children,omitempty"`
	Keys     []*APIKey  `json:"keys,omitempty"`
}

// GetAPIKeys returns the list of API keys associated with the current user.
func (c *Client) GetAPIKeys() (out *APIKeys, err error) {
	err = c.get(context.Background(), "api_keys", &out)
	return
}

// GetAPIKeysWithContext performs the same operation as GetAPIKeys, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetAPIKeysWithContext(ctx context.Context) (out *APIKeys, err error) {
	err = c.get(ctx, "api_keys", &out)
	return
}
