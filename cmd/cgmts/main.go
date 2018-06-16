package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	var data []byte
	for _, s := range os.Args[1:] {
		b, err := strconv.ParseUint(s, 16, 8)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, byte(b))
	}
	t := decodeCGMTime(data)
	fmt.Println(t)
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
