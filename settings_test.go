package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func TestSettings(t *testing.T) {
	cases := []struct {
		data   []byte
		family Family
		b      SettingsInfo
	}{
		{
			parseBytes("15 00 01 01 01 00 96 00 8C 00 00 00 00 01 00 00 64 01 04 00 14 00"),
			22,
			SettingsInfo{
				AutoOff:              0,
				InsulinAction:        4 * time.Hour,
				InsulinConcentration: 100,
				MaxBolus:             Insulin(15000),
				MaxBasal:             Insulin(3500),
				RFEnabled:            true,
				TempBasalType:        Absolute,
				SelectedPattern:      0,
			},
		},
		{
			parseBytes("19 00 04 01 05 01 00 8c 00 50 00 00 00 00 00 01 64 01 04 00 14 00 64 01 01 00"),
			23,
			SettingsInfo{
				AutoOff:              0,
				InsulinAction:        4 * time.Hour,
				InsulinConcentration: 100,
				MaxBolus:             Insulin(14000),
				MaxBasal:             Insulin(2000),
				RFEnabled:            false,
				TempBasalType:        Absolute,
				SelectedPattern:      0,
			},
		},
		{
			parseBytes("12 0F 02 01 01 00 64 00 78 00 00 00 00 01 00 00 64 01 00"),
			12,
			SettingsInfo{
				AutoOff:              15 * time.Hour,
				InsulinAction:        6 * time.Hour,
				InsulinConcentration: 100,
				MaxBolus:             Insulin(10000),
				MaxBasal:             Insulin(3000),
				RFEnabled:            true,
				TempBasalType:        Absolute,
				SelectedPattern:      0,
			},
		},
		{
			parseBytes("12 0D 02 01 01 00 64 00 50 00 00 00 00 00 00 00 64 01 00"),
			12,
			SettingsInfo{
				AutoOff:              13 * time.Hour,
				InsulinAction:        6 * time.Hour,
				InsulinConcentration: 100,
				MaxBolus:             Insulin(10000),
				MaxBasal:             Insulin(2000),
				RFEnabled:            false,
				TempBasalType:        Absolute,
				SelectedPattern:      0,
			},
		},
		{
			parseBytes("12 0D 02 01 01 00 64 00 50 00 00 00 00 01 00 00 64 01 01"),
			12,
			SettingsInfo{
				AutoOff:              13 * time.Hour,
				InsulinAction:        8 * time.Hour,
				InsulinConcentration: 100,
				MaxBolus:             Insulin(10000),
				MaxBasal:             Insulin(2000),
				RFEnabled:            true,
				TempBasalType:        Absolute,
				SelectedPattern:      0,
			},
		},
		{
			parseBytes("15 00 03 00 0A 01 7D 01 44 00 00 00 00 01 00 00 64 01 03 00 14 00"),
			22,
			SettingsInfo{
				AutoOff:              0,
				InsulinAction:        3 * time.Hour,
				InsulinConcentration: 100,
				MaxBolus:             Insulin(12500),
				MaxBasal:             Insulin(8100),
				RFEnabled:            true,
				TempBasalType:        Absolute,
				SelectedPattern:      0,
			},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			b, err := decodeSettings(c.data, c.family)
			if err != nil {
				t.Errorf("decodeSettings(% X, %d) returned %+v, want %+v", c.data, c.family, err, c.b)
				return
			}
			if !reflect.DeepEqual(b, c.b) {
				t.Errorf("decodeSettings(% X, %d) == %+v, want %+v", c.data, c.family, b, c.b)
			}
		})
	}
}
