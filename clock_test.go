package medtronic

import (
	"testing"
	"time"
)

func TestClock(t *testing.T) {
	cases := []struct {
		data []byte
		t    time.Time
	}{
		{
			parseBytes("07 17 2D 01 07 E0 0B 06"),
			parseTime("2016-11-06T23:45:01"),
		},
		{
			parseBytes("07 09 16 3B 07 E1 0C 1D"),
			parseTime("2017-12-29T09:22:59"),
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			tt := decodeClock(c.data)
			if !tt.Equal(c.t) {
				t.Errorf("decodeClock(% X) == %v, want %v", c.data, tt, c.t)
			}
		})
	}
}
