package medtronic

import (
	"testing"
)

func TestBattery(t *testing.T) {
	cases := []struct {
		data []byte
		b    BatteryInfo
	}{
		{
			parseBytes("03 00 00 96"),
			BatteryInfo{Voltage: 1500, LowBattery: false},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			b := decodeBatteryInfo(c.data)
			if b != c.b {
				t.Errorf("decodeBatteryInfo(% X) == %+v, want %+v", c.data, b, c.b)
			}
		})
	}
}
