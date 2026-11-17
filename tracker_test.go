package easypost

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

func (c *ClientTests) TestTrackerCreate() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	tracker, err := client.CreateTracker(
		&CreateTrackerOptions{
			TrackingCode: "EZ1000000001",
		},
	)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&Tracker{}), reflect.TypeOf(tracker))
	assert.True(strings.HasPrefix(tracker.ID, "trk_"))
	assert.Equal("pre_transit", tracker.Status)
}

func (c *ClientTests) TestTrackerRetrieve() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	tracker, err := client.CreateTracker(
		&CreateTrackerOptions{
			TrackingCode: "EZ1000000001",
		},
	)
	require.NoError(err)

	// Test trackers cycle through their "dummy" statuses automatically, the created and retrieved objects may differ
	retrievedTracker, err := client.GetTracker(tracker.ID)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&Tracker{}), reflect.TypeOf(retrievedTracker))
	assert.Equal(tracker.ID, retrievedTracker.ID)
}

func (c *ClientTests) TestTrackerAll() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	trackers, err := client.ListTrackers(
		&ListTrackersOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	trackersList := trackers.Trackers

	assert.LessOrEqual(len(trackersList), c.fixture.pageSize())
	assert.NotNil(trackers.HasMore)
	for _, tracker := range trackersList {
		assert.Equal(reflect.TypeOf(&Tracker{}), reflect.TypeOf(tracker))
	}
}

func (c *ClientTests) TestTrackersGetNextPage() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	firstPage, err := client.ListTrackers(
		&ListTrackersOptions{
			PageSize: c.fixture.pageSize(),
		},
	)
	require.NoError(err)

	nextPage, err := client.GetNextTrackerPageWithPageSize(firstPage, c.fixture.pageSize())
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Trackers) <= c.fixture.pageSize())

			lastIDOfFirstPage := firstPage.Trackers[len(firstPage.Trackers)-1].ID
			firstIDOfSecondPage := nextPage.Trackers[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIDOfSecondPage)
		}
	}()
	if err != nil {
		var endOfPaginationErr *EndOfPaginationError
		if errors.As(err, &endOfPaginationErr) {
			assert.Equal(err.Error(), endOfPaginationErr.Error())
			return
		}
		require.NoError(err)
	}
}

func (c *ClientTests) TestTrackerAllWithTrackingCodes() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	trackingCodes := []string{"EZ1000000001", "EZ2000000002"}
	for _, code := range trackingCodes {
		_, err := client.CreateTracker(&CreateTrackerOptions{TrackingCode: code})
		require.NoError(err)
	}

	trackers, err := client.ListTrackers(
		&ListTrackersOptions{
			TrackingCodes: trackingCodes,
			PageSize:      c.fixture.pageSize(),
		},
	)
	require.NoError(err)
	require.NotEmpty(trackers.Trackers)

	found := map[string]bool{}
	for _, tracker := range trackers.Trackers {
		assert.Equal(reflect.TypeOf(&Tracker{}), reflect.TypeOf(tracker))
		assert.Contains(trackingCodes, tracker.TrackingCode)
		found[tracker.TrackingCode] = true
	}
	for _, code := range trackingCodes {
		assert.True(found[code], "expected tracking code %s in list results", code)
	}
}

func (c *ClientTests) TestTrackersGetNextPageWithTrackingCodes() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	trackingCodes := []string{"EZ1000000001", "EZ2000000002"}
	for _, code := range trackingCodes {
		_, err := client.CreateTracker(&CreateTrackerOptions{TrackingCode: code})
		require.NoError(err)
	}

	// Wide fixed window so the cassette URL is stable on replay and still
	// includes the test trackers created earlier in this test.
	start := NewDateTime(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := NewDateTime(2030, time.January, 1, 0, 0, 0, 0, time.UTC)

	firstPage, err := client.ListTrackers(
		&ListTrackersOptions{
			TrackingCodes: trackingCodes,
			StartDateTime: &start,
			EndDateTime:   &end,
			PageSize:      1,
		},
	)
	require.NoError(err)
	require.NotEmpty(firstPage.Trackers)
	assert.Equal(trackingCodes, firstPage.TrackingCodes)
	assert.Equal(&start, firstPage.StartDateTime)
	assert.Equal(&end, firstPage.EndDateTime)
	for _, tracker := range firstPage.Trackers {
		assert.Contains(trackingCodes, tracker.TrackingCode)
	}

	nextPage, err := client.GetNextTrackerPageWithPageSize(firstPage, 1)
	defer func() {
		if err == nil {
			assert.True(len(nextPage.Trackers) <= 1)
			assert.Equal(trackingCodes, nextPage.TrackingCodes)
			assert.Equal(&start, nextPage.StartDateTime)
			assert.Equal(&end, nextPage.EndDateTime)
			for _, tracker := range nextPage.Trackers {
				assert.Contains(trackingCodes, tracker.TrackingCode)
			}

			lastIDOfFirstPage := firstPage.Trackers[len(firstPage.Trackers)-1].ID
			firstIdOfSecondPage := nextPage.Trackers[0].ID

			assert.NotEqual(lastIDOfFirstPage, firstIdOfSecondPage)
		}
	}()
	if err != nil {
		var endOfPaginationErr *EndOfPaginationError
		if errors.As(err, &endOfPaginationErr) {
			assert.Equal(err.Error(), endOfPaginationErr.Error())
			return
		}
		require.NoError(err)
	}
}

func (c *ClientTests) TestTrackerRetrieveBatch() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	tracker, err := client.CreateTracker(
		&CreateTrackerOptions{
			TrackingCode: "EZ1000000001",
		},
	)
	require.NoError(err)

	trackers, err := client.RetrieveTrackerBatch(
		&ListTrackersOptions{
			TrackingCodes: []string{tracker.TrackingCode},
		},
	)
	require.NoError(err)

	trackersList := trackers.Trackers

	for _, tracker := range trackersList {
		assert.Equal(reflect.TypeOf(&Tracker{}), reflect.TypeOf(tracker))
	}
}

func (c *ClientTests) TestTrackerDelete() {
	client := c.TestClient()
	require := c.Require()

	tracker, err := client.CreateTracker(
		&CreateTrackerOptions{
			TrackingCode: "EZ1000000001",
		},
	)
	require.NoError(err)

	err = client.DeleteTracker(tracker.ID)
	require.NoError(err)
}
