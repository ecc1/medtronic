package packet

import (
	"testing"
)

func TestCRC8(t *testing.T) {
	cases := []struct {
		msg []byte
		sum byte
	}{
		{parseBytes("00 01 02 03 04 05 06 07 08 09"), 0x98},
		{parseBytes("A7 12 89 86 5D 00"), 0xBE},
		{parseBytes("A7 12 89 86 06 00"), 0x15},
		{parseBytes("A7 12 89 86 15 09"), 0x56},
		{parseBytes("A7 12 89 86 8D 00"), 0xB0},
		{parseBytes("A7 12 89 86 8D 09 03 37 32 32 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00"), 0x39},
	}
	for _, c := range cases {
		sum := CRC8(c.msg)
		if sum != c.sum {
			t.Errorf("CRC8(% X) == %02X, want %02X", c.msg, sum, c.sum)
		}
	}
}
