package medtronic

import (
	"testing"
)

func TestReservoir(t *testing.T) {
	cases := []struct {
		data   []byte
		family Family
		i      Insulin
	}{
		{
			[]byte{0x02, 0x05, 0x46},
			22,
			135000,
		},
		{
			[]byte{0x04, 0x00, 0x00, 0x0C, 0xA3},
			23,
			80875,
		},
	}
	for _, c := range cases {
		i, err := decodeReservoir(c.data, c.family)
		if err != nil {
			t.Errorf("decodeReservoir(% X, %d) returned %v, want %v", c.data, c.family, err, c.i)
			continue
		}
		if i != c.i {
			t.Errorf("decodeReservoir(% X, %d) == %v, want %v", c.data, c.family, i, c.i)
		}
	}
}
