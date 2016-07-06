package medtronic

import (
	"time"
)

const (
	TimeLayout = "2006-01-02 15:04:05" // ISO 8601-ish
)

// Convert a multiple of half-hours to a Duration.
func scheduleToDuration(t uint8) time.Duration {
	return time.Duration(t) * 30 * time.Minute
}

// Convert a time to a Duration representing the offset since 00:00:00.
func sinceMidnight(t time.Time) time.Duration {
	year, month, day := t.Date()
	midnight := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return t.Sub(midnight)
}

// Decode a 5-byte timestamp from a pump history record.
func decodeTimestamp(data []byte) time.Time {
	sec := int(data[0] & 0x3F)
	min := int(data[1] & 0x3F)
	hour := int(data[2] & 0x1F)
	day := int(data[3] & 0x1F)
	// The 4-bit month value is encoded in the high 2 bits of the first 2 bytes.
	month := time.Month(int(data[0]>>6)<<2 | int(data[1]>>6))
	year := 2000 + int(data[4]&0x7F)
	return time.Date(year, month, day, hour, min, sec, 0, time.UTC)
}

// Decode a 2-byte date from a pump history record.
func decodeDate(data []byte) time.Time {
	day := int(data[0] & 0x1F)
	month := time.Month(int(data[0]>>5)<<1 + int(data[1]>>7))
	year := 2000 + int(data[1]&0x7F)
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// TimeNow returns a UTC time with the same Date amd Clock value as time.Now.
func TimeNow() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	return time.Date(year, month, day, hour, min, sec, 0, time.UTC)
}
