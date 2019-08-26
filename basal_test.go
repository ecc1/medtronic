package medtronic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestEncodeBasalRate(t *testing.T) {
	cases := []struct {
		family Family
		rate   Insulin
		actual Insulin
	}{
		{22, 1000, 1000},
		{22, 2550, 2550},
		{23, 575, 575},
		{23, 2575, 2550},
		{23, 11250, 11200},
	}
	log.SetOutput(ioutil.Discard)
	for _, c := range cases {
		name := fmt.Sprintf("%d_%d", c.family, c.rate)
		t.Run(name, func(t *testing.T) {
			r, err := encodeBasalRate("basal", c.rate, c.family)
			if err != nil {
				t.Errorf("encodeBasalRate(%d, %d) raised error (%v)", c.rate, c.family, err)
			}
			a := Insulin(r) * milliUnitsPerStroke(23)
			if a != c.actual {
				t.Errorf("encodeBasalRate(%v, %d) == %d, want %d", c.rate, c.family, a, c.actual)
			}
		})
	}
}

func TestBasalRates(t *testing.T) {
	cases := []struct {
		family Family
		data   []byte
		sched  BasalRateSchedule
	}{
		{
			22,
			parseBytes("28 00 00 30 00 12"),
			BasalRateSchedule{
				{parseTD("00:00"), 1000},
				{parseTD("09:00"), 1200},
			},
		},
		{
			22,
			parseBytes("28 00 00 40 01 08 28 00 2C"),
			BasalRateSchedule{
				{parseTD("00:00"), 1000},
				{parseTD("04:00"), 8000},
				{parseTD("22:00"), 1000},
			},
		},
		{
			23,
			parseBytes("20 00 00 26 00 0D 2C 00 13 26 00 1C"),
			BasalRateSchedule{
				{parseTD("00:00"), 800},
				{parseTD("06:30"), 950},
				{parseTD("09:30"), 1100},
				{parseTD("14:00"), 950},
			},
		},
		{
			22,
			parseBytes("28 00 00 28 00 06 2C 00 0C 30 00 14 30 00 2C"),
			BasalRateSchedule{
				{parseTD("00:00"), 1000},
				{parseTD("03:00"), 1000},
				{parseTD("06:00"), 1100},
				{parseTD("10:00"), 1200},
				{parseTD("22:00"), 1200},
			},
		},
		{
			22,
			parseBytes("00 00 00 04 00 02 08 00 04 0C 00 06 10 00 08 14 00 0A 18 00 0C 1C 00 0E 20 00 10 24 00 12 28 00 14 2C 00 16 30 00 18 34 00 1A 38 00 1C 3C 00 1E 40 00 20 44 00 22 48 00 24 4C 00 26 50 00 28 54 00 2A 58 00 2C 5C 00 2E"),
			BasalRateSchedule{
				{parseTD("00:00"), 0},
				{parseTD("01:00"), 100},
				{parseTD("02:00"), 200},
				{parseTD("03:00"), 300},
				{parseTD("04:00"), 400},
				{parseTD("05:00"), 500},
				{parseTD("06:00"), 600},
				{parseTD("07:00"), 700},
				{parseTD("08:00"), 800},
				{parseTD("09:00"), 900},
				{parseTD("10:00"), 1000},
				{parseTD("11:00"), 1100},
				{parseTD("12:00"), 1200},
				{parseTD("13:00"), 1300},
				{parseTD("14:00"), 1400},
				{parseTD("15:00"), 1500},
				{parseTD("16:00"), 1600},
				{parseTD("17:00"), 1700},
				{parseTD("18:00"), 1800},
				{parseTD("19:00"), 1900},
				{parseTD("20:00"), 2000},
				{parseTD("21:00"), 2100},
				{parseTD("22:00"), 2200},
				{parseTD("23:00"), 2300},
			},
		},
		{
			23,
			parseBytes("42 00 00 40 00 02 3A 00 04 3E 00 06 36 00 08 46 00 0A 4A 00 0C 4C 00 0E 4E 00 10 4C 00 12 4E 00 14 50 00 16 50 00 18 4E 00 1A 4A 00 1C 4A 00 1E 4A 00 20 4C 00 22 4A 00 24 4A 00 26 48 00 28 46 00 2A 4A 00 2C 4A 00 2E"),
			BasalRateSchedule{
				{parseTD("00:00"), 1650},
				{parseTD("01:00"), 1600},
				{parseTD("02:00"), 1450},
				{parseTD("03:00"), 1550},
				{parseTD("04:00"), 1350},
				{parseTD("05:00"), 1750},
				{parseTD("06:00"), 1850},
				{parseTD("07:00"), 1900},
				{parseTD("08:00"), 1950},
				{parseTD("09:00"), 1900},
				{parseTD("10:00"), 1950},
				{parseTD("11:00"), 2000},
				{parseTD("12:00"), 2000},
				{parseTD("13:00"), 1950},
				{parseTD("14:00"), 1850},
				{parseTD("15:00"), 1850},
				{parseTD("16:00"), 1850},
				{parseTD("17:00"), 1900},
				{parseTD("18:00"), 1850},
				{parseTD("19:00"), 1850},
				{parseTD("20:00"), 1800},
				{parseTD("21:00"), 1750},
				{parseTD("22:00"), 1850},
				{parseTD("23:00"), 1850},
			},
		},
	}
	for _, c := range cases {
		name := fmt.Sprintf("%d_%d", c.family, len(c.sched))
		t.Run("decode_"+name, func(t *testing.T) {
			s := decodeBasalRateSchedule(c.data)
			if !reflect.DeepEqual(s, c.sched) {
				t.Errorf("decodeBasalRateSchedule(% X) == %+v, want %+v", c.data, s, c.sched)
			}
		})
		t.Run("encode_"+name, func(t *testing.T) {
			data, err := encodeBasalRateSchedule(c.sched, c.family)
			if err != nil {
				t.Errorf("encodeBasalRateSchedule(%+v) raised error (%v)", c.sched, err)
			}
			if !bytes.Equal(data, c.data) {
				t.Errorf("encodeBasalRateSchedule(%+v) == % X, want % X", c.sched, data, c.data)
			}
		})
	}
}

func TestBasalRateAt(t *testing.T) {
	cases := []struct {
		sched  BasalRateSchedule
		at     time.Time
		target BasalRate
	}{
		{
			BasalRateSchedule{
				{parseTD("00:00"), 1000},
			},
			parseTime("2016-11-06T23:00:00"),
			BasalRate{parseTD("00:00"), 1000},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			target := c.sched.BasalRateAt(c.at)
			if !reflect.DeepEqual(target, c.target) {
				t.Errorf("%v.BasalRateAt(%v) == %+v, want %+v", c.sched, c.at, target, c.target)
			}
		})
	}
}
