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
	CGMStatus        CGMRecordType = 0x0B
	CGMTimeChange    CGMRecordType = 0x0C
	CGMSync          CGMRecordType = 0x0D
	CGMCalBG         CGMRecordType = 0x0E
	CGMCalFactor     CGMRecordType = 0x0F
	CGMEvent10       CGMRecordType = 0x10
	CGMEvent13       CGMRecordType = 0x13

	// Synthetic record type.
	// Single bytes with this value or greater represent glucose readings.
	CGMGlucose CGMRecordType = 0x20
)

type (
	// CGMRecord represents a CGM record.
	CGMRecord struct {
		Type    CGMRecordType
		Data    []byte
		Time    time.Time
		Glucose int    `json:",omitempty"`
		Value   string `json:",omitempty"`
	}

	// CGMHistory represents a sequence of CGM records.
	CGMHistory []CGMRecord

	cgmDecoder func(*CGMRecord)

	decodeInfo struct {
		length  int
		decoder cgmDecoder
	}
)

// 1-byte records containing just a type code need not be specified here.
// Timestamp decoding is done whenever the record is at least 5 bytes long,
// so there is no need to include that in the decoder.
var cgmDecodeInfo = map[CGMRecordType]decodeInfo{
	CGMCal:           {2, decodeCGMCal},
	CGMPacket:        {2, nil},
	CGMError:         {2, nil},
	CGMDataHigh:      {2, nil},
	CGMTimestamp:     {5, decodeCGMTimestamp},
	CGMBatteryChange: {5, nil},
	CGMStatus:        {5, decodeCGMStatus},
	CGMTimeChange:    {5, nil},
	CGMSync:          {5, decodeCGMSync},
	CGMCalBG:         {6, decodeCGMCalBG},
	CGMCalFactor:     {7, nil},
	CGMEvent10:       {8, nil},
	CGMGlucose:       {1, decodeCGMGlucose},
}

// Decode a 4-byte timestamp from a glucose history record.
func decodeCGMTime(data []byte) time.Time {
	sec := 0
	min := int(data[1] & 0x3F)
	hour := int(data[0] & 0x1F)
	day := int(data[2] & 0x1F)
	// The 4-bit month value is encoded in the high 2 bits of the first 2 bytes.
	month := time.Month(int(data[0]>>6)<<2 | int(data[1]>>6))
	year := 2000 + int(data[3]&0x7F)
	return time.Date(year, month, day, hour, min, sec, 0, time.Local)
}

func (t CGMRecordType) isRelative() bool {
	switch t {
	case CGMWeakSignal, CGMCal, CGMPacket, CGMError, CGMDataLow, CGMDataHigh, CGMGlucose:
		return true
	default:
		return false
	}
}

func decodeCGMCal(r *CGMRecord) {
	switch r.Data[1] {
	case 0:
		r.Value = "bgNow"
	case 1:
		r.Value = "waiting"
	case 2:
		r.Value = "error"
	default:
		r.Value = "unknown"
	}
}

func decodeCGMTimestamp(r *CGMRecord) {
	switch (r.Data[3] >> 5) & 0x3 {
	case 0:
		r.Value = "lastRF"
	case 1:
		r.Value = "pageEnd"
	case 2:
		r.Value = "gap"
	default:
		r.Value = "unknown"
	}
}

func decodeCGMStatus(r *CGMRecord) {
	switch (r.Data[3] >> 5) & 0x3 {
	case 0:
		r.Value = "off"
	case 1:
		r.Value = "on"
	case 2:
		r.Value = "lost"
	default:
		r.Value = "unknown"
	}
}

func decodeCGMSync(r *CGMRecord) {
	switch (r.Data[3] >> 5) & 0x3 {
	case 1:
		r.Value = "new"
	case 2:
		r.Value = "old"
	default:
		r.Value = "find"
	}
}

func decodeCGMCalBG(r *CGMRecord) {
	r.Glucose = int(r.Data[5])
}

func decodeCGMGlucose(r *CGMRecord) {
	r.Glucose = 2 * int(r.Data[0])
}

// DecodeCGMRecord decodes a CGM history record based on its type.
func DecodeCGMRecord(data []byte) (CGMRecord, error) {
	if len(data) == 0 {
		return CGMRecord{}, fmt.Errorf("DecodeCGMRecord: len(data) == 0")
	}
	t := CGMRecordType(data[0])
	if t >= CGMGlucose {
		t = CGMGlucose
	}
	n := 1
	var decode cgmDecoder
	d, found := cgmDecodeInfo[t]
	if found {
		n = d.length
		decode = d.decoder
	}
	if n > len(data) {
		return CGMRecord{}, fmt.Errorf("DecodeCGMRecord: expected %d-byte record but len(data) = %d", n, len(data))
	}
	r := CGMRecord{Type: t, Data: data[:n]}
	if n >= 5 {
		r.Time = decodeCGMTime(r.Data[1:5])
	}
	if decode != nil {
		decode(&r)
	}
	return r, nil
}

// DecodeCGMHistory decodes the records in a page of CGM data and
// returns them in reverse chronological order (most recent first).
// If a non-zero time is given, it is used as the initial timestamp.
func DecodeCGMHistory(data []byte, t time.Time) (CGMHistory, time.Time, error) {
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
	var err error
	if t.IsZero() {
		t, results, err = initialTimestamp(data)
		if err != nil {
			return results, t, err
		}
	}
	results = nil
	for len(data) != 0 {
		var r CGMRecord
		r, err = DecodeCGMRecord(data)
		if err != nil {
			break
		}
		if r.Type == CGMTimestamp {
			t = r.Time
		} else if r.Type.isRelative() {
			r.Time = t
			t = t.Add(-5 * time.Minute)
		}
		results = append(results, r)
		data = data[len(r.Data):]
	}
	return results, t, err
}

// ErrorNeedsTimestamp indicates that no initial timestamp was found.
var ErrorNeedsTimestamp = fmt.Errorf("CGM history needs timestamp")

func initialTimestamp(data []byte) (time.Time, CGMHistory, error) {
	var results CGMHistory
	numRelative := 0
	var err error
	for len(data) != 0 {
		var r CGMRecord
		r, err = DecodeCGMRecord(data)
		if err != nil {
			return time.Time{}, results, err
		}
		results = append(results, r)
		data = data[len(r.Data):]
		if r.Type.isRelative() {
			numRelative++
			continue
		}
		if r.Type == CGMTimestamp {
			if numRelative == 0 || r.hasOffsetTimestamp() {
				delta := time.Duration(numRelative) * 5 * time.Minute
				return r.Time.Add(delta), results, nil
			}
		}
		if r.Type != CGMDataEnd && r.Type != CGMEvent13 {
			break
		}
	}
	return time.Time{}, results, ErrorNeedsTimestamp
}

func (r CGMRecord) hasOffsetTimestamp() bool {
	return r.Value == "lastRF" || r.Value == "pageEnd"
}

func reverseBytes(a []byte) {
	for i, j := 0, len(a)-1; i < len(a)/2; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

// ReverseCGMHistory reverses a slice of CGM history records.
func ReverseCGMHistory(a CGMHistory) {
	for i, j := 0, len(a)-1; i < len(a)/2; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}
