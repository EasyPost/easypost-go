package easypost_test

import (
	"github.com/EasyPost/easypost-go/v2"
	"reflect"
	"strconv"
	"strings"
)

func (c *ClientTests) TestTrackerCreate() {
	client := c.TestClient()
	assert := c.Assert()

	tracker, _ := client.CreateTracker(
		&easypost.CreateTrackerOptions{
			TrackingCode: "EZ1000000001",
		},
	)

	assert.Equal(reflect.TypeOf(&easypost.Tracker{}), reflect.TypeOf(tracker))
	assert.True(strings.HasPrefix(tracker.ID, "trk_"))
	assert.Equal("pre_transit", tracker.Status)
}

func (c *ClientTests) TestTrackerRetrieve() {
	client := c.TestClient()
	assert := c.Assert()

	tracker, _ := client.CreateTracker(
		&easypost.CreateTrackerOptions{
			TrackingCode: "EZ1000000001",
		},
	)

	// Test trackers cycle through their "dummy" statuses automatically, the created and retrieved objects may differ
	retrievedTracker, _ := client.GetTracker(tracker.ID)

	assert.Equal(reflect.TypeOf(&easypost.Tracker{}), reflect.TypeOf(retrievedTracker))
	assert.Equal(tracker.ID, retrievedTracker.ID)
}

func (c *ClientTests) TestTrackerAll() {
	client := c.TestClient()
	assert := c.Assert()

	trackers, _ := client.ListTrackers(
		&easypost.ListTrackersOptions{
			PageSize: c.fixture.pageSize(),
		},
	)

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

	// This endpoint returns nothing so we assert the function returns true
	assert.True(response)
	require.NoError(err)
}
