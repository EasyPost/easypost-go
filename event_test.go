package easypost

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func (c *ClientTests) TestEventAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	events, err := client.ListEvents(
		&ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	eventsList := events.Events

	assert.LessOrEqual(len(eventsList), c.fixture.pageSize())
	assert.NotNil(events.HasMore)
	for _, event := range eventsList {
		assert.Equal(reflect.TypeOf(&Event{}), reflect.TypeOf(event))
	}
}

func (c *ClientTests) TestEventRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	events, err := client.ListEvents(
		&ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	event, err := client.GetEvent(events.Events[0].ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&Event{}), reflect.TypeOf(event))
	assert.True(strings.HasPrefix(event.ID, "evt_"))
}

func (c *ClientTests) TestEventRetrieveAllPayloads() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	// Create a webhook to receive the event
	webhook, webhookErr := client.CreateWebhook(
		&CreateUpdateWebhookOptions{
			URL: c.fixture.WebhookUrl(),
		},
	)
	require.NoError(webhookErr)

	// Create a batch to trigger an event
	_, batchErr := client.CreateBatch(
		c.fixture.OneCallBuyShipment(),
	)
	require.NoError(batchErr)

	currentDir, _ := os.Getwd()
	cassettePath := filepath.Join(currentDir, "cassettes", "TestEventRetrieveAllPayloads.yaml")
	if _, err := os.Stat(cassettePath); errors.Is(err, os.ErrNotExist) {
		time.Sleep(5 * time.Second) // Wait enough time for the event to be created
	}

	// Retrieve all events and extract the newest one
	events, err := client.ListEvents(
		&ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)
	event := events.Events[0]

	// Retrieve all payloads for the event
	eventPayloads, err := client.ListEventPayloads(event.ID)
	require.NoError(err)

	assert.NotNil(eventPayloads)

	// Delete the webhook
	_ = client.DeleteWebhook(webhook.ID)
}

func (c *ClientTests) TestEventRetrievePayload() {
	client := c.TestClient()
	require := c.Require()

	// Create a webhook to receive the event
	webhook, webhookErr := client.CreateWebhook(
		&CreateUpdateWebhookOptions{
			URL: c.fixture.WebhookUrl(),
		},
	)
	require.NoError(webhookErr)

	// Create a batch to trigger an event
	_, batchErr := client.CreateBatch(
		c.fixture.OneCallBuyShipment(),
	)
	require.NoError(batchErr)

	currentDir, _ := os.Getwd()
	cassettePath := filepath.Join(currentDir, "cassettes", "TestEventRetrieveAllPayloads.yaml")
	if _, err := os.Stat(cassettePath); errors.Is(err, os.ErrNotExist) {
		time.Sleep(5 * time.Second) // Wait enough time for the event to be created
	}

	events, err := client.ListEvents(
		&ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)
	event := events.Events[0]

	// Payload does not exist due to queueing, so this will throw an error
	_, err = client.GetEventPayload(event.ID, "payload_11111111111111111111111111111111")
	require.Error(err)

	_ = client.DeleteWebhook(webhook.ID)
}

func (c *ClientTests) TestEventsGetNextPage() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListEvents(
		&ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextEventPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Events) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Events[len(firstPage.Events)-1].ID
			firstIdOfSecondPage := nextPage.Events[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIdOfSecondPage)
		}
	}()
	if err != nil {
		assert.Equal(NoPagesLeftToRetrieve, err.Error())
		return
	}
}

func TestEventUnmarshalJSONWithoutResult(t *testing.T) {
	data := []byte(`{"id":"evt_123","object":"Event","description":"test event"}`)

	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("unmarshal event failed: %v", err)
	}

	if event.ID != "evt_123" {
		t.Fatalf("unexpected event ID: %s", event.ID)
	}

	if event.Result != nil {
		t.Fatalf("expected nil result when result field is absent, got %T", event.Result)
	}
}

func TestEventPayloadUnmarshalJSONRequestBodyString(t *testing.T) {
	data := []byte(`{"id":"payload_123","object":"Payload","request_body":"{\"id\":\"evt_123\",\"object\":\"Event\"}"}`)

	var payload EventPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal event payload failed: %v", err)
	}

	event, ok := payload.RequestBody.(*Event)
	if !ok {
		t.Fatalf("expected request_body to decode into *Event, got %T", payload.RequestBody)
	}

	if event.ID != "evt_123" {
		t.Fatalf("unexpected nested event ID: %s", event.ID)
	}
}
