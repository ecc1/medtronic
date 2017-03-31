package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func TestCarbRatios(t *testing.T) {
	cases := []struct {
		data  []byte
		units CarbUnitsType
		newer bool
		sched CarbRatioSchedule
	}{
		{
			[]byte{0, 6, 18, 8},
			Grams,
			false,
			[]CarbRatio{
				{durationToTimeOfDay(0), 60, Grams},
				{durationToTimeOfDay(9 * time.Hour), 80, Grams},
			},
		},
	}
	for _, c := range cases {
		s := decodeCarbRatioSchedule(c.data, c.units, c.newer)
		if !reflect.DeepEqual(s, c.sched) {
			t.Errorf("decodeCarbRatioSchedule(% X, %v, %v) == %v, want %v", c.data, c.units, c.newer, s, c.sched)
		}
	}
}

func TestCarbRatioAt(t *testing.T) {
	cases := []struct {
		sched  CarbRatioSchedule
		at     time.Time
		target CarbRatio
	}{
		{
			[]CarbRatio{
				{durationToTimeOfDay(0), 60, Grams},
			},
			parseTime("2016-11-06T23:00:00"),
			CarbRatio{durationToTimeOfDay(0), 60, Grams},
		},
	}
	for _, c := range cases {
		target := c.sched.CarbRatioAt(c.at)
		if !reflect.DeepEqual(target, c.target) {
			t.Errorf("%v.CarbRatioAt(%v) == %v, want %v", c.sched, c.at, target, c.target)
		}
	}
}
