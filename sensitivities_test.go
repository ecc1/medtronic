package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func TestInsulinSensitivitys(t *testing.T) {
	cases := []struct {
		data  []byte
		units GlucoseUnitsType
		sched InsulinSensitivitySchedule
	}{
		{
			[]byte{0x00, 0x28},
			MgPerDeciLiter,
			[]InsulinSensitivity{
				{durationToTimeOfDay(0), 40, MgPerDeciLiter},
			},
		},
	}
	for _, c := range cases {
		s := decodeInsulinSensitivitySchedule(c.data, c.units)
		if !reflect.DeepEqual(s, c.sched) {
			t.Errorf("decodeInsulinSensitivitySchedule(% X, %v) == %+v, want %+v", c.data, c.units, s, c.sched)
		}
	}
}

func TestInsulinSensitivityAt(t *testing.T) {
	cases := []struct {
		sched  InsulinSensitivitySchedule
		at     time.Time
		target InsulinSensitivity
	}{
		{
			[]InsulinSensitivity{
				{durationToTimeOfDay(0), 40, MgPerDeciLiter},
			},
			parseTime("2016-11-06T23:00:00"),
			InsulinSensitivity{durationToTimeOfDay(0), 40, MgPerDeciLiter},
		},
	}
	for _, c := range cases {
		target := c.sched.InsulinSensitivityAt(c.at)
		if !reflect.DeepEqual(target, c.target) {
			t.Errorf("%v.InsulinSensitivityAt(%v) == %+v, want %+v", c.sched, c.at, target, c.target)
		}
	}
}
