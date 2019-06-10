package medtronic

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeOfDay(t *testing.T) {
	cases := []struct {
		s   string
		t   TimeOfDay
		err error
	}{
		{"00:00", 0, nil},
		{"12:00", Duration(12 * time.Hour).TimeOfDay(), nil},
		{"23:59", Duration(24*time.Hour - 1*time.Minute).TimeOfDay(), nil},
		{"01:02:03", 0, fmt.Errorf("")},
		{"24:00", 0, fmt.Errorf("")},
		{"01:60", 0, fmt.Errorf("")},
	}
	for _, c := range cases {
		t.Run(c.s, func(t *testing.T) {
			if c.err == nil {
				s := c.t.String()
				if s != c.s {
					t.Errorf("%v.String() == %v, want %v", c.t, s, c.s)
				}
			}
			td, err := ParseTimeOfDay(c.s)
			if err == nil {
				if c.err == nil {
					if td == c.t {
						return
					}
					t.Errorf("ParseTimeOfDay(%s) == %v, want %v", c.s, td, c.t)
				} else {
					t.Errorf("ParseTimeOfDay(%s) == %v, want error", c.s, td)
				}
			} else {
				if c.err != nil {
					return
				}
				t.Errorf("ParseTimeOfDay(%s) == %v, want %v", c.s, err, c.t)
			}
		})
	}

}

func TestHalfHours(t *testing.T) {
	cases := []struct {
		h uint8
		d time.Duration
	}{
		{0, 0},
		{1, 30 * time.Minute},
		{3, 90 * time.Minute},
		{4, 2 * time.Hour},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("from%d", c.h), func(t *testing.T) {
			d := halfHoursToDuration(c.h)
			if d != Duration(c.d) {
				t.Errorf("halfHoursToDuration(%d) == %v, want %v", c.h, d, c.d)
			}
		})
		t.Run(fmt.Sprintf("to%d", c.h), func(t *testing.T) {
			h := TimeOfDay(c.d).HalfHours()
			if h != c.h {
				t.Errorf("HalfHours(%v) == %d, want %d", c.d, h, c.h)
			}
		})
	}
}

func TestSinceMidnight(t *testing.T) {
	cases := []struct {
		t time.Time
		d TimeOfDay
	}{
		{parseTime("2015-01-01T09:00"), Duration(9 * time.Hour).TimeOfDay()},
		{parseTime("2016-03-15T10:00:00.5"), Duration(10*time.Hour + 500*time.Millisecond).TimeOfDay()},
		{parseTime("2016-06-15T20:30"), Duration(20*time.Hour + 30*time.Minute).TimeOfDay()},
		{parseTime("2010-11-30T23:59:59.999"), Duration(24*time.Hour - time.Millisecond).TimeOfDay()},
		// DST changes
		{parseTime("2016-03-13T01:00"), Duration(1 * time.Hour).TimeOfDay()},
		{parseTime("2016-03-13T03:00"), Duration(3 * time.Hour).TimeOfDay()},
		{parseTime("2016-03-13T12:00"), Duration(12 * time.Hour).TimeOfDay()},
		{parseTime("2016-11-06T01:00"), Duration(1 * time.Hour).TimeOfDay()},
		{parseTime("2016-11-06T02:00"), Duration(2 * time.Hour).TimeOfDay()},
		{parseTime("2016-11-06T03:00"), Duration(3 * time.Hour).TimeOfDay()},
		{parseTime("2016-11-06T23:00"), Duration(23 * time.Hour).TimeOfDay()},
		{parseTime("2016-11-06T23:30"), Duration(23*time.Hour + 30*time.Minute).TimeOfDay()},
	}
	for _, c := range cases {
		t.Run(c.t.Format(time.Kitchen), func(t *testing.T) {
			d := SinceMidnight(c.t)
			if d != c.d {
				// Print TimeOfDay as underlying time.Duration.
				t.Errorf("sinceMidnight(%v) == %v, want %v", c.t, time.Duration(d), time.Duration(c.d))
			}
		})
	}
}

func TestDecodeTime(t *testing.T) {
	cases := []struct {
		b []byte
		t time.Time
	}{
		{parseBytes("1F 40 00 01 05"), parseTime("2005-01-01T00:00:31")},
		{parseBytes("09 A2 0A 15 10"), parseTime("2016-02-21T10:34:09")},
		{parseBytes("42 22 54 65 10"), parseTime("2016-04-05T20:34:02")},
		{parseBytes("79 23 0C 12 10"), parseTime("2016-04-18T12:35:57")},
		{parseBytes("75 B7 13 04 10"), parseTime("2016-06-04T19:55:53")},
		{parseBytes("5D B3 0F 06 10"), parseTime("2016-06-06T15:51:29")},
		{parseBytes("40 94 12 0F 10"), parseTime("2016-06-15T18:20:00")},
		{parseBytes("B1 34 87 6B 12"), parseTime("2018-08-11T07:52:49")},
	}
	for _, c := range cases {
		t.Run(c.t.Format(time.Kitchen), func(t *testing.T) {
			ts := time.Time(decodeTime(c.b))
			if !ts.Equal(c.t) {
				t.Errorf("decodeTime(% X) == %v, want %v", c.b, ts, c.t)
			}
		})
	}
}

func TestDecodeDate(t *testing.T) {
	cases := []struct {
		b []byte
		t time.Time
	}{
		{parseBytes("BF 0F"), parseTime("2015-10-31T00:00")},
		{parseBytes("78 10"), parseTime("2016-06-24T00:00")},
	}
	for _, c := range cases {
		t.Run(c.t.Format("2006-01-02"), func(t *testing.T) {
			ts := time.Time(decodeDate(c.b))
			if !ts.Equal(c.t) {
				t.Errorf("decodeDate(% X) == %v, want %v", c.b, ts, c.t)
			}
		})
	}
}
