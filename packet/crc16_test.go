package packet

import (
	"testing"
)

func TestCRC16(t *testing.T) {
	cases := []struct {
		msg []byte
		sum uint16
	}{
		{parseBytes("02 06 06 03"), 0x41CD},
		{parseBytes("02 09 00 05 0D 02 03"), 0x71DA},
		{parseBytes("A8 0F 25 C1 23 0D 19 1C 50 00 8F 00 90 00 34 34 99 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00"), 0xDEE7},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			sum := CRC16(c.msg)
			if sum != c.sum {
				t.Errorf("CRC16(% X) == %X, want %X", c.msg, sum, c.sum)
			}
		})
	}
}
