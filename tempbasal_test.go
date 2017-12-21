package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func insulinPointer(n int) *Insulin {
	v := Insulin(n)
	return &v
}

func TestTempBasal(t *testing.T) {
	cases := []struct {
		data []byte
		b    TempBasalInfo
	}{
		{
			[]byte{0x06, 0x00, 0x00, 0x00, 0x8C, 0x00, 0x1E},
			TempBasalInfo{
				Duration: 30 * time.Minute,
				Type:     Absolute,
				Rate:     insulinPointer(3500),
			},
		},
		{
			[]byte{0x06, 0x00, 0x00, 0x00, 0x37, 0x00, 0x17},
			TempBasalInfo{
				Duration: 23 * time.Minute,
				Type:     Absolute,
				Rate:     insulinPointer(1375),
			},
		},
		{
			[]byte{0x06, 0x00, 0x00, 0x05, 0x50, 0x00, 0x1E},
			TempBasalInfo{
				Duration: 30 * time.Minute,
				Type:     Absolute,
				Rate:     insulinPointer(34000),
			},
		},
	}
	for _, c := range cases {
		b, err := decodeTempBasal(c.data)
		if err != nil {
			t.Errorf("decodeTempBasal(% X) returned %v, want %v", c.data, err, c.b)
			continue
		}
		if !reflect.DeepEqual(b, c.b) {
			t.Errorf("decodeTempBasal(% X) == %v, want %v", c.data, b, c.b)
		}
	}
}
