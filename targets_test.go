package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func TestGlucoseTargets(t *testing.T) {
	cases := []struct {
		name   string
		data   []byte
		units  GlucoseUnitsType
		family Family
		sched  GlucoseTargetSchedule
	}{
		{
			"1_target",
			[]byte{0x00, 0x50, 0x78},
			MgPerDeciLiter,
			22,
			GlucoseTargetSchedule{
				{parseTD("00:00"), 80, 120, MgPerDeciLiter},
			},
		},
		{
			"1_target_mmol",
			[]byte{0x00, 0x2C, 0x43},
			MMolPerLiter,
			22,
			GlucoseTargetSchedule{
				{parseTD("00:00"), 4400, 6700, MMolPerLiter},
			},
		},
		{
			"8_targets_x23",
			[]byte{0x00, 0x64, 0x64, 0x02, 0x65, 0x65, 0x04, 0x66, 0x66, 0x06, 0x67, 0x67, 0x08, 0x68, 0x68, 0x0A, 0x69, 0x69, 0x0C, 0x6A, 0x6A, 0x0E, 0x6B, 0x6B, 0x00, 0x00, 0x00},
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
			"8_targets_x12",
			[]byte{0x00, 0x64, 0x02, 0x65, 0x04, 0x66, 0x06, 0x67, 0x08, 0x68, 0x0A, 0x69, 0x0C, 0x6A, 0x0E, 0x6B, 0x00, 0x00},
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
		t.Run(c.name, func(t *testing.T) {
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
				{durationToTimeOfDay(0), 80, 120, MgPerDeciLiter},
			},
			parseTime("2016-11-06T23:00:00"),
			GlucoseTarget{durationToTimeOfDay(0), 80, 120, MgPerDeciLiter},
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
