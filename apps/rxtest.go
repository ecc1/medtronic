package main

import (
	"log"

	"github.com/ecc1/cc1100"
)

const (
	DefaultFrequency = 916600000
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.LUTC)
	r, err := cc1100.Open()
	if err != nil {
		log.Fatal(err)
	}
	err = r.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = r.SetFrequency(DefaultFrequency)
	if err != nil {
		log.Fatal(err)
	}
	for packet := range r.Incoming() {
		log.Printf("% X (RSSI = %d)\n", packet.Data, packet.Rssi)
	}
}
