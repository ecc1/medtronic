package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func TestGlucoseTargets(t *testing.T) {
	cases := []struct {
		data   []byte
		units  GlucoseUnitsType
		family Family
		sched  GlucoseTargetSchedule
	}{
		{
			parseBytes("00 50 78"),
			MgPerDeciLiter,
			22,
			GlucoseTargetSchedule{
				{parseTD("00:00"), 80, 120, MgPerDeciLiter},
			},
		},
		{
			parseBytes("00 2C 43"),
			MMolPerLiter,
			22,
			GlucoseTargetSchedule{
				{parseTD("00:00"), 4400, 6700, MMolPerLiter},
			},
		},
		{
			parseBytes("00 64 64 02 65 65 04 66 66 06 67 67 08 68 68 0A 69 69 0C 6A 6A 0E 6B 6B 00 00 00"),
			MgPerDeciLiter,
			23,
			GlucoseTargetSchedule{
				{parseTD("00:00"), 100, 100, MgPerDeciLiter},
				{parseTD("01:00"), 101, 101, MgPerDeciLiter},
				{parseTD("02:00"), 102, 102, MgPerDeciLiter},
				{parseTD("03:00"), 103, 103, MgPerDeciLiter},
				{parseTD("04:00"), 104, 104, MgPerDeciLiter},
				{parseTD("05:00"), 105, 105, MgPerDeciLiter},
				{parseTD("06:00"), 106, 106, MgPerDeciLiter},
				{parseTD("07:00"), 107, 107, MgPerDeciLiter},
			},
		},
		{
			parseBytes("00 64 02 65 04 66 06 67 08 68 0A 69 0C 6A 0E 6B 00 00"),
			MgPerDeciLiter,
			12,
			GlucoseTargetSchedule{
				{parseTD("00:00"), 100, 100, MgPerDeciLiter},
				{parseTD("01:00"), 101, 101, MgPerDeciLiter},
				{parseTD("02:00"), 102, 102, MgPerDeciLiter},
				{parseTD("03:00"), 103, 103, MgPerDeciLiter},
				{parseTD("04:00"), 104, 104, MgPerDeciLiter},
				{parseTD("05:00"), 105, 105, MgPerDeciLiter},
				{parseTD("06:00"), 106, 106, MgPerDeciLiter},
				{parseTD("07:00"), 107, 107, MgPerDeciLiter},
			},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			s := decodeGlucoseTargetSchedule(c.data, c.units, c.family)
			if !reflect.DeepEqual(s, c.sched) {
				t.Errorf("decodeGlucoseTargetSchedule(% X, %v, %d) == %+v, want %+v", c.data, c.units, c.family, s, c.sched)
			}
		})
	}
}

func TestGlucoseTargetAt(t *testing.T) {
	cases := []struct {
		sched  GlucoseTargetSchedule
		at     time.Time
		target GlucoseTarget
	}{
		{
			GlucoseTargetSchedule{
				{TimeOfDay(0), 80, 120, MgPerDeciLiter},
			},
			parseTime("2016-11-06T23:00:00"),
			GlucoseTarget{TimeOfDay(0), 80, 120, MgPerDeciLiter},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			target := c.sched.GlucoseTargetAt(c.at)
			if !reflect.DeepEqual(target, c.target) {
				t.Errorf("%v.GlucoseTargetAt(%v) == %+v, want %+v", c.sched, c.at, target, c.target)
			}
		})
	}
}
