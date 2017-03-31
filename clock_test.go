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
			[]byte{0x07, 0x17, 0x2D, 0x01, 0x07, 0xE0, 0x0B, 0x06},
			parseTime("2016-11-06T23:45:01"),
		},
		{
			[]byte{0x07, 0x09, 0x16, 0x3B, 0x07, 0xE1, 0x0C, 0x1D},
			parseTime("2017-12-29T09:22:59"),
		},
	}
	for _, c := range cases {
		tt := decodeClock(c.data)
		if !tt.Equal(c.t) {
			t.Errorf("decodeClock(% X) == %v, want %v", c.data, tt, c.t)
		}
	}
}
