package easypost_test

import (
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
)

func (c *ClientTests) TestEventAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	events, err := client.ListEvents(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	eventsList := events.Events

	assert.LessOrEqual(len(eventsList), c.fixture.pageSize())
	assert.NotNil(events.HasMore)
	for _, event := range eventsList {
		assert.Equal(reflect.TypeOf(&easypost.Event{}), reflect.TypeOf(event))
	}
}

func (c *ClientTests) TestEventRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	events, err := client.ListEvents(
		&easypost.ListOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	event, err := client.GetEvent(events.Events[0].ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Event{}), reflect.TypeOf(event))
	assert.True(strings.HasPrefix(event.ID, "evt_"))
}
