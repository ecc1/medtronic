package medtronic

import (
	"testing"
)

func TestReservoir(t *testing.T) {
	cases := []struct {
		data  []byte
		newer bool
		i     Insulin
	}{
		{
			[]byte{0x02, 0x05, 0x46},
			false,
			135000,
		},
		{
			[]byte{0x04, 0x00, 0x00, 0x0C, 0xA3},
			true,
			80875,
		},
	}
	for _, c := range cases {
		i, err := decodeReservoir(c.data, c.newer)
		if err != nil {
			t.Errorf("decodeReservoir(%X, %v) returned %v, want %v", c.data, c.newer, err, c.i)
			continue
		}
		if i != c.i {
			t.Errorf("decodeReservoir(%X, %v) == %v, want %v", c.data, c.newer, i, c.i)
		}
	}
}
