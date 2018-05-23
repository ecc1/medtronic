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
			[]byte{0x15, 0x00, 0x01, 0x01, 0x01, 0x00, 0x96, 0x00, 0x8C, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x64, 0x01, 0x04, 0x00, 0x14, 0x00},
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
			[]byte{0x19, 0x00, 0x04, 0x01, 0x05, 0x01, 0x00, 0x8c, 0x00, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x64, 0x01, 0x04, 0x00, 0x14, 0x00, 0x64, 0x01, 0x01, 0x00},
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
			[]byte{0x12, 0x0F, 0x02, 0x01, 0x01, 0x00, 0x64, 0x00, 0x78, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x64, 0x01, 0x00},
			12,
			SettingsInfo{
				AutoOff:              15 * time.Hour,
				InsulinConcentration: 100,
				MaxBolus:             Insulin(10000),
				MaxBasal:             Insulin(3000),
				RFEnabled:            true,
				TempBasalType:        Absolute,
				SelectedPattern:      0,
			},
		},
	}
	for _, c := range cases {
		b, err := decodeSettings(c.data, c.family)
		if err != nil {
			t.Errorf("decodeSettings(% X, %d) returned %+v, want %+v", c.data, c.family, err, c.b)
			continue
		}
		if !reflect.DeepEqual(b, c.b) {
			t.Errorf("decodeSettings(% X, %d) == %+v, want %+v", c.data, c.family, b, c.b)
		}
	}
}
