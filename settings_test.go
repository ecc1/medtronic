package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func TestSettings(t *testing.T) {
	cases := []struct {
		data  []byte
		newer bool
		b     SettingsInfo
	}{
		{
			[]byte("\x15\x00\x01\x01\x01\x00\x96\x00\x8C\x00\x00\x00\x00\x01\x00\x00\x64\x01\x04\x00\x14\x00"),
			false,
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
			[]byte("\x19\x00\x04\x01\x05\x01\x00\x8c\x00\x50\x00\x00\x00\x00\x00\x01\x64\x01\x04\x00\x14\x00\x64\x01\x01\x00"),
			true,
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
	}
	for _, c := range cases {
		b, err := decodeSettings(c.data, c.newer)
		if err != nil {
			t.Errorf("decodeSettings(% X, %v) returned %+v, want %+v", c.data, c.newer, err, c.b)
			continue
		}
		if !reflect.DeepEqual(b, c.b) {
			t.Errorf("decodeSettings(% X, %v) == %+v, want %+v", c.data, c.newer, b, c.b)
		}
	}
}
