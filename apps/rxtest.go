package main

import (
	"log"

	"github.com/ecc1/rfm69"
)

const (
	DefaultFrequency = 916600000
)

func main() {
	r, err := rfm69.Open()
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
