package medtronic

import (
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	cases := []struct {
		t uint8
		d time.Duration
	}{
		{0, 0},
		{1, 30 * time.Minute},
		{3, 90 * time.Minute},
		{4, 2 * time.Hour},
	}
	for _, c := range cases {
		d := scheduleToDuration(c.t)
		if d != c.d {
			t.Errorf("scheduleToDuration(%d) == %v, want %v", c.t, d, c.d)
		}
	}
}

func parseTime(s string) time.Time {
	const layout = "2006-01-02T15:04:05.999999999"
	t, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

func TestSinceMidnight(t *testing.T) {
	cases := []struct {
		t time.Time
		d time.Duration
	}{
		{parseTime("2015-01-01T09:00:00"), 9 * time.Hour},
		{parseTime("2016-03-15T10:00:00.5"), 10*time.Hour + 500*time.Millisecond},
		{parseTime("2016-06-15T20:30:00"), 20*time.Hour + 30*time.Minute},
		{parseTime("2010-11-30T23:59:59.999"), 24*time.Hour - time.Millisecond},
	}
	for _, c := range cases {
		d := sinceMidnight(c.t)
		if d != c.d {
			t.Errorf("sinceMidnight(%v) == %v, want %v", c.t, d, c.d)
		}
	}
}

func TestDecodeTimestamp(t *testing.T) {
	cases := []struct {
		b []byte
		t time.Time
	}{
		{[]byte{0x1F, 0x40, 0x00, 0x01, 0x05}, parseTime("2005-01-01T00:00:31")},
		{[]byte{0x75, 0xB7, 0x13, 0x04, 0x10}, parseTime("2016-06-04T19:55:53")},
		{[]byte{0x5D, 0xB3, 0x0F, 0x06, 0x10}, parseTime("2016-06-06T15:51:29")},
		{[]byte{0x40, 0x94, 0x12, 0x0F, 0x10}, parseTime("2016-06-15T18:20:00")},
	}
	for _, c := range cases {
		ts := decodeTimestamp(c.b)
		if !ts.Equal(c.t) {
			t.Errorf("decodeTimestamp(% X) == %v, want %v", c.b, ts, c.t)
		}
	}
}

func TestDecodeDate(t *testing.T) {
	cases := []struct {
		b []byte
		t time.Time
	}{
		{[]byte{0xBF, 0x0F}, parseTime("2015-10-31T00:00:00")},
		{[]byte{0x78, 0x10}, parseTime("2016-06-24T00:00:00")},
	}
	for _, c := range cases {
		ts := decodeDate(c.b)
		if !ts.Equal(c.t) {
			t.Errorf("decodeDate(% X) == %v, want %v", c.b, ts, c.t)
		}
	}
}
