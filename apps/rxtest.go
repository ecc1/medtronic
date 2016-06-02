package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ecc1/cc1100"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s frequency\n", os.Args[0])
	}
	frequency := getFrequency(os.Args[1])
	log.Printf("setting frequency to %d\n", frequency)
	r, err := cc1100.Open()
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
	MHz := 0.0
	n, err := fmt.Sscanf(s, "%f", &MHz)
	if err == nil && n == 1 && 860.0 <= MHz && MHz <= 920.0 {
		return uint32(MHz * 1000000.0)
	}
	Hz := uint32(0)
	n, err = fmt.Sscanf(s, "%d", &Hz)
	if err == nil && n == 1 && 860000000 <= Hz && Hz <= 920000000 {
		return Hz
	}
	if err != nil {
		log.Fatalf("Argument (%s): %v\n", s, err)
	}
	log.Fatalf("Argument (%s) should be the pump frequency in MHz or Hz\n", s)
	panic("unreachable")
}
