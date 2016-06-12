package main

import (
	"log"

	"github.com/ecc1/rfm69"
)

func main() {
	r := rfm69.Open()
	if r.Error() != nil {
		log.Fatal(r.Error())
	}

	log.Printf("Resetting radio")
	r.Reset()
	dumpRF(r)

	freq := uint32(916600000)
	log.Println("")
	log.Printf("Initializing radio to %d MHz", freq)
	r.InitRF(freq)
	dumpRF(r)

	log.Println("")
	freq += 500000
	log.Printf("Changing frequency to %d", freq)
	r.SetFrequency(freq)
	dumpRF(r)

	bw := uint32(100000)
	log.Println("")
	log.Printf("Changing channel bandwidth to %d Hz", bw)
	r.SetChannelBw(bw)
	dumpRF(r)

	log.Println("")
	log.Printf("Sleeping")
	r.Sleep()
	dumpRF(r)
}

func dumpRF(r *rfm69.Radio) {
	if r.Error() != nil {
		log.Fatal(r.Error())
	}
	log.Printf("Mode: %s", r.State())
	log.Printf("Frequency: %d Hz", r.Frequency())
	mod := r.ReadModulationType()
	switch mod {
	case rfm69.ModulationTypeFSK:
		log.Printf("Modulation type: FSK")
	case rfm69.ModulationTypeOOK:
		log.Printf("Modulation type: OOK")
	default:
		log.Panicf("Unknown modulation mode %X", mod)
	}
	log.Printf("Bitrate: %d baud", r.ReadBitrate())
	log.Printf("Channel BW: %d Hz", r.ChannelBw())
}
