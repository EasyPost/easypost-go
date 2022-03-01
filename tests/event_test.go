package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
	"reflect"
	"strings"
)

func (c *ClientTests) TestEventAll() {
	client := c.TestClient()
	assert := c.Assert()

	events, _ := client.ListEvents(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)

	eventsList := events.Events

	assert.LessOrEqual(len(eventsList), c.fixture.pageSize())
	assert.NotNil(events.HasMore)
	for _, event := range eventsList {
		assert.Equal(reflect.TypeOf(&easypost.Event{}), reflect.TypeOf(event))
	}
}

func (c *ClientTests) TestEventRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	events, _ := client.ListEvents(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)

	event, _ := client.GetEvent(events.Events[0].ID)
	assert.Equal(reflect.TypeOf(&easypost.Event{}), reflect.TypeOf(event))
	assert.True(strings.HasPrefix(event.ID, "evt_"))
}
