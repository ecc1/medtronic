package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func TestDecodeCGMTime(t *testing.T) {
	cases := []struct {
		b []byte
		t time.Time
	}{
		{parseBytes("8D 9B 1D 0C"), parseTime("2012-10-29T13:27")},
		{parseBytes("0B AE 0A 0E"), parseTime("2014-02-10T11:46")},
		{parseBytes("4F 5B 13 8F"), parseTime("2015-05-19T15:27")},
		{parseBytes("14 B6 28 10"), parseTime("2016-02-08T20:54")},
	}
	for _, c := range cases {
		t.Run(c.t.Format(time.Kitchen), func(t *testing.T) {
			ts := time.Time(decodeCGMTime(c.b))
			if !ts.Equal(c.t) {
				t.Errorf("decodeCGMTime(% X) == %v, want %v", c.b, ts, c.t)
			}
		})
	}
}

func TestInitialTimestamp(t *testing.T) {
	cases := []struct {
		data []byte
		ts   time.Time
		err  error
	}{
		{parseBytes(""), time.Time{}, ErrorNeedsTimestamp},
		{parseBytes("01011053b3940810531111AE"), time.Time{}, ErrorNeedsTimestamp},
		{parseBytes("1013b39408534232"), parseTime("2016-10-19T21:06"), nil},
		{parseBytes("1053b3940853535A"), time.Time{}, ErrorNeedsTimestamp},
	}
	for _, c := range cases {
		reverseBytes(c.data)
		t.Run("", func(t *testing.T) {
			ts, _, err := initialTimestamp(c.data)
			if ts != c.ts {
				t.Errorf("initialTimestamp returned %v, want %v", ts, c.ts)
			}
			if err != c.err {
				t.Errorf("initialTimestamp raised %v, want %v", err, c.err)
			}
		})
	}
}

func TestDecodeCGMRecord(t *testing.T) {
	cases := []struct {
		data []byte
		r    CGMRecord
	}{
		// github.com/ps2/rileylink_ios/MinimedKitTests/GlucoseEvents/
		{parseBytes("0300"), CGMRecord{
			Type:  CGMCal,
			Value: "bgNow",
		}},
		{parseBytes("0301"), CGMRecord{
			Type:  CGMCal,
			Value: "waiting",
		}},
		{parseBytes("0302"), CGMRecord{
			Type:  CGMCal,
			Value: "error",
		}},
		{parseBytes("0501"), CGMRecord{
			Type: CGMError,
		}},
		{parseBytes("0a0bae0a0e"), CGMRecord{
			Type: CGMBatteryChange,
			Time: parseTime("2014-02-10T11:46"),
		}},
		{parseBytes("0c0ad23e0e"), CGMRecord{
			Type: CGMTimeChange,
			Time: parseTime("2014-03-30T10:18"),
		}},
		{parseBytes("0e4f5b138fa0"), CGMRecord{
			Type:    CGMCalBG,
			Time:    parseTime("2015-05-19T15:27"),
			Glucose: 160,
		}},
		{parseBytes("0f4f67130f128c"), CGMRecord{
			Type: CGMCalFactor,
			Time: parseTime("2015-05-19T15:39"),
		}},
		{parseBytes("0402"), CGMRecord{
			Type: CGMPacket,
		}},
		{parseBytes("06"), CGMRecord{
			Type:    CGMDataLow,
			Glucose: 40,
		}},
		{parseBytes("07ff"), CGMRecord{
			Type:    CGMDataHigh,
			Glucose: 400,
		}},
		{parseBytes("0814b62810"), CGMRecord{
			Type:  CGMTimestamp,
			Time:  parseTime("2016-02-08T20:54"),
			Value: "pageEnd",
		}},
		{parseBytes("088d9b5d0c"), CGMRecord{
			Type:  CGMTimestamp,
			Time:  parseTime("2012-10-29T13:27"),
			Value: "gap",
		}},
		{parseBytes("0b0baf0a0e"), CGMRecord{
			Type:  CGMStatus,
			Time:  parseTime("2014-02-10T11:47"),
			Value: "off",
		}},
		{parseBytes("0b0baf2a0e"), CGMRecord{
			Type:  CGMStatus,
			Time:  parseTime("2014-02-10T11:47"),
			Value: "on",
		}},
		{parseBytes("0b0baf4a0e"), CGMRecord{
			Type:  CGMStatus,
			Time:  parseTime("2014-02-10T11:47"),
			Value: "lost",
		}},
		{parseBytes("0d4d44330f"), CGMRecord{
			Type:  CGMSync,
			Time:  parseTime("2015-05-19T13:04"),
			Value: "new",
		}},
		{parseBytes("0d4d44530f"), CGMRecord{
			Type:  CGMSync,
			Time:  parseTime("2015-05-19T13:04"),
			Value: "old",
		}},
		{parseBytes("0d4d44730f"), CGMRecord{
			Type:  CGMSync,
			Time:  parseTime("2015-05-19T13:04"),
			Value: "find",
		}},
		{parseBytes("35"), CGMRecord{
			Type:    CGMGlucose,
			Glucose: 106,
		}},
	}
	for _, c := range cases {
		c.r.Data = c.data
		t.Run("", func(t *testing.T) {
			r, err := DecodeOneRecord(t, c.data)
			if err != nil && err != ErrorNeedsTimestamp {
				t.Errorf("DecodeCGMRecord raised %v", err)
			}
			if !reflect.DeepEqual(r, c.r) {
				t.Errorf("DecodeCGMRecord returned %+v, want %+v", r, c.r)
			}
		})
	}
}

func DecodeOneRecord(t *testing.T, data []byte) (CGMRecord, error) {
	reverseBytes(data)
	h, _, err := DecodeCGMHistory(data, time.Time{})
	var r CGMRecord
	if len(h) != 1 {
		t.Errorf("DecodeCGMHistory returned %d records", len(h))
	}
	if len(h) != 0 {
		r = h[0]
	}
	return r, err
}
