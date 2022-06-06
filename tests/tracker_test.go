package easypost_test

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/EasyPost/easypost-go/v2"
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

func (c *ClientTests) TestTrackerCreateList() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	trackingCodeParam := make(map[string]interface{})
	trackingCodes := []string{"EZ1000000001", "EZ1000000002", "EZ1000000003"}

	for i := range trackingCodes {
		trackingCodeParam[strconv.Itoa(i)] = map[string]string{
			"tracking_code": trackingCodes[i],
		}
	}

	response, err := client.CreateTrackerList(trackingCodeParam)
	require.NoError(err)

	// This endpoint returns nothing, so we assert the function returns true
	assert.True(response)
}
