package easypost

import (
	"context"
	"net/http"
)

// APIKey represents a single API key.
type APIKey struct {
	ID        string    `json:"id,omitempty" url:"id,omitempty"`
	Object    string    `json:"object,omitempty" url:"object,omitempty"`
	Mode      string    `json:"mode,omitempty" url:"mode,omitempty"`
	CreatedAt *DateTime `json:"created_at,omitempty" url:"created_at,omitempty"`
	Key       string    `json:"key,omitempty" url:"key,omitempty"`
	Active    bool      `json:"active,omitempty" url:"active,omitempty"`
}

// APIKeys contains information about a list of API keys for the given user and
// any child users.
type APIKeys struct {
	ID       string     `json:"id,omitempty" url:"id,omitempty"`
	Children []*APIKeys `json:"children,omitempty" url:"children,omitempty"`
	Keys     []*APIKey  `json:"keys,omitempty" url:"keys,omitempty"`
}

// GetAPIKeys returns the list of API keys associated with the current user.
func (c *Client) GetAPIKeys() (out *APIKeys, err error) {
	return c.GetAPIKeysWithContext(context.Background())
}

// GetAPIKeysWithContext performs the same operation as GetAPIKeys, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetAPIKeysWithContext(ctx context.Context) (out *APIKeys, err error) {
	err = c.do(ctx, http.MethodGet, "api_keys", nil, &out)
	return
}

// GetAPIKeysForUser returns the list of API keys associated with the given user ID.
func (c *Client) GetAPIKeysForUser(userID string) (out *APIKeys, err error) {
	return c.GetAPIKeysForUserWithContext(context.Background(), userID)
}

// GetAPIKeysForUserWithContext performs the same operation as GetAPIKeysForUser, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetAPIKeysForUserWithContext(ctx context.Context, userID string) (out *APIKeys, err error) {
	keys, err := c.GetAPIKeysWithContext(ctx)

	if err != nil {
		return nil, err
	}

	if keys.ID == userID {
		return keys, nil
	}

	for _, child := range keys.Children {
		if child.ID == userID {
			return child, nil
		}
	}

	return nil, newFilteringError(NoUserFoundForId)
}

// CreateAPIKey creates an API key for a child or referral customer user.
func (c *Client) CreateAPIKey(mode string) (out *APIKey, err error) {
	return c.CreateAPIKeyWithContext(context.Background(), mode)
}

// CreateAPIKeyWithContext performs the same operation as CreateAPIKey, but allows
// specifying a context that can interrupt the request.
func (c *Client) CreateAPIKeyWithContext(ctx context.Context, mode string) (out *APIKey, err error) {
	params := map[string]interface{}{"mode": mode}
	err = c.do(ctx, http.MethodPost, "api_keys", params, &out)
	return
}

// DeleteAPIKey deletes an API key for a child or referral customer user.
func (c *Client) DeleteAPIKey(id string) error {
	return c.DeleteAPIKeyWithContext(context.Background(), id)
}

// DeleteAPIKeyWithContext performs the same operation as DeleteAPIKey, but allows
// specifying a context that can interrupt the request.
func (c *Client) DeleteAPIKeyWithContext(ctx context.Context, id string) error {
	return c.do(ctx, http.MethodDelete, "api_keys/"+id, nil, nil)
}

// EnableAPIKey enables a child or referral customer API key.
func (c *Client) EnableAPIKey(id string) (out *APIKey, err error) {
	return c.EnableAPIKeyWithContext(context.Background(), id)
}

// EnableAPIKeyWithContext performs the same operation as EnableAPIKey, but allows
// specifying a context that can interrupt the request.
func (c *Client) EnableAPIKeyWithContext(ctx context.Context, id string) (out *APIKey, err error) {
	err = c.do(ctx, http.MethodPost, "api_keys/"+id+"/enable", nil, &out)
	return
}

// DisableAPIKey disables a child or referral customer API key.
func (c *Client) DisableAPIKey(id string) (out *APIKey, err error) {
	return c.DisableAPIKeyWithContext(context.Background(), id)
}

// DisableAPIKeyWithContext performs the same operation as DisableAPIKey, but allows
// specifying a context that can interrupt the request.
func (c *Client) DisableAPIKeyWithContext(ctx context.Context, id string) (out *APIKey, err error) {
	err = c.do(ctx, http.MethodPost, "api_keys/"+id+"/disable", nil, &out)
	return
}
