package main

import (
	"log"

	"github.com/ecc1/medtronic/cc1101"
)

func main() {
	r := cc1101.Open()
	if r.Error() != nil {
		log.Fatal(r.Error())
	}

	log.Printf("Resetting radio")
	r.Reset()
	r.DumpRF()

	freq := uint32(916600000)
	log.Println("")
	log.Printf("Initializing radio to %d MHz", freq)
	r.InitRF(freq)
	r.DumpRF()

	log.Println("")
	freq += 500000
	log.Printf("Changing frequency to %d", freq)
	r.SetFrequency(freq)
	r.DumpRF()
}
