package packet

import (
	"testing"
)

func TestCRC16(t *testing.T) {
	cases := []struct {
		msg []byte
		sum uint16
	}{
		{[]byte{0x02, 0x06, 0x06, 0x03}, 0x41CD},
		{[]byte{0x02, 0x09, 0x00, 0x05, 0x0D, 0x02, 0x03}, 0x71DA},
	}
	for _, c := range cases {
		sum := CRC16(c.msg)
		if sum != c.sum {
			t.Errorf("CRC16(% X) == %X, want %X", c.msg, sum, c.sum)
		}
	}
}
