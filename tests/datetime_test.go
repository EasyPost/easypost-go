package easypost_test

import (
	"github.com/EasyPost/easypost-go/v4"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
)

func (c *ClientTests) TestDateTime() {
	assert := c.Assert()

	year := 2019
	month := time.Month(1)
	day := 1
	hour := 0
	minute := 0
	second := 0
	nanosecond := 0
	location := time.UTC

	_time := time.Date(year, month, day, hour, minute, second, nanosecond, location)

	datetime := easypost.DateTime(_time)
	assert.Equal(_time, datetime.AsTime())

	datetimeFromTime := easypost.DateTimeFromTime(_time)
	assert.Equal(_time, datetimeFromTime.AsTime())

	datetimeFromNumbers := easypost.NewDateTime(year, month, day, hour, minute, second, nanosecond, location)
	assert.Equal(_time, datetimeFromNumbers.AsTime())
}

func (c *ClientTests) TestDateTimeJSON() {
	client := c.TestClient()
	assert, require := c.Assert(), c.Require()

	shipment, err := client.CreateShipment(c.fixture.OneCallBuyShipment())
	require.NoError(err)

	// fixture uses a DateTime object behind the scenes
	pickupData := c.fixture.BasicPickup()
	pickupData.Shipment = shipment

	// will test serializing (during call) and deserializing (parsing response) the DateTime object
	pickup, err := client.CreatePickup(pickupData)
	require.NoError(err)

	assert.Equal(reflect.TypeOf(&easypost.Pickup{}), reflect.TypeOf(pickup))
	assert.NotNil(pickup.MaxDatetime)
	assert.NotNil(pickup.MinDatetime)

	maxDateTimeTime := pickup.MaxDatetime.AsTime()
	minDateTimeTime := pickup.MinDatetime.AsTime()
	assert.NotNil(maxDateTimeTime)
	assert.NotNil(minDateTimeTime)

	maxDatetimeString := pickup.MaxDatetime.String()
	minDatetimeString := pickup.MinDatetime.String()
	assert.NotNil(maxDatetimeString)
	assert.NotNil(minDatetimeString)
}

// / TestDateTimeUrlQueryStringParameterInclusion tests that DateTime-type parameters are encoded properly (RFC3339) for use in HTTP payloads
func (c *ClientTests) TestDateTimeUrlQueryStringParameterInclusion() {
	assert := c.Assert()

	now := time.Now().UTC()
	past := now.AddDate(0, -4, 0)

	start := easypost.DateTimeFromTime(past)
	end := easypost.DateTimeFromTime(now)

	options := easypost.ListOptions{
		StartDateTime: &start,
		EndDateTime:   &end,
		PageSize:      100,
	}

	// This is the internals of the `convertOptsToURLValues` method in the `Client` struct (can't be used directly because it's private)
	values, _ := query.Values(options)

	assert.Equal(past.Format(time.RFC3339), values.Get("start_datetime"))
	assert.Equal(now.Format(time.RFC3339), values.Get("end_datetime"))
}
