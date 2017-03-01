package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	switch len(os.Args) {
	case 3:
		t := parseTime(os.Args[1] + " " + os.Args[2])
		fmt.Printf(" time: %v\n", t)
		fmt.Printf("bytes: % X\n", encodeTime(t))
	default:
		var data []byte
		for _, s := range os.Args[1:] {
			b, err := strconv.ParseUint(s, 16, 8)
			if err != nil {
				log.Fatal(err)
			}
			data = append(data, byte(b))
		}
		t := parseTimestamp(data)
		fmt.Printf("           time: %v\n", t)
		fmt.Printf(" original bytes: % X\n", data)
		fmt.Printf("canonical bytes: % X\n", encodeTime(t))
	}
}

// Parse a 5-byte timestamp from a pump history record.
func parseTimestamp(data []byte) time.Time {
	sec := int(data[0] & 0x3F)
	min := int(data[1] & 0x3F)
	hour := int(data[2] & 0x1F)
	day := int(data[3] & 0x1F)
	// The 4-bit month value is encoded in the high 2 bits of the first 2 bytes.
	month := time.Month(int(data[0]>>6)<<2 | int(data[1]>>6))
	year := 2000 + int(data[4]&0x7F)
	return time.Date(year, month, day, hour, min, sec, 0, time.Local)
}

func parseTime(date string) time.Time {
	const layout = "1/2/2006 15:04:05"
	t, err := time.ParseInLocation(layout, date, time.Local)
	if err != nil {
		log.Fatalf("Cannot parse %s: %v", date, err)
	}
	return t
}

func encodeTime(t time.Time) []byte {
	data := make([]byte, 5)
	data[0] = byte(t.Second())
	data[1] = byte(t.Minute())
	data[2] = byte(t.Hour())
	data[3] = byte(t.Day())
	mm := byte(t.Month())
	data[0] |= (mm >> 2) << 6
	data[1] |= (mm & 0x3) << 6
	data[4] = byte(t.Year() - 2000)
	return data
}
