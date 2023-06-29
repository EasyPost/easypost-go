package easypost

import "time"

type DateTime time.Time

func (dt *DateTime) UnmarshalJSON(b []byte) (err error) {
	// if string is empty, return nil
	if string(b) == "null" {
		return
	}

	var t time.Time

	// try to parse
	// 2006-01-02
	t, err = time.Parse(`"2006-01-02"`, string(b))
	if err == nil {
		*dt = DateTime(t)
		return
	}

	// try to parse
	// 2006-01-02T15:04:05Z
	t, err = time.Parse(`"2006-01-02T15:04:05Z"`, string(b))
	if err == nil {
		*dt = DateTime(t)
		return
	}

	// try to parse as RFC3339 (default for time.Time)
	// 2006-01-02T15:04:05Z07:00
	t, err = time.Parse(`"`+time.RFC3339+`"`, string(b))
	if err == nil {
		*dt = DateTime(t)
		return
	}

	// try to parse as RFC3339Nano
	// 2006-01-02T15:04:05.999999999Z07:00
	t, err = time.Parse(`"`+time.RFC3339Nano+`"`, string(b))
	if err == nil {
		*dt = DateTime(t)
		return
	}

	// try to parse as RFC1123
	// Mon, 02 Jan 2006 15:04:05 MST
	t, err = time.Parse(`"`+time.RFC1123+`"`, string(b))
	if err == nil {
		*dt = DateTime(t)
		return
	}

	// try to parse as RFC1123Z
	// Mon, 02 Jan 2006 15:04:05 -0700
	t, err = time.Parse(`"`+time.RFC1123Z+`"`, string(b))
	if err == nil {
		*dt = DateTime(t)
		return
	}

	// try to parse as RFC822
	// 02 Jan 06 15:04 MST
	t, err = time.Parse(`"`+time.RFC822+`"`, string(b))
	if err == nil {
		*dt = DateTime(t)
		return
	}

	// try to parse as RFC822Z
	// 02 Jan 06 15:04 -0700
	t, err = time.Parse(`"`+time.RFC822Z+`"`, string(b))
	if err == nil {
		*dt = DateTime(t)
		return
	}

	// try to parse as RFC850
	// Monday, 02-Jan-06 15:04:05 MST
	t, err = time.Parse(`"`+time.RFC850+`"`, string(b))
	if err == nil {
		*dt = DateTime(t)
		return
	}

	// last ditch effort, fallback to whatever the JSON marshaller thinks
	err = t.UnmarshalJSON(b)
	if err != nil {
		return err
	}
	*dt = DateTime(t)
	return nil
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	return time.Time(dt).MarshalJSON()
}

func (dt DateTime) String() string {
	return time.Time(dt).String()
}

func (dt DateTime) AsTime() time.Time {
	return time.Time(dt)
}

// NewDateTime returns the DateTime corresponding to
//
//	yyyy-mm-dd hh:mm:ss + nsec nanoseconds
//
// in the appropriate zone for that time in the given location.
//
// The month, day, hour, min, sec, and nsec values may be outside
// their usual ranges and will be normalized during the conversion.
// For example, October 32 converts to November 1.
//
// A daylight savings time transition skips or repeats times.
// For example, in the United States, March 13, 2011 2:15am never occurred,
// while November 6, 2011 1:15am occurred twice. In such cases, the
// choice of time zone, and therefore the time, is not well-defined.
// NewDateTime returns a time that is correct in one of the two zones involved
// in the transition, but it does not guarantee which.
//
// NewDateTime panics if loc is nil.
func NewDateTime(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) DateTime {
	timeDate := time.Date(year, month, day, hour, min, sec, nsec, loc)
	return DateTimeFromTime(timeDate)
}

// DateTimeFromTime construct a DateTime from a time.Time
func DateTimeFromTime(t time.Time) DateTime {
	return DateTime(t)
}
