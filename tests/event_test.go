package easypost_test

import "github.com/EasyPost/easypost-go"

func (c *ClientTests) TestEventsGet() {
	client := c.TestClient()
	assert := c.Assert()

	events, err := client.ListEvents(nil)
	assert.NoError(err)
	if assert.Len(events.Events, 7) {
		assert.Equal(events.Events[0].Description, "tracker.updated")
		assert.Equal(
			events.Events[0].PreviousAttributes,
			map[string]interface{}{"status": "out_for_delivery"},
		)
		assert.Equal(events.Events[5].Description, "refund.successful")
		assert.Equal(events.Events[5].Status, "completed")
		assert.Equal(events.Events[6].Description, "payment.failed")
		assert.Len(events.Events[6].CompletedURLs, 1)

		event, err := client.GetEvent(events.Events[0].ID)
		assert.NoError(err)
		assert.Equal(events.Events[0], event)
	}
}

func (c *ClientTests) TestEventsGetPayloads() {
	client := c.TestClient()
	assert := c.Assert()

	events, err := client.ListEvents(nil)
	assert.NoError(err)
	if assert.Len(events.Events, 7) {
		payloads, err := client.ListEventPayloads(events.Events[0].ID)
		assert.NoError(err)
		if assert.Len(payloads, 1) {
			if assert.IsType((*easypost.Event)(nil), payloads[0].RequestBody) {
				event := payloads[0].RequestBody.(*easypost.Event)
				assert.Equal("tracker.updated", event.Description)
				if assert.IsType((*easypost.Tracker)(nil), event.Result) {
					tracker := event.Result.(*easypost.Tracker)
					assert.Equal("9470100897846040813025", tracker.TrackingCode)
				}
			}
		}
		payloads, err = client.ListEventPayloads(events.Events[5].ID)
		assert.NoError(err)
		if assert.Len(payloads, 1) {
			if assert.IsType((*easypost.Event)(nil), payloads[0].RequestBody) {
				event := payloads[0].RequestBody.(*easypost.Event)
				assert.Equal("refund.successful", event.Description)
				if assert.IsType((*easypost.Refund)(nil), event.Result) {
					refund := event.Result.(*easypost.Refund)
					assert.Equal("USPS", refund.Carrier)
				}
			}
		}
		payloads, err = client.ListEventPayloads(events.Events[6].ID)
		assert.NoError(err)
		if assert.Len(payloads, 1) {
			if assert.IsType((*easypost.Event)(nil), payloads[0].RequestBody) {
				event := payloads[0].RequestBody.(*easypost.Event)
				assert.Equal("payment.failed", event.Description)
				if assert.IsType((*easypost.PaymentLog)(nil), event.Result) {
					payment := event.Result.(*easypost.PaymentLog)
					assert.Equal("1000.00000", payment.Amount)
				}
			}
		}
	}

}
