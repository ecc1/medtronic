package medtronic

import (
	"fmt"
	"testing"
)

func TestReservoir(t *testing.T) {
	cases := []struct {
		data   []byte
		family Family
		i      Insulin
	}{
		{
			parseBytes("02 05 46"),
			22,
			135000,
		},
		{
			parseBytes("04 00 00 0C A3"),
			23,
			80875,
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%d", c.i), func(t *testing.T) {
			i, err := decodeReservoir(c.data, c.family)
			if err != nil {
				t.Errorf("decodeReservoir(% X, %d) returned %v, want %v", c.data, c.family, err, c.i)
				return
			}
			if i != c.i {
				t.Errorf("decodeReservoir(% X, %d) == %v, want %v", c.data, c.family, i, c.i)
			}
		})
	}
}
