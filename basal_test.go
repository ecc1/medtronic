package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func TestBasalRates(t *testing.T) {
	cases := []struct {
		data  []byte
		sched BasalRateSchedule
	}{
		{
			[]byte{0x06, 0x28, 0x00, 0x00, 0x30, 0x00, 0x12},
			[]BasalRate{
				{parseTD("00:00"), 1000},
				{parseTD("09:00"), 1200},
			},
		},
		{
			[]byte{0x0C, 0x20, 0x00, 0x00, 0x26, 0x00, 0x0D, 0x2C, 0x00, 0x13, 0x26, 0x00, 0x1C},
			[]BasalRate{
				{parseTD("00:00"), 800},
				{parseTD("06:30"), 950},
				{parseTD("09:30"), 1100},
				{parseTD("14:00"), 950},
			},
		},
	}
	for _, c := range cases {
		s := decodeBasalRateSchedule(c.data)
		if !reflect.DeepEqual(s, c.sched) {
			t.Errorf("decodeBasalRateSchedule(% X) == %+v, want %+v", c.data, s, c.sched)
		}
	}
}

func TestBasalRateAt(t *testing.T) {
	cases := []struct {
		sched  BasalRateSchedule
		at     time.Time
		target BasalRate
	}{
		{
			[]BasalRate{
				{parseTD("00:00"), 1000},
			},
			parseTime("2016-11-06T23:00:00"),
			BasalRate{parseTD("00:00"), 1000},
		},
	}
	for _, c := range cases {
		target := c.sched.BasalRateAt(c.at)
		if !reflect.DeepEqual(target, c.target) {
			t.Errorf("%v.BasalRateAt(%v) == %+v, want %+v", c.sched, c.at, target, c.target)
		}
	}
}

func parseTD(s string) TimeOfDay {
	t, err := parseTimeOfDay(s)
	if err != nil {
		panic(err)
	}
	return t
}
