package medtronic

import (
	"reflect"
	"testing"
	"time"
)

// Ratio values are 10x grams/unit or 100x units/exchange (see carbratios.go).
func TestCarbRatios(t *testing.T) {
	cases := []struct {
		name   string
		data   []byte
		units  CarbUnitsType
		family Family
		sched  CarbRatioSchedule
	}{
		{
			"2_carbratios",
			parseBytes("00 06 12 08"),
			Grams,
			22,
			CarbRatioSchedule{
				{parseTD("00:00"), 60, Grams},
				{parseTD("09:00"), 80, Grams},
			},
		},
		{
			"8_carbratios",
			parseBytes("00 0A 02 0B 04 0C 06 0D 08 0E 0A 0F 0C 10 0E 11 00 00 00 00"),
			Grams,
			12,
			CarbRatioSchedule{
				{parseTD("00:00"), 100, Grams},
				{parseTD("01:00"), 110, Grams},
				{parseTD("02:00"), 120, Grams},
				{parseTD("03:00"), 130, Grams},
				{parseTD("04:00"), 140, Grams},
				{parseTD("05:00"), 150, Grams},
				{parseTD("06:00"), 160, Grams},
				{parseTD("07:00"), 170, Grams},
			},
		},
		{
			"12_carbratios",
			parseBytes("00 05 01 06 02 07 03 08 04 09 05 0A 06 0B 07 0C 00 00"),
			Grams,
			22,
			CarbRatioSchedule{
				{parseTD("00:00"), 50, Grams},
				{parseTD("00:30"), 60, Grams},
				{parseTD("01:00"), 70, Grams},
				{parseTD("01:30"), 80, Grams},
				{parseTD("02:00"), 90, Grams},
				{parseTD("02:30"), 100, Grams},
				{parseTD("03:00"), 110, Grams},
				{parseTD("03:30"), 120, Grams},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := decodeCarbRatioSchedule(c.data, c.units, c.family)
			if !reflect.DeepEqual(s, c.sched) {
				t.Errorf("decodeCarbRatioSchedule(% X, %v, %d) == %+v, want %+v", c.data, c.units, c.family, s, c.sched)
			}
		})
	}
}

func TestCarbRatioAt(t *testing.T) {
	cases := []struct {
		sched  CarbRatioSchedule
		at     time.Time
		target CarbRatio
	}{
		{
			CarbRatioSchedule{
				{parseTD("00:00"), 60, Grams},
			},
			parseTime("2016-11-06T23:00:00"),
			CarbRatio{parseTD("00:00"), 60, Grams},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			target := c.sched.CarbRatioAt(c.at)
			if !reflect.DeepEqual(target, c.target) {
				t.Errorf("%v.CarbRatioAt(%v) == %+v, want %+v", c.sched, c.at, target, c.target)
			}
		})
	}
}
