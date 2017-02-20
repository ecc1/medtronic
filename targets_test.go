package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func TestGlucoseTargets(t *testing.T) {
	cases := []struct {
		data  []byte
		units GlucoseUnitsType
		sched GlucoseTargetSchedule
	}{
		{
			[]byte{0, 80, 120},
			MgPerDeciLiter,
			[]GlucoseTarget{
				{durationToTimeOfDay(0), 80, 120, MgPerDeciLiter},
			},
		},
		{
			[]byte{0, 44, 67},
			MMolPerLiter,
			[]GlucoseTarget{
				{durationToTimeOfDay(0), 4400, 6700, MMolPerLiter},
			},
		},
	}
	for _, c := range cases {
		s := decodeGlucoseTargetSchedule(c.data, c.units)
		if !reflect.DeepEqual(s, c.sched) {
			t.Errorf("decodeGlucoseTargetSchedule(%X, %v) == %v, want %v", c.data, c.units, s, c.sched)
		}
	}
}

func TestGlucoseTargetAt(t *testing.T) {
	cases := []struct {
		sched  GlucoseTargetSchedule
		at     time.Time
		target GlucoseTarget
	}{
		{
			[]GlucoseTarget{
				{durationToTimeOfDay(0), 80, 120, MgPerDeciLiter},
			},
			parseTime("2016-11-06T23:00:00"),
			GlucoseTarget{durationToTimeOfDay(0), 80, 120, MgPerDeciLiter},
		},
	}
	for _, c := range cases {
		target := c.sched.GlucoseTargetAt(c.at)
		if !reflect.DeepEqual(target, c.target) {
			t.Errorf("%v.GlucoseTargetAt(%v) == %v, want %v", c.sched, c.at, target, c.target)
		}
	}
}
