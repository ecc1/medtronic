package medtronic

import (
	"fmt"
	"time"
)

// CGMRecordType represents a CGM record type.
type CGMRecordType byte

//go:generate stringer -type CGMRecordType

// Events stored in the pump's CGM history pages.
const (
	CGMDataEnd       CGMRecordType = 0x01
	CGMWeakSignal    CGMRecordType = 0x02
	CGMCal           CGMRecordType = 0x03
	CGMPacket        CGMRecordType = 0x04
	CGMError         CGMRecordType = 0x05
	CGMDataLow       CGMRecordType = 0x06
	CGMDataHigh      CGMRecordType = 0x07
	CGMTimestamp     CGMRecordType = 0x08
	CGMBatteryChange CGMRecordType = 0x0A
	CGMSensorStatus  CGMRecordType = 0x0B
	CGMTimeChange    CGMRecordType = 0x0C
	CGMSync          CGMRecordType = 0x0D
	CGMCalBGForGH    CGMRecordType = 0x0E
	CGMCalFactor     CGMRecordType = 0x0F
	CGMEvent10       CGMRecordType = 0x10

	// Synthetic record type.
	// Single bytes with this value or greater represent glucose readings.
	CGMGlucose CGMRecordType = 0x20
)

// (Only records with length > 1 need to be included here.)
var cgmRecordLength = map[CGMRecordType]int{
	CGMCal:           2,
	CGMPacket:        2,
	CGMError:         2,
	CGMDataHigh:      2,
	CGMTimestamp:     5,
	CGMBatteryChange: 5,
	CGMSensorStatus:  5,
	CGMTimeChange:    5,
	CGMSync:          5,
	CGMCalBGForGH:    6,
	CGMCalFactor:     7,
	CGMEvent10:       8,
}

type (
	// CGMRecord represents a CGM record.
	CGMRecord struct {
		Type    CGMRecordType
		Data    []byte
		Time    Time
		Glucose int `json:",omitempty"`
	}

	// CGMHistory represents a sequence of CGM records.
	CGMHistory []CGMRecord
)

func decodeCGMTimestamp(r *CGMRecord) error {
	r.Time = decodeCGMTime(r.Data[1:])
	return nil
}

// Decode a 4-byte timestamp from a glucose history record.
func decodeCGMTime(data []byte) Time {
	sec := 0
	min := int(data[1] & 0x3F)
	hour := int(data[0] & 0x1F)
	day := int(data[2] & 0x1F)
	// The 4-bit month value is encoded in the high 2 bits of the first 2 bytes.
	month := time.Month(int(data[0]>>6)<<2 | int(data[1]>>6))
	year := 2000 + int(data[3]&0x7F)
	return Time(time.Date(year, month, day, hour, min, sec, 0, time.Local))
}

func decodeCGMGlucose(r *CGMRecord) error {
	r.Glucose = 2 * int(r.Data[0])
	return nil
}

func decodeCGMCalBGForGH(r *CGMRecord) error {
	r.Glucose = int(r.Data[5])
	return nil
}

// DecodeCGMRecord decodes a CGM history record based on its type.
func DecodeCGMRecord(data []byte) (CGMRecord, error) {
	if len(data) == 0 {
		return CGMRecord{}, fmt.Errorf("empty CGM record")
	}
	t := CGMRecordType(data[0])
	if t >= CGMGlucose {
		t = CGMGlucose
	}
	n, found := cgmRecordLength[t]
	if !found {
		n = 1
	}
	var err error
	r := CGMRecord{Type: t, Data: data[:n]}
	switch t {
	case CGMGlucose:
		err = decodeCGMGlucose(&r)
	case CGMCalBGForGH:
		err = decodeCGMCalBGForGH(&r)
	default:
		if n >= 5 {
			err = decodeCGMTimestamp(&r)
		}
	}
	return r, err
}

// DecodeCGMHistory decodes the records in a page of CGM data and
// returns them in reverse chronological order (most recent first),
// to match the order of the CGM pages themselves.
func DecodeCGMHistory(data []byte) (CGMHistory, error) {
	reverseBytes(data)
	var i int
	var b byte
	for i, b = range data {
		if b != 0 {
			break
		}
	}
	data = data[i:]
	var results CGMHistory
	var r CGMRecord
	var err error
	for len(data) > 0 {
		r, err = DecodeCGMRecord(data)
		if err != nil {
			break
		}
		results = append(results, r)
		data = data[len(r.Data):]
	}
	addTimestamps(results)
	return results, err
}

func reverseBytes(a []byte) {
	for i, j := 0, len(a)-1; i < len(a)/2; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func addTimestamps(results CGMHistory) {
	// Start at the end to find the earliest timestamp.
	var prev int
	var ts time.Time
	i := len(results) - 1
	for i >= 0 {
		r := results[i]
		t := time.Time(r.Time)
		if !t.IsZero() {
			prev = i
			ts = t
		} else if prev != 0 {
			t = ts.Add(time.Duration(prev-i) * 5 * time.Minute)
			results[i].Time = Time(t)
		}
		i--
	}
}
