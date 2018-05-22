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
				{parseTD("00:00"), 40, MgPerDeciLiter},
			},
		},
		{
			[]byte{0x00, 0x14, 0x02, 0x19, 0x04, 0x1E, 0x06, 0x23, 0x08, 0x28, 0x0A, 0x2D, 0x0C, 0x32, 0x0E, 0x37, 0x00, 0x00, 0x00},
			MgPerDeciLiter,
			[]InsulinSensitivity{
				{parseTD("00:00"), 20, MgPerDeciLiter},
				{parseTD("01:00"), 25, MgPerDeciLiter},
				{parseTD("02:00"), 30, MgPerDeciLiter},
				{parseTD("03:00"), 35, MgPerDeciLiter},
				{parseTD("04:00"), 40, MgPerDeciLiter},
				{parseTD("05:00"), 45, MgPerDeciLiter},
				{parseTD("06:00"), 50, MgPerDeciLiter},
				{parseTD("07:00"), 55, MgPerDeciLiter},
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
