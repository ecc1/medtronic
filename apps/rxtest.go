package main

import (
	"log"
	"os"
	"strconv"

	"github.com/ecc1/cc1101"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s frequency\n", os.Args[0])
	}
	frequency := getFrequency(os.Args[1])
	log.Printf("setting frequency to %d\n", frequency)
	r, err := cc1101.Open()
	if err != nil {
		log.Fatal(err)
	}
	err = r.Init(frequency)
	if err != nil {
		log.Fatal(err)
	}
	for packet := range r.Incoming() {
		log.Printf("% X (RSSI = %d)\n", packet.Data, packet.Rssi)
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
	log.Fatalf("%s: invalid pump frequency\n", s)
	panic("unreachable")
}
