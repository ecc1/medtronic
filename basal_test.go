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
			[]byte{0x01, 0x28, 0x00, 0x00, 0x30, 0x00, 0x12},
			[]BasalRate{
				{parseTD("00:00"), 1000},
				{parseTD("09:00"), 1200},
			},
		},
		{
			[]byte{0x01, 0x20, 0x00, 0x00, 0x26, 0x00, 0x0D, 0x2C, 0x00, 0x13, 0x26, 0x00, 0x1C},
			[]BasalRate{
				{parseTD("00:00"), 800},
				{parseTD("06:30"), 950},
				{parseTD("09:30"), 1100},
				{parseTD("14:00"), 950},
			},
		},
		{
			[]byte{0x01, 0x28, 0x00, 0x00, 0x28, 0x00, 0x06, 0x2C, 0x00, 0x0C, 0x30, 0x00, 0x14, 0x30, 0x00, 0x2C, 0x00, 0x00, 0x00},
			[]BasalRate{
				{parseTD("00:00"), 1000},
				{parseTD("03:00"), 1000},
				{parseTD("06:00"), 1100},
				{parseTD("10:00"), 1200},
				{parseTD("22:00"), 1200},
			},
		},
		{
			[]byte{0x01, 0x14, 0x00, 0x00, 0x15, 0x00, 0x02, 0x16, 0x00, 0x04, 0x17, 0x00, 0x06, 0x18, 0x00, 0x08, 0x19, 0x00, 0x0A, 0x1A, 0x00, 0x0C, 0x1B, 0x00, 0x0E, 0x1C, 0x00, 0x10, 0x1D, 0x00, 0x12, 0x1E, 0x00, 0x14, 0x1F, 0x00, 0x16, 0x20, 0x00, 0x18, 0x21, 0x00, 0x1A, 0x22, 0x00, 0x1C, 0x23, 0x00, 0x1E, 0x24, 0x00, 0x20, 0x25, 0x00, 0x22, 0x26, 0x00, 0x24, 0x27, 0x00, 0x26, 0x28, 0x00, 0x28, 0x2C},
			[]BasalRate{
				{parseTD("00:00"), 500},
				{parseTD("01:00"), 525},
				{parseTD("02:00"), 550},
				{parseTD("03:00"), 575},
				{parseTD("04:00"), 600},
				{parseTD("05:00"), 625},
				{parseTD("06:00"), 650},
				{parseTD("07:00"), 675},
				{parseTD("08:00"), 700},
				{parseTD("09:00"), 725},
				{parseTD("10:00"), 750},
				{parseTD("11:00"), 775},
				{parseTD("12:00"), 800},
				{parseTD("13:00"), 825},
				{parseTD("14:00"), 850},
				{parseTD("15:00"), 875},
				{parseTD("16:00"), 900},
				{parseTD("17:00"), 925},
				{parseTD("18:00"), 950},
				{parseTD("19:00"), 975},
				{parseTD("20:00"), 1000},
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
				{parseTD("00:00"), 1000},
			},
			parseTime("2016-11-06T23:00:00"),
			BasalRate{parseTD("00:00"), 1000},
		},
	}
	for _, c := range cases {
		target := c.sched.BasalRateAt(c.at)
		if !reflect.DeepEqual(target, c.target) {
			t.Errorf("%v.BasalRateAt(%v) == %+v, want %+v", c.sched, c.at, target, c.target)
		}
	}
}

func parseTD(s string) TimeOfDay {
	t, err := parseTimeOfDay(s)
	if err != nil {
		panic(err)
	}
	return t
}
