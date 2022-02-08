package easypost_test

import (
	"strconv"

	"github.com/EasyPost/easypost-go"
)

func (c *ClientTests) TestTrackerValues() {
	client := c.TestClient()
	assert := c.Assert()
	type Value struct{ Code, Status string }
	values := []Value{
		{Code: "EZ1000000001", Status: "pre_transit"},
		{Code: "EZ2000000002", Status: "in_transit"},
		{Code: "EZ3000000003", Status: "out_for_delivery"},
		{Code: "EZ4000000004", Status: "delivered"},
		{Code: "EZ5000000005", Status: "return_to_sender"},
		{Code: "EZ6000000006", Status: "failure"},
		{Code: "EZ7000000007", Status: "unknown"},
	}
	for i := range values {
		opts := easypost.CreateTrackerOptions{TrackingCode: values[i].Code}
		tracker, err := client.CreateTracker(&opts)
		if assert.NoError(err) {
			assert.Equal(values[i].Status, tracker.Status)
			assert.NotEmpty(tracker.TrackingDetails)
			if values[i].Status == "delivered" {
				assert.Equal("John Tester", tracker.SignedBy)
			}
		}
	}
}

func (c *ClientTests) TestTrackerCreateList() {
	client := c.TestClient()
	assert := c.Assert()

	trackingCodeParam := make(map[string]interface{})
	trackingCodes := []string{"EZ1000000001", "EZ1000000002", "EZ1000000003"}

	for i := range trackingCodes {
		trackingCodeParam[strconv.Itoa(i)] = map[string]string{
			"tracking_code": trackingCodes[i],
			"carrier":       "USPS",
		}
	}
	response, err := client.CreateTrackerList(trackingCodeParam)
	assert.NoError(err)
	assert.True(response)
}
