package medtronic

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	JsonTimeLayout = time.RFC3339          // Nightscout-interoperable
	UserTimeLayout = "2006-01-02 15:04:05" // human-readable
)

// Convert a multiple of half-hours to a Duration.
func halfHoursToDuration(t uint8) time.Duration {
	return time.Duration(t) * 30 * time.Minute
}

// TimeOfDay represents a value between 0 and 24 hours.
type TimeOfDay time.Duration

// Convert a duration to a time of day.
func durationToTimeOfDay(d time.Duration) TimeOfDay {
	if d < 0 || 24*time.Hour <= d {
		log.Panicf("duration %v is not a valid time of day", d)
	}
	return TimeOfDay(d)
}

// Convert a time of day to a string of the form HH:MM.
func (t TimeOfDay) String() string {
	d := time.Duration(t)
	hour := d / time.Hour
	min := (d % time.Hour) / time.Minute
	return fmt.Sprintf("%02d:%02d", hour, min)
}

// ParseTimeOfDay parses a string of the form HH:MM
// and returns a time of day.
func parseTimeOfDay(s string) (TimeOfDay, error) {
	if len(s) == 5 && s[2] == ':' {
		hour, hErr := strconv.Atoi(s[0:2])
		min, mErr := strconv.Atoi(s[3:5])
		if hErr == nil && 0 <= hour && hour <= 23 && mErr == nil && 0 <= min && min <= 59 {
			d := time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute
			return durationToTimeOfDay(d), nil
		}
	}
	return 0, fmt.Errorf("parseTimeOfDay: %q must be of the form HH:MM", s)
}

// Convert a multiple of half-hours to a time of day.
func halfHoursToTimeOfDay(t uint8) TimeOfDay {
	return durationToTimeOfDay(halfHoursToDuration(t))
}

// Convert a time to a time of day.
func sinceMidnight(t time.Time) TimeOfDay {
	hour, min, sec := t.Clock()
	h, m, s := time.Duration(hour), time.Duration(min), time.Duration(sec)
	n := time.Duration(t.Nanosecond())
	d := h*time.Hour + m*time.Minute + s*time.Second + n*time.Nanosecond
	return durationToTimeOfDay(d)
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
	return time.Date(year, month, day, hour, min, sec, 0, time.Local)
}

// Decode a 2-byte date from a pump history record.
func decodeDate(data []byte) time.Time {
	day := int(data[0] & 0x1F)
	month := time.Month(int(data[0]>>5)<<1 + int(data[1]>>7))
	year := 2000 + int(data[1]&0x7F)
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}
