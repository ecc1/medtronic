package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func TestInsulinSensitivities(t *testing.T) {
	cases := []struct {
		data  []byte
		units GlucoseUnitsType
		sched InsulinSensitivitySchedule
	}{
		{
			parseBytes("00 28"),
			MgPerDeciLiter,
			InsulinSensitivitySchedule{
				{parseTD("00:00"), 40, MgPerDeciLiter},
			},
		},
		{
			parseBytes("00 32 2C 3C"),
			MgPerDeciLiter,
			InsulinSensitivitySchedule{
				{parseTD("00:00"), 50, MgPerDeciLiter},
				{parseTD("22:00"), 60, MgPerDeciLiter},
			},
		},
		{
			parseBytes("00 14 02 19 04 1E 06 23 08 28 0A 2D 0C 32 0E 37 00 00 00"),
			MgPerDeciLiter,
			InsulinSensitivitySchedule{
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
		t.Run("", func(t *testing.T) {
			s := decodeInsulinSensitivitySchedule(c.data, c.units)
			if !reflect.DeepEqual(s, c.sched) {
				t.Errorf("decodeInsulinSensitivitySchedule(% X, %v) == %+v, want %+v", c.data, c.units, s, c.sched)
			}
		})
	}
}

func TestInsulinSensitivityAt(t *testing.T) {
	cases := []struct {
		sched  InsulinSensitivitySchedule
		at     time.Time
		target InsulinSensitivity
	}{
		{
			InsulinSensitivitySchedule{
				{TimeOfDay(0), 40, MgPerDeciLiter},
			},
			parseTime("2016-11-06T23:00:00"),
			InsulinSensitivity{TimeOfDay(0), 40, MgPerDeciLiter},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			target := c.sched.InsulinSensitivityAt(c.at)
			if !reflect.DeepEqual(target, c.target) {
				t.Errorf("%v.InsulinSensitivityAt(%v) == %+v, want %+v", c.sched, c.at, target, c.target)
			}
		})
	}
}
