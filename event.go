package easypost

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

// Event objects contain details about changes to EasyPost objects
type Event struct {
	ID                 string                 `json:"id,omitempty" url:"id,omitempty"`
	UserID             string                 `json:"user_id,omitempty" url:"user_id,omitempty"`
	Object             string                 `json:"object,omitempty" url:"object,omitempty"`
	Mode               string                 `json:"mode,omitempty" url:"mode,omitempty"`
	CreatedAt          *DateTime              `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt          *DateTime              `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	Description        string                 `json:"description,omitempty" url:"description,omitempty"`
	PreviousAttributes map[string]interface{} `json:"previous_attributes,omitempty" url:"previous_attributes,omitempty"`
	// Result will be populated with the relevant object type, i.e.
	// *Batch, *Insurance, *PaymentLog, *Refund, *Report, *Tracker or *ScanForm.
	// It will be nil if no 'result' field is present, which is the case for
	// the ListEvents and GetEvents methods. The RequestBody field of the
	// EventPayload type will generally be an instance of *Event with this field
	// present. Having the field here also enables re-using this type to
	// implement a webhook handler.
	Result        interface{} `json:"result,omitempty" url:"result,omitempty"`
	Status        string      `json:"status,omitempty" url:"status,omitempty"`
	PendingURLs   []string    `json:"pending_urls,omitempty" url:"pending_urls,omitempty"`
	CompletedURLs []string    `json:"completed_urls,omitempty" url:"completed_urls,omitempty"`
}

// EventPayload represents the result of a webhook call.
type EventPayload struct {
	ID             string            `json:"id,omitempty" url:"id,omitempty"`
	Object         string            `json:"object,omitempty" url:"object,omitempty"`
	CreatedAt      *DateTime         `json:"created_at,omitempty" url:"created_at,omitempty"`
	UpdatedAt      *DateTime         `json:"updated_at,omitempty" url:"updated_at,omitempty"`
	RequestURL     string            `json:"request_url,omitempty" url:"request_url,omitempty"`
	RequestHeaders map[string]string `json:"request_headers,omitempty" url:"request_headers,omitempty"`
	// RequestBody is the raw request body that was sent to the webhook. This is
	// expected to be an Event object. It may either be encoded in the API
	// response as a string (with JSON delimiters escaped) or as base64. The
	// UnmarshalJSON method will attempt to convert it to an *Event type, but it
	// may be set to a default type if decoding to an object fails.
	RequestBody     interface{}       `json:"request_body,omitempty" url:"request_body,omitempty"`
	ResponseHeaders map[string]string `json:"response_headers,omitempty" url:"response_headers,omitempty"`
	ResponseBody    string            `json:"response_body,omitempty" url:"response_body,omitempty"`
	ResponseCode    int               `json:"response_code,omitempty" url:"response_code,omitempty"`
	TotalTime       int               `json:"total_time,omitempty" url:"total_time,omitempty"`
}

// ListEventsResult holds the results from the list events API.
type ListEventsResult struct {
	Events []*Event `json:"events,omitempty" url:"events,omitempty"`
	PaginatedCollection
}

type listEventPayloadsResult struct {
	Payloads *[]*EventPayload `json:"payloads,omitempty" url:"payloads,omitempty"`
}

// UnmarshalJSON will attempt to convert an event response to an *Event type,
// but it may be set to a default type if decoding to an object fails.
func (e *Event) UnmarshalJSON(data []byte) (err error) {
	var buf json.RawMessage
	event := Event{Result: &buf}

	type nonUnmarshaler *Event
	if err = json.Unmarshal(data, nonUnmarshaler(&event)); err != nil {
		return err
	}

	if event.Result, err = UnmarshalJSONObject(buf); err == nil {
		*e = event
	}

	return err
}

// UnmarshalJSON will attempt to convert an event payload response to an *EventPayload type,
// but it may be set to a default type if decoding to an object fails.
func (e *EventPayload) UnmarshalJSON(data []byte) (err error) {
	var s string
	payload := EventPayload{RequestBody: &s}

	type nonUnmarshaler *EventPayload
	if err = json.Unmarshal(data, nonUnmarshaler(&payload)); err != nil {
		return err
	}

	// Attempt to base64 decode the body. Ignore errors.
	if buf, err := base64.StdEncoding.DecodeString(s); err == nil {
		s = string(buf)
	}

	// try to decode RequestBody to an object, but if we can't, then just
	// set it to the string.
	if payload.RequestBody, err = UnmarshalJSONObject([]byte(s)); err != nil {
		payload.RequestBody = s
	}

	*e = payload
	return nil
}

// ListEvents provides a paginated result of Event objects.
func (c *Client) ListEvents(opts *ListOptions) (out *ListEventsResult, err error) {
	return c.ListEventsWithContext(context.Background(), opts)
}

// ListEventsWithContext performs the same operation as ListEvents, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListEventsWithContext(ctx context.Context, opts *ListOptions) (out *ListEventsResult, err error) {
	err = c.do(ctx, http.MethodGet, "events", opts, &out)
	return
}

// GetNextEventPage returns the next page of events
func (c *Client) GetNextEventPage(collection *ListEventsResult) (out *ListEventsResult, err error) {
	return c.GetNextEventPageWithContext(context.Background(), collection)
}

// GetNextEventPageWithPageSize returns the next page of events with a specific page size
func (c *Client) GetNextEventPageWithPageSize(collection *ListEventsResult, pageSize int) (out *ListEventsResult, err error) {
	return c.GetNextEventPageWithPageSizeWithContext(context.Background(), collection, pageSize)
}

// GetNextEventPageWithContext performs the same operation as GetNextEventPage, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextEventPageWithContext(ctx context.Context, collection *ListEventsResult) (out *ListEventsResult, err error) {
	return c.GetNextEventPageWithPageSizeWithContext(ctx, collection, 0)
}

// GetNextEventPageWithPageSizeWithContext performs the same operation as GetNextEventPageWithPageSize, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetNextEventPageWithPageSizeWithContext(ctx context.Context, collection *ListEventsResult, pageSize int) (out *ListEventsResult, err error) {
	if len(collection.Events) == 0 {
		err = newEndOfPaginationError()
		return
	}
	lastID := collection.Events[len(collection.Events)-1].ID
	params, err := nextPageParameters(collection.HasMore, lastID, pageSize)
	if err != nil {
		return
	}
	return c.ListEventsWithContext(ctx, params)
}

// GetEvent retrieves a previously-created event by its ID.
func (c *Client) GetEvent(eventID string) (out *Event, err error) {
	return c.GetEventWithContext(context.Background(), eventID)
}

// GetEventWithContext performs the same operation as GetEvent, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetEventWithContext(ctx context.Context, eventID string) (out *Event, err error) {
	err = c.do(ctx, http.MethodGet, "events/"+eventID, nil, &out)
	return
}

// ListEventPayloads retrieves the payload results of a previous webhook call.
func (c *Client) ListEventPayloads(eventID string) (out []*EventPayload, err error) {
	return c.ListEventPayloadsWithContext(context.Background(), eventID)
}

// ListEventPayloadsWithContext performs the same operation as GetEventPayload, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListEventPayloadsWithContext(ctx context.Context, eventID string) (out []*EventPayload, err error) {
	err = c.do(ctx, http.MethodGet, "events/"+eventID+"/payloads", nil, &listEventPayloadsResult{Payloads: &out})
	return
}

// GetEventPayload retrieves a specific payload result of a previous webhook call.
func (c *Client) GetEventPayload(eventID, payloadID string) (out *EventPayload, err error) {
	return c.GetEventPayloadWithContext(context.Background(), eventID, payloadID)
}

// GetEventPayloadWithContext performs the same operation as GetEventPayload, but
// allows specifying a context that can interrupt the request.
func (c *Client) GetEventPayloadWithContext(ctx context.Context, eventID, payloadID string) (out *EventPayload, err error) {
	err = c.do(ctx, http.MethodGet, "events/"+eventID+"/payloads/"+payloadID, nil, &out)
	return
}
