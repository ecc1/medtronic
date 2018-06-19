package medtronic

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

const (
	// JSONTimeLayout specifies the format for JSON time values.
	JSONTimeLayout = time.RFC3339
	// UserTimeLayout specifies a consistent, human-readable format for local time.
	UserTimeLayout = "2006-01-02 15:04:05"
)

type (
	// Duration allows custom JSON marshaling for time.Duration values.
	Duration time.Duration
	// TimeOfDay represents a value between 0 and 24 hours.
	TimeOfDay time.Duration
)

// hourstoDuration converts n hours to a Duration.
func hoursToDuration(n uint8) Duration {
	return Duration(time.Duration(n) * time.Hour)
}

// halfHoursToDuration converts n half-hours to a Duration.
func halfHoursToDuration(n uint8) Duration {
	return Duration(time.Duration(n) * 30 * time.Minute)
}

// minutesToDuration converts n minutes to a Duration.
func minutesToDuration(n uint8) Duration {
	return Duration(time.Duration(n) * time.Minute)
}

// TimeOfDay converts a duration to a time of day.
func (d Duration) TimeOfDay() TimeOfDay {
	td := time.Duration(d)
	if td < 0 || 24*time.Hour <= td {
		log.Panicf("duration %v is not a valid time of day", td)
	}
	return TimeOfDay(td)
}

// Convert a time of day into half-hours.
func (t TimeOfDay) HalfHours() uint8 {
	n := time.Duration(t) / (30 * time.Minute)
	if n > 255 {
		log.Panicf("time of day %v is too large", t)
	}
	return uint8(n)
}

// Convert a time of day to a string of the form HH:MM.
func (t TimeOfDay) String() string {
	d := time.Duration(t)
	hour := d / time.Hour
	min := (d % time.Hour) / time.Minute
	return fmt.Sprintf("%02d:%02d", hour, min)
}

// ParseTimeOfDay parses a string of the form HH:MM into a time of day.
func ParseTimeOfDay(s string) (TimeOfDay, error) {
	if len(s) == 5 && s[2] == ':' {
		hour, hErr := strconv.Atoi(s[0:2])
		min, mErr := strconv.Atoi(s[3:5])
		if hErr == nil && 0 <= hour && hour <= 23 && mErr == nil && 0 <= min && min <= 59 {
			d := time.Duration(hour)*time.Hour + time.Duration(min)*time.Minute
			return Duration(d).TimeOfDay(), nil
		}
	}
	return 0, fmt.Errorf("parseTimeOfDay: %q must be of the form HH:MM", s)
}

// halfHoursToTimeOfDay converts n half-hours to a time of day.
func halfHoursToTimeOfDay(n uint8) TimeOfDay {
	return Duration(time.Duration(n) * 30 * time.Minute).TimeOfDay()
}

// SinceMidnight converts a time to a time of day.
func SinceMidnight(t time.Time) TimeOfDay {
	hour, min, sec := t.Clock()
	h, m, s := time.Duration(hour), time.Duration(min), time.Duration(sec)
	n := time.Duration(t.Nanosecond())
	d := h*time.Hour + m*time.Minute + s*time.Second + n*time.Nanosecond
	return Duration(d).TimeOfDay()
}

// Decode a 5-byte timestamp from a pump history record.
func decodeTime(data []byte) time.Time {
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
