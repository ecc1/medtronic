package medtronic

import (
	"reflect"
	"testing"
	"time"
)

func cgmTS(t time.Time) CGMRecord {
	return CGMRecord{
		Type: CGMTimestamp,
		Time: t,
	}
}

func cgmBG() CGMRecord {
	return CGMRecord{
		Type:    CGMGlucose,
		Glucose: 100,
	}
}

func cgmBGT(t time.Time) CGMRecord {
	return CGMRecord{
		Type:    CGMGlucose,
		Time:    t,
		Glucose: 100,
	}
}

func TestAddTimestamps(t *testing.T) {
	var (
		r0 = cgmTS(parseTime("2018-05-01T00:00"))
		r1 = cgmBG()
		r2 = cgmTS(parseTime("2018-05-01T01:00"))
		t1 = cgmBGT(parseTime("2018-05-01T00:05"))
		t2 = cgmBGT(parseTime("2018-05-01T00:10"))
		t3 = cgmBGT(parseTime("2018-05-01T00:15"))
		t4 = cgmBGT(parseTime("2018-05-01T01:05"))
		t5 = cgmBGT(parseTime("2018-05-01T01:10"))
	)
	cases := []struct {
		before CGMHistory
		after  CGMHistory
	}{
		{nil, nil},
		{CGMHistory{r0}, CGMHistory{r0}},
		{CGMHistory{r1}, CGMHistory{r1}},
		{CGMHistory{r1, r1, r1}, CGMHistory{r1, r1, r1}},
		{CGMHistory{r1, r0}, CGMHistory{t1, r0}},
		{CGMHistory{r1, r1, r1, r0}, CGMHistory{t3, t2, t1, r0}},
		{CGMHistory{r1, r1, r2, r1, r0}, CGMHistory{t5, t4, r2, t1, r0}},
	}
	for _, c := range cases {
		before := c.before[:]
		addTimestamps(before)
		if !reflect.DeepEqual(before, c.after) {
			t.Errorf("addTimestamps(%+v) == %+v, want %+v", c.before, before, c.after)
		}
	}
}
