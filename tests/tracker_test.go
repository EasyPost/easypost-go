package easypost_test

import (
	"errors"
	"reflect"
	"strings"

	"github.com/EasyPost/easypost-go/v5"
)

func (c *ClientTests) TestTrackerCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	tracker, err := client.CreateTracker(
		&easypost.CreateTrackerOptions{
			TrackingCode: "EZ1000000001",
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Tracker{}), reflect.TypeOf(tracker))
	assert.True(strings.HasPrefix(tracker.ID, "trk_"))
	assert.Equal("pre_transit", tracker.Status)
}

func (c *ClientTests) TestTrackerRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	tracker, err := client.CreateTracker(
		&easypost.CreateTrackerOptions{
			TrackingCode: "EZ1000000001",
		},
	)
	require.NoError(err)

	// Test trackers cycle through their "dummy" statuses automatically, the created and retrieved objects may differ
	retrievedTracker, err := client.GetTracker(tracker.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Tracker{}), reflect.TypeOf(retrievedTracker))
	assert.Equal(tracker.ID, retrievedTracker.ID)
}

func (c *ClientTests) TestTrackerAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	trackers, err := client.ListTrackers(
		&easypost.ListTrackersOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	trackersList := trackers.Trackers

	assert.LessOrEqual(len(trackersList), c.fixture.pageSize())
	assert.NotNil(trackers.HasMore)
	for _, tracker := range trackersList {
		assert.Equal(reflect.TypeOf(&easypost.Tracker{}), reflect.TypeOf(tracker))
	}
}

func (c *ClientTests) TestTrackersGetNextPage() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListTrackers(
		&easypost.ListTrackersOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextTrackerPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Trackers) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Trackers[len(firstPage.Trackers)-1].ID
			firstIdOfSecondPage := nextPage.Trackers[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIdOfSecondPage)
		}
	}()
	if err != nil {
		var endOfPaginationErr *easypost.EndOfPaginationError
		if errors.As(err, &endOfPaginationErr) {
			assert.Equal(err.Error(), endOfPaginationErr.Error())
			return
		}
		require.NoError(err)
	}
}
