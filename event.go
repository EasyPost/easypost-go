package easypost

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

// Event objects contain details about changes to EasyPost objects
type Event struct {
	ID                 string                 `json:"id,omitempty"`
	UserID             string                 `json:"user_id,omitempty"`
	Object             string                 `json:"object,omitempty"`
	Mode               string                 `json:"mode,omitempty"`
	CreatedAt          *time.Time             `json:"created_at,omitempty"`
	UpdatedAt          *time.Time             `json:"updated_at,omitempty"`
	Description        string                 `json:"description,omitempty"`
	PreviousAttributes map[string]interface{} `json:"previous_attributes,omitempty"`
	// Result will be populated with the relevant object type, i.e.
	// *Batch, *Insurance, *PaymentLog, *Refund, *Report, *Tracker or *ScanForm.
	// It will be nil if no 'result' field is present, which is the case for
	// the ListEvents and GetEvents methods. The RequestBody field of the
	// EventPayload type will generally be an instance of *Event with this field
	// present. Having the field here also enables re-using this type to
	// implement a webhook handler.
	Result        interface{} `json:"result,omitempty"`
	Status        string      `json:"status,omitempty"`
	PendingURLs   []string    `json:"pending_urls,omitempty"`
	CompletedURLs []string    `json:"completed_urls,omitempty"`
}

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

// EventPayload represents the result of a webhook call.
type EventPayload struct {
	ID             string            `json:"id,omitempty"`
	Object         string            `json:"object,omitempty"`
	CreatedAt      *time.Time        `json:"created_at,omitempty"`
	UpdatedAt      *time.Time        `json:"updated_at,omitempty"`
	RequestURL     string            `json:"request_url,omitempty"`
	RequestHeaders map[string]string `json:"request_headers,omitempty"`
	// RequestBody is the raw request body that was sent to the webhook. This is
	// expected to be an Event object. It may either be encoded in the API
	// response as a string (with JSON delimiters escaped) or as base64. The
	// UnmarshalJSON method will attempt to convert it to an *Event type, but it
	// may be set to a default type if decoding to an object fails.
	RequestBody     interface{}       `json:"request_body,omitempty"`
	ResponseHeaders map[string]string `json:"response_headers,omitempty"`
	ResponseBody    string            `json:"response_body,omitempty"`
	ResponseCode    int               `json:"response_code,omitempty"`
	TotalTime       int               `json:"total_time,omitempty"`
}

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
	return c.ListEventsWithContext(context.Background(), opts)
}

// ListEventsWithContext performs the same operation as ListEventes, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListEventsWithContext(ctx context.Context, opts *ListOptions) (out *ListEventsResult, err error) {
	err = c.do(ctx, http.MethodGet, "events", c.convertOptsToURLValues(opts), &out)
	return
}

// GetEvent retrieves a previously-created event by its ID.
func (c *Client) GetEvent(eventID string) (out *Event, err error) {
	err = c.get(context.Background(), "events/"+eventID, &out)
	return
}

// GetEventWithContext performs the same operation as GetEvent, but allows
// specifying a context that can interrupt the request.
func (c *Client) GetEventWithContext(ctx context.Context, eventID string) (out *Event, err error) {
	err = c.get(ctx, "events/"+eventID, &out)
	return
}

type listEventPayloadsResult struct {
	Payloads *[]*EventPayload `json:"payloads,omitempty"`
}

// GetEventPayload retrieves the payload results of a previous webhook call.
func (c *Client) ListEventPayloads(eventID string) (out []*EventPayload, err error) {
	return c.ListEventPayloadsWithContext(context.Background(), eventID)
}

// GetEventPayloadWithContext performs the same operation as GetEventPaylod, but
// allows specifying a context that can interrupt the request.
func (c *Client) ListEventPayloadsWithContext(ctx context.Context, eventID string) (out []*EventPayload, err error) {
	err = c.get(
		ctx,
		"events/"+eventID+"/payloads",
		&listEventPayloadsResult{Payloads: &out},
	)
	return
}
