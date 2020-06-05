package easypost

import (
	"context"
	"net/http"
	"time"
)

// Event objects contain details about changes to EasyPost objects
type Event struct {
	ID                 string                 `json:"id,omitempty"`
	Object             string                 `json:"object,omitempty"`
	Mode               string                 `json:"mode,omitempty"`
	CreatedAt          *time.Time             `json:"created_at,omitempty"`
	UpdatedAt          *time.Time             `json:"updated_at,omitempty"`
	Description        string                 `json:"description,omitempty"`
	PreviousAttributes map[string]interface{} `json:"previous_attributes,omitempty"`
	Result             map[string]interface{} `json:"result,omitempty"`
	Status             string                 `json:"status,omitempty"`
	PendingURLs        []string               `json:"pending_urls,omitempty"`
	CompletedURLs      []string               `json:"completed_urls,omitempty"`
}

// ListEventsResult holds the results from the list events API.
type ListEventsResult struct {
	Events []*Event `json:"events,omitempty"`
	// HasMore indicates if there are more responses to be fetched. If True,
	// additional responses can be fetched by updating the ListEventsOptions
	// parameter's AfterID field with the ID of the last item in this object's
	// Events field.
	HasMore bool `json:"has_more,omitempty"`
}

// ListEvents provides a paginated result of Event objects.
func (c *Client) ListEvents(opts *ListOptions) (out *ListEventsResult, err error) {
	return c.ListEventsWithContext(nil, opts)
}

// ListEventsWithContext performs the same operation as ListEventes, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListEventsWithContext(ctx context.Context, opts *ListOptions) (out *ListEventsResult, err error) {
	err = c.do(ctx, http.MethodGet, "events", c.convertOptsToURLValues(opts), &out)
	return
}

// GetEvent retrieves a previously-created event by its ID.
func (c *Client) GetEvent(eventID string) (out *Event, err error) {
	err = c.get(nil, "events/"+eventID, &out)
	return
}

// GetEventWithContext performs the same operation as GetEvent, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetEventWithContext(ctx context.Context, eventID string) (out *Event, err error) {
	err = c.get(ctx, "events/"+eventID, &out)
	return
}
