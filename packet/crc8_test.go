package packet

import (
	"testing"
)

func TestCrc8(t *testing.T) {
	cases := []struct {
		msg []byte
		sum byte
	}{
		{[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, 0x98},
		{[]byte("0123456789"), 0xFD},
		{[]byte{0xFF, 0xFF, 0xFF, 0xFF}, 0xA3},
		{[]byte{0x01, 0x07, 0x00, 0x10, 0x04}, 0x44},
	}
	for _, c := range cases {
		sum := Crc8(c.msg)
		if sum != c.sum {
			t.Errorf("crc8(%X) == %X, want %X", c.msg, sum, c.sum)
		}
	}
}
