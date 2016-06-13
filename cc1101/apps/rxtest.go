package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ecc1/medtronic/cc1101"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s frequency", os.Args[0])
	}
	frequency := getFrequency(os.Args[1])
	log.Printf("setting frequency to %d", frequency)
	r := cc1101.Open()
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	r.Init(frequency)
	for r.Error() == nil {
		data, rssi := r.Receive(time.Hour)
		log.Printf("% X (RSSI = %d)", data, rssi)
	}
}

func getFrequency(s string) uint32 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}
	if 860.0 <= f && f <= 920.0 {
		return uint32(f * 1000000.0)
	}
	if 860000000.0 <= f && f <= 920000000.0 {
		return uint32(f)
	}
	log.Fatalf("%s: invalid pump frequency", s)
	panic("unreachable")
}
