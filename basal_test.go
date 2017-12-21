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
				{durationToTimeOfDay(0), 1000},
				{durationToTimeOfDay(9 * time.Hour), 1200},
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
				{durationToTimeOfDay(0), 1000},
			},
			parseTime("2016-11-06T23:00:00"),
			BasalRate{durationToTimeOfDay(0), 1000},
		},
	}
	for _, c := range cases {
		target := c.sched.BasalRateAt(c.at)
		if !reflect.DeepEqual(target, c.target) {
			t.Errorf("%v.BasalRateAt(%v) == %+v, want %+v", c.sched, c.at, target, c.target)
		}
	}
}
