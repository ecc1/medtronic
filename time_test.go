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
	t, err := time.Parse(time.RFC3339Nano, s)
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
		{parseTime("2015-01-01T09:00:00Z"), 9 * time.Hour},
		{parseTime("2016-03-15T10:00:00.5Z"), 10*time.Hour + 500*time.Millisecond},
		{parseTime("2016-06-15T20:30:00-04:00"), 20*time.Hour + 30*time.Minute},
		{parseTime("2010-11-30T23:59:59.999+06:00"), 24*time.Hour - time.Millisecond},
	}
	for _, c := range cases {
		d := sinceMidnight(c.t)
		if d != c.d {
			t.Errorf("sinceMidnight(%v) == %v, want %v", c.t, d, c.d)
		}
	}
}
