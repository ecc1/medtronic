package main

import (
	"log"

	"github.com/ecc1/cc1101"
)

func main() {
	r, err := cc1101.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Resetting radio")
	err = r.Reset()
	if err != nil {
		log.Fatal(err)
	}
	r.DumpRF()

	freq := uint32(916600000)
	log.Println("")
	log.Printf("Initializing radio to %d MHz", freq)
	err = r.InitRF(freq)
	if err != nil {
		log.Fatal(err)
	}
	r.DumpRF()

	log.Println("")
	freq += 500000
	log.Printf("Changing frequency to %d", freq)
	err = r.SetFrequency(freq)
	if err != nil {
		log.Fatal(err)
	}
	r.DumpRF()
}
