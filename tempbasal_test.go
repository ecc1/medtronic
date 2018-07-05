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
			parseBytes("06 00 00 00 8C 00 1E"),
			TempBasalInfo{
				Duration: 30 * time.Minute,
				Type:     Absolute,
				Rate:     insulinPointer(3500),
			},
		},
		{
			parseBytes("06 00 00 00 37 00 17"),
			TempBasalInfo{
				Duration: 23 * time.Minute,
				Type:     Absolute,
				Rate:     insulinPointer(1375),
			},
		},
		{
			parseBytes("06 00 00 05 50 00 1E"),
			TempBasalInfo{
				Duration: 30 * time.Minute,
				Type:     Absolute,
				Rate:     insulinPointer(34000),
			},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			b, err := decodeTempBasal(c.data)
			if err != nil {
				t.Errorf("decodeTempBasal(% X) returned %+v, want %+v", c.data, err, c.b)
				return
			}
			if !reflect.DeepEqual(b, c.b) {
				t.Errorf("decodeTempBasal(% X) == %+v, want %+v", c.data, b, c.b)
			}
		})
	}
}
