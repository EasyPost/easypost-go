package easypost

import (
	"context"
	"net/http"
)

// A EmbeddablesSession object represents an object containing an AccountLink.
type EmbeddablesSession struct {
	Object    string    `json:"object,omitempty" url:"object,omitempty"`
	SessionId string    `json:"session_id,omitempty" url:"session_id,omitempty"`
	CreatedAt *DateTime `json:"created_at,omitempty" url:"created_at,omitempty"`
	ExpiresAt *DateTime `json:"expires_at,omitempty" url:"expires_at,omitempty"`
}

// EmbeddablesSessionParameters is used to specify parameters for creating an AccountLink.
type EmbeddablesSessionParameters struct {
	OriginHost string `json:"origin_host,omitempty" url:"origin_host,omitempty"`
	UserId     string `json:"user_id,omitempty" url:"user_id,omitempty"`
}

// CreateEmbeddablesSession creates a temporary session for initializing embeddable components.
func (c *Client) CreateEmbeddablesSession(in *EmbeddablesSessionParameters) (out *EmbeddablesSession, err error) {
	return c.CreateEmbeddablesSessionWithContext(context.Background(), in)
}

// CreateEmbeddablesSessionWithContext performs the same operation as CreateEmbeddablesSession, but allows specifying a context that can interrupt the request.
func (c *Client) CreateEmbeddablesSessionWithContext(ctx context.Context, in *EmbeddablesSessionParameters) (out *EmbeddablesSession, err error) {
	err = c.do(ctx, http.MethodPost, "embeddables/session", in, &out)
	return
}
