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
			[]byte{3, 0, 0, 150},
			BatteryInfo{Voltage: 1500, LowBattery: false},
		},
	}
	for _, c := range cases {
		b := decodeBatteryInfo(c.data)
		if b != c.b {
			t.Errorf("decodeBatteryInfo(% X) == %v, want %v", c.data, b, c.b)
		}
	}
}
