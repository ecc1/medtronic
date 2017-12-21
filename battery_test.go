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
			[]byte{0x03, 0x00, 0x00, 0x96},
			BatteryInfo{Voltage: 1500, LowBattery: false},
		},
	}
	for _, c := range cases {
		b := decodeBatteryInfo(c.data)
		if b != c.b {
			t.Errorf("decodeBatteryInfo(% X) == %+v, want %+v", c.data, b, c.b)
		}
	}
}
